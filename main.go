package main

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "github.com/InstaUpload/common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getUserService(ctx context.Context, addr string) (pb.UserServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	userService := pb.NewUserServiceClient(conn)
	return userService, conn, nil
}

func run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         "5000",
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 40,
		IdleTimeout:  time.Minute,
	}
	log.Println("Http server running in port :5000")

	return srv.ListenAndServe()
}

func main() {
	ctx := context.Background()
	userService, conn, err := getUserService(ctx, "localhost:5003")
	if err != nil {
		log.Fatal("can not get user service")
	}
	defer conn.Close()
	handler := &Handler{userClient: userService}
	mux := handler.mount()
	if err := run(mux); err != nil {
		log.Fatal("Failed to start server")
	}
}
