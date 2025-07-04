package auth

import (
	"context"
	"time"

	"github.com/aslon1213/go-pos-erp/pkg/configs"
	"github.com/aslon1213/go-pos-erp/pkg/middleware"
	models "github.com/aslon1213/go-pos-erp/pkg/repository"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"

	pasetoware "github.com/gofiber/contrib/paseto"
	"github.com/rs/zerolog/log"
)

// AuthControllers handles authentication-related operations
type AuthControllers struct {
	UserCollection       *mongo.Collection
	ActivitiesCollection *mongo.Collection
	SecretSymmetricKey   string
}

// New initializes a new AuthControllers instance
func New(db *mongo.Database) *AuthControllers {

	users_collection := db.Collection("users")
	users_collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"username": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	)

	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	return &AuthControllers{
		UserCollection:       users_collection,
		ActivitiesCollection: db.Collection("activities"),
		SecretSymmetricKey:   config.Server.SecretSymmetricKey,
	}
}

// User represents a user in the system
type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Info handles user info
// @Summary Get user info
// @Security BearerAuth
// @Description Get user info
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "message"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/auth/me [get]
func (a *AuthControllers) InfoMe(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}

// Login handles user login
// @Summary Login a user
// @Description Authenticate user and return a token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body LoginInput true "User credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/login [post]
func (a *AuthControllers) Login(c *fiber.Ctx) error {
	var user_to_check LoginInput

	if err := c.BodyParser(&user_to_check); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if user_to_check.Username == "" || user_to_check.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}
	middleware.SetActionType(c, middleware.ActivityTypeLogin)
	middleware.SetUser(c, user_to_check.Username)
	middleware.LogActivity(c)
	pass := user_to_check.Password

	// check in the database if the user exists
	user_db := models.User{}
	err := a.UserCollection.FindOne(c.Context(), bson.M{"username": user_to_check.Username}).Decode(&user_db)

	if err != nil {
		log.Error().Err(err).Msg("Failed to find user")
		// no user found
		middleware.SetData(c, fiber.Map{
			"error": "User not found",
		})
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	equal := bcrypt.CompareHashAndPassword([]byte(user_db.Password), []byte(pass))
	if equal != nil {
		log.Warn().Msg("Unauthorized access attempt")
		middleware.SetData(c, fiber.Map{
			"error": "Invalid password",
		})
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create token and encrypt it
	encryptedToken, err := pasetoware.CreateToken([]byte(a.SecretSymmetricKey), user_to_check.Username, 48*time.Hour, pasetoware.PurposeLocal)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create token")
		middleware.DontLogActivity(c)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	payload, err := pasetoware.NewPayload(
		encryptedToken,
		48*time.Hour,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create payload")
		middleware.DontLogActivity(c)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Info().Str("username", user_to_check.Username).Msg("User logged in successfully")
	return c.JSON(payload)
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User credentials"
// @Success 201 {string} string "Created"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/register [post]
func (a *AuthControllers) Register(c *fiber.Ctx) error {
	// Parse user input
	var user models.UserRegisterInput

	err := c.BodyParser(&user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse user")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// validate the user

	// log the action
	middleware.SetActionType(c, middleware.ActivityTypeRegister)
	middleware.SetUser(c, user.Username)
	middleware.SetData(
		c,
		map[string]string{
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"phone":    user.Phone,
			"branch":   user.Branch,
		},
	)
	middleware.LogActivity(c)

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		middleware.DontLogActivity(c)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Create a new user document
	newUser := models.NewUser(user.Username, string(hashedPassword), user.Email, user.Role, user.Phone, user.Branch)

	// Insert the new user into the database
	_, err = a.UserCollection.InsertOne(c.Context(), newUser)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register user")
		middleware.DontLogActivity(c)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Info().Str("username", user.Username).Msg("User registered successfully")
	return c.SendStatus(fiber.StatusCreated)
}
