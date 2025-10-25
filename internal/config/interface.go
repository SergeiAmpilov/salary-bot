// internal/config/interface.go
package config

// ConfigReader описывает контракт для чтения конфигурации
type ConfigReader interface {
	Read() (*Config, error)
}
