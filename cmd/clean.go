/*
Copyright © 2023 Tony

*/
package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"os"

	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans temporary deb files",
	Long: `Cleans temporary deb files
Usage: ezdeb clean`,
	Run: func(cmd *cobra.Command, args []string) {
		// delete .deb files in os.TempDir() directory
		fmt.Println("Cleaning temporary deb files...")
		tempPath := filepath.Join(os.TempDir(), "ezdeb")

		// check if tempPath exists
		if _, err := os.Stat(tempPath); os.IsNotExist(err) {
			fmt.Println("No temporary files to delete")
			return
		}

		err := filepath.Walk(tempPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Print(Red, "Error: failed to access path", Reset)
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".deb") {
				err := os.Remove(path)
				if err != nil {
					fmt.Print(Red, "Error: failed to delete file" + path, Reset)
					return err
				}
				fmt.Println("Deleted: " + path)
			}
			return nil
		})
		if err != nil {
			fmt.Print(Red, "Error: failed to delete temporary files", Reset)
			return
		}
		fmt.Println("Done!")
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
