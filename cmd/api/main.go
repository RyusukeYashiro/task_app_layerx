package main

import (
	"context"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	echomw "github.com/oapi-codegen/echo-middleware"

	api "github.com/ryusuke/task_app_layerx/internal/presentation/http/echo"
)

func main() {
	e := echo.New()
	e.Use(echoMw.Recover(), echoMw.Logger())

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("api/openapi.yaml")
	if err != nil { log.Fatalf("load openapi: %v", err) }

	// 開発中は認証を通す（あとでJWT検証に置換）
	e.Use(echomw.OapiRequestValidatorWithOptions(doc, &echomw.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
				// TODO: 本実装では Authorization ヘッダーのJWTを検証する
				return nil
			},
		},
	}))

	api.RegisterHandlersWithBaseURL(e, api.NewHandler(), "/api/v1")

	port := os.Getenv("APP_PORT")
	if port == "" { port = "8080" }
	log.Printf("listening on :%s", port)
	log.Fatal(e.Start(":" + port))
}