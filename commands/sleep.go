package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSleepCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sleep",
		Short: "View sleep records",
	}
	cmd.AddCommand(newSleepListCmd())
	cmd.AddCommand(newSleepGetCmd())
	return cmd
}

func newSleepListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sleep records (includes naps)",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := buildListParams(cmd)
			all, _ := cmd.InheritedFlags().GetBool("all")

			if all {
				records, err := state.svc.Sleep.ListAll(cmd.Context(), params)
				if err != nil {
					return err
				}
				return state.formatter.Format(cmd.OutOrStdout(), records)
			}

			page, err := state.svc.Sleep.List(cmd.Context(), params)
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

func newSleepGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a sleep record by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sleep, err := state.svc.Sleep.Get(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(cmd.OutOrStdout(), sleep)
		},
	}
}
