/*
Copyright © 2023 Tony

*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"io/ioutil"

	"github.com/spf13/cobra"
)

var heldPkgNames []string

func listHeldPkgs(dirPath string) {
	// Create the folder if it doesnt exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			fmt.Println(Red, "Failed to create held folder", Reset)
			return
		}
	}

	// List all the package config files in the folder
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(Red, "Failed to read held folder", Reset)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
            fileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
            heldPkgNames = append(heldPkgNames, fileName)
        }
	}
}

func isHeldPkg(pkg string) (bool, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(Red, "Failed to get user home directory", Reset)
		return false, err
	}

	heldDirPath := filepath.Join(homeDir, ".ezdeb", "held")

	listHeldPkgs(heldDirPath)

	for _, heldPkg := range heldPkgNames {
		if heldPkg == pkg {
			return true, nil
		}
	}
	return false, nil
}

// holdCmd represents the hold command
var holdCmd = &cobra.Command{
	Use:   "hold",
	Short: "Hold packages from updating",
	Long: `Usage: ezdeb hold <package_name>`,
	Run: func(cmd *cobra.Command, args []string) {
		// init logging
		logger, err := InitLogger()
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(args) < 1 {
			fmt.Println(Red, "Please provide a package name", Reset)
			return
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(Red, "Failed to get user home directory", Reset)
			return
		}

		installDirPath := filepath.Join(homeDir, ".ezdeb", "packages")
		heldDirPath := filepath.Join(homeDir, ".ezdeb", "held")

		err = listPkgConfigs(installDirPath)
		if err != nil {
			fmt.Println(Yellow, "Failed to list package configs", Reset)
			fmt.Println(Yellow, "Install a package first before holding...", Reset)
			return
		}

		listHeldPkgs(heldDirPath)

		for _, pkg := range args {
			// if pkg is not installed skip
			if !isInstalled(pkg) {
				fmt.Println(Red, "Package", pkg, "not installed\n", Reset)
				continue
			}
			if check, _ := isHeldPkg(pkg); check {
				fmt.Println(Red, "Package", pkg, "already held\n", Reset)
				continue
			} else {
				// if it is not held then check if config file exists in packages folder
				for _, check2 := range pkgNames {
					// if it exists in packages folder then create a file in held folder
					if pkg == check2 {
						fileName := filepath.Join(heldDirPath, pkg+".json")
						_, err := os.Create(fileName)
						if err != nil {
							fmt.Println(Red, "Failed to create held file", Reset)
							continue
						} else {
							fmt.Println(Green, "Package", pkg, "held\n", Reset)
							logger.Infof("hold: %v", pkg)
							continue
						}
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(holdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// holdCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// holdCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
