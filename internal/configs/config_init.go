package configs

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvConfigPathName  = "CONFIG_PATH"
	FlagConfigPathName = "config"
	ConfigPathConst    = "configs/config_prod.yaml"
)

var configPath string
var configMap *Config
var once sync.Once

// GetConfig возвращает конфигурацию приложения
func GetConfig() *Config {
	// Запоминаем путь к файлу конфигурации приложения
	once.Do(func() {
		//// Логируем путь перед загрузкой
		//slog.Info("КОНФИГУРАЦИЯ. Загрузка конфигурации из файла", slog.String("path", configPath))

		// Загружаем конфигурацию
		configMap = &Config{}
		aa := getConfigPath(ConfigPathConst)
		slog.Debug("КОНФИГУРАЦИЯ. Загрузка конфигурации из файла", slog.String("CONFIG", aa))
		if err := cleanenv.ReadConfig(aa, configMap); err != nil {
			helpText := "КОНФИГУРАЦИЯ. Ошибка считывания данных из файла"
			help, _ := cleanenv.GetDescription(configMap, &helpText)
			slog.Error(help)
			slog.Error(err.Error())
			os.Exit(1)
		}
	})
	return configMap
}

func getConfigPath(configPathConst string) string {

	var configPath string
	var varConfig string
	// Проверяем указанный путь к файлу конфигурации через флаг
	flag.StringVar(&configPath, FlagConfigPathName, "", "Файл конфигурации приложения")
	flag.Parse()
	if configPath != "" {
		// 1) если передали через флаг
		varConfig = "флаг"
	} else if env := os.Getenv(EnvConfigPathName); env != "" {
		// 2) иначе, если задано через переменную окружения
		configPath = env
		varConfig = "переменная окружения"
	} else {
		// 3) ни того ни другого — берём путь по умолчанию
		configPath = configPathConst
		varConfig = "Константа"
	}

	log.Printf("КОНФИГУРАЦИЯ. Источник: %s. Путь до файла: %s", varConfig, configPath)
	return configPath
}
