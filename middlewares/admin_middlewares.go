package middlewares

import (
	"github.com/firdisml/go-mongo-rest/configs"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func Authenticated() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(configs.Env("JWT_SECRET")),
	})
}
