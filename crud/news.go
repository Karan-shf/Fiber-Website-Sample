package crud

import (
	"context"
	"net/http"
	"tarafdari-sample/database"
	"tarafdari-sample/models"
	"time"

	jtoken "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gofiber/fiber/v2"
)

func Add_News(c *fiber.Ctx) error {

	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)

	_, adminObj, is_admin, err := Get_Ent(claims)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Token"})
	}

	if !is_admin {
		return c.Status(http.StatusMethodNotAllowed).JSON(fiber.Map{"error": "you must be an admin to access this method"})
	}

	var new_news models.News

	if err := c.BodyParser(&new_news); err != nil {
		return err
	}

	if new_news.Body == "" || new_news.Title == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "all fields must be filled"})
	}

	new_news.TimeOfCast = time.Now()
	new_news.Author = *adminObj

	insertResult, err := database.NewsCollection.InsertOne(context.Background(), new_news)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add news"})
	}

	new_news.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(http.StatusCreated).JSON(new_news)

}

func Get_All_News(c *fiber.Ctx) error {

	var news_list []models.News

	cursor, err := database.NewsCollection.Find(context.Background(), bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var news models.News
		if err := cursor.Decode(&news); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error while decoding cursor": err.Error()})
		}
		news_list = append(news_list, news)
	}

	return c.Status(http.StatusOK).JSON(news_list)
}

func Get_News_By_ID(c *fiber.Ctx) error {

	var news models.News

	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	filter := bson.M{"_id": objectID}

	errr := database.NewsCollection.FindOne(context.Background(), filter).Decode(&news)

	if errr != nil {
		if errr == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "news not found"})
		}
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error while fetching news from database:" + err.Error()})
	}

	return c.Status(http.StatusOK).JSON(news)
}

func Delete_News(c *fiber.Ctx) error {

	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)

	_, adminObj, is_admin, err := Get_Ent(claims)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Token"})
	}

	if !is_admin {
		return c.Status(http.StatusMethodNotAllowed).JSON(fiber.Map{"error": "you must be an admin to access this method"})
	}

	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	filter := bson.M{"_id": objectID}

	var deleting_news models.News

	errr := database.NewsCollection.FindOne(context.Background(), filter).Decode(&deleting_news)

	if errr != nil {
		if errr == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "news not found"})
		}
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error while fetching news from database:" + errr.Error()})
	}

	if adminObj.ID != deleting_news.Author.ID {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "only author of news can delete it"})
	}

	deleteResult, err := database.NewsCollection.DeleteOne(context.Background(), filter)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error while deleting from database"})
	}

	if deleteResult.DeletedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "news was not found"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "news deleted succesfully"})
}
