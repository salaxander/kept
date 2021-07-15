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
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/pterm/pterm"
	"github.com/salaxander/kepctl/pkg/kep"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a KEP",
	Long:  `Get a KEP by providing an individual KEP number.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		k := kep.Get(args[0])
		if open {
			openbrowser(k.URL)
			return
		}

		headerStyle := pterm.NewStyle(pterm.FgLightCyan, pterm.Bold)
		pterm.DefaultSection.WithStyle(headerStyle).WithIndentCharacter("\u2638\ufe0f ").Printfln("KEP %s", k.IssueNumber)
		pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
			{
				Level: 0,
				Text:  pterm.Sprintf("Title: %s", k.Title),
			},
			{
				Level: 0,
				Text:  pterm.Sprintf("URL: %s", k.URL),
			},
		}).Render()
	},
}

var open bool

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVarP(&open, "open", "o", false, "Open the KEP in your default web browser.")
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
