package commands

import (
	"github.com/spf13/cobra"
)

func newUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "View user profile and measurements",
	}
	cmd.AddCommand(newUserProfileCmd())
	cmd.AddCommand(newUserMeasurementsCmd())
	return cmd
}

func newUserProfileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "profile",
		Short: "Get your WHOOP profile (name, email)",
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, err := state.svc.User.GetProfile(cmd.Context())
			if err != nil {
				return err
			}
			return state.formatter.Format(cmd.OutOrStdout(), profile)
		},
	}
}

func newUserMeasurementsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "measurements",
		Short: "Get your body measurements (height, weight, max HR)",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := state.svc.User.GetBodyMeasurement(cmd.Context())
			if err != nil {
				return err
			}
			return state.formatter.Format(cmd.OutOrStdout(), m)
		},
	}
}
