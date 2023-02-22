package main

import (
	"log"
	"os"
	"time"

	"github.com/firdisml/go-mongo-rest/configs"
	"github.com/firdisml/go-mongo-rest/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	recaptcha "github.com/jansvabik/fiber-recaptcha"
)

func main() {

	//Recaptcha
	recaptcha.SecretKey = configs.Env("RECAPTCHA_SECRET")

	//Fiber
	app := fiber.New()

	//Limiter
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))

	//Cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
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
