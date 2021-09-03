package app

type AppConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseUrl       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
}
