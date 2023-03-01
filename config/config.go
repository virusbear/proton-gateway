package config

import (
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
	"io"
)

type Config struct {
	Serial  SerialConfig   `yaml:"serial"`
	Mqtt    MqttConfig     `yaml:"mqtt"`
	Devices []DeviceConfig `yaml:"devices" default:"[]"`
}

type SerialConfig struct {
	Port     string `yaml:"port"`
	BaudRate uint   `yaml:"baudrate" default:"115200"`
}

type MqttConfig struct {
	Host string `yaml:"host" default:"localhost"`
	Port uint16 `yaml:"port" default:"1883"`
}

type DeviceConfig struct {
	Type string `yaml:"type"`
	Mac  string `yaml:"mac"`
}

func Load(reader io.Reader) (*Config, error) {
	config := Config{}

	if err := defaults.Set(&config); err != nil {
		return nil, err
	}

	if err := yaml.NewDecoder(reader).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
