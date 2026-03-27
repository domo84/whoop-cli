package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/magnusmv/whoop-cli/internal/auth"
)

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
	}
	cmd.AddCommand(newAuthLoginCmd())
	cmd.AddCommand(newAuthLogoutCmd())
	cmd.AddCommand(newAuthStatusCmd())
	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate with WHOOP via OAuth2",
		Long: `Opens a browser window to complete OAuth2 authentication.

Set WHOOP_CLIENT_ID and WHOOP_CLIENT_SECRET environment variables (or add
them to ~/.whoop-cli/config.yaml) before running this command.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if state.cfg.ClientID == "" {
				return fmt.Errorf("WHOOP_CLIENT_ID is not set.\n\nGet your credentials at https://developer.whoop.com and set:\n  export WHOOP_CLIENT_ID=<your-client-id>\n  export WHOOP_CLIENT_SECRET=<your-client-secret>")
			}

			flow := auth.NewPKCEFlow(state.cfg)
			token, err := flow.Login(cmd.Context())
			if err != nil {
				return fmt.Errorf("login failed: %w", err)
			}

			store := auth.NewFileTokenStore(state.cfgMgr.Dir())
			if err := store.Save(token); err != nil {
				return fmt.Errorf("saving token: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), color.GreenString("✓ Successfully authenticated with WHOOP."))
			return nil
		},
	}
}

func newAuthLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Revoke authentication and delete stored credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			store := auth.NewFileTokenStore(state.cfgMgr.Dir())
			if !store.Exists() {
				fmt.Fprintln(cmd.OutOrStdout(), "Not authenticated.")
				return nil
			}

			token, err := store.Load()
			if err == nil && state.cfg.ClientID != "" {
				flow := auth.NewPKCEFlow(state.cfg)
				if rErr := flow.Revoke(cmd.Context(), token); rErr != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Warning: could not revoke token remotely: %v\n", rErr)
				}
			}

			if err := store.Delete(); err != nil {
				return fmt.Errorf("deleting credentials: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), color.YellowString("Logged out. Credentials removed."))
			return nil
		},
	}
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			store := auth.NewFileTokenStore(state.cfgMgr.Dir())
			if !store.Exists() {
				fmt.Fprintln(cmd.OutOrStdout(), color.RedString("✗ Not authenticated. Run: whoop auth login"))
				return nil
			}

			token, err := store.Load()
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), color.RedString("✗ Credentials file is unreadable: %v", err))
				return nil
			}

			if token.Valid() {
				fmt.Fprintf(cmd.OutOrStdout(), "%s Token is valid (expires: %s)\n",
					color.GreenString("✓"),
					token.Expiry.Format("2006-01-02 15:04:05 MST"))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s Token expired — will refresh automatically on next request.\n",
					color.YellowString("⚠"))
			}
			return nil
		},
	}
}
