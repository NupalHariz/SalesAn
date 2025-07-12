package config

import (
	"time"

	"github.com/reyhanmichiels/go-pkg/v2/auth"
	"github.com/reyhanmichiels/go-pkg/v2/log"
	"github.com/reyhanmichiels/go-pkg/v2/parser"
	"github.com/reyhanmichiels/go-pkg/v2/rate_limiter"
	"github.com/reyhanmichiels/go-pkg/v2/redis"
	"github.com/reyhanmichiels/go-pkg/v2/sql"
	"github.com/reyhanmichiels/go-pkg/v2/translator"
)

type Application struct {
	Meta        ApplicationMeta
	Gin         GinConfig
	Log         log.Config
	SQL         sql.Config
	Auth        auth.Config
	Redis       redis.Config
	Translator  translator.Config
	RateLimiter rate_limiter.Config
	Parser      parser.Options
}

type ApplicationMeta struct {
	Title       string
	Description string
	Host        string
	BasePath    string
	Version     string
}

type GinConfig struct {
	Port            string
	Mode            string
	LogRequest      bool
	LogResponse     bool
	Timeout         time.Duration
	ShutdownTimeout time.Duration
	CORS            CORSConfig
	Meta            ApplicationMeta
	Swagger         SwaggerConfig
	Dummy           DummyConfig
}

type CORSConfig struct {
	Mode string
}
type SwaggerConfig struct {
	Enabled   bool
	Path      string
	BasicAuth BasicAuthConf
}

type PlatformConfig struct {
	Enabled   bool
	Path      string
	BasicAuth BasicAuthConf
}

type DummyConfig struct {
	Enabled bool
	Path    string
}

type BasicAuthConf struct {
	Username string
	Password string
}

func Init() Application {
	return Application{}
}
