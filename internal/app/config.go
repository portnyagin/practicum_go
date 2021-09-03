package app

type AppConfig struct {
	Server_address string `env:"SERVER_ADDRESS" envDefault:":8080"`
	Base_URL       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
}
