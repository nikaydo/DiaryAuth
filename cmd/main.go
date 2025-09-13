package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikaydo/DiaryAuth/internal/config"
	"github.com/nikaydo/DiaryAuth/internal/database"
	authService "github.com/nikaydo/DiaryAuth/internal/grpc"
	"github.com/nikaydo/DiaryContract/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func unaryLogger(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	start := time.Now()

	resp, err = handler(ctx, req)

	duration := time.Since(start)

	st := status.Convert(err)
	code := st.Code()

	if err != nil {
		log.Printf(" %s | %s | %+v | %s",
			info.FullMethod, code, err, duration)
	} else {
		log.Printf(" %s | %s | %s",
			info.FullMethod, code, duration)
	}

	return resp, err
}

func main() {
	log.Println("Loading env...")
	env, err := config.ReadEnv()
	if err != nil {
		log.Fatalln("Config error loading: ", err)
		return
	}
	log.Println("Env load successful")

	log.Println("Database Init...")
	db, err := database.InitBD(env)
	if err != nil {
		log.Fatalln("Database error loading: ", err)
		return
	}
	log.Println("Database Init successful")

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", env.Host, env.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryLogger))

	auth.RegisterAuthServer(grpcServer, &authService.Auth{DB: db, Env: env})

	log.Println("gRPC server running...")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	log.Println("gRPC server started")

	<-stop
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped")
}
