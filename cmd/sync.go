/*
Copyright © 2023 Tony

*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the packageList from the remote repository",
	Long: `Sync the packageList from the remote repository
Usage: ezdeb sync`,
	Run: func(cmd *cobra.Command, args []string) {

		/*
			Download the packageList from the remote repository
		*/

		response, err := http.Get("https://gitlab.com/Charlie-117/ezdeb/-/raw/master/pkglist/pkglist.json")
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
			return
		}

		defer response.Body.Close()

		// Create the output file
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Failed to get home directory: %v\n", err)
			return
		}

		dirPath := filepath.Join(homeDir, ".ezdeb")
		filePath := filepath.Join(homeDir, ".ezdeb", "pkglist.json")

		// Create the directory if it does not exist
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err := os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				fmt.Printf("Failed to create directory: %v\n", err)
				return
			}
		}


		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Failed to create file: %v\n", err)
			return
		}
		defer file.Close()

		// Copy the contents of the response body to the output file
		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Printf("Failed to copy contents of response body to file: %v\n", err)
			return
		}

		fmt.Printf("Successfully synced the packageList from the remote repository\n")

	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
