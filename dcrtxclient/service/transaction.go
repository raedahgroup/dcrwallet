package service

import (
	pb "github.com/raedahgroup/dcrtxmatcher/pkg/matcherrpc"
)

type (
	TransactionService struct {
		client pb.SplitTicketMatcherServiceClient
	}
)

func NewTransactionService(conn *grpc.ClientConn) *TransactionService {
	return &TransactionService{
		client: pb.NewSplitTicketMatcherServiceClient(conn),
	}
}

func (t *TransactionService) Join() {

}

func (t *TransactionService) GenerateTicket() {

}
