package config

import "time"

// env: "local" # local, dev, prod
// storage_path: "./storage/storage.db"
// http_server:
//   address: "localhost:8080"
//   timeout: 4s
//   idle_timeout: 60s

type Config struct {
	Env         string `yaml:"env" env-default:"local" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-default:"./storage/storage.db" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}
