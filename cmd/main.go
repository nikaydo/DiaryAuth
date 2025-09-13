package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nikaydo/DiaryAuth/internal/config"
	"github.com/nikaydo/DiaryAuth/internal/database"
	authService "github.com/nikaydo/DiaryAuth/internal/grpc"
	"github.com/nikaydo/DiaryContract/gen/auth"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Loading env...")
	env, err := config.ReadEnv()
	if err != nil {
		log.Fatalln("Config error loading: ", err)
		return
	}
	log.Println("Env load sucessful")
	log.Println("Database Init...")
	db, err := database.InitBD(env)
	if err != nil {
		log.Fatalln("Database error loading: ", err)
		return
	}
	log.Println("Database Init sucessful")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", env.Host, env.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Grpc server running...")
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServer(grpcServer, &authService.Auth{DB: db, Env: env})
	log.Println("Grpc server running sucessful")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	<-stop
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped")
}
