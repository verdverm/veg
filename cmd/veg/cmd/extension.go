package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/verdverm/veg/cmd/veg/flags"

	"github.com/verdverm/veg/pkg/extension"

	"github.com/verdverm/veg/pkg/runtime"

)

var extensionLong = `run the extension server`

func ExtensionPersistentPreRun(args []string) (err error) {

	err = runtime.EnsureInfra()

	return err
}

func ExtensionRun(args []string) (err error) {

	err = extension.Run(args, flags.RootPflags)

	return err
}

var ExtensionCmd = &cobra.Command{

	Use: "extension [args]",

	Short: "run the extension server",

	Long: extensionLong,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ExtensionPersistentPreRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {


		var err error

		// Argument Parsing

		err = ExtensionRun(args)
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

	ohelp := ExtensionCmd.HelpFunc()
	ousage := ExtensionCmd.UsageFunc()

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
	ExtensionCmd.SetHelpFunc(thelp)
	ExtensionCmd.SetUsageFunc(tusage)

}
