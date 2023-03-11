/*
Copyright © 2023 Tony

*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"bufio"
	"time"

	"github.com/spf13/cobra"
)

func formatTime(logLine string) string {
    timeString := logLine[6:25]
    t, _ := time.Parse("2006-01-02T15:04:05", timeString)
    return t.Format("2006-01-02 15:04:05")
}

func formatAction(logLine string) string {
	if strings.Contains(logLine, "uninstall") {
		return "uninstall"
	} else if strings.Contains(logLine, "\"install") {
		return "install"
	} else if strings.Contains(logLine, "update") {
		return "update"
	} else {
		return "unknown"
	}
}

func formatPkg(logLine string) string {
	if strings.Contains(logLine, "uninstall") {
		return fmt.Sprintf("%s%s", "\"", logLine[60:])
	} else if strings.Contains(logLine, "\"install") {
		return fmt.Sprintf("%s%s", "\"", logLine[58:])
	} else if strings.Contains(logLine, "update") {
		return fmt.Sprintf("%s%s", "\"", logLine[57:])
	} else {
		return "unknown"
	}
}

func readLog(action string) {
	// read the log file
	// return error if failed
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("failed to get home directory")
	}

	logFile := filepath.Join(homeDir, ".ezdeb", "ezdeb.log")

	file, err := os.Open(logFile)
	if err != nil {
		fmt.Printf("failed to open log file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// if action is all
	// read all logs
	if action == "all" {
		for scanner.Scan() {
			logLine := 	scanner.Text()
			fmt.Printf("time: %s action: %s package: %s\n", formatTime(logLine), formatAction(logLine), formatPkg(logLine))
		}
	} else if action == "install" {
		for scanner.Scan() {
			logLine := 	scanner.Text()
			if strings.Contains(logLine, "\"install") {
				fmt.Printf("time: %s action: install package: %s\n", formatTime(logLine), formatPkg(logLine))
			}
		}
	} else if action == "uninstall" {
		for scanner.Scan() {
			logLine := 	scanner.Text()
			if strings.Contains(logLine, "uninstall") {
				fmt.Printf("time: %s action: uninstall package: %s\n", formatTime(logLine), formatPkg(logLine))
			}
		}
	} else if action == "update" {
		for scanner.Scan() {
			logLine := 	scanner.Text()
			if strings.Contains(logLine, "update") {
				fmt.Printf("time: %s action: update package: %s\n", formatTime(logLine), formatPkg(logLine))
			}
		}
	}
}

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show logs",
	Long: `Usage: ezdeb logs [flags]`,
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("action").Value.String() == "" {
			readLog("all")
		} else if cmd.Flag("action").Value.String() == "install" {
			readLog("install")
		} else if cmd.Flag("action").Value.String() == "uninstall" {
			readLog("uninstall")
		} else if cmd.Flag("action").Value.String() == "update" {
			readLog("update")
		} else {
			fmt.Println("Invalid action, use -h to see available actions")
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	logsCmd.Flags().StringP("action", "a", "", "Show logs for a specific action (install, remove, update)")
}