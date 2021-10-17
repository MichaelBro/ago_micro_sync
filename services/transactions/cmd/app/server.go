package app

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	transactionsV1Pb "transactions/pkg/proto/v1"
	"transactions/pkg/transactions"
)

type Server struct {
	transactionsV1Pb.UnimplementedTransactionsServiceServer
	transactionsSvc *transactions.Service
}

func New(transactionsSvc *transactions.Service) *Server {
	return &Server{transactionsSvc: transactionsSvc}
}

func (s *Server) Transactions(ctx context.Context, request *transactionsV1Pb.TransactionsRequest) (*transactionsV1Pb.TransactionsResponse, error) {
	records, err := s.transactionsSvc.Transactions(ctx, request.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	data := &transactionsV1Pb.TransactionsResponse{
		Transactions: make([]*transactionsV1Pb.TransactionResponse, len(records)),
	}

	for i, record := range records {
		data.Transactions[i] = &transactionsV1Pb.TransactionResponse{
			Id:       record.ID,
			UserId:   record.UserID,
			Category: record.Category,
			Amount:   record.Amount,
			Created:  &timestamppb.Timestamp{Seconds: record.Created},
		}
	}

	return data, nil
}
