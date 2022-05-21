package main

import "github.com/goohopeteam/auth-service/internal/app"

func main() {
	app := app.Init()
	app.RunMigrations()
	app.Run()
}
