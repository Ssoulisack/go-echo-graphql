package main

import (
	"fmt"
	"log"
	"my-graphql-project/api/graph/resolver"
	"my-graphql-project/api/routes"
	"my-graphql-project/bootstrap"
)

func main() {
	app := bootstrap.App()
	globalEnv := app.Env
	echo := app.Echo
	db := app.DB

	resolver := resolver.InitializeResolver(db)

	routes.EchoSetup(echo, resolver)

	log.Fatal(echo.Start(fmt.Sprintf(":%d", globalEnv.App.Port)))
}
