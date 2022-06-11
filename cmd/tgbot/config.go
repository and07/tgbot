package main

type Configuration struct {
	Port    string `envconfig:"PORT" required:"true"`
	Token   string `envconfig:"TOKEN" required:"true"`
	AppID   int    `envconfig:"APP_ID" required:"true"`
	AppHash string `envconfig:"APP_HASH" required:"true"`
}
