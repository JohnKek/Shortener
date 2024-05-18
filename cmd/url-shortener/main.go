package main

import (
	"UrlShort/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// TODO: init configs: cleanenv
	// TODO: init loggers: slog
	// TODO: init storage: sqllite
	// TODO: init router: chi, render

}
