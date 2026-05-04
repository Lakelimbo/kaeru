package config

// for environment variables
const (
	ENVIRONMENT = "KAERU_ENVIRONMENT"
	JWT_SECRET  = "KAERU_JWT_SECRET"
)

func (c *Config) IsDev() bool {
	return c.Environment == "development"
}

func (c *Config) IsProd() bool {
	return c.Environment == "production"
}
