package cmds

import (
	"context"
	"fmt"

	draftconfig "github.com/go-go-golems/draft-review/pkg/config"
	draftdb "github.com/go-go-golems/draft-review/pkg/db"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	pkgerrors "github.com/pkg/errors"
)

type MigrateStatusCommand struct {
	*cmds.CommandDescription
}

var _ cmds.BareCommand = &MigrateStatusCommand{}

func NewMigrateStatusCommand() (*MigrateStatusCommand, error) {
	sqlSection, err := draftconfig.NewSQLConnectionSection()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create sql section")
	}

	commandSettingsSection, err := cli.NewCommandSettingsSection()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create command settings section")
	}

	description := cmds.NewCommandDescription(
		"status",
		cmds.WithShort("Show migration status"),
		cmds.WithLong(`Show which embedded Draft Review migrations are applied in the configured database.`),
		cmds.WithSections(sqlSection, commandSettingsSection),
	)

	return &MigrateStatusCommand{CommandDescription: description}, nil
}

func (c *MigrateStatusCommand) Run(ctx context.Context, parsedValues *values.Values) error {
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

	statuses, err := draftdb.MigrationStatuses(ctx, db.Pool())
	if err != nil {
		return pkgerrors.Wrap(err, "failed to load migration statuses")
	}

	for _, status := range statuses {
		state := "pending"
		if status.Applied {
			state = "applied"
		}
		fmt.Printf("%s\t%s\n", state, status.Name)
	}

	return nil
}
