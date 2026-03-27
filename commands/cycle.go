package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newCycleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cycle",
		Short: "View physiological cycles",
	}
	cmd.AddCommand(newCycleListCmd())
	cmd.AddCommand(newCycleGetCmd())
	cmd.AddCommand(newCycleSleepCmd())
	return cmd
}

func newCycleListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List physiological cycles",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := buildListParams(cmd)
			all, _ := cmd.InheritedFlags().GetBool("all")

			if all {
				records, err := state.svc.Cycle.ListAll(cmd.Context(), params)
				if err != nil {
					return err
				}
				return state.formatter.Format(cmd.OutOrStdout(), records)
			}

			page, err := state.svc.Cycle.List(cmd.Context(), params)
			if err != nil {
				return err
			}
			if page.NextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Note: more pages available. Use --all to fetch everything.\n")
			}
			return state.formatter.Format(cmd.OutOrStdout(), page)
		},
	}
	addListFlags(cmd)
	return cmd
}

func newCycleGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a cycle by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid cycle ID %q: must be a number", args[0])
			}
			cycle, err := state.svc.Cycle.Get(cmd.Context(), id)
			if err != nil {
				return err
			}
			return state.formatter.Format(cmd.OutOrStdout(), cycle)
		},
	}
}

func newCycleSleepCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sleep <cycleId>",
		Short: "Get sleep data for a cycle",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid cycle ID %q: must be a number", args[0])
			}
			sleep, err := state.svc.Cycle.GetSleep(cmd.Context(), id)
			if err != nil {
				return err
			}
			return state.formatter.Format(cmd.OutOrStdout(), sleep)
		},
	}
}
