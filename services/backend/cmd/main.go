package main

import (
	"backend/cmd/app"
	"backend/pkg/auth"
	transactionsV1Pb "backend/pkg/proto/v1"
	"backend/pkg/transactions"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort             = "9999"
	defaultHost             = "0.0.0.0"
	defaultTransactionsAddr = "localhost:11111"
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

	tokenLifeTime := time.Hour

	publicKey, err := ioutil.ReadFile("keys/public.key")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	transactionsAddr, ok := os.LookupEnv("APP_TRANSACTIONS_ADDR")
	if !ok {
		transactionsAddr = defaultTransactionsAddr
	}

	if err := execute(net.JoinHostPort(host, port), publicKey, tokenLifeTime, transactionsAddr); err != nil {
		os.Exit(1)
	}
}

func execute(addr string, publicKey []byte, tokenLifeTime time.Duration, transactionsAddr string) error {
	authSvc := auth.NewService(publicKey, tokenLifeTime)

	conn, err := grpc.Dial(transactionsAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()
	client := transactionsV1Pb.NewTransactionsServiceClient(conn)
	transactionsSvc := transactions.NewService(client)

	mux := chi.NewRouter()

	application := app.NewServer(authSvc, transactionsSvc, mux)
	err = application.Init()
	if err != nil {
		log.Print(err)
		return err
	}

	server := &http.Server{
		Addr:    addr,
		Handler: application,
	}
	return server.ListenAndServe()
}
