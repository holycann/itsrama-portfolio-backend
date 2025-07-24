package configs

type ServerConfig struct {
	ProductionDomain string
	Host             string
	Port             int
	ReadTimeout      int
	WriteTimeout     int
	ShutdownTimeout  int
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		ProductionDomain: getEnv("PRODUCTION_DOMAIN", ""),
		Host:             getEnv("SERVER_HOST", "0.0.0.0"),
		Port:             getEnvAsInt("SERVER_PORT", 8181),
		ReadTimeout:      getEnvAsInt("SERVER_READ_TIMEOUT", 15),
		WriteTimeout:     getEnvAsInt("SERVER_WRITE_TIMEOUT", 15),
		ShutdownTimeout:  getEnvAsInt("SERVER_SHUTDOWN_TIMEOUT", 30),
	}
}