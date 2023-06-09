/*
Copyright © 2023 Tony

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"net/http"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
)

const (
	Green = "\033[32m"
	Red = "\033[31m"
	Cyan = "\033[36m"
	Blue = "\033[34m"
	Magneta = "\033[35m"
	Yellow = "\033[33m"
	Reset = "\033[0m"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ezdeb",
	Short: "Manage .deb packages with ease",
	Long: `ezdeb is a tool to manage .deb packages sourced from GitHub and other websites.`,
	Run: func(cmd *cobra.Command, args []string) {
		if (len(args) == 0) {
			fmt.Println("No arguments provided. Run ezdeb --help for more information.")
		}
	 },
}

func checkAppUpdate(url string, content string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	fileContent := strings.Trim(string(body), "\n")

	return fileContent == content, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	// check if application is up to date
	// compare commit hash from version file in repository
	// if different then show msg alerting user to update
	check, err := checkAppUpdate("https://gitlab.com/Charlie-117/ezdeb/-/raw/master/release/version", "90985c299a5f5e28a44e7f7b7a3d68c5118cb5ed")
	if err != nil {
		fmt.Println(Red, "\n\nError checking for App update: " + err.Error() + Reset)
	}
	if !check {
		fmt.Println(Yellow, "\n\n******\n\nAn update is available for EZDEB, please refer to guide for upgrading.\n\n******", Reset)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get home directory: %v\n", err)
		return
	}

	// check if packageList is updated within 24 hrs or not
	listPath := filepath.Join(homeDir, ".ezdeb", "pkglist.json")
	fileInfo, err := os.Stat(listPath)

	listModTime := fileInfo.ModTime()
	listAge := time.Since(listModTime)

	if listAge > 24 * time.Hour {
		fmt.Println(Yellow, "\n\n******\n\nPackage list is older than 24 hours. Run 'ezdeb sync' to update the package list.\n\n******", Reset)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configPath := filepath.Join(home, ".ezdeb")

	// Search config in home/.ezdeb directory with name "pkglist" (without extension).
	viper.AddConfigPath(configPath)
	viper.SetConfigType("json")
	viper.SetConfigName("pkglist")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	// else print a msg and sync it
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(Yellow, "******\n\nPackage list does not exist.\nSyncing package list from repository.\n\n******", Reset)
		syncCmd.Run(rootCmd, []string{})
		panic(fmt.Errorf("Run the command again, if it doesn't work then contact us with the debug message"))
	}
}

// logger function
func InitLogger() (*logrus.Logger, error) {
	// create logger object
	logger := logrus.New()

	// get home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	logFile := filepath.Join(homeDir, ".ezdeb", "ezdeb.log")
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// set log output to file
	logger.SetOutput(f)

	return logger, nil
}