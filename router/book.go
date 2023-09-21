package router

import (
	"golang-mongodb/common"
	"golang-mongodb/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateBookDto struct {
	Title  string `json:"title" bson:"title"`
	Author string `json:"author" bson:"author"`
	Year   string `json:"year" bson:"year"`
}

type UpdateBookDto struct {
	Title  string `json:"title,omitempty" bson:"title,omitempty"`
	Author string `json:"author,omitempty" bson:"author,omitempty"`
	Year   string `json:"year,omitempty" bson:"year,omitempty"`
}

func AddBookGroup(app *fiber.App) {
	bookGroup := app.Group("/books")

	bookGroup.Post("/", Create)
	bookGroup.Patch("/:id", Update)
	bookGroup.Get("/", FindAll)
	bookGroup.Get("/:id", FindOne)
	bookGroup.Delete("/:id", Delete)
}

func Create(c *fiber.Ctx) error {
	collection := common.GetDBCollection("books")

	b := new(CreateBookDto)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	book, err := collection.InsertOne(c.Context(), b)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "failed to created book",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data":    book,
		"message": "book added",
	})
}

func Update(c *fiber.Ctx) error {
	collection := common.GetDBCollection("books")

	b := new(UpdateBookDto)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	book, err := collection.UpdateOne(c.Context(), bson.M{"_id": objectId}, bson.M{"$set": b})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "failed to update book",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data":    book,
		"message": "book updated",
	})
}

func FindAll(c *fiber.Ctx) error {
	collection := common.GetDBCollection("books")

	books := make([]models.Book, 0)
	cursor, err := collection.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for cursor.Next(c.Context()) {
		book := models.Book{}
		err := cursor.Decode(&book)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		books = append(books, book)
	}

	return c.Status(200).JSON(fiber.Map{
		"data":    books,
		"message": "all books found",
	})
}

func FindOne(c *fiber.Ctx) error {
	collection := common.GetDBCollection("books")

	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	book := models.Book{}
	err = collection.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&book)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data":    book,
		"message": "book found",
	})
}

func Delete(c *fiber.Ctx) error {
	collection := common.GetDBCollection("books")

	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	result, err := collection.DeleteOne(c.Context(), bson.M{"_id": objectId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "failed to delete book",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data":    result,
		"message": "book deleted",
	})
}
