package cmd

import (
	"github.com/spf13/cobra"
)

// configureCmd returns configure cobra Command.
func toolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Tools mostly used for debugging",
	}

	return cmd
}
