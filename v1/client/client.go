package main

import (
	"context"
	"fmt"
	cnf "github.com/cermu/gRPC-blog-app/conf"
	"github.com/cermu/gRPC-blog-app/v1/pb/blog"
	"google.golang.org/grpc"
	"log"
)

func main() {
	// get file name and line number when the code crashes
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("INFO | Starting gRPC blog App client...")

	// dialing connection to a gRPC server
	gRPCHost := fmt.Sprintf("%s:%d", cnf.GetAppConfigs().GrpcHost, cnf.GetAppConfigs().ApplicationPort)
	opts := grpc.WithInsecure()
	conn, err := grpc.Dial(gRPCHost, opts)

	if err != nil {
		log.Fatalf("ERROR | Dialing gRPC server failed with message: %v\n", err.Error())
	}

	// close connection after use
	defer conn.Close()

	// Implementing BlogServiceClient interface from blog.pb.go
	c := blog.NewBlogServiceClient(conn)
	//createAuthor(c)
	// fetchAuthor(c)
	updateAuthor(c)
}

// createAuthor private function that calls CreateAuthor gRPC method on gRPC server
func createAuthor(c blog.BlogServiceClient) {
	log.Println("INFO | Starting a CreateAuthor unary gRPC")

	address := &blog.Address{
		City:          "Nairobi",
		Country:       "Kenya",
		ZipCode:       "00101",
		PostalAddress: "243",
	}

	author := &blog.Author{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "jdoe@gmail.com",
		Address:   address,
	}

	createAuthorReq := &blog.CreateAuthorRequest{
		Author: author,
	}
	response, err := c.CreateAuthor(context.Background(), createAuthorReq)

	if err != nil {
		log.Println(err)
	}
	log.Printf("INFO | Response from CreateAuthor gRPC: %v\n", response.GetAuthor())
}

// fetchAuthor private function that calls CreateAuthor gRPC method on gRPC server
func fetchAuthor(c blog.BlogServiceClient) {
	log.Println("INFO | Starting a FetchAuthor unary gRPC")

	fetchAuthorReq := &blog.FetchAuthorRequest{
		AuthorId: "5ffd6e13de19e0a9724b9b20",
	}

	response, err := c.FetchAuthor(context.Background(), fetchAuthorReq)
	if err != nil {
		log.Println(err)
	}
	log.Printf("INFO | Response from FetchAuthor gRPC: %v\n", response.GetAuthor())
}

// updateAuthor private function that calls UpdateAuthor gRPC method on gRPC server
func updateAuthor(c blog.BlogServiceClient) {
	log.Println("INFO | Starting an UpdateAuthor unary gRPC")

	address := &blog.Address{
		Id:            "5ffd6b06bd02eeada4eadd0f",
		City:          "Mombasa",
		Country:       "Kenya",
		ZipCode:       "00570",
		PostalAddress: "324",
	}

	author := &blog.Author{
		Id:        "5ffd6b06bd02eeada4eadd10",
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "janedoe@yahoo.com",
		Address:   address,
	}

	updateAuthorReq := &blog.UpdateAuthorRequest{
		Author: author,
	}
	response, err := c.UpdateAuthor(context.Background(), updateAuthorReq)

	if err != nil {
		log.Println(err)
	}
	log.Printf("INFO | Response from UpdateAuthor gRPC: %v\n", response.GetAuthor())
}
