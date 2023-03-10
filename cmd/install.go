/*
Copyright © 2023 Tony

*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/google/go-github/v50/github"
	"github.com/schollz/progressbar/v3"
)

// Create variables for global usage
var ghuser string
var ghrepo string
var pkgurl string
var debName string
var debSize int64

func isInstalled(packageName string) bool {
	// check if the package is installed in the system
	// return true if installed
	// return false if not installed
	isPkgInstalled, _ := exec.Command("dpkg", "-s", packageName).CombinedOutput()

		if strings.Contains(string(isPkgInstalled), "Status: install ok installed") {
			return true
		}

		return false
}

func searchPkgDetails(pkgName string) bool {
	// search for the package in the config file
	// if found return true and set the global variables
	// return false if not found
	packages := viper.Get("packages").([]interface{})
	for _, pkg := range packages {
		pkgMap := pkg.(map[string]interface{})
		if pkgMap["name"].(string) == pkgName {
			if pkgMap["source"].(string) == "github" {
				ghuser = pkgMap["ghuser"].(string)
				ghrepo = pkgMap["ghrepo"].(string)
				return true
			} else if pkgMap["source"].(string) == "website" {
				pkgurl = pkgMap["url"].(string)
				return true
			}
		}
	}
	return false
}

func fetchGithubRelease(ghuser string, ghrepo string) (string, error) {
	// download the latest .deb file from github release
	client := github.NewClient(nil)
	ctx := context.Background()

	// get latest release
	ghRelease, _, err := client.Repositories.GetLatestRelease(ctx, ghuser, ghrepo)
	if err != nil {
		return "", fmt.Errorf("failed to get latest release: %v", err)
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
		return "", fmt.Errorf("no .deb file asset found in release")
	}

	// download deb file
	resp, err := http.Get(asset.GetBrowserDownloadURL())
	if err != nil {
		return "", fmt.Errorf("failed to get download link: %v", err)
	}
	defer resp.Body.Close()

	// create progress bar and set it to the number of bytes downloaded
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	// Create os.TempDir()/ezdeb if it doesn't exist
	if _, err := os.Stat(filepath.Join(os.TempDir(), "ezdeb")); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(os.TempDir(), "ezdeb"), 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create temporary directory: %v", err)
		}
	}

	// create temporary file
	debFileLoc := filepath.Join(os.TempDir(), "ezdeb", asset.GetName())
	f, err := os.Create(debFileLoc)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer f.Close()

	// write content to file
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to temporary file: %v", err)
	}

	// grab the downloaded file size
	fileInfo, err := f.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}
	debSize = fileInfo.Size()

	// Return the location of the downloaded file
	return debFileLoc, nil

}

func fetchPackage(url string) (string, error) {
	// download the package from the url

	// if the url is a dynamic url i.e it keeps changing the .deb name then
	// we need to search for the package in the page and get the url of the .deb file
	if !strings.Contains(url, ".deb") {
		pageResp, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("failed to get download link: %v", err)
		}
		defer pageResp.Body.Close()

		// Read the response body into a buffer
		body, err := ioutil.ReadAll(pageResp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to get download link: %v", err)
		}

		// Find the URL of the .deb package
		re := regexp.MustCompile(`"([^"]*\.deb)"`)
		matches := re.FindSubmatch(body)
		if len(matches) < 2 {
			return "", fmt.Errorf("No .deb package found in response body")
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

    // Download the .deb package
    resp, err := http.Get(url)
    if err != nil {
        return "", fmt.Errorf("failed to download package: %v", err)
    }
    defer resp.Body.Close()

	// create progress bar and set it to the number of bytes downloaded
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	// Create os.TempDir()/ezdeb if it doesn't exist
	if _, err := os.Stat(filepath.Join(os.TempDir(), "ezdeb")); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(os.TempDir(), "ezdeb"), 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create temporary directory: %v", err)
		}
	}

    // Create a temporary file to store the downloaded .deb package
    debFileLoc := filepath.Join(os.TempDir(), "ezdeb", debName)
	f, err := os.Create(debFileLoc)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer f.Close()

    // Copy the contents of the downloaded .deb package to the temporary file
    _, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to write to temporary file: %v", err)
    }

	// grab the downloaded file size
	fileInfo, err := f.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}
	debSize = fileInfo.Size()

    // Return the location of the downloaded file
    return debFileLoc, nil
}

func installPackage(location string) error {
    // Run apt-get command to resolve dependencies and install deb file
    cmd := exec.Command("sudo", "apt-get", "install", "-y", location)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
    err := cmd.Run()
    if err != nil {
        return fmt.Errorf("failed to solve dependencies %v", err)
    }

    return nil
}

func storePackageDetails(packageName string, packageVersion string) error {
	// Store the package name and version in the package.json file
	// Get the home directory of the user
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dirPath := filepath.Join(homeDir, ".ezdeb", "packages")

	// Create the directory if it does not exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Create the package config file
	filePath := filepath.Join(dirPath, packageName+".json")

	// insert info into the package.json file
	pkgConfig := viper.New()
	pkgConfig.Set("name", packageName)
	pkgConfig.Set("version", packageVersion)
	pkgConfig.Set("size", debSize)
	err = pkgConfig.WriteConfigAs(filePath)
	if err != nil {
		return err
	}

	return nil
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a package",
	Long: `Usage: ezdeb install <package_name>`,
	Run: func(cmd *cobra.Command, args []string) {

		/*

		Ensure argument is provided
		check if package is already installed
		check if package exists in the list and fetch details
		check if source is gh or url
		download gh release or url file
		install package with apt
		on successful installation store pkg details, separate file for every pkg
		show success or error msg

		*/

		if (len(args) < 1) {
			fmt.Println("Please provide a package name")
			return
		}

		for _, pkg := range args {

			fmt.Printf("\n\nInstalling package %v\n", pkg)

			if isInstalled(pkg) {
				fmt.Printf("\n\nPackage %v is already installed\n", pkg)
				continue
			}

			if !searchPkgDetails(pkg) {
				fmt.Printf("\n\nPackage %v not found\n", pkg)
				continue
			}

			if ghuser != "" && ghrepo != "" {
				if location, err := fetchGithubRelease(ghuser, ghrepo); err == nil {
					if err = installPackage(location); err == nil {
						if err = storePackageDetails(pkg, debName); err == nil {
							fmt.Printf("\n\nPackage %v installed successfully\n", pkg)
						} else {
							fmt.Printf("\n\nPackage %v successfully installed but not logged", pkg)
							continue
						}
					} else {
						fmt.Printf("\n\nFailed to install package %v\n", pkg)
						continue
					}
				} else {
					fmt.Printf("\n\nFailed to fetch package %v\n", pkg)
					continue
				}
			} else if pkgurl != "" {
				if location, err := fetchPackage(pkgurl); err == nil {
					if err = installPackage(location); err == nil {
						if err = storePackageDetails(pkg, debName); err == nil {
							fmt.Printf("\n\nPackage %v installed successfully\n", pkg)
						} else {
							fmt.Printf("\n\nPackage %v successfully installed but not logged", pkg)
							continue
						}
					} else {
						fmt.Printf("\n\nFailed to install package %v\n", pkg)
						continue
					}
				} else {
					fmt.Printf("\n\nFailed to fetch package %v\n", pkg)
					continue
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
