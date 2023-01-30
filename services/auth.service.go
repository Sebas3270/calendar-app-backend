package services

import (
	"context"
	"log"

	"github.com/Sebas3270/calendar-app-backend/db"
	"github.com/Sebas3270/calendar-app-backend/helpers"
	"github.com/Sebas3270/calendar-app-backend/middlewares"
	"github.com/Sebas3270/calendar-app-backend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func getUserById(id string) (models.User, bool) {

	var user models.User

	objID, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{
		{Key: "_id", Value: objID},
	}

	if err := db.UserCollection.FindOne(context.TODO(), filter).Decode(&user); err != nil {

		// This error means your query did not match any documents.
		if err == mongo.ErrNoDocuments {
			return user, false
		}
	}

	return user, true

}

func getUserByEmail(email string) (models.User, bool) {

	var user models.User

	filter := bson.D{
		{Key: "email", Value: email},
	}

	if err := db.UserCollection.FindOne(context.TODO(), filter).Decode(&user); err != nil {

		// This error means your query did not match any documents.
		if err == mongo.ErrNoDocuments {
			return user, false
		}
	}

	return user, true

}

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	errors := helpers.ValidateUserStruct(*user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})

	}

	if _, found := getUserByEmail(user.Email); found == true {
		return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
			"error": "Email already registered",
		})
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(passwordHashed)
	user.Id = primitive.NewObjectID()

	res, err := db.UserCollection.InsertOne(context.TODO(), user)

	insertedId := res.InsertedID.(primitive.ObjectID).Hex()
	token, err := helpers.GenerateJWT(insertedId, user.Name)

	if err != nil {
		log.Fatal("error here", err)
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"name": user.Name,
			"id":   res.InsertedID,
		},
		"token": token,
	})
}

func Login(c *fiber.Ctx) error {

	logInBody := new(models.LogIn)
	if err := c.BodyParser(logInBody); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := helpers.ValidateLogInStruct(*logInBody)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	// filter := bson.D{
	// 	{Key: "email", Value: logInBody.Email},
	// }

	// var user models.User

	// err := db.UserCollection.FindOne(context.TODO(), filter).Decode(&user)
	// if err != nil {

	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	// 		"errors": "User not found",
	// 	})

	// }

	user, found := getUserByEmail(logInBody.Email)
	if found == false {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Incorrect credentials",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(logInBody.Password))

	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Incorrect credentials",
		})

	}

	insertedId := user.Id.Hex()
	token, err := helpers.GenerateJWT(insertedId, user.Name)

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"name": user.Name,
			"id":   user.Id,
		},
		"token": token,
	})
}

func RenewToken(c *fiber.Ctx) error {

	userId, userName, err := middlewares.GetTokenInfo(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	newToken, err := helpers.GenerateJWT(userId, userName)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	user, found := getUserById(userId)
	if found == false {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Incorrect credentials",
		})
	}

	return c.JSON(fiber.Map{
		"token": newToken,
		"user": fiber.Map{
			"name": user.Name,
			"id":   user.Id,
		},
	})
}
