package main

import (
	"os"

	draftreviewcmds "github.com/go-go-golems/draft-review/cmd/draft-review/cmds"
	draftreviewdoc "github.com/go-go-golems/draft-review/cmd/draft-review/doc"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds/logging"
	"github.com/go-go-golems/glazed/pkg/help"
	help_cmd "github.com/go-go-golems/glazed/pkg/help/cmd"
	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "draft-review",
	Short:   "Draft Review backend and developer tooling",
	Long:    "Draft Review backend and developer tooling for the local Go, PostgreSQL, and Keycloak-based development workflow.",
	Version: version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return logging.InitLoggerFromCobra(cmd)
	},
}

func main() {
	cobra.CheckErr(logging.AddLoggingSectionToRootCommand(rootCmd, "draft-review"))

	helpSystem := help.NewHelpSystem()
	cobra.CheckErr(draftreviewdoc.AddDocToHelpSystem(helpSystem))
	help_cmd.SetupCobraRootCommand(helpSystem, rootCmd)

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
		Long:  "Database migration commands for applying and inspecting the embedded Draft Review PostgreSQL schema.",
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
		Long:  "Seed data commands for loading deterministic local Draft Review development fixtures.",
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
