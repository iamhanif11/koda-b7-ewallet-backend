package service

import (
	"context"
	"fmt"

	"github.com/iamhanif11/ewallet-backend/internal/dto"
	"github.com/iamhanif11/ewallet-backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type TransactionService struct {
	db                    *pgxpool.Pool
	transactionRepository *repository.TransactionRepository
	rdb                   *redis.Client
}

func NewTransactionService(transactionRepository *repository.TransactionRepository, db *pgxpool.Pool, rdb *redis.Client) *TransactionService {
	return &TransactionService{
		db:                    db,
		transactionRepository: transactionRepository,
		rdb:                   rdb,
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

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	go func() {
		bgCtx := context.Background()
		senderCacheKey := fmt.Sprintf("user:%d:dashboard", senderId)
		receiverCacheKey := fmt.Sprintf("user:%d:dashboard", req.ReceiverId)

		ts.rdb.Del(bgCtx, senderCacheKey, receiverCacheKey)
	}()

	return nil
}

func (ts *TransactionService) TopUp(ctx context.Context, userId int, req dto.TopUpRequest) error {

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	err = ts.transactionRepository.TopUp(ctx, tx, userId, req.Amount, req.PaymentMethodId)

	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	go func() {
		dashboardCacheKey := fmt.Sprintf("user:%d:dashboard", userId)
		ts.rdb.Del(context.Background(), dashboardCacheKey)
	}()

	return nil
}

func (ts *TransactionService) GetTransactionHistory(ctx context.Context, userId int, search string, page, limit int) (dto.TransactionHistoryResponse, error) {
	offset := (page - 1) * limit

	histories, err := ts.transactionRepository.GetTransactionHistoryById(ctx, ts.db, userId, search, limit, offset)

	if err != nil {
		return dto.TransactionHistoryResponse{}, err
	}

	items := make([]dto.TransactionHistoryItem, 0, len(histories))
	for _, history := range histories {

		items = append(items, dto.TransactionHistoryItem{
			Id:        history.Id,
			Type:      history.Type,
			Amount:    history.Amount,
			Status:    history.Status,
			CreatedAt: history.CreatedAt,
			Fullname:  *history.Fullname,
			Picture:   *history.Picture,
			Phone:     *history.Phone,
		})
	}

	return dto.TransactionHistoryResponse{
		Items: items,
		Pages: dto.PaginationResponse{
			Page:  page,
			Limit: limit,
		},
	}, nil
}
