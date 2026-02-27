package main

import (
	"task4/internal/app"
	"task4/internal/config"
)

func main() {
	cfg := config.MustLoad("config.env")
	pool := app.ConnectToDB(cfg)
	defer pool.Close()

	BookUseCase := app.SetupUseCases(cfg, pool)

	router := app.SetupRouter(BookUseCase, cfg)

	srv := app.SetupServer(cfg, router)

	app.StartServer(srv)

	app.HandleQuit(srv)
}
