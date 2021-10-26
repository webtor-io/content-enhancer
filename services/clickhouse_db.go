package services

import (
	"database/sql"
	"sync"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/urfave/cli"
)

type DBProvider interface {
	Get() (*sql.DB, error)
}

const (
	clickHouseDSN = "clickhouse-dsn"
)

func RegisterClickHouseDBFlags(f []cli.Flag) []cli.Flag {
	return append(f,
		cli.StringFlag{
			Name:   clickHouseDSN,
			Usage:  "clickhouse dsn",
			EnvVar: "CLICKHOUSE_DSN",
		})
}

type ClickHouseDB struct {
	dsn  string
	err  error
	db   *sql.DB
	once sync.Once
}

func NewClickHouseDB(c *cli.Context) *ClickHouseDB {
	return &ClickHouseDB{
		dsn: c.String(clickHouseDSN),
	}
}

func (s *ClickHouseDB) Get() (*sql.DB, error) {
	s.once.Do(func() {
		s.db, s.err = sql.Open("clickhouse", s.dsn)
	})
	return s.db, s.err
}

func (s *ClickHouseDB) Close() {
	if s.db != nil {
		s.db.Close()
	}
}
