package cmd

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/verdverm/veg/pkg/runtime"

	libagentcmd "github.com/verdverm/veg/pkg/agent/cmd"

	"github.com/verdverm/veg/cmd/veg/cmd/agent"

	"github.com/verdverm/veg/cmd/veg/flags"

)

var agentLong = `build, use, evaluate, and serve agentic systems`

func init() {

	flags.SetupAgentPflags(AgentCmd.PersistentFlags(), &(flags.AgentPflags))

}

func AgentPersistentPreRun(args []string) (err error) {

	err = runtime.EnsureInfra()

	return err
}

func AgentRun(args []string) (err error) {

	err = libagentcmd.Main(args, flags.RootPflags, flags.AgentPflags, flags.Agent__ChatPflags)

	return err
}

var AgentCmd = &cobra.Command{

	Use: "agent [...target] [% ...cue]",

	Short: "build, chat with, and serve agentic systems",

	Long: agentLong,

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		glob := toComplete + "*"
		matches, _ := filepath.Glob(glob)
		return matches, cobra.ShellCompDirectiveDefault
	},

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = AgentPersistentPreRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {


		var err error

		// Argument Parsing

		err = AgentRun(args)
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

	ohelp := AgentCmd.HelpFunc()
	ousage := AgentCmd.UsageFunc()

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
	AgentCmd.SetHelpFunc(thelp)
	AgentCmd.SetUsageFunc(tusage)

	AgentCmd.AddCommand(cmdagent.ListCmd)
	AgentCmd.AddCommand(cmdagent.ChatCmd)
	AgentCmd.AddCommand(cmdagent.ServeCmd)
	AgentCmd.AddCommand(cmdagent.BulkCmd)
	AgentCmd.AddCommand(cmdagent.EvalCmd)

}
