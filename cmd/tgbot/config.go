package main

type Configuration struct {
	Port  string `envconfig:"PORT" required:"true"`
	Token string `envconfig:"TOKEN" required:"true"`
}
