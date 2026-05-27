package service

import (
	"context"

	"github.com/iamhanif11/ewallet-backend/internal/dto"
	"github.com/iamhanif11/ewallet-backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionService struct {
	db                    *pgxpool.Pool
	transactionRepository *repository.TransactionRepository
}

func NewTransactionService(transactionRepository *repository.TransactionRepository, db *pgxpool.Pool) *TransactionService {
	return &TransactionService{
		db:                    db,
		transactionRepository: transactionRepository,
	}
}

func (ts *TransactionService) FindReceivers(ctx context.Context, userId int, search string, page, limit int) (dto.ReceiverListResponse, error) {
	offset := (page - 1) * limit

	receivers, err := ts.transactionRepository.FindReceivers(ctx, ts.db, userId, search, limit, offset)
	if err != nil {
		return dto.ReceiverListResponse{}, err
	}

	items := make([]dto.ReceiverResponse, 0, len(receivers))
	for _, receiver := range receivers {
		items = append(items, dto.ReceiverResponse{
			Id:       receiver.Id,
			Picture:  receiver.Picture,
			Receiver: receiver.Receiver,
			Phone:    receiver.Phone,
		})
	}

	return dto.ReceiverListResponse{
		Items: items,
		Pages: dto.PaginationResponse{
			Page:  page,
			Limit: limit,
		},
	}, nil
}

func (ts *TransactionService) Transfer(ctx context.Context, senderId int, req dto.TransferRequest) error {
	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	err = ts.transactionRepository.Transfer(ctx, tx, senderId, req.ReceiverId, req.Amount, req.Notes)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
