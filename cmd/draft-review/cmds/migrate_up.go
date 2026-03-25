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

type MigrateUpCommand struct {
	*cmds.CommandDescription
}

var _ cmds.BareCommand = &MigrateUpCommand{}

func NewMigrateUpCommand() (*MigrateUpCommand, error) {
	sqlSection, err := draftconfig.NewSQLConnectionSection()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create sql section")
	}

	commandSettingsSection, err := cli.NewCommandSettingsSection()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create command settings section")
	}

	description := cmds.NewCommandDescription(
		"up",
		cmds.WithShort("Apply embedded PostgreSQL migrations"),
		cmds.WithLong(`Apply embedded Draft Review PostgreSQL migrations to the configured database.`),
		cmds.WithSections(sqlSection, commandSettingsSection),
	)

	return &MigrateUpCommand{CommandDescription: description}, nil
}

func (c *MigrateUpCommand) Run(ctx context.Context, parsedValues *values.Values) error {
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

	return db.Migrate(ctx)
}
