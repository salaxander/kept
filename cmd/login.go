package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/salaxander/kepctl/pkg/auth"
	"github.com/salaxander/kepctl/pkg/util"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to GitHub with kepctl",
	Long: `Login to GitHub with kepctl. This will allow you to edit milstones and labels if you're authorized to do so, and
	prevent rate limiting from the GitHub API.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Get the device auth info from GitHub.
		resp, _ := auth.Login()

		pterm.Println(`You will find an auth code displayed below. A GitHub auth window will open in your browser
			for you to enter the code. Once you've entered the code, return to the terminal!`)

		// Display the code users need to enter in their browser.
		pterm.FgGreen.Println(resp.UserCode)

		// Open a browser window to the auth page.
		util.Openbrowser(resp.VerificationURI)

		// Poll the GitHub auth service until the user has authenticated using the code.
		// Recieve the users token back on the channel and set it in the config file.
		ch := make(chan string)
		go func() {
			auth.PollAuthRequest(resp, ch)
		}()
		authSpinner, _ := pterm.DefaultSpinner.Start("Waiting for authorization from GitHub...")
		token := <-ch
		viper.Set("authToken", token)
		viper.WriteConfig()
		authSpinner.Success("GitHub authorized")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
