package main

import (
	"log"
	"os"

	"github.com/firdisml/go-mongo-rest/configs"
	"github.com/firdisml/go-mongo-rest/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recaptcha "github.com/jansvabik/fiber-recaptcha"
)

func main() {

	//Recaptcha
	recaptcha.SecretKey = configs.Env("RECAPTCHA_SECRET")

	//Fiber
	app := fiber.New()

	//Cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	//Mongo
	configs.ConnectMongo()

	//Storage
	configs.ConnectStorage()

	//Routes
	routes.PromotionRoutes(app)
	routes.AdminRoutes(app)

	//Get Port
	port := os.Getenv("PORT")

	//Set Port
	if port == "" {
		port = "3000"
	}

	//Listen
	log.Fatal(app.Listen("0.0.0.0:" + port))
}
