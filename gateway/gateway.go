package gateway

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/sakiib/grpc-gateway-demo/insecure"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	pb "github.com/sakiib/grpc-gateway-demo/gen/go/proto"
)

// Run runs the gRPC-Gateway, dialling the provided address.
func Run(dialAddr string) error {
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, err := grpc.DialContext(
		context.Background(),
		dialAddr,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(insecure.CertPool, "")),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterBookServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	gatewayAddr := "0.0.0.0:" + port
	gwServer := &http.Server{
		Addr: gatewayAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/v1") {
				gwmux.ServeHTTP(w, r)
				return
			}
		}),
	}
	// Empty parameters mean use the TLS Config specified with the server.
	if strings.ToLower(os.Getenv("SERVE_HTTP")) == "true" {
		log.Info("Serving gRPC-Gateway on http://", gatewayAddr)
		return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServe())
	}

	gwServer.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{insecure.Cert},
	}
	log.Info("Serving gRPC-Gateway on https://", gatewayAddr)
	return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServe())
}