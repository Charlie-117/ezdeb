/*
Copyright © 2023 Tony

*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var searchRsltCount = 0

func searchName(searchTerm string, packages []interface{}) bool {
	pkgFound := false
	for _, pkg := range packages {
		pkgMap := pkg.(map[string]interface{})
		if strings.Contains(strings.ToLower(pkgMap["name"].(string)), searchTerm) {
			searchRsltCount++
			fmt.Println(Cyan, pkgMap["name"], Reset, " - ", pkgMap["description"], "\n")
			pkgFound = true
		}
	}
	return pkgFound
}

func searchDesc(searchTerm string, packages []interface{}) bool {
	pkgFound := false
	for _, pkg := range packages {
		pkgMap := pkg.(map[string]interface{})
		if strings.Contains(strings.ToLower(pkgMap["description"].(string)), searchTerm) {
			searchRsltCount++
			fmt.Println(Cyan, pkgMap["name"], Reset, " - ", pkgMap["description"], "\n")
			pkgFound = true
		}
	}
	return pkgFound
}


// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a package",
	Long: `Search for a package
Usage: ezdeb search <search_term>`,
	Run: func(cmd *cobra.Command, args []string) {
		// exit if no args is provided
		if len(args) < 1 {
			fmt.Println(Red, "Please provide a search term", Reset)
			return
		}

		fmt.Println("Searching through packages...\n")
		packages := viper.Get("packages").([]interface{})
		pkgFound := false

		// search for the the whole args as one in package names
		searchTerm := strings.ToLower(strings.Join(args, " "))
		if !(searchName(searchTerm, packages)) {
			// if no package is found then search for each arg in package names
			for _, searchTerm := range args {
				searchTerm = strings.ToLower(searchTerm)
				if (searchName(searchTerm, packages)) {
					pkgFound = true
				}
			}
			if !pkgFound {
				// if no package is found then search for the whole args as one in package descriptions
				searchTerm := strings.ToLower(strings.Join(args, " "))
				if !(searchDesc(searchTerm, packages)) {
					// if no package is found then search for each arg in package descriptions
					for _, searchTerm := range args {
						searchTerm = strings.ToLower(searchTerm)
						searchDesc(searchTerm, packages)
					}
				}
			}
		}

		fmt.Println(Green, "Found", searchRsltCount, "results", Reset)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
