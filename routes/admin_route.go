package routes

import (
	"github.com/firdisml/go-mongo-rest/controllers"
	"github.com/firdisml/go-mongo-rest/middlewares"
	"github.com/gofiber/fiber/v2"
)

func AdminRoutes(app *fiber.App) {
	//Group
	admin := app.Group("/api/admin")

	//Routes
	admin.Post("/", controllers.SignUpAdmin)
	admin.Post("/login", controllers.SignInAdmin)
	admin.Get("/get", middlewares.Authenticated(), controllers.GetAdmin)
}
