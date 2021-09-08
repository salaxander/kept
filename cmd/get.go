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

	"github.com/salaxander/kept/pkg/kep"
	"github.com/salaxander/kept/pkg/util"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a KEP",
	Long:  `Get a KEP by providing an individual KEP number.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pterm.Println("")
		getSpinner, _ := pterm.DefaultSpinner.Start("Getting KEP...")
		k := kep.Get(args[0])
		if open {
			util.Openbrowser(k.URL)
			return
		}

		headerStyle := pterm.NewStyle(pterm.FgLightCyan, pterm.Bold)
		pterm.DefaultSection.WithStyle(headerStyle).WithIndentCharacter("\u2638\ufe0f ").Printfln("KEP %s", k.IssueNumber)
		pterm.Printfln("Title: %s", k.Title)
		if k.LatestMilestone != "" {
			pterm.Printfln("Milestone: %s", k.LatestMilestone)
		}
		if k.SIG != "" {
			pterm.Println("SIG: %s", k.SIG)
		}
		if k.Stage != "" {
			pterm.Println("Stage: %s", k.Stage)
		}
		if k.Tracked {
			pterm.Println("Tracked: yes")
		}
		pterm.Printfln("URL: %s", k.URL)

		pterm.Println("")
		getSpinner.Success("Found KEP!")
	},
}

var open bool

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVarP(&open, "open", "o", false, "Open the KEP in your default web browser.")
}
