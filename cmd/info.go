/*
Copyright © 2023 Tony

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show information about a particular package",
	Long: `Usage: ezdeb info <package_name>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(Red, "Please provide a package name", Reset)
			return
		}

		if len(args) > 1 {
			fmt.Println(Red, "Please provide only one package name", Reset)
			return
		}

		pkgName := args[0]

		packages := viper.Get("packages").([]interface{})
		for _, pkg := range packages {
			pkgMap := pkg.(map[string]interface{})
			if pkgMap["name"].(string) == pkgName {
				fmt.Println("Package name: ", pkgMap["name"])
				fmt.Println("Package description: ", pkgMap["description"])
				fmt.Println("Package source: ", pkgMap["source"])
				if pkgMap["source"] == "github" {
					fmt.Println("Package Repository: ", pkgMap["ghuser"], "/", pkgMap["ghrepo"])
				}
				if pkgMap["source"] == "website" {
					fmt.Println("Package link: ", pkgMap["link"])
				}
				if isInstalled(pkgName) {
					fmt.Println("Installed: Yes")
				} else {
					fmt.Println("Installed: No")
				}
				if held, err := isHeldPkg(pkgName); err == nil && held {
					fmt.Println("Held: Yes")
				} else {
					fmt.Println("Held: No")
				}
				return
			}
		}

		fmt.Println("Package not found")
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
