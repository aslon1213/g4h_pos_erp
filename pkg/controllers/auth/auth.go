package auth

import (
	"context"
	"time"

	"github.com/aslon1213/go-pos-erp/pkg/configs"
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
	UserCollection     *mongo.Collection
	SecretSymmetricKey string
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
		UserCollection:     users_collection,
		SecretSymmetricKey: config.Server.SecretSymmetricKey,
	}
}

// User represents a user in the system
type LoginInput struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

// Login handles user login
// @Summary Login a user
// @Description Authenticate user and return a token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body User true "User credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/login [post]
func (a *AuthControllers) Login(c *fiber.Ctx) error {
	var user_to_check LoginInput

	err := c.BodyParser(&user_to_check)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse user")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	pass := user_to_check.Password

	// check in the database if the user exists
	user_db := models.User{}
	err = a.UserCollection.FindOne(c.Context(), bson.M{"username": user_to_check.Username}).Decode(&user_db)

	if err != nil {
		log.Error().Err(err).Msg("Failed to find user")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	equal := bcrypt.CompareHashAndPassword([]byte(user_db.Password), []byte(pass))
	if equal != nil {
		log.Warn().Msg("Unauthorized access attempt")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create token and encrypt it
	encryptedToken, err := pasetoware.CreateToken([]byte(a.SecretSymmetricKey), user_to_check.Username, 48*time.Hour, pasetoware.PurposeLocal)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create token")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	payload, err := pasetoware.NewPayload(
		encryptedToken,
		48*time.Hour,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create payload")
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

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Create a new user document
	newUser := models.NewUser(user.Username, string(hashedPassword), user.Email, user.Role, user.Phone, user.Branch)

	// Insert the new user into the database
	_, err = a.UserCollection.InsertOne(c.Context(), newUser)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register user")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Info().Str("username", user.Username).Msg("User registered successfully")
	return c.SendStatus(fiber.StatusCreated)
}
