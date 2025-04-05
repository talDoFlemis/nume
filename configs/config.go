package configs

import (
	"bytes"
	_ "embed"
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

//go:embed base.yaml
var baseConfig []byte

type CORSCfg struct {
	Origins []string `mapstructure:"origins" validate:"required,dive,http_url"`
	Methods []string `mapstructure:"methods" validate:"required,dive,oneof=GET POST PUT PATCH DELETE OPTIONS"`
	Headers []string `mapstructure:"headers" validate:"required,dive,oneof=Origin Accept Content-Type Authorization X-CSRF-Token"`
}

type HTTPCfg struct {
	Port                     int     `mapstructure:"port"                        validate:"required,min=1,max=65535"`
	ApiPrefix                string  `mapstructure:"api-prefix"                  validate:"required"`
	IP                       string  `mapstructure:"ip"                          validate:"required,ip"`
	CORS                     CORSCfg `mapstructure:"cors"                        validate:"required"`
	ShutdownTimeoutInSeconds int     `mapstructure:"shutdown-timeout-in-seconds" validate:"required,gt=10,lt=600"`
	ReadTimeoutInSeconds     int     `mapstructure:"read-timeout-in-seconds"     validate:"required,gt=10,lt=600"`
	WriteTimeoutInSeconds    int     `mapstructure:"write-timeout-in-seconds"    validate:"required,gt=10,lt=600"`
	IdleTimeoutInSeconds     int     `mapstructure:"idle-timeout-in-seconds"     validate:"required,gt=10,lt=600"`
}

type AppCfg struct {
	Name        string `mapstructure:"name"        validate:"required"`
	Version     string `mapstructure:"version"     validate:"required"`
	Environment string `mapstructure:"environment" validate:"required,oneof=develop prod local"`
}

type IBSUCfg struct {
	EnableSpoof  bool `mapstructure:"enable-spoof"`
	MinNFIQScore int  `mapstructure:"min-nfiq-score" validate:"required,min=1,max=100"`
}

type LoggerCfg struct {
	Level      string `mapstructure:"level"       validate:"required,oneof=DEBUG INFO WARN ERROR"`
	EnableJSON bool   `mapstructure:"enable-json"`
}

type Config struct {
	HTTP   HTTPCfg   `mapstructure:"http"   validate:"required"`
	App    AppCfg    `mapstructure:"app"    validate:"required"`
	Logger LoggerCfg `mapstructure:"logger" validate:"required"`
}

func LoadConfig() (*Config, error) {
	var cfg *Config

	// Load base config
	viper.SetConfigName("base")
	viper.SetConfigType("yaml")

	err := viper.ReadConfig(bytes.NewReader(baseConfig))
	if err != nil {
		slog.Error("failed to read base config", slog.Any("err", err))
		return nil, err
	}

	viper.SetEnvPrefix("NUME")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", ""))
	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)
	if err != nil {
		slog.Error("failed to unmarshal config", slog.Any("err", err))
		return nil, err
	}

	validate := validator.New()
	err = validate.Struct(cfg)
	if err != nil {
		slog.Error("config validation failed", slog.Any("err", err))
		return nil, err
	}

	return cfg, nil
}
