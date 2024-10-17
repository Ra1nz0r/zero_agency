package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerHost               string        `mapstructure:"ADDRESS"`                // адрес сервера
	LogLevel                 string        `mapstructure:"LOG_LEVEL"`              // уровень логирования
	MigrationPath            string        `mapstructure:"MIGRATION_PATH"`         // путь до папки с миграциями
	PaginationLimit          int32         `mapstructure:"PAGINATION_LIMIT"`       // размер пагинации по-умолчанию
	DatabaseUser             string        `mapstructure:"DB_USER"`                // имя пользователя датабазы
	DatabasePassword         string        `mapstructure:"DB_PASSWORD"`            // пароль пользователя датабазы
	DatabaseHost             string        `mapstructure:"DB_HOST"`                // адрес для подключения к датабазе
	DatabasePort             string        `mapstructure:"DB_PORT"`                // порт для подключения к датабазе
	DatabaseName             string        `mapstructure:"DB_NAME"`                // имя датабазы
	DatabaseDriver           string        `mapstructure:"DB_DRIVER"`              // драйвер датабазы
	DatabaseMaxOpenConns     int           `mapstructure:"DB_MAX_OPEN_CONNS"`      // максимальное количество открытых подключений
	DatabaseMaxIdleConns     int           `mapstructure:"DB_MAX_IDLE_CONNS"`      // максимальное количество подключений в пуле бездействующих подключений
	DatabaseMaxLifetimeInMin time.Duration `mapstructure:"DB_MAX_LIFETIME_IN_MIN"` // максимальное количество времени, в течение которого соединение может быть повторно использовано
}

// LoadConfig загружает из файла '.env' переменные окружения.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
