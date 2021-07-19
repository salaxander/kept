/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/salaxander/kepctl/pkg/kep"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get a list of KEPs",
	Long:  `Get a list of KEPs with options to filter by tracked status and milestones.`,
	Run: func(cmd *cobra.Command, args []string) {
		pterm.Println("")
		listSpinner, _ := pterm.DefaultSpinner.Start("Fetching KEPs...")
		keps := kep.List(milestone, sig, stage, tracked)
		tableData := [][]string{{"KEP Number", "Title", "URL"}}
		for i := range keps {
			kep := []string{keps[i].IssueNumber, keps[i].Title, keps[i].URL}
			tableData = append(tableData, kep)
		}
		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		listSpinner.Success("Retrieved KEPs.")
	},
}

var milestone string
var stage string
var sig string

var all bool
var tracked bool

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&milestone, "milestone", "m", "", "Milestone to filter KEPs by.")
	listCmd.Flags().StringVarP(&stage, "stage", "", "", "Stage to filter KEPs by (alpha|beta|stable).")
	listCmd.Flags().StringVarP(&sig, "sig", "", "", "SIG to filter KEPs by.")

	listCmd.Flags().BoolVarP(&all, "all", "a", false, "Show all KEPs, including closed.")
	listCmd.Flags().BoolVarP(&tracked, "tracked", "t", false, "Filter for tracked KEPs only.")
}
