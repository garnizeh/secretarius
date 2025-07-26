package config

type DBConfig struct {
	User          string
	Password      string
	HostReadWrite string
	HostReadOnly  string
	Host          string
	Name          string
	Schema        string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		User:          getEnv("DB_USER", "englog"),
		Password:      getEnv("DB_PASSWORD", "englog_dev_password"),
		HostReadWrite: getEnv("DB_HOST_READ_WRITE", "localhost:5432"),
		HostReadOnly:  getEnv("DB_HOST_READ_ONLY", "localhost:5432"),
		Name:          getEnv("DB_NAME", "englog"),
		Schema:        getEnv("DB_SCHEMA", "englog"),
	}
}
