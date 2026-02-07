package config_accrual

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"
)

type Config struct {
	Opts  *Options
	Token *Token
}

type Options struct {
	ServerHost   string `env:"RUN_ADDRESS"`
	DataBaseHost string `env:"DATABASE_URI"`
}

type Token struct {
	SecretKey  string        `env:"SECRET_KEY" default:"super-secret-key"`
	CookieName string        `env:"COOKIE_NAME" default:"url_shortener_session"`
	CookieTTL  time.Duration `env:"COOKIE_TTL" default:"24h"`
}

func NewConfig() (*Config, error) {
	opts, err := newOpts()
	if err != nil {
		return nil, err
	}
	return &Config{
		Opts: opts,
	}, nil
}

func newOpts() (*Options, error) {
	opts, ok := parseEnv()
	if ok {
		return opts, nil
	}
	var addr = flag.String("a", "localhost:8080", "server host")
	var psqlHost = flag.String("d", "", "psql data")
	flag.Parse()

	if opts.ServerHost == "" {
		opts.ServerHost = *addr
	}

	if opts.DataBaseHost == "" {
		opts.DataBaseHost = *psqlHost
	}
	if _, err := url.Parse("https://" + opts.ServerHost); err != nil {
		return nil, fmt.Errorf("incorrect parametr `-a` %s", opts.ServerHost)
	}

	if _, err := url.Parse("https://" + opts.DataBaseHost); err != nil {
		return nil, fmt.Errorf("incorrect parametr `-d` %s", opts.DataBaseHost)
	}

	return opts, nil
}

func parseEnv() (*Options, bool) {
	envAddr := os.Getenv("RUN_ADDRESS")
	envPsqlDsn := os.Getenv("DATABASE_URI")

	opts := &Options{}
	if _, err := url.Parse("http://" + envAddr); err == nil {
		opts.ServerHost = envAddr
	}
	
	if _, err := url.Parse(envPsqlDsn); err == nil {
		opts.DataBaseHost = envPsqlDsn
	}

	sucessAll := opts.ServerHost != "" && opts.DataBaseHost != ""
	return opts, sucessAll
}