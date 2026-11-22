package cmdchat

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	libagentcmd "github.com/verdverm/veg/pkg/agent/cmd"

	"github.com/verdverm/veg/cmd/veg/flags"

)

var infoLong = `get info for session`

func InfoRun(args []string) (err error) {

	err = libagentcmd.ChatInfo(args, flags.RootPflags, flags.AgentPflags, flags.Agent__ChatPflags)

	return err
}

var InfoCmd = &cobra.Command{

	Use: "info [...target] [% ...cue]",

	Short: "get info for session",

	Long: infoLong,

	Run: func(cmd *cobra.Command, args []string) {


		var err error

		// Argument Parsing

		err = InfoRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := InfoCmd.HelpFunc()
	ousage := InfoCmd.UsageFunc()

	help := func(cmd *cobra.Command, args []string) {


		if extra(cmd) {
			return
		}
		ohelp(cmd, args)
	}
	usage := func(cmd *cobra.Command) error {
		if extra(cmd) {
			return nil
		}
		return ousage(cmd)
	}

	thelp := func(cmd *cobra.Command, args []string) {
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		return usage(cmd)
	}
	InfoCmd.SetHelpFunc(thelp)
	InfoCmd.SetUsageFunc(tusage)

}
