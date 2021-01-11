package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Author struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"first_name"`
	LastName  string             `bson:"last_name"`
	Email     string             `bson:"email"`
	AddressID []address          `bson:"address_id"`
	BlogID    []BlogItem         `bson:"blog_id"`
	Created   time.Time          `bson:"created"`
	Updated   time.Time          `bson:"dated,omitempty"`
}

type address struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	City          string             `bson:"city,omitempty"`
	Country       string             `bson:"country,omitempty"`
	ZipCode       string             `bson:"zip_code,omitempty"`
	PostalAddress string             `bson:"postal_address,omitempty"`
	Created       time.Time          `bson:"created"`
	Updated       time.Time          `bson:"dated,omitempty"`
}
