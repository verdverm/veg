package cmd

import (
	"github.com/spf13/cobra"

	"github.com/verdverm/veg/cmd/veg/cmd/datamodel"

	"github.com/verdverm/veg/cmd/veg/flags"

)

var datamodelLong = `Data models are values or objects used in many veg processes and modules.
The "datamodel" command helps you manage them and track their change history.
At their core, they represent the models that make up your application.
The intention is to define a data model for your entire application once,
then use this source of truth to generate code from database to server to client.

Veg's schema for datamodels is minimal and flexible, allowing you to define the
shape based on your application. You can have multiple datamodels as well.
You can also control and where and how history should be tracked. This history
is included during code generation so that database migrations and functions
for converting between versions can be created.

# Examples Datamodels

-- config.cue --
package datamodel

import (
	"github.com/verdverm/veg/schemas/dm"
	"github.com/verdverm/veg/schemas/dm/fields"
)

// Track an entire oject
Config: dm.Object & {

	host: fields.String & { Default: "8080" }

	database: {
		host:   fields.String
		port:   fields.String
		dbconn: fields.String
	}
}

-- database.cue --
package datamodel

import (
	"github.com/verdverm/veg/schemas/dm/sql"
	"github.com/verdverm/veg/schemas/dm/fields"
)

// Traditional database model which maps onto tables & columns
Datamodel: sql.Datamodel & {
	// implied through definition, duplicated here for example clarity
	$veg: metadata: {
		id:   "datamodel-abc123"
		name: "MyDatamodel"
	}

	Models: {
		User: {
			Fields: {
				ID:        fields.UUID
				CreatedAt: fields.Datetime
				UpdatedAt: fields.Datetime
				DeletedAt: fields.Datetime

				email:    fields.Email
				username: fields.String
				password: fields.Password
				verified: fields.Bool
				active:   fields.Bool

				persona: fields.Enum & {
					Vals: ["guest", "user", "admin", "owner"]
					Default: "user"
				}
			}
		}
	}
}

# Example Usage   (dm is short for datamodel)

  $ veg dm list   (print known data models)
  NAME         TYPE       VERSION  STATUS  ID
  Config       object     -        ok      Config
  MyDatamodel  datamodel  -        ok      datamodel-abc123

  $ veg dm tree   (print the structure of the datamodels)

  $ veg dm diff   (prints a tree based diff of the datamodel)

  $ veg dm checkpoint -m "a message about this checkpoint"

  $ veg dm log    (prints the log of changes from latest to oldest)

  You can also use the -d & -e flags to subselect datamodels and nested values

# Learn more:
  - https://github.com/verdverm/veg/wiki/getting-started/data-layer/
  - https://github.com/verdverm/veg/wiki/data-modeling/`

func init() {

	flags.SetupDatamodelPflags(DatamodelCmd.PersistentFlags(), &(flags.DatamodelPflags))

}

var DatamodelCmd = &cobra.Command{

	Use: "datamodel",

	Aliases: []string{
		"dm",
	},

	Short: "manage, diff, and migrate your data models",

	Long: datamodelLong,
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := DatamodelCmd.HelpFunc()
	ousage := DatamodelCmd.UsageFunc()

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
	DatamodelCmd.SetHelpFunc(thelp)
	DatamodelCmd.SetUsageFunc(tusage)

	DatamodelCmd.AddCommand(cmddatamodel.CheckpointCmd)
	DatamodelCmd.AddCommand(cmddatamodel.DiffCmd)
	DatamodelCmd.AddCommand(cmddatamodel.TreeCmd)
	DatamodelCmd.AddCommand(cmddatamodel.ListCmd)
	DatamodelCmd.AddCommand(cmddatamodel.LogCmd)

}
