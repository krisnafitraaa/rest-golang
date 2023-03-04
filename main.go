package main

import (
	"github.com/fundraising/rest-api/database"
	"github.com/fundraising/rest-api/routes"
)

func main() {
	database.Initialize()
	e := routes.DefineRoutes()

	e.Logger.Fatal(e.Start(":8000"))
}
