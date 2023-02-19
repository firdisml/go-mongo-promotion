package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/firdisml/go-mongo-rest/configs"
	"github.com/firdisml/go-mongo-rest/models"
	"github.com/firdisml/go-mongo-rest/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var admin_collection *mongo.Collection = configs.GetCollection(configs.Database, "admins")

func SignUpAdmin(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var admin models.SignUpAdmin
	defer cancel()

	if parse_error := c.BodyParser(&admin); parse_error != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.AdminResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    &fiber.Map{"data": parse_error.Error()}})
	}

	if validation_error := validate.Struct(&admin); validation_error != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.AdminResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    &fiber.Map{"data": validation_error.Error()}})
	}

	hashed_password, password_hashing_error := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if password_hashing_error != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.AdminResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    &fiber.Map{"data": password_hashing_error.Error()}})
	}

	admin_id := primitive.NewObjectID()

	admin_created_date := time.Now()
	admin_updated_date := time.Now()

	new_admin := models.SignUpAdmin{
		Id:       admin_id,
		Email:    strings.ToLower(admin.Email),
		Password: string(hashed_password),
		Name:     admin.Name,
		Created:  admin_created_date,
		Updated:  admin_updated_date,
		Status:   admin.Status,
	}

	insert_result, insert_error := admin_collection.InsertOne(ctx, new_admin)
	if insert_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AdminResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": insert_error.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.AdminResponse{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    &fiber.Map{"mongo_data": insert_result, "admin_id": admin_id}})
}

func SignInAdmin(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var admin_entered models.SignInAdmin
	var admin_stored models.SignUpAdmin
	defer cancel()

	if parse_error := c.BodyParser(&admin_entered); parse_error != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.AdminResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    &fiber.Map{"data": parse_error.Error()}})
	}

	if validation_error := validate.Struct(&admin_entered); validation_error != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.AdminResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    &fiber.Map{"data": validation_error.Error()}})
	}

	if find_error := admin_collection.FindOne(ctx, bson.M{"email": admin_entered.Email}).Decode(&admin_stored); find_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AdminResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": find_error.Error()}})
	}

	password_compare_error := bcrypt.CompareHashAndPassword([]byte(admin_stored.Password), []byte(admin_entered.Password))
	if password_compare_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AdminResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": password_compare_error.Error()}})
	}

	token_byte := jwt.New(jwt.SigningMethodHS256)

	token_claims := token_byte.Claims.(jwt.MapClaims)

	token_claims["sub"] = admin_stored.Id.Hex()
	token_claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	token_string, token_string_error := token_byte.SignedString([]byte(configs.Env("JWT_SECRET")))
	if token_string_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AdminResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": token_string_error.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.AdminResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    &fiber.Map{"token": token_string}})
}

func GetAdmin(c *fiber.Ctx) error {
	return c.Status(http.StatusCreated).JSON(responses.AdminResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    &fiber.Map{"token": "gay"}})
}
