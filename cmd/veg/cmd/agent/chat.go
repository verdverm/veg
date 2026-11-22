package cmdagent

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/spf13/cobra"

	libagentcmd "github.com/verdverm/veg/pkg/agent/cmd"

	cmdchat "github.com/verdverm/veg/cmd/veg/cmd/agent/chat"

	"github.com/verdverm/veg/cmd/veg/flags"

)

var chatLong = `chat with an agent`

func init() {

	flags.SetupAgent__ChatPflags(ChatCmd.PersistentFlags(), &(flags.Agent__ChatPflags))

}

func ChatRun(args []string) (err error) {

	err = libagentcmd.Chat(args, flags.RootPflags, flags.AgentPflags, flags.Agent__ChatPflags)

	return err
}

var ChatCmd = &cobra.Command{

	Use: "chat [...target] [% ...cue]",

	Short: "chat with an agent",

	Long: chatLong,

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		glob := toComplete + "*"
		matches, _ := filepath.Glob(glob)
		return matches, cobra.ShellCompDirectiveDefault
	},

	Run: func(cmd *cobra.Command, args []string) {


		var err error

		// Argument Parsing

		err = ChatRun(args)
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

	ohelp := ChatCmd.HelpFunc()
	ousage := ChatCmd.UsageFunc()

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
	ChatCmd.SetHelpFunc(thelp)
	ChatCmd.SetUsageFunc(tusage)

	ChatCmd.AddCommand(cmdchat.InfoCmd)
	ChatCmd.AddCommand(cmdchat.ListCmd)
	ChatCmd.AddCommand(cmdchat.DeleteCmd)

}
