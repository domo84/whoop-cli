package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newRecoveryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recovery",
		Short: "View recovery scores",
	}
	cmd.AddCommand(newRecoveryListCmd())
	cmd.AddCommand(newRecoveryGetCmd())
	return cmd
}

func newRecoveryListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List recovery scores",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := buildListParams(cmd)
			all, _ := cmd.InheritedFlags().GetBool("all")

			if all {
				records, err := state.svc.Recovery.ListAll(cmd.Context(), params)
				if err != nil {
					return err
				}
				return state.formatter.Format(cmd.OutOrStdout(), records)
			}

			page, err := state.svc.Recovery.List(cmd.Context(), params)
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

func newRecoveryGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <cycleId>",
		Short: "Get recovery score for a cycle ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid cycle ID %q: must be a number", args[0])
			}
			recovery, err := state.svc.Recovery.Get(cmd.Context(), id)
			if err != nil {
				return err
			}
			return state.formatter.Format(cmd.OutOrStdout(), recovery)
		},
	}
}
