package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/magnusmv/whoop-cli/internal/api"
	"github.com/magnusmv/whoop-cli/internal/auth"
	"github.com/magnusmv/whoop-cli/internal/config"
	"github.com/magnusmv/whoop-cli/internal/output"
)

// services groups all Whoop API service interfaces.
type services struct {
	User     api.UserService
	Cycle    api.CycleService
	Sleep    api.SleepService
	Recovery api.RecoveryService
	Workout  api.WorkoutService
}

// rootState holds shared state accessible to all subcommands.
type rootState struct {
	cfgMgr    config.Manager
	cfg       *config.Config
	svc       *services
	formatter output.Formatter
}

var state = &rootState{}

// authExemptCommands lists command use-paths that don't require a valid token.
var authExemptCommands = map[string]bool{
	"whoop auth login":  true,
	"whoop auth logout": true,
	"whoop auth setup":  true,
	"whoop auth status": true,
	"whoop completion":  true,
	"whoop help":        true,
}

// NewRootCmd creates and returns the root cobra command.
func NewRootCmd() *cobra.Command {
	var cfgDir string

	root := &cobra.Command{
		Use:   "whoop",
		Short: "Whoop CLI — access your WHOOP data from the terminal",
		Long: `whoop-cli communicates with the official WHOOP API v2.

Before using most commands, authenticate with:
  whoop auth login`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// If state has already been injected (e.g. in tests), skip setup.
			if state.svc != nil {
				return nil
			}

			// Initialise config manager.
			mgr, err := config.NewManager(cfgDir)
			if err != nil {
				return err
			}
			state.cfgMgr = mgr

			cfg, err := mgr.Load()
			if err != nil {
				return err
			}
			state.cfg = cfg

			// Determine output format from persistent flag.
			fmtStr, _ := cmd.Flags().GetString("output")
			if fmtStr == "" {
				fmtStr = cmd.InheritedFlags().Lookup("output").Value.String()
			}
			if fmtStr == "" {
				fmtStr = cfg.OutputFormat
			}
			state.formatter = output.New(fmtStr)

			// Skip token check for auth commands.
			if authExemptCommands[cmd.CommandPath()] {
				return nil
			}

			// Load token and build API client.
			store := auth.NewFileTokenStore(mgr.Dir())
			if !store.Exists() {
				return fmt.Errorf("not authenticated — run: whoop auth login")
			}
			token, err := store.Load()
			if err != nil {
				return err
			}

			ts := auth.NewTokenSource(cmd.Context(), cfg, token, store)
			client := api.NewClient(cmd.Context(), ts)
			state.svc = &services{
				User:     api.NewUserService(client),
				Cycle:    api.NewCycleService(client),
				Sleep:    api.NewSleepService(client),
				Recovery: api.NewRecoveryService(client),
				Workout:  api.NewWorkoutService(client),
			}
			return nil
		},
	}

	root.PersistentFlags().StringVar(&cfgDir, "config-dir", "", "config directory (default: ~/.whoop-cli)")
	root.PersistentFlags().StringP("output", "o", "table", "output format: table|json")
	root.PersistentFlags().Bool("all", false, "fetch all pages (list commands only)")

	root.AddCommand(newAuthCmd())
	root.AddCommand(newUserCmd())
	root.AddCommand(newCycleCmd())
	root.AddCommand(newSleepCmd())
	root.AddCommand(newRecoveryCmd())
	root.AddCommand(newWorkoutCmd())

	return root
}

// Execute is the main entry point called from main.go.
func Execute() {
	ctx := context.Background()
	if err := NewRootCmd().ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// listFlags attaches --limit, --start, --end to a command.
func addListFlags(cmd *cobra.Command) {
	cmd.Flags().Int("limit", 25, "number of records per page (max 25)")
	cmd.Flags().String("start", "", "start date filter (YYYY-MM-DD or RFC3339)")
	cmd.Flags().String("end", "", "end date filter (YYYY-MM-DD or RFC3339)")
}

// buildListParams reads list flags and returns a ListParams.
func buildListParams(cmd *cobra.Command) api.ListParams {
	limit, _ := cmd.Flags().GetInt("limit")
	startStr, _ := cmd.Flags().GetString("start")
	endStr, _ := cmd.Flags().GetString("end")

	p := api.ListParams{Limit: limit}
	if startStr != "" {
		if t := parseDate(startStr); t != nil {
			p.Start = t
		}
	}
	if endStr != "" {
		if t := parseDate(endStr); t != nil {
			p.End = t
		}
	}
	return p
}

// isFetchAll returns true if --all flag was set.
func isFetchAll(cmd *cobra.Command) bool {
	all, _ := cmd.Flags().GetBool("all")
	if !all {
		all = cmd.InheritedFlags().Lookup("all").Value.String() == "true"
	}
	return all
}
