package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newWorkoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workout",
		Short: "View workouts",
	}
	cmd.AddCommand(newWorkoutListCmd())
	cmd.AddCommand(newWorkoutGetCmd())
	return cmd
}

func newWorkoutListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workouts",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := buildListParams(cmd)
			all, _ := cmd.InheritedFlags().GetBool("all")

			if all {
				records, err := state.svc.Workout.ListAll(cmd.Context(), params)
				if err != nil {
					return err
				}
				return state.formatter.Format(cmd.OutOrStdout(), records)
			}

			page, err := state.svc.Workout.List(cmd.Context(), params)
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

func newWorkoutGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a workout by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workout, err := state.svc.Workout.Get(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(cmd.OutOrStdout(), workout)
		},
	}
}
