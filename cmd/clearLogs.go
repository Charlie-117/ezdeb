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
}
