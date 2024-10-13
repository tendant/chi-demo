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

type AppConfig struct {
	Server
	Metrics MetricsConfig
	AppEnv  string `env:"APP_ENV" env-default:"dev"` // "dev", "prodction"
}
