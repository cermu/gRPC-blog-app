package main

import (
	"fmt"
	cnf "github.com/cermu/gRPC-blog-app/conf"
	utl "github.com/cermu/gRPC-blog-app/utils"
	"github.com/cermu/gRPC-blog-app/v1/pb/blog"
	_ "github.com/cermu/gRPC-blog-app/v1/pb/blog"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// get file name and line number when the code crashes
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// create a listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cnf.GetAppConfigs().ApplicationPort))
	if err != nil {
		log.Fatalf("ERROR | Failed to create a listener with message: %v\n", err.Error())
	}

	// create a gRPC server
	blogServer := &blog.Server{} // can be written as: new(blog.Server), new returns pointers

	opts := []grpc.ServerOption{}
	gRPCServer := grpc.NewServer(opts...)
	blog.RegisterBlogServiceServer(gRPCServer, blogServer)

	// start the server in a go routine
	go func() {
		log.Printf("INFO | Starting gRPC server on port: %d\n", cnf.GetAppConfigs().ApplicationPort)
		// bind the listener to gRPC server
		if err := gRPCServer.Serve(lis); err != nil {
			log.Fatalf("ERROR | Failed to start gRPC server with message: %v\n", err.Error())
		}
	}()

	// shutting down the application gracefully
	// wait for ctrl + c  signal to exit
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// block until a signal is received
	receivedSignal := <-ch

	log.Printf("INFO | Stopping gRPC server, signal: %v received.\n", receivedSignal.String())
	gRPCServer.Stop()

	if lis != nil {
		log.Println("INFO | Closing the listener...")
		_ = lis.Close()
	}
	utl.MongoClientDisconnect()
	log.Println("INFO | Application has shutdown.")
}
