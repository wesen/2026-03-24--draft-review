package cmds

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	draftconfig "github.com/go-go-golems/draft-review/pkg/config"
	draftdb "github.com/go-go-golems/draft-review/pkg/db"
	"github.com/go-go-golems/draft-review/pkg/server"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type ServeCommand struct {
	*cmds.CommandDescription
	version string
}

type ServeSettings struct {
	ListenHost string `glazed:"listen-host"`
	ListenPort int    `glazed:"listen-port"`
}

var _ cmds.BareCommand = &ServeCommand{}

func NewServeCommand(version string) (*ServeCommand, error) {
	defaultSection, err := schema.NewSection(
		schema.DefaultSlug,
		"Server Settings",
		schema.WithFields(
			fields.New(
				"listen-host",
				fields.TypeString,
				fields.WithHelp("Host interface to bind"),
				fields.WithDefault("0.0.0.0"),
			),
			fields.New(
				"listen-port",
				fields.TypeInteger,
				fields.WithHelp("Port to listen on"),
				fields.WithDefault(8080),
			),
		),
	)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create default section")
	}

	sqlSection, err := draftconfig.NewSQLConnectionSection()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create sql section")
	}

	backendSection, err := draftconfig.NewBackendSection()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create backend section")
	}

	commandSettingsSection, err := cli.NewCommandSettingsSection()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to create command settings section")
	}

	description := cmds.NewCommandDescription(
		"serve",
		cmds.WithShort("Start the Draft Review backend"),
		cmds.WithLong(`Start the Draft Review Go backend.

The backend serves:
- health and runtime info routes
- future author and reader JSON endpoints on /api/*
- a PostgreSQL-backed persistence layer

Examples:
  draft-review serve --dsn postgres://draft_review:draft_review@127.0.0.1:5432/draft_review?sslmode=disable
  draft-review serve --host 127.0.0.1 --port 8080 --auto-migrate
`),
		cmds.WithSections(defaultSection, sqlSection, backendSection, commandSettingsSection),
	)

	return &ServeCommand{
		CommandDescription: description,
		version:            version,
	}, nil
}

func (c *ServeCommand) Run(ctx context.Context, parsedValues *values.Values) error {
	settings := &ServeSettings{}
	if err := parsedValues.DecodeSectionInto(schema.DefaultSlug, settings); err != nil {
		return pkgerrors.Wrap(err, "failed to decode serve settings")
	}

	sqlSettings, err := draftconfig.LoadSQLConnectionSettings(parsedValues)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to load sql connection settings")
	}

	backendSettings, err := draftconfig.LoadBackendSettings(parsedValues)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to load backend settings")
	}

	connectionString, err := sqlSettings.ConnectionString()
	if err != nil {
		return pkgerrors.Wrap(err, "failed to compute connection string")
	}

	db, err := draftdb.Open(ctx, connectionString)
	if err != nil {
		return pkgerrors.Wrap(err, "failed to open application database")
	}
	if db != nil {
		defer db.Close()
	}

	if db != nil && backendSettings.AutoMigrate {
		if err := db.Migrate(ctx); err != nil {
			return pkgerrors.Wrap(err, "failed to apply embedded migrations")
		}
	}

	serverCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpServer, err := server.NewHTTPServer(serverCtx, server.Options{
		Host:                settings.ListenHost,
		Port:                settings.ListenPort,
		Version:             c.version,
		Database:            db,
		FrontendDevProxyURL: backendSettings.FrontendDevProxyURL,
	})
	if err != nil {
		return pkgerrors.Wrap(err, "failed to create http server")
	}

	go func() {
		<-serverCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Warn().Err(err).Msg("failed to shutdown server cleanly")
		}
	}()

	log.Info().
		Str("address", httpServer.Addr).
		Bool("database_configured", db != nil).
		Bool("auto_migrate", backendSettings.AutoMigrate).
		Str("frontend_dev_proxy_url", backendSettings.FrontendDevProxyURL).
		Msg("Starting Draft Review backend")

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server exited with error: %w", err)
	}

	return nil
}
