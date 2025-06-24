package models

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID       bson.ObjectID `bson:"_id" json:"id"`
	Email    string        `bson:"email" unique:"true" json:"email"`
	Username string        `bson:"username" unique:"true" json:"username"`
	Password string        `bson:"password" json:"password"`
	Role     string        `bson:"role" json:"role"`
	Phone    string        `bson:"phone" json:"phone"`
	Branch   string        `bson:"branch" json:"branch"`
}
