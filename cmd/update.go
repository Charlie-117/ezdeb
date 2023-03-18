/*
Copyright © 2023 Tony

*/
package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/google/go-github/v50/github"
)

var pkgNames []string

func listPkgConfigs(folderPath string) error {
	// List all the package config files in the folder

    files, err := ioutil.ReadDir(folderPath)
    if err != nil {
        return fmt.Errorf("failed to read directory: %v", err)
    }

    for _, file := range files {
        if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
            fileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
            pkgNames = append(pkgNames, fileName)
        }
    }

	if len(pkgNames) == 0 {
		return fmt.Errorf("No packages installed")
	}

    return nil
}

func checkIfInstalled(pkg string) bool {
	// check if the package is installed in the system
	isPkgInstalled, _ := exec.Command("dpkg", "-s", pkg).CombinedOutput()

	if strings.Contains(string(isPkgInstalled), "Status: install ok installed") {
		return true
	}

	return false
}

func checkUpdateGh(pkg string, ghuser string, ghrepo string) (bool, error) {
	// create a GH client and check if the package names or size don't match
	// with the one in the config file of the pkg

	client := github.NewClient(nil)
	ctx := context.Background()

	// get latest release
	ghRelease, _, err := client.Repositories.GetLatestRelease(ctx, ghuser, ghrepo)
	if err != nil {
		return false, fmt.Errorf("failed to get latest release: %v", err)
	}

	// find deb file asset
	var asset *github.ReleaseAsset

	// first search for .deb file with amd64 or x86_64 in name to avoid arm builds
	for _, a := range ghRelease.Assets {
		if filepath.Ext(a.GetName()) == ".deb" && (strings.Contains(a.GetName(), "amd64") || strings.Contains(a.GetName(), "x86_64")) {
			asset = a
			debName = a.GetName()
			break
		}
	}

	// if no .deb file was found with arch in name then search for .deb files
	if asset == nil {
		for _, a := range ghRelease.Assets {
			if filepath.Ext(a.GetName()) == ".deb" {
				asset = a
				debName = a.GetName()
				break
			}
		}
	}

	// if no .deb file was found then return error
	if asset == nil {
		return false, fmt.Errorf("no .deb file asset found in release")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	dirPath := filepath.Join(homeDir, ".ezdeb", "packages")

	// Check if directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err != nil {
			return false, err
		}
	}

	// get info from the package.json file
	pkgConfig := viper.New()
	pkgConfig.SetConfigName(pkg)
	pkgConfig.SetConfigType("json")
	pkgConfig.AddConfigPath(dirPath)
	err = pkgConfig.ReadInConfig()
	if err != nil {
		return false, err
	}

	pkgSize := pkgConfig.GetInt("size")
	pkgName := pkgConfig.GetString("version")

	if (asset.GetSize() != pkgSize) || (pkgName != asset.GetName()) {
		return true, nil
	} else {
		return false, nil
	}
}

