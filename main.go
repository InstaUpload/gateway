package main

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "github.com/InstaUpload/common/api"
	_ "github.com/InstaUpload/gateway/docs"
	"github.com/InstaUpload/gateway/utils"
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
		Addr:         utils.GetEnvString("HTTP_SERVER_PORT", ":5000"),
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 40,
		IdleTimeout:  time.Minute,
	}
	log.Println("Http server running in port :5000")

	return srv.ListenAndServe()
}

//	@title						InstaUpload
//	@version					0.1
//	@description				This is swagger api page for InstaUpload gateway service.
//	@contact.name				Sahaj
//	@contact.email				gpt.sahaj28@gmail.com
//	@host						localhost:5000
//	@BasePath					/v1
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
func main() {
	ctx := context.Background()
	// TODO: Get user service address from env
	userService, conn, err := getUserService(ctx, "localhost:5003")
	if err != nil {
		log.Fatal("can not get user service")
	}
	defer conn.Close()
	handler := Handler{userClient: userService}
	mux := handler.mount()
	if err := run(mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
