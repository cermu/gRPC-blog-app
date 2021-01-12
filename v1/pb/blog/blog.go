package blog

import (
	"context"
	"fmt"
	utl "github.com/cermu/gRPC-blog-app/utils"
	"github.com/cermu/gRPC-blog-app/v1/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type Server struct {
}

// CreateAuthor method
// implementing BlogServiceServer interface from blogApp.pb.go
func (*Server) CreateAuthor(_ context.Context, req *CreateAuthorRequest) (*CreateAuthorResponse, error) {
	log.Printf("INFO | CreateAuthor method was invoked with request: %v\n", req)

	// Fetch data from request
	creationTime := time.Now()
	author := req.GetAuthor()

	// extract address to be saved first in the DB
	addressData := models.Address{
		City:          author.GetAddress().GetCity(),
		Country:       author.GetAddress().GetCountry(),
		ZipCode:       author.GetAddress().GetZipCode(),
		PostalAddress: author.GetAddress().GetPostalAddress(),
		Created:       creationTime,
	}

	// Insert Address
	result, err := utl.GetMongoDB().Collection("addresses").InsertOne(context.Background(), addressData)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("ERROR | The following error occurred while inserting address into database: %v", err.Error()))
	}

	// fetch the inserted address ID
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("ERROR | An error occurred while converting inserted address id to object id"))
	}

	// Insert Author
	var addressSlice []string
	addressSlice = append(addressSlice, oid.Hex())
	authorData := models.Author{
		FirstName: author.GetFirstName(),
		LastName:  author.GetLastName(),
		Email:     author.GetEmail(),
		AddressID: addressSlice,
		Created:   creationTime,
	}
	result, err = utl.GetMongoDB().Collection("authors").InsertOne(context.Background(), authorData)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("ERROR | The following error occurred while inserting author into database: %v", err.Error()))
	}

	// fetch the inserted author ID
	oid, ok = result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("ERROR | An error occurred while converting inserted author id to object id"))
	}

	return &CreateAuthorResponse{
		Author: &Author{
			Id:        oid.Hex(),
			FirstName: author.GetFirstName(),
			LastName:  author.GetLastName(),
			Email:     author.GetEmail(),
			Address: &Address{
				City:          author.GetAddress().GetCity(),
				Country:       author.GetAddress().GetCountry(),
				ZipCode:       author.GetAddress().GetZipCode(),
				PostalAddress: author.GetAddress().GetPostalAddress(),
				Created:       creationTime.String(),
			},
			Created: creationTime.String(),
		},
	}, nil
}

// FetchAuthor method
// implementing BlogServiceServer interface from blogApp.pb.go
func (*Server) FetchAuthor(_ context.Context, req *FetchAuthorRequest) (*FetchAuthorResponse, error) {
	log.Printf("INFO | FetchAuthor method was invoked with request: %v\n", req)

	// retrieve author id from request
	authorId := req.GetAuthorId()

	// convert the id to oid
	oid, err := primitive.ObjectIDFromHex(authorId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("ERROR | Converting string author id to oid failed: %v", err))
	}

	//fetch data from DB
	addressData := &models.Address{}
	authorData := &models.Author{}

	filter := bson.M{"_id": oid}
	results := utl.GetMongoDB().Collection("authors").FindOne(context.Background(), filter)

	if err := results.Decode(authorData); err != nil {
		return nil, status.Errorf(codes.NotFound,
			fmt.Sprintf("WARNING | Cannot find author with specified ID: %v", err))
	}

	// fetch address
	for _, v := range authorData.AddressID {
		// convert fetched address id to oid
		oid, err = primitive.ObjectIDFromHex(v)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				fmt.Sprintf("ERROR | Converting string fetched address id to oid failed: %v", err))
		}

		addrResults := utl.GetMongoDB().Collection("addresses").FindOne(context.Background(), bson.M{"_id": oid})
		if err := addrResults.Decode(addressData); err != nil {
			return nil, status.Errorf(codes.NotFound,
				fmt.Sprintf("WARNING | Cannot find address with fetched ID: %v", err))
		}
	}

	return &FetchAuthorResponse{
		Author: &Author{
			Id: authorData.ID.Hex(),
			FirstName: authorData.FirstName,
			LastName: authorData.LastName,
			Email: authorData.Email,
			Address: &Address{
				City: addressData.City,
				Country: addressData.Country,
				ZipCode: addressData.ZipCode,
				PostalAddress: addressData.PostalAddress,
				Created: addressData.Created.String(),
				Updated: addressData.Updated.String(),
			},
			Created: authorData.Created.String(),
			Updated: authorData.Updated.String(),
		},
	}, nil
}
