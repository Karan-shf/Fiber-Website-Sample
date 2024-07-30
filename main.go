package main

import (
	"context"
	"log"
	"tarafdari-sample/config"
	"tarafdari-sample/crud"
	"tarafdari-sample/database"
	"tarafdari-sample/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("error while loading .env file:", err)
	}

	database.DeployCollections()
	config.LoadSecretKey()

	defer database.Client.Disconnect(context.Background())

	app := fiber.New()

	Auth := middlewares.AuthMiddleware(config.JWT_SECRET_KEY)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, http://localhost:3000",
		AllowHeaders: "Origin,Content-Type,Accept",
		AllowMethods: "*",
	}))
	app.Use(logger.New())

	app.Post("/user/register", crud.Add_User)

	app.Post("/user/login", crud.Login_User)

	app.Get("/user", Auth, crud.Get_User)

	app.Post("/admin/register", crud.Add_Admin)

	app.Post("/admin/login", crud.Login_Admin)

	app.Post("/news", Auth, crud.Add_News)

	app.Get("/news", Auth, crud.Get_All_News)

	log.Fatal(app.Listen(":8080"))
}
