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
	createAuthor(c)
}

// createAuthor private function that calls CreateAuthor gRPC method on gRPC server
func createAuthor(c blog.BlogServiceClient) {
	log.Println("INFO | Starting a CreateAuthor unary gRPC")

	address := &blog.Address{
		City: "Nairobi",
		Country: "Kenya",
		ZipCode: "00101",
		PostalAddress: "243",
	}

	author := &blog.Author{
		FirstName: "John",
		LastName: "Doe",
		Email: "jdoe@gmail.com",
		Address: address,
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
