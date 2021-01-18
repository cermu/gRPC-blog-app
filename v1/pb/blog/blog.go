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

	//// fetch author details to be returned within Blog message
	//var authorSlice []*Author
	//for _, email := range blog.GetWriterEmail() {
	//	filter := bson.M{"email": email}
	//	authorData := &models.Author{}
	//	addressData := &models.Address{}
	//
	//	results := utl.GetMongoDB().Collection("authors").FindOne(ctx, filter)
	//	if err := results.Decode(authorData); err != nil {
	//		return nil, status.Errorf(codes.NotFound,
	//			fmt.Sprintf("WARNING | Cannot find author with specified ID: %v", err))
	//	}
	//
	//	for _, addressId := range authorData.AddressID {
	//		addrOid, err := primitive.ObjectIDFromHex(addressId)
	//		if err != nil {
	//			return nil, status.Errorf(codes.InvalidArgument,
	//				fmt.Sprintf("ERROR | Converting string address id to oid failed: %v", err))
	//		}
	//
	//		filter := bson.M{"_id": addrOid}
	//		addrResults := utl.GetMongoDB().Collection("addresses").FindOne(ctx, filter)
	//		if err := addrResults.Decode(addressData); err != nil {
	//			return nil, status.Errorf(codes.NotFound,
	//				fmt.Sprintf("WARNING | Cannot find address with specified ID: %v", err))
	//		}
	//	}
	//	// For every author found, append them in a slice
	//	authorSlice = append(authorSlice, authorStructToPbAuthorMessage(authorData, addressData))
	//}
	return &CreateBlogResponse{
		Blog: &Blog{
			Id:          oid.Hex(),
			Title:       blog.GetTitle(),
			Content:     blog.GetContent(),
			WriterEmail: blog.GetWriterEmail(),
			Created:     creationTime.String(),
		},
	}, nil
}

// FetchBlog method
// implementing BlogServiceServer interface from blogApp.pb.go
func (*Server) FetchBlog(_ context.Context, req *ReadBlogRequest) (*ReadBlogResponse, error) {
	log.Printf("INFO | FetchBLog method was invoked with request: %v\n", req)

	// retrieve blog id from request
	blogId := req.GetBlogId()

	// convert blogId to oid
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("ERROR | Converting string blog id to oid failed: %v", err))
	}

	// fetch data from DB
	blogItem := &models.BlogItem{}

	ctx := context.Background()
	filter := bson.M{"_id": oid}
	results := utl.GetMongoDB().Collection("blog").FindOne(ctx, filter)

	if err := results.Decode(blogItem); err != nil {
		return nil, status.Errorf(codes.NotFound,
			fmt.Sprintf("WARNING | Cannot find blog item with specified ID: %v", err))
	}

	var authorList []*Author
	for _, authorEmail := range blogItem.WriterEmail {
		authorData := &models.Author{}
		addressData := &models.Address{}

		filter := bson.M{"email": authorEmail}
		authorResults := utl.GetMongoDB().Collection("authors").FindOne(ctx, filter)
		if err := authorResults.Decode(authorData); err != nil {
			return nil, status.Errorf(codes.NotFound,
				fmt.Sprintf("WARNING | Cannot find author with fetched email address: %v", err))
		}

		// fetch address details
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
					fmt.Sprintf("WARNING | Cannot find address with fetched address ID: %v", err))
			}
		}

		// append authors found in a slice
		authorList = append(authorList, authorStructToPbAuthorMessage(authorData, addressData))
	}
	return &ReadBlogResponse{
		Blog: &Blog{
			Id:            blogItem.ID.Hex(),
			Title:         blogItem.Title,
			Content:       blogItem.Content,
			WriterEmail:   blogItem.WriterEmail,
			AuthorDetails: authorList,
			Created:       blogItem.Created.String(),
			Updated:       blogItem.Updated.String(),
		},
	}, nil
}

// UpdateBlog method
// implementing BlogServiceServer interface from blogApp.pb.go
func (*Server) UpdateBlog(_ context.Context, req *UpdateBlogRequest) (*UpdateBlogResponse, error) {
	log.Printf("INFO | UpdateBlog method was invoked with request: %v\n", req)

	// retrieve data from request
	updateTime := time.Now()
	blogMessage := req.GetBlog()

	oid, err := primitive.ObjectIDFromHex(blogMessage.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("ERROR | Converting string blog id to oid failed: %v", err))
	}
	filter := bson.M{"_id": oid}

	// populate BlogItem struct to be used to update
	// a blog in DB
	blogItem := &models.BlogItem{}
	blogItem.ID = oid
	blogItem.Title = blogMessage.GetTitle()
	blogItem.Content = blogMessage.GetContent()
	blogItem.WriterEmail = blogMessage.GetWriterEmail()
	blogItem.Updated = updateTime

	doc := bson.M{"$set": blogItem}

	_, updateErr := utl.GetMongoDB().Collection("blog").UpdateOne(context.Background(), filter, doc)
	if updateErr != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("ERROR | Blog update failed with message: %v", updateErr))
	}
	return &UpdateBlogResponse{
		Blog: &Blog{
			Id:          blogItem.ID.Hex(),
			Title:       blogItem.Title,
			Content:     blogItem.Content,
			WriterEmail: blogItem.WriterEmail,
			Created:     blogItem.Created.String(),
			Updated:     blogItem.Updated.String(),
		},
	}, nil
}

// DeleteBlog method
// implementing BlogServiceServer interface from blogApp.pb.go
func (*Server) DeleteBlog(_ context.Context, req *DeleteBlogRequest) (*DeleteBlogResponse, error) {
	log.Printf("INFO | DeleteBlog method was invoked with request: %v\n", req)

	// retrieve blog id from the request
	blogId := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("ERROR | Converting string blog id to oid failed with message: %v", err))
	}

	filter := bson.M{"_id": oid}
	result, delErr := utl.GetMongoDB().Collection("blog").DeleteOne(context.Background(), filter)
	if delErr != nil {
		return nil, status.Errorf(codes.Internal,
			fmt.Sprintf("WARNING | Deleting blog failied with message: %v", delErr))
	}

	if result.DeletedCount > 0 {
		return &DeleteBlogResponse{
			DeleteResponse: fmt.Sprintf("Blog: %v was deleted successfully", blogId),
		}, nil
	}
	return &DeleteBlogResponse{
		DeleteResponse: "Nothing to delete",
	}, nil
}
