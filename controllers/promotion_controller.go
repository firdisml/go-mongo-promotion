package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/firdisml/go-mongo-rest/configs"
	"github.com/firdisml/go-mongo-rest/models"
	"github.com/firdisml/go-mongo-rest/responses"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var promotion_collection *mongo.Collection = configs.GetCollection(configs.Database, "promotions")
var validate = validator.New()

func CreatePromotion(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var promotion models.Promotion
	defer cancel()

	file_header, file_header_error := c.FormFile("image")
	if file_header_error != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PromotionResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    &fiber.Map{"data": file_header_error.Error()}})
	}

	if parse_error := c.BodyParser(&promotion); parse_error != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PromotionResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    &fiber.Map{"data": parse_error.Error()}})
	}

	if validation_error := validate.Struct(&promotion); validation_error != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PromotionResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    &fiber.Map{"data": validation_error.Error()}})
	}

	file, file_error := file_header.Open()
	if file_error != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PromotionResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    &fiber.Map{"data": file_error.Error()}})
	}

	defer file.Close()

	promotion_id := primitive.NewObjectID()

	promotion_created_date := primitive.NewDateTimeFromTime(time.Now())

	file_size := file_header.Size

	file_buffer := make([]byte, file_size)

	file.Read(file_buffer)

	file_upload_result := configs.UploadFile(configs.Storage, file_size, file_buffer, promotion_id)

	new_promotion := models.Promotion{
		Id:          promotion_id,
		Title:       promotion.Title,
		Category:    promotion.Category,
		Description: promotion.Description,
		Shop:        promotion.Shop,
		State:       promotion.State,
		Link:        promotion.Link,
		Created:     promotion_created_date,
		Start:       promotion.Start,
		End:         promotion.End,
		Visible:     promotion.Visible,
	}

	insert_result, insert_error := promotion_collection.InsertOne(ctx, new_promotion)
	if insert_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PromotionResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": insert_error.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.PromotionResponse{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    &fiber.Map{"mongo_data": insert_result, "image_data": file_upload_result, "promotion_id": promotion_id}})

}

func GetPromotions(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	promotions := []models.Promotion{}
	promotion_skip_string := c.Query("skip")
	promotion_limit_string := c.Query("limit")
	promotion_search := c.Query("search")
	defer cancel()

	promotion_skip_int64, promotion_skip_convert_error := strconv.ParseInt(promotion_skip_string, 10, 64)
	if promotion_skip_convert_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PromotionResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": promotion_skip_convert_error.Error()}})
	}

	promotion_limit_int64, promotion_limit_convert_error := strconv.ParseInt(promotion_limit_string, 10, 64)
	if promotion_limit_convert_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PromotionResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": promotion_limit_convert_error.Error()}})
	}

	find_options := options.Find()
	find_options.SetSkip(promotion_skip_int64)
	find_options.SetLimit(promotion_limit_int64)

	filter := bson.M{"$text": bson.M{"$search": promotion_search}}

	find_cursor, find_error := promotion_collection.Find(ctx, filter, find_options)
	if find_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PromotionResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": find_error.Error()}})
	}

	if cursor_error := find_cursor.All(context.TODO(), &promotions); cursor_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PromotionResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": cursor_error.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.PromotionResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    &fiber.Map{"data": promotions}})
}

func GetUniquePromotion(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	promotion_id := c.Params("id")
	var promotion models.Promotion
	defer cancel()

	promotion_object_id, object_error := primitive.ObjectIDFromHex(promotion_id)
	if object_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PromotionResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": object_error.Error()}})
	}

	if find_error := promotion_collection.FindOne(ctx, bson.M{"id": promotion_object_id}).Decode(&promotion); find_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PromotionResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": find_error.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.PromotionResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    &fiber.Map{"data": promotion}})

}

func EditUniquePromotion(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	promotion_id := c.Params("id")
	var promotion models.Promotion
	defer cancel()

	promotion_object_id, object_error := primitive.ObjectIDFromHex(promotion_id)
	if object_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PromotionResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": object_error.Error()}})
	}

	if parse_error := c.BodyParser(&promotion); parse_error != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PromotionResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    &fiber.Map{"data": parse_error.Error()}})
	}

	if validation_error := validate.Struct(&promotion); validation_error != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.PromotionResponse{
			Status:  http.StatusBadRequest,
			Message: "Error",
			Data:    &fiber.Map{"data": validation_error.Error()}})
	}

	update_promotion := bson.M{
		"title":       promotion.Title,
		"category":    promotion.Category,
		"description": promotion.Description,
		"shop":        promotion.Shop,
		"state":       promotion.State,
		"link":        promotion.Link,
		"created":     promotion.Created,
		"start":       promotion.Start,
		"end":         promotion.End,
		"visible":     promotion.Visible,
	}

	update_result, update_error := promotion_collection.UpdateOne(ctx, bson.M{"id": promotion_object_id}, bson.M{"$set": update_promotion})
	if update_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PromotionResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": update_error.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.PromotionResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    &fiber.Map{"data": update_result}})
}

func DeleteUniquePromotion(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	promotion_id := c.Params("id")
	defer cancel()

	promotion_object_id, object_error := primitive.ObjectIDFromHex(promotion_id)
	if object_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PromotionResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": object_error.Error()}})
	}

	delete_result, delete_error := promotion_collection.DeleteOne(ctx, bson.M{"id": promotion_object_id})
	if delete_error != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.PromotionResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error",
			Data:    &fiber.Map{"data": delete_error.Error()}})
	}

	return c.Status(http.StatusOK).JSON(
		responses.PromotionResponse{
			Status:  http.StatusOK,
			Message: "Success",
			Data:    &fiber.Map{"data": delete_result}},
	)
}
