/*
Copyright © 2023 Tony

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// clearLogsCmd represents the clearLogs command
var clearLogsCmd = &cobra.Command{
	Use:   "clearLogs",
	Short: "Clear logs",
	Long: `Clear logs
Usage: ezdeb clearLogs`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Clearings logs...")

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(Red, "failed to get home directory", Reset)
			return
		}

		logFile := filepath.Join(homeDir, ".ezdeb", "ezdeb.log")

		// replace the logs file with empty file
		err = os.Truncate(logFile, 0)
		if err != nil {
			fmt.Println(Red, "Failed to clear log file", Reset)
			return
		}

		fmt.Println("Logs cleared")
	},
}

func init() {
	rootCmd.AddCommand(clearLogsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clearLogsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clearLogsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
