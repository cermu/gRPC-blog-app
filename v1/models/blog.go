package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
	AuthorID []string           `bson:"author_id"`
	Created  time.Time          `bson:"created,omitempty"`
	Updated  time.Time          `bson:"dated,omitempty"`
}
