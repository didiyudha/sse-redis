package main

import (
	"fmt"
	"log"

	"github.com/didiyudha/sse-redis/config"
	"github.com/didiyudha/sse-redis/domain/product/handler"
	"github.com/didiyudha/sse-redis/domain/product/repository"
	"github.com/didiyudha/sse-redis/domain/product/usecase"
	"github.com/didiyudha/sse-redis/internal/platform/redis"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	godotenv.Load(".env")
	config.LoadEnv()

	redisConn, err := redis.NewRedis(config.Cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}

	productCache := repository.NewProductCache(redisConn)
	productUseCase := usecase.NewProductUseCase(productCache)
	productHandler := handler.NewProductHandler(productUseCase)

	e := echo.New()
	e.Debug = true

	e.POST("/products", productHandler.Store)
	e.GET("/products/streams/:productId", productHandler.Streams)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Cfg.Port)))
}
