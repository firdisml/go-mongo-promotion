package routes

import (
	"github.com/firdisml/go-mongo-rest/controllers"
	"github.com/gofiber/fiber/v2"
)

func PromotionRoutes(app *fiber.App) {
	//Group
	promotion := app.Group("/api/promotion")

	//Routes
	promotion.Post("/", controllers.CreatePromotion)
	promotion.Get("/:id", controllers.GetUniquePromotion)
	promotion.Put("/:id", controllers.EditUniquePromotion)
	promotion.Delete("/:id", controllers.DeleteUniquePromotion)
}
