package cmd

import (
	"bufio"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Open session and set up token",
	Run: func(cmd *cobra.Command, args []string) {
		var email string
		if !cmd.Flag("login").Changed {
			np.FEEDBACK.Print("Enter your email: ")
			var err error
			email, err = bufio.NewReader(os.Stdin).ReadString('\n')
			email = strings.TrimRight(email, "\r\n")
			exitOnErr(err)
		} else {
			email = cmd.Flag("login").Value.String()
		}
		var password string
		if !cmd.Flag("password").Changed {
			np.FEEDBACK.Print("Enter your password: ")
			passwordB, err := terminal.ReadPassword(int(syscall.Stdin))
			exitOnErr(err)
			password = string(passwordB)
		} else {
			password = cmd.Flag("password").Value.String()
		}
		exitOnErr(ChkitClient.Login(email, password))
		np.FEEDBACK.Printf("Succesfull login!")
		exitOnErr(ChkitClient.SaveTokens())
	},
}

func init() {
	loginCmd.PersistentFlags().StringP("login", "l", "", "User login (email)")
	loginCmd.PersistentFlags().StringP("password", "p", "", "User password")
}
