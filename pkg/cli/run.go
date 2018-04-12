package cli

import (
	"fmt"

	"github.com/containerum/chkit/pkg/configuration"

	"github.com/containerum/chkit/pkg/cli/login"
	"github.com/containerum/chkit/pkg/cli/mode"
	"github.com/containerum/chkit/pkg/cli/prerun"

	"github.com/containerum/chkit/pkg/chkitErrors"
	. "github.com/containerum/chkit/pkg/context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// Version -- chkit version
	Version = "3.0.0-alpha"
)

const (
	// FlagConfigFile -- context config data key
	FlagConfigFile = "config"
	// FlagAPIaddr -- API address context key
	FlagAPIaddr = "apiaddr"
)

const (
	// ErrFatalError -- unrecoverable fatal error
	ErrFatalError chkitErrors.Err = "fatal error"
)

var runContext = struct {
	ConfigFile    string
	APIaddr       string
	Username      string
	Pass          string
	DebugRequests bool
}{}

var Root = &cobra.Command{
	Use:     "chkit",
	Short:   "chkit is a containerum.io terminal client",
	Version: Version,
	PreRun: func(*cobra.Command, []string) {
		prerun.PreRun()
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Hello, %q!\nUsing %q as default namespace\n",
			Context.Client.Username,
			Context.Namespace)
		if err := mainActivity(); err != nil {
			logrus.Fatalf("error in main activity: %v", err)
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		if !Context.Changed {
			return
		}
		if err := configuration.SaveConfig(); err != nil {
			fmt.Printf("Unable to save config file: %v\n", err)
		}
	},
}

func init() {
	Context.Client.APIaddr = mode.API_ADDR
	Root.AddCommand(
		login.Command,
		Get,
		Delete,
	)
	Root.PersistentFlags().
		StringVarP(&Context.Namespace, "namespace", "n", Context.Namespace, "")
	Root.PersistentFlags().
		BoolVarP(&Context.Quiet, "quiet", "q", Context.Quiet, "quiet mode")
}
