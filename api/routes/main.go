package routes

import (
	"my-graphql-project/api/graph"
	"my-graphql-project/api/graph/resolver"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
)

func EchoSetup(app *echo.Echo, resolver *resolver.Resolver) {
	api := app.Group("/api/v1", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	})
	// Define routes inside the `api` group
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "API is running"})
	})
	// GraphQL endpoint
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	api.POST("/graphql", echo.WrapHandler(srv))
	api.GET("/", echo.WrapHandler(playground.Handler("GraphQL Playground", "/graphql")))

}
