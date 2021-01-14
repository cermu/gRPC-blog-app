package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
	WriterEmail []string        `bson:"authors"`
	Created  time.Time          `bson:"created,omitempty"`
	Updated  time.Time          `bson:"updated,omitempty"`
}
