package main

import (
	"os"

	draftreviewcmds "github.com/go-go-golems/draft-review/cmd/draft-review/cmds"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "draft-review",
	Short:   "Draft Review backend and developer tooling",
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return logging.InitLoggerFromCobra(cmd)
	},
}

func main() {
	cobra.CheckErr(logging.AddLoggingSectionToRootCommand(rootCmd, "draft-review"))

	serveCmd, err := draftreviewcmds.NewServeCommand(version)
	cobra.CheckErr(err)
	cobraServeCmd, err := cli.BuildCobraCommandFromCommand(serveCmd,
		cli.WithParserConfig(cli.CobraParserConfig{
			AppName:           "draft-review",
			ShortHelpSections: []string{"default", "sql-connection", "auth", "backend"},
			MiddlewaresFunc:   cli.CobraCommandDefaultMiddlewares,
		}),
	)
	cobra.CheckErr(err)
	rootCmd.AddCommand(cobraServeCmd)

	migrateRootCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration commands",
	}

	migrateUpCmd, err := draftreviewcmds.NewMigrateUpCommand()
	cobra.CheckErr(err)
	cobraMigrateUpCmd, err := cli.BuildCobraCommandFromCommand(migrateUpCmd,
		cli.WithParserConfig(cli.CobraParserConfig{
			AppName:           "draft-review",
			ShortHelpSections: []string{"sql-connection"},
			MiddlewaresFunc:   cli.CobraCommandDefaultMiddlewares,
		}),
	)
	cobra.CheckErr(err)
	migrateRootCmd.AddCommand(cobraMigrateUpCmd)

	migrateStatusCmd, err := draftreviewcmds.NewMigrateStatusCommand()
	cobra.CheckErr(err)
	cobraMigrateStatusCmd, err := cli.BuildCobraCommandFromCommand(migrateStatusCmd,
		cli.WithParserConfig(cli.CobraParserConfig{
			AppName:           "draft-review",
			ShortHelpSections: []string{"sql-connection"},
			MiddlewaresFunc:   cli.CobraCommandDefaultMiddlewares,
		}),
	)
	cobra.CheckErr(err)
	migrateRootCmd.AddCommand(cobraMigrateStatusCmd)
	rootCmd.AddCommand(migrateRootCmd)

	seedRootCmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed data commands",
	}

	seedDevCmd, err := draftreviewcmds.NewSeedDevCommand()
	cobra.CheckErr(err)
	cobraSeedDevCmd, err := cli.BuildCobraCommandFromCommand(seedDevCmd,
		cli.WithParserConfig(cli.CobraParserConfig{
			AppName:           "draft-review",
			ShortHelpSections: []string{"sql-connection"},
			MiddlewaresFunc:   cli.CobraCommandDefaultMiddlewares,
		}),
	)
	cobra.CheckErr(err)
	seedRootCmd.AddCommand(cobraSeedDevCmd)
	rootCmd.AddCommand(seedRootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
