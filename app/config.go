package app

type Server struct {
	Host string `env:"HOST" env-default:"localhost"`
	Port int    `env:"PORT" env-default:"4000"`
}

type MetricsConfig struct {
	Enabled bool   `env:"METRICS_ENABLED" env-default:"true"`
	Host    string `env:"METRICS_HOST" env-default:"localhost"`
	Port    int    `env:"METRICS_PORT" env-default:"9100"`
}

type TimeoutConfig struct {
	ReadTimeout    int `env:"SERVER_READ_TIMEOUT" env-default:"30"`
	WriteTimeout   int `env:"SERVER_WRITE_TIMEOUT" env-default:"30"`
	IdleTimeout    int `env:"SERVER_IDLE_TIMEOUT" env-default:"120"`
	HandlerTimeout int `env:"SERVER_HANDLER_TIMEOUT" env-default:"25"`
}

type AppConfig struct {
	Server
	Metrics  MetricsConfig
	Timeouts TimeoutConfig
	AppEnv   string `env:"APP_ENV" env-default:"dev"` // "dev", "prodction"
}
