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

// CreateBlog method
// implementing BlogServiceServer interface from blogApp.pb.go
func (*Server) CreateBlog(_ context.Context, req *CreateBlogRequest) (*CreateBlogResponse, error) {
	log.Printf("INFO | CreateBlog method was invoked with request: %v\n", req)

	// fetch data from request
	creationTime := time.Now()
	blog := req.GetBlog()

	// insert data in the DB
	blogData := &models.BlogItem{
		Title:       blog.GetTitle(),
		Content:     blog.GetContent(),
		WriterEmail: blog.GetWriterEmail(),
		Created:     creationTime,
	}

	ctx := context.Background()
	result, err := utl.GetMongoDB().Collection("blog").InsertOne(ctx, blogData)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("ERROR The following error occurred while inserting author into database: %v", err))
	}

	// fetch id of the inserted blog
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("ERROR | An error occurred while converting inserted author id to object id"))
	}

	// fetch author details to be returned within Blog message
	var authorSlice []*Author
	for _, email := range blog.GetWriterEmail() {
		filter := bson.M{"email": email}
		authorData := &models.Author{}
		addressData := &models.Address{}

		results := utl.GetMongoDB().Collection("authors").FindOne(ctx, filter)
		if err := results.Decode(authorData); err != nil {
			return nil, status.Errorf(codes.NotFound,
				fmt.Sprintf("WARNING | Cannot find author with specified ID: %v", err))
		}

		for _, addressId := range authorData.AddressID {
			addrOid, err := primitive.ObjectIDFromHex(addressId)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument,
					fmt.Sprintf("ERROR | Converting string address id to oid failed: %v", err))
			}

			filter := bson.M{"_id": addrOid}
			addrResults := utl.GetMongoDB().Collection("addresses").FindOne(ctx, filter)
			if err := addrResults.Decode(addressData); err != nil {
				return nil, status.Errorf(codes.NotFound,
					fmt.Sprintf("WARNING | Cannot find address with specified ID: %v", err))
			}
		}
		// For every author found, append them in a slice
		authorSlice = append(authorSlice, authorStructToPbAuthorMessage(authorData, addressData))
	}
	return &CreateBlogResponse{
		Blog: &Blog{
			Id:            oid.Hex(),
			Title:         blog.GetTitle(),
			Content:       blog.GetContent(),
			WriterEmail:   blog.GetWriterEmail(),
			AuthorDetails: authorSlice,
			Created:       creationTime.String(),
		},
	}, nil
}
