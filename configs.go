package main

import "github.com/kelseyhightower/envconfig"

// Config returns the local environment variables
type Config struct {
	MQTTURL    string `envconfig:"MQTT_URL" required:"true"`
	MQTTTopic  string `envconfig:"MQTT_TOPIC" required:"true"`
	MQTTClient string `envconfig:"MQTT_CLIENT" required:"true"`
}

func loadConfigs() (*Config, error) {
	var config Config
	err := envconfig.Process("AS", &config)
	return &config, err
}
