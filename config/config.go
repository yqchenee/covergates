package config

import (
	log "github.com/sirupsen/logrus"
)

var config *Config

func GetConfig() *Config {
	if config == nil {
		log.Fatal("Config not initialize")
	}
	return config
}

type Config struct {
	Server Server
	Gitea  Gitea
	Github Github
}

type Server struct {
	Secret string `default:"secret"`
}

type Gitea struct {
	Server       string
	ClientID     string
	ClientSecret string
	SkipVerity   bool
	Scope        []string `default:"repo,repo:status,user:email,read:org"`
}

type Github struct {
	Server       string `default:"https://github.com"`
	APIServer    string
	ClientID     string
	ClientSecret string
	SkipVerity   bool
	Scope        []string `default:"repo,repo:status,user:email,read:org"`
}
