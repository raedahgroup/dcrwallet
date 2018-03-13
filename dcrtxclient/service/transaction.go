package service

import (
	"bytes"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/wire"
	pb "github.com/raedahgroup/dcrtxmatcher/pkg/api/matcherrpc"
	"google.golang.org/grpc"

	"golang.org/x/net/context"
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

func (t *TransactionService) JoinTransaction(tx *wire.MsgTx, voteAddress dcrutil.Address, ticketPrice dcrutil.Amount) (*wire.MsgTx, error) {
	canContribute := ticketPrice.MulF64(0.5) // Half of ticket price

	joinReq := &pb.FindMatchesRequest{
		Amount: uint64(canContribute),
	}

	findRes, err := t.client.FindMatches(context.Background(), joinReq)
	if err != nil {
		return nil, err
	}

	generateReq := &pb.GenerateTicketRequest{
		SessionId: findRes.SessionId,
		CommitmentOutput: &pb.TxOut{
			Value:  uint64(tx.TxOut[1].Value),
			Script: tx.TxOut[1].PkScript,
		},
		ChangeOutput: &pb.TxOut{
			Value:  uint64(tx.TxOut[2].Value),
			Script: tx.TxOut[2].PkScript,
		},
		VoteAddress: voteAddress.String(),
	}

	genRes, err := t.client.GenerateTicket(context.Background(), generateReq)
	if err != nil {
		return nil, err
	}

	buffTx := bytes.NewBuffer(nil)
	buffTx.Grow(tx.SerializeSize())
	err = tx.BtcEncode(buffTx, 0)
	if err != nil {
		return nil, err
	}

	publishReq := &pb.PublishTicketRequest{
		SessionId:            findRes.SessionId,
		SplitTx:              buffTx.Bytes(),
		SplitTxOutputIndex:   genRes.OutputIndex,
		TicketInputScriptsig: tx.TxIn[0].SignatureScript,
	}

	publishRes, err := t.client.PublishTicket(context.Background(), publishReq)
	if err != nil {
		return nil, err
	}

	var ticket *wire.MsgTx
	rbuf := bytes.NewReader(publishRes.TicketTx)
	err = ticket.BtcDecode(rbuf, 0)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}
