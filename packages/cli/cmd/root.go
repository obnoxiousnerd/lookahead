/*
Copyright © 2020 Pranav Karawale

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"lookahead.web.app/cli/internal/constants"
	"lookahead.web.app/cli/internal/store"
	"lookahead.web.app/cli/internal/version"

	"github.com/spf13/viper"
	"lookahead.web.app/cli/internal/config"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "look",
	Short:   "CLI interface of Lookahead (visit https://lookahead.web.app for more info)",
	Version: version.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SetVersionTemplate("Lookahead CLI version " + version.Version)
	if err := rootCmd.Execute(); err != nil {
		// fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lookahead.yaml)")
	//Initialize default config
	viper.SetDefault("numberOfEntries", 5)

	shouldUpdateStore := false
	//updateStoreOnTheseCommands array containing commands which will trigger
	//local store update
	updateStoreOnTheseCommands := []string{"list", "edit", "delete", "new"}

	for _, command := range updateStoreOnTheseCommands {
		if len(os.Args) >= 2 && os.Args[1] == command {
			shouldUpdateStore = true
		}
	}

	if shouldUpdateStore {
		//Sync the store
		store.Store.Sync(false)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home := constants.CONFIG_PATH

		// Search config in home directory with name ".lookahead" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".lookahead")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// This will be uncommented later for debugging purposes
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
		config.CheckConfigValidity()
	}
}
