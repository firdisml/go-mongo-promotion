package routes

import (
	"github.com/firdisml/go-mongo-rest/controllers"
	"github.com/firdisml/go-mongo-rest/middlewares"
	"github.com/gofiber/fiber/v2"
	recaptcha "github.com/jansvabik/fiber-recaptcha"
)

func PromotionRoutes(app *fiber.App) {
	//Group
	promotion := app.Group("/api/promotions")

	//Routes
	promotion.Post("/", recaptcha.Middleware, controllers.CreatePromotion)
	promotion.Get("/visible", controllers.GetPromotionsVisible)
	promotion.Get("/hidden", middlewares.Authenticated(), controllers.GetPromotionsHidden)
	promotion.Get("/:id", middlewares.Authenticated(), controllers.GetUniquePromotion)
	promotion.Patch("/:id", middlewares.Authenticated(), controllers.EditUniquePromotion)
	promotion.Delete("/:id", middlewares.Authenticated(), controllers.DeleteUniquePromotion)
}
