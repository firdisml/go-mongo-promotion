package routes

import (
	"github.com/firdisml/go-mongo-rest/controllers"
	"github.com/gofiber/fiber/v2"
)

func AdminRoutes(app *fiber.App) {
	//Group
	promotion := app.Group("/api/admin")

	//Routes
	promotion.Post("/", controllers.SignUpAdmin)
	promotion.Post("/login", controllers.SignInAdmin)
}
