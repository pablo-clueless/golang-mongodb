package router

import (
	"golang-mongodb/common"
	"golang-mongodb/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/asaskevich/govalidator.v9"
)

type CreateUserDto struct {
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func AddAuthGroup(app *fiber.App) {
	authGroup := app.Group("/auth")

	authGroup.Post("/signup", Signup)
	authGroup.Post("/signin", Signin)
}

func Signup(c *fiber.Ctx) error {
	collection := common.GetDBCollection("users")

	var user *models.User
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}
	user.Email = common.NormalizeEmail(user.Email)
	if !govalidator.IsEmail(user.Email) {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid email",
		})
	}

	exists, err := collection.Find(c.Context(), bson.M{"email": user.Email})
	if exists == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "user exists",
		})
	}
	if err == mongo.ErrNilDocument {
		if strings.TrimSpace(user.Password) == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}
	user.Password = common.EncryptPassword(user.Password)
	user.CreatedAt = time.Now()
	result, err := collection.InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "unable to create user",
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data":    result,
		"message": "user created",
	})
}

func Signin(c *fiber.Ctx) error {
	collection := common.GetDBCollection("users")

	var credentials *models.User
	var user *models.User
	err := c.BodyParser(&credentials)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	credentials.Email = common.NormalizeEmail(credentials.Email)
	if !govalidator.IsEmail(credentials.Email) {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid email",
		})
	}

	e := collection.FindOne(c.Context(), bson.M{"email": credentials.Email}).Decode(&user)
	if e != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "user not found",
			"message": err.Error(),
		})
	}

	err = common.VerifyPassword(user.Password, credentials.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "invalid password",
			"message": err.Error(),
		})
	}

	token, err := common.NewToken(user.ID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "unable to log user in",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data":  user,
		"token": token,
	})
}
