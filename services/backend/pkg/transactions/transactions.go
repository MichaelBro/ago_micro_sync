package transactions

import (
	transactionsV1Pb "backend/pkg/proto/v1"
	"context"
	"encoding/json"
	"log"
)

type Service struct {
	client transactionsV1Pb.TransactionsServiceClient
	url    string
}

func NewService(client transactionsV1Pb.TransactionsServiceClient) *Service {
	return &Service{client: client}
}

func (s *Service) Transactions(ctx context.Context, userID int64) ([]byte, error) {
	resp, err := s.client.Transactions(ctx, &transactionsV1Pb.TransactionsRequest{Id: userID})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return data, nil
}
