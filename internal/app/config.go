package app

type AppConfig struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
	FileStorage   string `env:"FILE_STORAGE_PATH" envDefault:"./data"`
}

func (c *AppConfig) validate() error {
	if c.BaseURL == "" {
		c.BaseURL = "http://localhost:8080/"
	}
	if c.BaseURL[len(c.BaseURL)-1:] != "/" {
		c.BaseURL += "/"
	}
	return nil
}
