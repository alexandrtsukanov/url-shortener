package main

import (
	"fmt"
	"url-shortener/internal/config"
)

func main() {
	// TODO init config: cleanenv
	cfg, err := config.MustLoad()

	fmt.Println(cfg)
	fmt.Println(err)

	// TODO init logger: slog
	// TODO init storage: sqlite
	// TODO init router: chi, chi render
	// TODO run server: cleanenv
}
