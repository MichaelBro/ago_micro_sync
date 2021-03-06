package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"transactions/cmd/app"
	transactionsV1Pb "transactions/pkg/proto/v1"
	"transactions/pkg/transactions"
)

const (
	defaultPort = "9999"
	defaultHost = "0.0.0.0"
	defaultDSN  = "postgres://app:pass@localhost:5532/db"
)

func main() {
	port, ok := os.LookupEnv("APP_PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("APP_HOST")
	if !ok {
		host = defaultHost
	}

	dsn, ok := os.LookupEnv("APP_DSN")
	if !ok {
		dsn = defaultDSN
	}

	if err := execute(net.JoinHostPort(host, port), dsn); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, dsn string) error {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Print(err)
		return err
	}
	defer pool.Close()

	transactionsSvc := transactions.NewService(pool)
	application := app.New(transactionsSvc)

	grpcServer := grpc.NewServer()
	transactionsV1Pb.RegisterTransactionsServiceServer(grpcServer, application)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return grpcServer.Serve(listener)
}
