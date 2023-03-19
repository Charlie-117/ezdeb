/*
Copyright © 2023 Tony

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available packages",
	Long: `List all available packages
Usage: ezdeb list`,
	Run: func(cmd *cobra.Command, args []string) {

		count := 0

		// List only installed packages if flag is set
		if cmd.Flag("installed").Value.String() == "true" {
			fmt.Println("Listing installed packages\n")

			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println(Red, "Error: failed to get user home directory", Reset)
				return
			}
			pkgPath := filepath.Join(homeDir, ".ezdeb", "packages")

			err = filepath.Walk(pkgPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					fmt.Println(Red, "Error: failed to access path", Reset)
					return err
				}
				if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
					// trim .json suffix from file name
					fileName := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
					fmt.Println(Cyan, fileName, Reset)
					count++
				}
				return nil
			})
			if err != nil {
				fmt.Print(Red, "Error: failed to list packages", Reset)
				return
			}

			if count == 0 {
				fmt.Println("No packages installed")
				return
			}
			fmt.Println(Green, "\nTotal number of installed packages:", count, Reset)
			return
		}

		// List only held packages if flag is set
		if cmd.Flag("held").Value.String() == "true" {
			fmt.Println("Listing held packages\n")

			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println(Red, "Failed to get user home directory", Reset)
				return
			}

			heldDirPath := filepath.Join(homeDir, ".ezdeb", "held")

			listHeldPkgs(heldDirPath)

			if len(heldPkgNames) == 0 {
				fmt.Println("No held packages")
				return
			}

			for _, heldPkg := range heldPkgNames {
				fmt.Println(Cyan, heldPkg, Reset)
			}

			fmt.Println(Green, "\nTotal number of held packages:", len(heldPkgNames), Reset)
			return

		}

		fmt.Println("Available packages:\n")

		packages := viper.Get("packages").([]interface{})
		for _, pkg := range packages {
			pkgMap := pkg.(map[string]interface{})

			fmt.Println(Cyan, pkgMap["name"], Reset, " - ", pkgMap["description"], "\n")
			count++
		}

		fmt.Println(Green, "\nTotal number of packages:", count, Reset)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listCmd.Flags().BoolP("installed", "i", false, "List only installed packages")
	listCmd.Flags().BoolP("held", "l", false, "List only held packages")
}