func checkUpdateUrl(pkg string, url string) (bool, error) {
	// if the url is a dynamic url i.e it keeps changing the .deb name then
	// we need to search for the package in the page and get the url of the .deb file
	if !strings.Contains(url, ".deb") {
		pageResp, err := http.Get(url)
		if err != nil {
			return false, fmt.Errorf("failed to get download link: %v", err)
		}
		defer pageResp.Body.Close()

		// Read the response body into a buffer
		body, err := ioutil.ReadAll(pageResp.Body)
		if err != nil {
			return false, fmt.Errorf("failed to get download link: %v", err)
		}

		// Find the URL of the .deb package
		re := regexp.MustCompile(`"([^"]*\.deb)"`)
		matches := re.FindSubmatch(body)
		if len(matches) < 2 {
			return false, fmt.Errorf("No .deb package found in response body")
		}

		// Construct the download URL of the .deb package
		url := fmt.Sprintf("%s/%s", url, matches[1])
		urlParts := strings.Split(url, "/")
		debName = urlParts[len(urlParts)-1]
	} else {
		// if the url is a .deb file then we just need to grab the name of the .deb file
		urlParts := strings.Split(url, "/")
		debName = urlParts[len(urlParts)-1]
	}

	// get the size of the deb file from url
	resp, err := http.Head(url)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return false, err
	}

	// check if the size of the deb file is the same as the one in the package.json file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	dirPath := filepath.Join(homeDir, ".ezdeb", "packages")

	// Check if directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err != nil {
			return false, err
		}
	}

	// get info from the package.json file
	pkgConfig := viper.New()
	pkgConfig.SetConfigName(pkg)
	pkgConfig.SetConfigType("json")
	pkgConfig.AddConfigPath(dirPath)
	err = pkgConfig.ReadInConfig()
	if err != nil {
		return false, err
	}

	pkgSize := pkgConfig.GetInt("size")

	if (size != pkgSize) {
		return true, nil
	} else {
		return false, nil
	}

}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all packages",
	Long: `Usage: ezdeb update`,
	Run: func(cmd *cobra.Command, args []string) {


		/*

		for every package config file in packages folder
		ensuure it is installed
		fetch the details
		check the source
		if gh then check deb name
		if url then check size
		if different then download
		if same then print no update

		*/

		// init logging
		logger, err := InitLogger()
		if err != nil {
			fmt.Println(err)
			return
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Failed to get user home directory")
			return
		}

		dirPath := filepath.Join(homeDir, ".ezdeb", "packages")

		err = listPkgConfigs(dirPath)
		if err != nil {
			fmt.Println("Failed to list package configs")
			return
		}

		for _, pkg := range pkgNames {
			if checkIfInstalled(pkg) {
				fmt.Printf("Checking update for %v\n", pkg)
				// fn from install.go
				if searchPkgDetails(pkg) {
					if ghuser != "" && ghrepo != "" {
						if check, err := checkUpdateGh(pkg, ghuser, ghrepo); err == nil && check {
							if cmd.Flag("check-only").Value.String() != "true" {
								if location, err := fetchGithubRelease(ghuser, ghrepo); err == nil {
									if err = installPackage(location); err == nil {
										if err = storePackageDetails(pkg, debName); err == nil {
											logger.Infof("update: %v", pkg)
											fmt.Printf("Package %v updated successfully\n\n", pkg)
											continue
										} else {
											fmt.Printf("Package %v successfully updated but not logged\n\n", pkg)
											continue
										}
									} else {
										fmt.Printf("Failed to update package %v\n\n", pkg)
										continue
									}
								} else {
									fmt.Printf("Failed to fetch package %v\n\n", pkg)
									continue
								}
							} else {
								fmt.Printf("Update available for package: %v\n\n", pkg)
							}
						} else {
							fmt.Printf("Package %v is up to date\n\n", pkg)
							continue
						}
					} else if pkgurl != "" {
						if check, err := checkUpdateUrl(pkg, pkgurl); err == nil && check {
							if cmd.Flag("check-only").Value.String() != "true" {
								if location, err := fetchPackage(pkgurl); err == nil {
									if err = installPackage(location); err == nil {
										if err = storePackageDetails(pkg, debName); err == nil {
											logger.Infof("update: %v", pkg)
											fmt.Printf("Package %v updated successfully\n\n", pkg)
										} else {
											fmt.Printf("Package %v successfully updated but not logged\n\n", pkg)
											continue
										}
									} else {
										fmt.Printf("Failed to updated package %v\n\n", pkg)
										continue
									}
								} else {
									fmt.Printf("Failed to fetch package %v\n\n", pkg)
									continue
								}
							} else {
								fmt.Printf("Update available for package: %v\n\n", pkg)
							}
						} else {
							fmt.Printf("Package %v is up to date\n\n", pkg)
							continue
						}
					} else {
						fmt.Printf("Failed to fetch details for Package %v\n\n", pkg)
						continue
					}
				} else {
					fmt.Printf("Package %v details not found\n\n", pkg)
					continue
				}
			} else {
				continue
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	updateCmd.Flags().BoolP("check-only", "c", false, "Only check for updates")
}