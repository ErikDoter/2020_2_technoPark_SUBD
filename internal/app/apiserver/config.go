package apiserver

import "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/store"

type Config struct {
	BindAddr string `toml:"bind_addr"`
	Store    *store.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8000",
		Store:    store.NewConfig(),
	}
}