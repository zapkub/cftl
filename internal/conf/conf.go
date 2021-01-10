package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/zapkub/cftl/internal/fsutil"
	"github.com/zapkub/cftl/internal/logger"
)

type c struct {
	Dev   bool
	Oauth struct {
		GithubClientID     string
		GithubClientSecret string
	}
	Address string

	// driver will fix to postgres
	// if want to change. do it from ENV
	DBDriver string `toml:"-"`

	DBUser     string
	DBPassword string
	DBPort     string
	DBHost     string
	DBName     string
}

// StatementTimeout is the value of the Postgres statement_timeout parameter.
// Statements that run longer than this are terminated.
// 10 minutes is the App Engine standard request timeout.
const StatementTimeout = 10 * time.Minute

// SourceTimeout is the value of the timeout for source.Client, which is used
// to fetch source code from third party URLs.
const SourceTimeout = 1 * time.Minute

func (c *c) DBConnInfo() string {
	// For the connection string syntax, see
	// https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING.
	// Set the statement_timeout config parameter for this session.
	// See https://www.postgresql.org/docs/current/runtime-config-client.html.
	timeoutOption := fmt.Sprintf("-c statement_timeout=%d", StatementTimeout/time.Millisecond)
	return fmt.Sprintf("user='%s' password='%s' host='%s' port=%s dbname='%s' sslmode=disable options='%s'",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, timeoutOption)
}

// C is configuration instance
var C = initconfig()

func initconfig() *c {

	var conf c
	appf := fsutil.Default.MustOpenFile("app.conf", os.O_RDONLY)
	appfb, err := ioutil.ReadAll(appf)
	if err != nil {
		logger.Fatalf(nil, "parse config file error: %v", err)
	}

	_, err = toml.Decode(string(appfb), &conf)
	if err != nil {
		logger.Fatalf(nil, "parse toml config error: %v", err)
	}

	fmt.Printf("dump configuration: %+v\n", conf)

	conf.DBDriver = GetEnv("CFTL_DATABASE_DRIVER", "postgres")

	return &conf
}

// GetEnv looks up the given key from the environment, returning its value if
// it exists, and otherwise returning the given fallback value.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
