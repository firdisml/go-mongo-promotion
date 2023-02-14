package main

import (
	"log"

	"github.com/firdisml/go-mongo-rest/configs"
	"github.com/firdisml/go-mongo-rest/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	//Fiber
	app := fiber.New()

	//Mongo
	configs.ConnectMongo()

	//Storage
	configs.ConnectStorage()

	//Routes
	routes.PromotionRoutes(app)

	//Listen
	log.Fatal(app.Listen(":3000"))
}
