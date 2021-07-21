package main

import (
	"flag"
	"fmt"
	"github.com/sakiib/grpc-gateway-demo/gateway"
	"github.com/sakiib/grpc-gateway-demo/gen/go/proto"
	"github.com/sakiib/grpc-gateway-demo/insecure"
	"github.com/sakiib/grpc-gateway-demo/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	port := flag.String("port", "", "server port")
	flag.Parse()

	log.Printf("starting the server on port: %s", *port)

	grpcServer := grpc.NewServer(grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)))
	bookServer := service.NewBookService(service.NewInMemStore())
	reflection.Register(grpcServer)

	pb.RegisterBookServiceServer(grpcServer, bookServer)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", *port))
	if err != nil {
		log.Fatalf("failed to listen with: %s", err.Error())
	}

	go func() {
		log.Fatal("grpc server crashed: ", grpcServer.Serve(listener))
	}()

	err = gateway.Run("dns:///" + fmt.Sprintf("0.0.0.0:%s", *port))
	log.Fatal("gateway crashed: ", err)
}