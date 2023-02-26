/*
Copyright © 2023 Tony

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available packages",
	Long: `Usage: ezdeb list`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available packages:\n")

		count := 0
		packages := viper.Get("packages").([]interface{})
		for _, pkg := range packages {
			pkgMap := pkg.(map[string]interface{})

			fmt.Println(pkgMap["name"], " - ", pkgMap["description"], "\n")
			count++
		}

		fmt.Println("\nTotal number of packages:", count)

		// TODO: 
			// add functionality to list only installed packages
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
}
