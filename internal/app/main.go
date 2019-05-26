package app

func New() *App {
	config := NewConfig()
	return &App{
		Config: config,
	}
}

type App struct {
	Config *Config
}
