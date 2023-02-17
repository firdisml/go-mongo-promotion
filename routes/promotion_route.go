package routes

import (
	"github.com/firdisml/go-mongo-rest/controllers"
	"github.com/gofiber/fiber/v2"
	recaptcha "github.com/jansvabik/fiber-recaptcha"
)

func PromotionRoutes(app *fiber.App) {
	//Group
	promotion := app.Group("/api/promotions")

	//Routes
	promotion.Post("/", recaptcha.Middleware, controllers.CreatePromotion)
	promotion.Get("/query", controllers.GetPromotions)
	promotion.Get("/:id", controllers.GetUniquePromotion)
	promotion.Put("/:id", controllers.EditUniquePromotion)
	promotion.Delete("/:id", controllers.DeleteUniquePromotion)
}
