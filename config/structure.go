package config

type Config struct {
	Name        string         `yaml:"name,omitempty"`
	Environment string         `yaml:"environment,omitempty"`
	OutputDir   string         `yaml:"output_dir,omitempty"`
	Server      ServerConfig   `yaml:"server,omitempty"`
	Database    DatabaseConfig `yaml:"database,omitempty"`
	Logging     LoggingConfig  `yaml:"logging,omitempty"`
}

type ServerConfig struct {
	Host              string   `yaml:"host,omitempty"`
	Port              uint     `yaml:"port,omitempty"`
	DockerHost        string   `yaml:"docker_host,omitempty"`
	ComposePath       string   `yaml:"compose_path,omitempty"`
	APISpecProduction bool     `yaml:"api_spec_production,omitempty"`
	CORSOrigins       []string `yaml:"cors_origins,omitempty"`
}

type DatabaseConfig struct {
	Path string `yaml:"path,omitempty"`
	WAL  bool   `yaml:"wal,omitempty"`
}

type LoggingConfig struct {
	Level string `yaml:"level,omitempty"`
}
