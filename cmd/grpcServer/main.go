package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vinicius-gregorio/go_grpc_bff/internal/database"
	"github.com/vinicius-gregorio/go_grpc_bff/internal/pb"
	"github.com/vinicius-gregorio/go_grpc_bff/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("Err to open db: %v", err)
	}
	defer db.Close()

	categoryDb := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryDb)
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
