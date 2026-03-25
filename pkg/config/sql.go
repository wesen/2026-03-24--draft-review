package config

import (
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/pkg/errors"
)

//go:embed "flags/sql-connection.yaml"
var sqlConnectionFlagsYAML []byte

const SQLConnectionSlug = "sql-connection"

type SQLConnectionSettings struct {
	Host       string `glazed:"host"`
	Port       int    `glazed:"port"`
	Database   string `glazed:"database"`
	User       string `glazed:"user"`
	Password   string `glazed:"password"` // #nosec G101 -- Password is part of the DB config model.
	Schema     string `glazed:"schema"`
	DBType     string `glazed:"db-type"`
	DSN        string `glazed:"dsn"`
	Driver     string `glazed:"driver"`
	SSLDisable bool   `glazed:"ssl-disable"`
}

func NewSQLConnectionSection(options ...schema.SectionOption) (schema.Section, error) {
	section, err := schema.NewSectionFromYAML(sqlConnectionFlagsYAML, options...)
	if err != nil {
		return nil, err
	}
	return section, nil
}

func LoadSQLConnectionSettings(parsedValues *values.Values) (*SQLConnectionSettings, error) {
	if parsedValues == nil {
		return nil, errors.New("parsed values are nil")
	}

	settings := &SQLConnectionSettings{}
	if err := parsedValues.DecodeSectionInto(SQLConnectionSlug, settings); err != nil {
		return nil, errors.Wrap(err, "failed to decode sql connection section")
	}

	return NormalizeSQLConnectionSettings(settings)
}

func NormalizeSQLConnectionSettings(settings *SQLConnectionSettings) (*SQLConnectionSettings, error) {
	if settings == nil {
		settings = &SQLConnectionSettings{}
	}

	if settings.Host == "" {
		settings.Host = strings.TrimSpace(os.Getenv("DRAFT_REVIEW_HOST"))
	}
	if settings.Port == 0 {
		settings.Port = envIntOr("DRAFT_REVIEW_PORT", 5432)
	}
	if settings.Database == "" {
		settings.Database = strings.TrimSpace(os.Getenv("DRAFT_REVIEW_DATABASE"))
	}
	if settings.User == "" {
		settings.User = strings.TrimSpace(os.Getenv("DRAFT_REVIEW_USER"))
	}
	if settings.Password == "" {
		settings.Password = strings.TrimSpace(os.Getenv("DRAFT_REVIEW_PASSWORD"))
	}
	if settings.Schema == "" {
		settings.Schema = strings.TrimSpace(os.Getenv("DRAFT_REVIEW_SCHEMA"))
	}
	if settings.DBType == "" {
		settings.DBType = strings.TrimSpace(os.Getenv("DRAFT_REVIEW_DB_TYPE"))
	}
	if settings.DSN == "" {
		settings.DSN = strings.TrimSpace(os.Getenv("DRAFT_REVIEW_DSN"))
	}
	if settings.Driver == "" {
		settings.Driver = strings.TrimSpace(os.Getenv("DRAFT_REVIEW_DRIVER"))
	}
	if !settings.SSLDisable {
		settings.SSLDisable = envBoolOr("DRAFT_REVIEW_SSL_DISABLE", false)
	}

	settings.Host = strings.TrimSpace(settings.Host)
	settings.Database = strings.TrimSpace(settings.Database)
	settings.User = strings.TrimSpace(settings.User)
	settings.Password = strings.TrimSpace(settings.Password)
	settings.Schema = strings.TrimSpace(settings.Schema)
	settings.DSN = strings.TrimSpace(settings.DSN)
	settings.DBType = strings.ToLower(strings.TrimSpace(settings.DBType))
	settings.Driver = strings.ToLower(strings.TrimSpace(settings.Driver))

	if settings.Port == 0 {
		settings.Port = 5432
	}
	if settings.DBType == "" {
		settings.DBType = "pgx"
	}
	if settings.Driver == "" {
		settings.Driver = "pgx"
	}

	return settings, nil
}

func envIntOr(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}

func (s *SQLConnectionSettings) ConnectionString() (string, error) {
	if s == nil {
		return "", errors.New("sql connection settings are nil")
	}
	if s.DSN != "" {
		return s.DSN, nil
	}
	if s.Host == "" || s.Database == "" || s.User == "" {
		return "", errors.New("dsn is empty and host/database/user are not fully configured")
	}

	sslMode := "require"
	if s.SSLDisable {
		sslMode = "disable"
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(s.User, s.Password),
		Host:   fmt.Sprintf("%s:%d", s.Host, s.Port),
		Path:   s.Database,
	}
	q := u.Query()
	q.Set("sslmode", sslMode)
	if s.Schema != "" {
		q.Set("search_path", s.Schema)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}
