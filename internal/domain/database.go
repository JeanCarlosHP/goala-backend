package domain

type Database interface {
	NewConnection(config *Config) error
}
