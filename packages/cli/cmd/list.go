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
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"lookahead.web.app/cli/internal/credential"
	"lookahead.web.app/cli/internal/logging"
	"lookahead.web.app/cli/internal/store"
	"lookahead.web.app/cli/internal/types"
)

func checkStringEmptyOrOnlySpaces(str string) bool {
	if strings.TrimSpace(str) == "" {
		return true
	}
	return false
}

func getLastPathOfDocId(docId string) string {
	arr := strings.Split(docId, "/")
	return arr[len(arr)-1]
}

func printWholeTodo(todo types.DataSchema) {
	fmt.Println("Todo id:", getLastPathOfDocId(todo.ID))
	color.HiCyan(todo.Title)
	if !checkStringEmptyOrOnlySpaces(todo.Content) {
		fmt.Println(todo.Content)
	}
	fmt.Println(todo.ToLastEditedHuman())
	fmt.Println()
}

func printOnlyId(todo map[string]interface{}) {
	// TODO(me): Implement print-ID-only functionality
}

// listCmd represents the show command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all notes and to-dos of the user",
	Run: func(cmd *cobra.Command, args []string) {
		//If user is logged out
		if credential.CheckIfUserLoggedIn() == false {
			logging.Error(1, "You should be logged in to run this command!!"+
				" Use `look login` to login")
		}
		entriesLimit := viper.GetInt("limitEntries")
		// -----------------------------------------
		// Implement get-by-id only
		// -----------------------------------------
		s := logging.DarkSpinner(" Fetching data")
		s.Start()
		documents, err := store.Store.GetAll()
		if err != nil {
			s.Stop()
			logging.Error(1, err.Error())
		}
		s.Stop()
		for i, todo := range documents {
			if entriesLimit > 0 && i >= entriesLimit {
				logging.Warn(
					"If you want to see more than %d entries, "+
						"please set the `limitEntries` flag in the configuration file"+
						" or use the `-l` flag in the command",
					entriesLimit,
				)
				break
			}
			printWholeTodo(todo)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	//Show this command as suggestion when these non-existent commands are used
	listCmd.SuggestFor = []string{"notes", "todos", "show"}
	listCmd.PersistentFlags().Int16P(
		"limitEntries",
		"l",
		0,
		"configure the number of entries the CLI shows on running `look list`")
	//Bind flag to viper for precedence
	viper.BindPFlag("limitEntries", listCmd.PersistentFlags().Lookup("limitEntries"))
}
