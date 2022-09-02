package config

//Config Server configuration struct
type Cfg struct {
	Server
}

type Server struct {
	Logrus struct {
		LogLevel string `env:"LOGSLEVEL" ` // info,debug
		JSON     string `env:"JSONLOGS" `  // log format in json
	}
	Address              string `env:"ADDRESS"`
	URLPostgres          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}
