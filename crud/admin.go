package crud

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"net/http"
	"tarafdari-sample/config"
	"tarafdari-sample/database"
	"tarafdari-sample/models"
	"time"

	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Add_Admin(c *fiber.Ctx) error {

	var new_admin models.Admin

	if err := c.BodyParser(&new_admin); err != nil {
		return err
	}

	if new_admin.Password == "" || new_admin.Username == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "all fields must be filled"})
	}

	hash := sha256.New()
	hash.Write([]byte(new_admin.Password))
	hashedInByte := hash.Sum(nil)
	new_admin.Password = hex.EncodeToString(hashedInByte)

	insertResult, err := database.AdminCollection.InsertOne(context.Background(), new_admin)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "username already taken"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register admin"})
	}

	new_admin.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(http.StatusCreated).JSON(new_admin)
}

func Login_Admin(c *fiber.Ctx) error {

	var loginRequest models.LoginRequest

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var admin models.Admin

	hash := sha256.New()
	hash.Write([]byte(loginRequest.Password))
	hashedInByte := hash.Sum(nil)
	loginRequest.Password = hex.EncodeToString(hashedInByte)

	filter := bson.M{"username": loginRequest.Username, "password": loginRequest.Password}
	err := database.AdminCollection.FindOne(context.Background(), filter).Decode(&admin)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "admin not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to find admin"})
	}

	claims := jtoken.MapClaims{
		"id":       admin.ID,
		"is_admin": true,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	// log.Println("secret key in function is:", config.JWT_SECRET_KEY)

	encoded_token, err := token.SignedString([]byte(config.JWT_SECRET_KEY))

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(models.LoginResponse{Message: "Login Succesfull", Token: encoded_token})

}
