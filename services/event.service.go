package services

import (
	"context"
	"errors"

	"github.com/Sebas3270/calendar-app-backend/db"
	"github.com/Sebas3270/calendar-app-backend/helpers"
	"github.com/Sebas3270/calendar-app-backend/middlewares"
	"github.com/Sebas3270/calendar-app-backend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getEventById(eventId string) (models.Event, error) {

	var event models.Event

	eventIdParsed, _ := primitive.ObjectIDFromHex(eventId)

	filter := bson.D{
		{Key: "_id", Value: eventIdParsed},
	}

	if err := db.EventsCollection.FindOne(context.TODO(), filter).Decode(&event); err != nil {
		return event, err
	}

	return event, nil

}

func parseBodyEvent(c *fiber.Ctx) (models.Event, error) {

	event := new(models.Event)
	if err := c.BodyParser(event); err != nil {
		return *event, err
	}

	err := helpers.ValidateEventStruct(*event)
	if err != nil {
		return *event, errors.New("Could not validate event")
	}

	return *event, nil

}

func checkEventCreator(c *fiber.Ctx, event models.Event) bool {
	userId, _, _ := middlewares.GetTokenInfo(c)
	return event.User == userId
}

func GetEvents(c *fiber.Ctx) error {

	userId, _, _ := middlewares.GetTokenInfo(c)
	var events []models.Event

	filter := bson.D{
		{Key: "user", Value: userId},
	}

	cursor, err := db.EventsCollection.Find(context.TODO(), filter)

	if err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": "Error getting events",
		})
	}

	if err = cursor.All(context.TODO(), &events); err != nil {
		panic(err)
	}

	return c.JSON(events)
}

func CreateEvent(c *fiber.Ctx) error {

	event, err := parseBodyEvent(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	userId, _, _ := middlewares.GetTokenInfo(c)
	event.User = userId
	event.Id = primitive.NewObjectID()

	_, err = db.EventsCollection.InsertOne(context.TODO(), event)

	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.JSON(event)
}

func UpdateEvent(c *fiber.Ctx) error {

	// Getting id from params
	params := c.AllParams()
	eventId, ok := params["id"]

	if ok == false {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Check request params",
		})
	}

	/*Checking if event exists*/
	event, err := getEventById(eventId)

	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event do not exist",
		})
	}

	// Checking the event was made by the person requesting the action
	if different := checkEventCreator(c, event); different == false {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You can not manage this event",
		})
	}

	// Making update request
	eventIdParsed, _ := primitive.ObjectIDFromHex(eventId)

	filter := bson.D{
		{Key: "_id", Value: eventIdParsed},
	}

	event, err = parseBodyEvent(c) //Parsing into event the body of request

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	updateEvent := bson.D{{Key: "$set", Value: event}}

	_, err = db.EventsCollection.UpdateOne(context.TODO(), filter, updateEvent)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": err,
		})
	}

	return c.JSON(event)
}

func DeleteEvent(c *fiber.Ctx) error {

	params := c.AllParams()
	eventId, ok := params["id"]

	if ok == false {
		c.Status(400).JSON(fiber.Map{
			"error": "Check request params",
		})
	}

	event, err := getEventById(eventId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Event do not exist",
		})
	}

	if different := checkEventCreator(c, event); different == false {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You can not manage this event",
		})
	}

	eventIdParsed, _ := primitive.ObjectIDFromHex(eventId)

	filter := bson.D{
		{Key: "_id", Value: eventIdParsed},
	}

	//TODO: Avoid to edit an event you did not do

	_, err = db.EventsCollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": err,
		})
	}

	return c.JSON(fiber.Map{
		"msg": "Event deleted successfully",
	})
}
