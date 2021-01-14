package main

import (
	"context"
	"fmt"
	cnf "github.com/cermu/gRPC-blog-app/conf"
	"github.com/cermu/gRPC-blog-app/v1/pb/blog"
	"google.golang.org/grpc"
	"io"
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
	// createAuthor(c)
	// fetchAuthor(c)
	// updateAuthor(c)
	// deleteAuthor(c)
	// allAuthors(c)
	// createBlog(c)
	fetchBlog(c)
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

// fetchAuthor private function that calls FetchAuthor gRPC method on gRPC server
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
		Id:            "5ffed8580dfec66815713941",
		City:          "Mombasa",
		Country:       "Kenya",
		ZipCode:       "00570",
		PostalAddress: "324",
	}

	author := &blog.Author{
		Id:        "5ffed8580dfec66815713942",
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

// deleteAuthor private function that calls DeleteAuthor gRPC method on gRPC server
func deleteAuthor(c blog.BlogServiceClient) {
	log.Println("INFO | Starting a DeleteAuthor unary gRPC")

	deleteAuthorRequest := &blog.DeleteAuthorRequest{
		AuthorId: "5ffd6e13de19e0a9724b9b20",
	}

	response, err := c.DeleteAuthor(context.Background(), deleteAuthorRequest)
	if err != nil {
		// log.Println(err)
		log.Fatal(err)
	}
	log.Printf("INFO | Response from DeleteAuthor gRPC: %v\n", response.GetDeleteResponse())
}

// allAuthors private function that streams AllAuthors gRPC method on gRPC server
func allAuthors(c blog.BlogServiceClient) {
	log.Printf("INFO | Starting an AllAuthors server streaming gRPC...\n")

	allAuthorReq := &blog.AllAuthorsRequest{}

	stream, err := c.AllAuthors(context.Background(), allAuthorReq)
	if err != nil {
		log.Fatal(err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
		}
		log.Printf("INFO | Response from AllAuthors server stream gRPC: %v\n", response.GetAuthor())
	}
	log.Println("INFO | Streaming Authors from AllAuthors gRPC has completed")
}

// createBlog private function that calls CreateBlog gRPC method on gRPC server
func createBlog(c blog.BlogServiceClient) {
	log.Println("INFO | Starting a CreateBlog unary gRPC")

	blogData := &blog.Blog{
		Title: "Golang in 30 mins",
		Content: "30 minutes of building a golang micro service",
		WriterEmail: []string{"jdoe@gmail.com"},
	}

	createBlogReq := &blog.CreateBlogRequest{
		Blog: blogData,
	}

	response, err := c.CreateBlog(context.Background(), createBlogReq)

	if err != nil {
		log.Println(err)
	}
	log.Printf("INFO | Response from CreateBlog gRPC: %v\n", response.GetBlog())
}

// fetchBlog private function that calls FetchBlog gRPC method on gRPC server
func fetchBlog(c blog.BlogServiceClient) {
	log.Println("INFO | Starting a FetchBlog unary gRPC")

	fetchBlogReq := &blog.ReadBlogRequest{
		BlogId: "600021457ca951f000beda40",
	}

	response, err := c.FetchBlog(context.Background(), fetchBlogReq)
	if err != nil {
		log.Println(err)
	}
	log.Printf("INFO | Response from FetchBlog gRPC: %v\n", response.GetBlog())
}
