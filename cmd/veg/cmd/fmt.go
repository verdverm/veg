package cmd

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/verdverm/veg/cmd/veg/cmd/fmt"
	"github.com/verdverm/veg/cmd/veg/flags"

	hfmt "github.com/verdverm/veg/pkg/fmt"
)

var fmtLong = `With veg fmt, you can
  1. format any language from a single tool
  2. run formatters as api servers for IDEs and veg
  3. manage the underlying formatter containers`

func init() {

	flags.SetupFmtFlags(FmtCmd.Flags(), &(flags.FmtFlags))

}

func FmtRun(files []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = hfmt.Run(files, flags.RootPflags, flags.FmtFlags)

	return err
}

var FmtCmd = &cobra.Command{

	Use: "fmt [filepaths or globs]",

	Short: "format any code and manage the formatters",

	Long: fmtLong,

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		glob := toComplete + "*"
		matches, _ := filepath.Glob(glob)
		return matches, cobra.ShellCompDirectiveDefault
	},

	Run: func(cmd *cobra.Command, args []string) {


		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'files'")
			cmd.Usage()
			os.Exit(1)
		}

		var files []string

		if 0 < len(args) {

			files = args[0:]

		}

		err = FmtRun(files)
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

	ohelp := FmtCmd.HelpFunc()
	ousage := FmtCmd.UsageFunc()

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
	FmtCmd.SetHelpFunc(thelp)
	FmtCmd.SetUsageFunc(tusage)

	FmtCmd.AddCommand(cmdfmt.InfoCmd)
	FmtCmd.AddCommand(cmdfmt.PullCmd)
	FmtCmd.AddCommand(cmdfmt.StartCmd)
	FmtCmd.AddCommand(cmdfmt.TestCmd)
	FmtCmd.AddCommand(cmdfmt.StopCmd)

}
