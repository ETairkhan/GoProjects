package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`    // Use primitive.ObjectID for the ID
	Name   string             `json:"name" bson:"name"` // Fix bson tags to match the field names
	Gender string             `json:"gender" bson:"gender"`
	Age    int                `json:"age" bson:"age"`
}
