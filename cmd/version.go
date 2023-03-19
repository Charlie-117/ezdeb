/*
Copyright © 2023 Tony

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of ezdeb",
	Long: `Print the version of ezdeb`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version 1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
