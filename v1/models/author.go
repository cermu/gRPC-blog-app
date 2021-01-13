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
	AddressID []string          `bson:"address_id,omitempty"`
	// BlogID    []string         `bson:"blog_id,omitempty"`
	Created   time.Time          `bson:"created,omitempty"`
	Updated   time.Time          `bson:"updated,omitempty"`
}

type Address struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	City          string             `bson:"city,omitempty"`
	Country       string             `bson:"country,omitempty"`
	ZipCode       string             `bson:"zip_code,omitempty"`
	PostalAddress string             `bson:"postal_address,omitempty"`
	Created       time.Time          `bson:"created,omitempty"`
	Updated       time.Time          `bson:"updated,omitempty"`
}
