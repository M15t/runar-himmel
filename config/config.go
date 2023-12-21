package config

type (
	// Configuration holds data necessery for configuring application
	Configuration struct {
		General
		Server
		DB
		JWT
	}

	// General holds general configurations
	General struct {
		Debug bool `env:"DEBUG" envDefault:"true"`
	}

	// Server holds http server configurations
	Server struct {
		// The port for the http server to listen on
		Port int `env:"SERVER_PORT" envDefault:"8080"`
		// ReadHeaderTimeout is the amount of time allowed to read request headers.
		ReadHeaderTimeout int `env:"SERVER_READ_HEADER_TIMEOUT" envDefault:"10"`
		// ReadTimeout is the maximum duration for reading the entire request, including the body
		ReadTimeout int `env:"SERVER_READ_TIMEOUT" envDefault:"30"`
		// WriteTimeout is the maximum duration before timing out writes of the response
		WriteTimeout int `env:"SERVER_WRITE_TIMEOUT" envDefault:"60"`
		// CORS settings
		AllowOrigins []string `env:"SERVER_ALLOW_ORIGINS" envDefault:"*"`
	}

	// DB holds DB configurations
	DB struct {
		Driver   string `env:"DB_DRIVER,notEmpty"`
		Host     string `env:"DB_HOST,notEmpty"`
		Port     int    `env:"DB_PORT,notEmpty"`
		Username string `env:"DB_USERNAME,notEmpty"`
		Password string `env:"DB_PASSWORD,notEmpty"`
		Database string `env:"DB_DATABASE,notEmpty"`
		Logging  int    `env:"DB_LOGGING" envDefault:"1"` // 0=discard, 1=silent, 2=error, 3=warn, 4=info
		Params   string `env:"DB_PARAMS"`
	}

	// JWT holds JWT configurations
	JWT struct {
		Secret               string `env:"JWT_SECRET,notEmpty"`
		Algorithm            string `env:"JWT_ALGORITHM" envDefault:"HS256"`
		DurationAccessToken  int    `env:"JWT_DURATION_ACCESS_TOKEN" envDefault:"3600"`   // 1 hour in second
		DurationRefreshToken int    `env:"JWT_DURATION_REFRESH_TOKEN" envDefault:"86400"` // 1 day in second
	}

	// App holds app specific configurations
	App struct {
		// more app specific configurations
	}
)

// LoadAll returns all configurations for the app
func LoadAll() (cfg Configuration, err error) {
	err = Load(&cfg)
	return
}
