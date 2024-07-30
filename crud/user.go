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

func Add_User(c *fiber.Ctx) error {

	var new_user models.User

	if err := c.BodyParser(&new_user); err != nil {
		return err
	}

	if new_user.Age == 0 || new_user.Firstname == "" || new_user.Lastname == "" || new_user.Password == "" || new_user.Username == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "all fields must be filled"})
	}

	hash := sha256.New()
	hash.Write([]byte(new_user.Password))
	hashedInByte := hash.Sum(nil)
	new_user.Password = hex.EncodeToString(hashedInByte)

	insertResult, err := database.UserCollection.InsertOne(context.Background(), new_user)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "username already taken"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
	}

	new_user.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(http.StatusCreated).JSON(new_user)

}

func Login_User(c *fiber.Ctx) error {

	var loginRequest models.LoginRequest

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User

	hash := sha256.New()
	hash.Write([]byte(loginRequest.Password))
	hashedInByte := hash.Sum(nil)
	loginRequest.Password = hex.EncodeToString(hashedInByte)

	filter := bson.M{"username": loginRequest.Username, "password": loginRequest.Password}
	err := database.UserCollection.FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to find user"})
	}

	claims := jtoken.MapClaims{
		"id":       user.ID,
		"is_admin": false,
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

func Get_User(c *fiber.Ctx) error {
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	// claim_id := claims["id"].(string)
	// id, _ := primitive.ObjectIDFromHex(claim_id)
	// is_admin := claims["is_admin"].(bool)
	// if !is_admin {
	// 	var user models.User
	// 	filter := bson.M{"_id": id}
	// 	database.UserCollection.FindOne(context.Background(), filter).Decode(&user)
	// 	return c.Status(http.StatusOK).JSON(user)
	// } else {
	// 	var admin models.Admin
	// 	filter := bson.M{"_id": id}
	// 	database.AdminCollection.FindOne(context.Background(), filter).Decode(&admin)
	// 	return c.Status(http.StatusOK).JSON(admin)
	// }
	userObj, adminObj, is_admin, err := Get_Ent(claims)

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Token"})
	}

	if is_admin {
		return c.Status(http.StatusOK).JSON(adminObj)
	}
	return c.Status(http.StatusOK).JSON(userObj)

}

func Get_Ent(claims jtoken.MapClaims) (*models.User, *models.Admin, bool, error) {

	claim_id := claims["id"].(string)
	id, _ := primitive.ObjectIDFromHex(claim_id)
	is_admin := claims["is_admin"].(bool)
	filter := bson.M{"_id": id}
	var user models.User
	var admin models.Admin
	var err error
	if !is_admin {
		err = database.UserCollection.FindOne(context.Background(), filter).Decode(&user)
		// return c.Status(http.StatusOK).JSON(user)
	} else {
		// filter := bson.M{"_id": id}
		err = database.AdminCollection.FindOne(context.Background(), filter).Decode(&admin)
		// return c.Status(http.StatusOK).JSON(admin)
	}
	return &user, &admin, is_admin, err

}
