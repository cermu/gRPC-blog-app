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
		Author: authorStructToPbAuthorMessage(authorData, addressData),
	}, nil
}

// UpdateAuthor method
// implementing BlogServiceServer interface from blogApp.pb.go
func (*Server) UpdateAuthor(_ context.Context, req *UpdateAuthorRequest) (*UpdateAuthorResponse, error) {
	log.Printf("INFO | UpdateAuthor method was invoked with request: %v\n", req)

	// retrieve data from request
	updatedTime := time.Now()
	author := req.GetAuthor()

	// fetch the address and author to be updated by their ids
	oid, err := primitive.ObjectIDFromHex(author.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("ERROR | Converting string author id to oid failed: %v", err))
	}

	addrOid, addrErr := primitive.ObjectIDFromHex(author.GetAddress().GetId())
	if addrErr != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("ERROR | Converting string address id to oid failed: %v", addrErr))
	}

	// address data
	addressData := &models.Address{}
	filter := bson.M{"_id": addrOid}
	addrResult := utl.GetMongoDB().Collection("addresses").FindOne(context.Background(), filter)
	if err := addrResult.Decode(addressData); err != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("WARNING | Cannot find address with specified address id: %v", err))
	}
	addressData.City = author.GetAddress().GetCity()
	addressData.Country = author.GetAddress().GetCountry()
	addressData.ZipCode = author.GetAddress().GetZipCode()
	addressData.PostalAddress = author.GetAddress().GetPostalAddress()
	addressData.Updated = updatedTime

	// update the address record in DB
	_, updateErr := utl.GetMongoDB().Collection("addresses").ReplaceOne(context.Background(), filter, addressData)
	if updateErr != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("ERROR | Address update failed with message: %v", updateErr))
	}

	// author data
	authorData := &models.Author{}
	filter = bson.M{"_id": oid}
	result := utl.GetMongoDB().Collection("authors").FindOne(context.Background(), filter)
	if err := result.Decode(authorData); err != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("WARNING | Cannot find author with specified author id: %v", err))
	}
	var addressSlice []string
	addressSlice = append(addressSlice, author.GetAddress().GetId())
	authorData.FirstName = author.GetFirstName()
	authorData.LastName = author.GetLastName()
	authorData.Email = author.GetEmail()
	authorData.AddressID = addressSlice
	authorData.Updated = updatedTime

	// update the author record in DB
	_, updateErr = utl.GetMongoDB().Collection("authors").ReplaceOne(context.Background(), filter, authorData)
	if updateErr != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("ERROR | Author update failed with message: %v", updateErr))
	}

	// return the updated author
	return &UpdateAuthorResponse{
		Author: authorStructToPbAuthorMessage(authorData, addressData),
	}, nil
}

// DeleteAuthor method
// implementing BlogServiceServer interface from blogApp.pb.go
func (*Server) DeleteAuthor(_ context.Context, req *DeleteAuthorRequest) (*DeleteAuthorResponse, error) {
	log.Printf("INFO | DeleteAuthor method was invoked with request: %v\n", req)

	// retrieve author id from request
	authorId := req.GetAuthorId()
	authorOid, err := primitive.ObjectIDFromHex(authorId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("ERROR | Converting string author id to oid failed: %v", err))
	}

	// fetch the address id to be deleted
	authorData := &models.Author{}
	filter := bson.M{"_id": authorOid}

	results := utl.GetMongoDB().Collection("authors").FindOne(context.Background(), filter)
	if err := results.Decode(authorData); err != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("WARNING | Cannot find author with specified author id: %v", err))
	}

	// delete the address first
	for _, v := range authorData.AddressID {
		addressOid, addrErr := primitive.ObjectIDFromHex(v)
		if addrErr != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				fmt.Sprintf("ERROR | Converting string address id to oid failed: %v", addrErr))
		}

		result, err := utl.GetMongoDB().Collection("addresses").DeleteOne(context.Background(), bson.M{"_id": addressOid})
		if err != nil {
			log.Printf("WARNING | Deleting address failied with message: %v\n", err)
			continue
		}

		if result.DeletedCount > 0 {
			log.Printf("INFO | Address: %v has been deleted successfully", v)
		}
	}

	// delete the author
	delResults, delErr := utl.GetMongoDB().Collection("authors").DeleteOne(context.Background(), filter)
	if delErr != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("WARNING | Deleting author failied with message: %v", delErr))
	}
	if delResults.DeletedCount > 0 {
		return &DeleteAuthorResponse{
			DeleteResponse: fmt.Sprintf("Author: %v has been deleted successfully", authorId),
		}, nil
	}

	return &DeleteAuthorResponse{
		DeleteResponse: "Nothing to delete",
	}, nil
}

// AllAuthors method, a server streaming scenario
func (*Server) AllAuthors(_ *AllAuthorsRequest, stream BlogService_AllAuthorsServer) error {
	log.Println("INFO | AllAuthors server streaming method was invoked")

	ctx := context.Background()
	filter := bson.M{} // pass an empty filter to fetch all records
	cur, err := utl.GetMongoDB().Collection("authors").Find(ctx, filter)
	if err != nil {
		return status.Errorf(codes.Internal,
			fmt.Sprintf("WARNING | Failed to fetch Authors from the DB: %v", err))
	}

	// close the cursor after use
	defer func () {
		if err := cur.Close(ctx); err != nil {
			log.Printf("WARNING | CLosing MongoDB cursor failed with message: %v\n", err)
		}
	}()

	for cur.Next(ctx) {
		authorData := &models.Author{}
		if err := cur.Decode(authorData); err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("ERROR | Decoding MongoDB object failed with message: %v", err))
		}

		// fetch addresses to be returned within the Author message response
		addressData := &models.Address{}
		for _, v := range authorData.AddressID {
			// convert fetched address id to oid
			oid, err := primitive.ObjectIDFromHex(v)
			if err != nil {
				return status.Errorf(codes.InvalidArgument,
					fmt.Sprintf("ERROR | Converting string fetched address id to oid failed: %v", err))
			}

			addrResults := utl.GetMongoDB().Collection("addresses").FindOne(context.Background(), bson.M{"_id": oid})
			if err := addrResults.Decode(addressData); err != nil {
				return status.Errorf(codes.NotFound,
					fmt.Sprintf("WARNING | Cannot find address with fetched ID: %v", err))
			}
		}

		// send Author response via a stream
		response := &AllAuthorsResponse{
			Author: authorStructToPbAuthorMessage(authorData, addressData),
		}

		streamErr := stream.Send(response)
		if streamErr != nil {
			log.Printf("WARNING | Sending server stream response failed with message: %v\n", streamErr)
		}
	}

	// handle cur.Err()
	if err := cur.Err(); err != nil {
		return status.Errorf(codes.Internal,
			fmt.Sprintf("ERROR | The following MongoDB cursor error has occurred: %v", err))
	}
	return nil
}

// authorStructToPbAuthorMessage private function that converts  Author struct
// to Author message in protobuf definition. It takes a pointer to Author and Address
// and returns a pointer to Author
func authorStructToPbAuthorMessage(author *models.Author, address *models.Address) *Author {
	return &Author{
		Id:        author.ID.Hex(),
		FirstName: author.FirstName,
		LastName:  author.LastName,
		Email:     author.Email,
		Address: &Address{
			Id:            address.ID.Hex(),
			City:          address.City,
			Country:       address.Country,
			ZipCode:       address.ZipCode,
			PostalAddress: address.PostalAddress,
			Created:       address.Created.String(),
			Updated:       address.Updated.String(),
		},
		Created: author.Created.String(),
		Updated: author.Updated.String(),
	}
}
