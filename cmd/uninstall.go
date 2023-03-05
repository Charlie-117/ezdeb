/*
Copyright © 2023 Tony

*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func isInstalledU(packageName string) bool {
	// check if the package is installed in the system
	// return true if installed
	// return false if not installed
	isPkgInstalled, _ := exec.Command("dpkg", "-s", packageName).CombinedOutput()

		if strings.Contains(string(isPkgInstalled), "Status: install ok installed") {
			return true
		}

		return false
}

func searchPkgDetailsU(pkgName string) bool {
	// search for the package in the config file
	// return false if not found
	packages := viper.Get("packages").([]interface{})
	for _, pkg := range packages {
		pkgMap := pkg.(map[string]interface{})
		if pkgMap["name"].(string) == pkgName {
				return true
		}
	}
	return false
}

func uninstallPkg(pkgName string) error {
	// uninstall the package
	// return error if failed
	cmd := exec.Command("sudo", "apt", "remove", "-y", pkgName)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func deletePkgConfig(pkgName string) error {
	// delete the package config file
	// return error if failed

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dirPath := filepath.Join(homeDir, ".ezdeb", "packages")
	filePath := filepath.Join(dirPath, pkgName+".json")

	cmd := exec.Command("rm", filePath)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall a package",
	Long: `Usage: ezdeb uninstall <package_name>`,
	Run: func(cmd *cobra.Command, args []string) {

		/*

		check if package is installed
		check if package exists in the list
		if package exists in list then uninstall package with apt remove
		on successful uninstall delete the pkg config file
		show success or error msg

		*/

		if (len(args) < 1) {
			fmt.Println("Please provide a package name")
			return
		}

		for _, pkg := range args {

			fmt.Printf("\n\nUninstalling package %v\n", pkg)

			if !(isInstalledU(pkg)) {
				fmt.Printf("\n\nPackage %v is not installed\n", pkg)
				continue
			}

			if !(searchPkgDetailsU(pkg)) {
				fmt.Printf("\n\nPackage %v was not installed with ezdeb\n", pkg)
				continue
			}

			if err := uninstallPkg(pkg); err != nil {
				fmt.Printf("\n\nFailed to uninstall package %v\n", pkg)
				continue
			} else {
				if err := deletePkgConfig(pkg); err != nil {
					fmt.Printf("\n\nPackage %v successfully uninstalled but config not removed\n", pkg)
				} else {
					fmt.Printf("\n\nPackage %v successfully uninstalled\n", pkg)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uninstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
