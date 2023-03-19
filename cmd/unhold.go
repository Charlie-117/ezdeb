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

func unholdPkg(pkgName string) error {
	// check if the package is held
	// return error if not held
	isHeld, err := isHeldPkg(pkgName)
	if err != nil {
		return err
	}

	if !isHeld {
		return fmt.Errorf("Package %s is not held", pkgName)
	}

	// delete the package config file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(Red, "Failed to get user home directory", Reset)
		return err
	}

	pkgPath := filepath.Join(homeDir, ".ezdeb", "held", pkgName+".json")

	err = os.Remove(pkgPath)
	if err != nil {
		fmt.Println(Red, "Failed to remove held package", Reset)
		return err
	}

	return nil
}

// unholdCmd represents the unhold command
var unholdCmd = &cobra.Command{
	Use:   "unhold",
	Short: "Unhold held packages",
	Long: `Usage: ezdeb unhold [flags] [package_name]`,
	Run: func(cmd *cobra.Command, args []string) {
		// init logging
		logger, err := InitLogger()
		if err != nil {
			fmt.Println(err)
			return
		}

		// if all flag is set then unhold all held packages
		if cmd.Flag("all").Value.String() == "true" {
			fmt.Println("Unholding all held packages\n")
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println(Red, "Error: failed to get user home directory", Reset)
				return
			}
			heldDirPath := filepath.Join(homeDir, ".ezdeb", "held")
			// delete the heldDirPath
			err = os.RemoveAll(heldDirPath)
			if err != nil {
				fmt.Println(Red, "Error: failed to remove held dir", Reset)
				return
			}
		}

		if ((cmd.Flag("all").Value.String() == "false") && (len(args) < 1)) {
			fmt.Println(Red, "Please provide a package name", Reset)
			return
		} else {
			for _, pkg := range args {
				// if pkg is not installed skip
				if !isInstalled(pkg) {
					fmt.Println(Red, "Package", pkg, "not installed\n", Reset)
					continue
				}
				// if pkg is not held skip
				if isHeld, _ := isHeldPkg(pkg); !isHeld {
					fmt.Println(Red, "Package", pkg, "not held\n", Reset)
					continue
				}
				// unhold pkg
				err := unholdPkg(pkg)
				if err != nil {
					fmt.Println(Red, "Error: failed to unhold package\n", Reset)
					return
				}
				fmt.Println(Green, "Package", pkg, "unheld\n", Reset)
				logger.Infof("unhold: %v", pkg)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(unholdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unholdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unholdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	unholdCmd.Flags().BoolP("all", "a", false, "Unhold all held packages")
}
