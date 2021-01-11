package utils

import (
	"context"
	"fmt"
	cnf "github.com/cermu/gRPC-blog-app/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var mongoDatabase *mongo.Database
var mongoClient *mongo.Client
var ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)

func init() {
	// Get filename and line number whenever the code crashes
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("INFO | Initializing mongo database connection...")

	// create a mongo client
	mongoURI := fmt.Sprintf("mongodb://%s:%d", cnf.GetAppConfigs().MongoDBHost, cnf.GetAppConfigs().MongoDBPort)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("ERROR | Failed to create mongo client with message: %v\n", err.Error())
	}

	mongoClient = client
	connectErr := mongoClient.Connect(ctx)
	if connectErr != nil {
		log.Fatalf("ERROR | Mongo client failed to connect with message: %v\n", connectErr.Error())
	}

	mongoDatabase = mongoClient.Database(cnf.GetAppConfigs().MongoDBName)
	log.Printf("INFO | Initializing mongo database connection \t [OK]")
}

// GetMongoDB public function that exposes a private mongo Database pointer
func GetMongoDB () *mongo.Database {
	return mongoDatabase
}

// MongoClientDisconnect public function which will be used to close the mongo client
func MongoClientDisconnect() {
	defer cancel()

	log.Println("INFO | Disconnecting mongo client...")
	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Fatalf("ERROR | Mongo client failed to disconnect with message: %v\v", err.Error())
	}
}
