package infra

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/dig"
)

type (
	Databases struct {
		dig.Out
		Pg *sql.DB
	}

	DatabaseCfgs struct {
		dig.In
		Pg *DatabaseCfg
	}

	DatabaseCfg struct {
		DBName string `envconfig:"DBNAME" required:"true" default:"dbname"`
		DBUser string `envconfig:"DBUSER" required:"true" default:"dbuser"`
		DBPass string `envconfig:"DBPASS" required:"true" default:"dbpass"`
		Host   string `envconfig:"HOST" required:"true" default:"localhost"`
		Port   string `envconfig:"PORT" required:"true" default:"9999"`

		MaxOpenConns    int           `envconfig:"MAX_OPEN_CONNS" default:"20" required:"true"`
		MaxIdleConns    int           `envconfig:"MAX_IDLE_CONNS" default:"5" required:"true"`
		ConnMaxLifetime time.Duration `envconfig:"CONN_MAX_LIFETIME" default:"15m" required:"true"`
	}

	JwtCfg struct {
		SecretKey string `envconfig:"JWT_SECRET_KEY" required:"true" default:"secret"`
	}

	DiscordCfg struct {
		DiscordToken   string `envconfig:"DISCORD_TOKEN" required:"true" default:"discord_token"`
		DiscordGuildID string `envconfig:"DISCORD_GUILD_ID" required:"true" default:"discord_guild_id"`
		DiscordRoleID  string `envconfig:"DISCORD_ROLE_ID" required:"true" default:"discord_role_id"`
	}
)

func NewDatabases(cfgs DatabaseCfgs) Databases {
	return Databases{
		Pg: openPostgres(cfgs.Pg),
	}
}

func openPostgres(p *DatabaseCfg) *sql.DB {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.DBUser, url.QueryEscape(p.DBPass), p.Host, p.Port, p.DBName,
	)

	// Open a connection to the database.
	db, err := sql.Open("postgres", conn)
	if err != nil {
		slog.Error(fmt.Sprintf("postgres: conn, err : %v", err.Error()))
	}

	db.SetConnMaxLifetime(p.ConnMaxLifetime)
	db.SetMaxIdleConns(p.MaxIdleConns)
	db.SetMaxOpenConns(p.MaxOpenConns)

	if err = db.Ping(); err != nil {
		slog.Error(fmt.Sprintf("postgres: ping, err : %v", err.Error()))
	}
	fmt.Println("Connected to the database")
	return db
}
