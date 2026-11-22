package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/verdverm/veg/pkg/verinfo"
	"github.com/verdverm/veg/pkg/container"
)

const versionMessage = `veg - the high code framework

Version:     %s
Commit:      %s

BuildDate:   %s
GoVersion:   %s
CueVersion:  %s
OS_Arch:     %s_%s
ConfigDir:   %s
CacheDir:    %s
Containers:  %s

Author:      verdverm
License:     Apache-2.0
Homepage:    https://github.com/verdverm/veg
GitHub:      https://github.com/verdverm/veg
`

var VersionLong = `Print the build version for veg`

var VersionCmd = &cobra.Command{

	Use: "version",

	Aliases: []string{
		"ver",
	},

	Short: "print the version",

	Long: VersionLong,

	Run: func(cmd *cobra.Command, args []string) {

		configDir, _ := os.UserConfigDir()
		cacheDir, _ := os.UserCacheDir()

		err := container.InitClient()
		if err != nil {
			fmt.Println(err)
		}

		rt, err := container.GetVersion()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf(
			versionMessage,
			verinfo.Version,
			verinfo.Commit,
			verinfo.BuildDate,
			verinfo.GoVersion,
			verinfo.CueVersion,
			verinfo.BuildOS,
			verinfo.BuildArch,
			filepath.Join(configDir,"veg"),
			filepath.Join(cacheDir,"veg"),
			rt,
		)
	},
}

func init() {
	help := VersionCmd.HelpFunc()
	usage := VersionCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		if VersionCmd.Name() == cmd.Name() {
		}
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		if VersionCmd.Name() == cmd.Name() {
		}
		return usage(cmd)
	}
	VersionCmd.SetHelpFunc(thelp)
	VersionCmd.SetUsageFunc(tusage)

}
