package cmds

import (
	"context"

	draftconfig "github.com/go-go-golems/draft-review/pkg/config"
	draftdb "github.com/go-go-golems/draft-review/pkg/db"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	pkgerrors "github.com/pkg/errors"
)

type SeedDevCommand struct {
	*cmds.CommandDescription
}

var _ cmds.BareCommand = &SeedDevCommand{}

func NewSeedDevCommand() (*SeedDevCommand, error) {
	sqlSection, err := draftconfig.NewSQLConnectionSection()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create sql section")
	}

	commandSettingsSection, err := cli.NewCommandSettingsSection()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create command settings section")
	}

	description := cmds.NewCommandDescription(
		"dev",
		cmds.WithShort("Seed development data"),
		cmds.WithLong(`Seed a small Draft Review development dataset into the configured database.`),
		cmds.WithSections(sqlSection, commandSettingsSection),
	)

	return &SeedDevCommand{CommandDescription: description}, nil
}

func (c *SeedDevCommand) Run(ctx context.Context, parsedValues *values.Values) error {
	sqlSettings, err := draftconfig.LoadSQLConnectionSettings(parsedValues)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to load sql connection settings")
	}

	connectionString, err := sqlSettings.ConnectionString()
	if err != nil {
		return pkgerrors.Wrap(err, "failed to compute connection string")
	}

	db, err := draftdb.Open(ctx, connectionString)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to open database")
	}
	if db == nil {
		return draftdb.ErrNotConfigured
	}
	defer db.Close()

	if err := db.Migrate(ctx); err != nil {
		return pkgerrors.Wrap(err, "failed to apply migrations before seeding")
	}

	return draftdb.SeedDev(ctx, db.Pool())
}
