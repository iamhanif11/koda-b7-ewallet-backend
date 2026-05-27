package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iamhanif11/ewallet-backend/internal/dto"
	"github.com/iamhanif11/ewallet-backend/internal/service"
	"github.com/iamhanif11/ewallet-backend/pkg"
)

type TransactionController struct {
	transactionService *service.TransactionService
}

func NewTransactionController(transactionService *service.TransactionService) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

// Find Receivers
//
//	@Summary		Search and Retriever Receiver List
//	@Description	List other users for transfer purpose with search and pagination
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			search			query	string		false	"Search keyword by name and phone"
//	@Param			page			query	int			false	"Page Number"
//	@Param			limit			query	int			false	"Max data per page"
//	@Success		200	{object}	dto.Response[dto.ReceiverListResponse]
//	@Failure		500	{object}	dto.ErrorResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Router			/transaction/receivers 	[get]
func (tc *TransactionController) FindReceivers(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication Failed",
			Success: false,
		})
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid page parameter",
			Success: false,
		})
		return
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid Limit Parameter",
			Success: false,
		})
		return
	}

	search := strings.TrimSpace(ctx.DefaultQuery("search", ""))

	res, err := tc.transactionService.FindReceivers(ctx.Request.Context(), userClaims.Id, search, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to Retrieve Receiver List",
			Success: false,
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response[dto.ReceiverListResponse]{
		Message: "Get Receivers Succesfully",
		Data:    res,
		Success: true,
	})
}

// Transfer
//
//	@Summary		Transfer Balance
//	@Description	Transfer balance to another user
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		dto.TransferRequest	true	"Transfer Request"
//	@Success		200		{object}	dto.Response[any]
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/transaction/transfer [post]
func (tc *TransactionController) Transfer(ctx *gin.Context) {
	claims, ok := ctx.Get("user")
	userClaims, ok := claims.(*pkg.Claims)

	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication Failed",
			Success: false,
		})
		return
	}

	var req dto.TransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid Request Body",
			Success: false,
		})
		return
	}

	if req.ReceiverId <= 0 {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Receiver Id is Required",
			Success: false,
		})
		return
	}

	if req.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Amount Must be Greater than 0",
			Success: false,
		})
		return
	}

	if userClaims.Id == req.ReceiverId {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Cannot Transfer to Yourself",
			Success: false,
		})
		return
	}

	err := tc.transactionService.Transfer(
		ctx.Request.Context(),
		userClaims.Id, req,
	)

	if err != nil {
		if err.Error() == "insufficient balance" {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Inssuficient Balance",
				Success: false,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Transfer Failed",
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response[any]{
		Message: "Transfer Success",
		Success: true,
	})

}

// Top Up
//
//	@Summary		Top Up Balance
//	@Description	Top up user wallet balance
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		dto.TopUpRequest	true	"Top Up Request"
//	@Success		200		{object}	dto.Response[any]
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/transaction/topup [post]
func (tc *TransactionController) Topup(ctx *gin.Context) {
	claims, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication Failed",
			Success: false,
		})
		return
	}

	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Invalid",
			Success: false,
		})
		return
	}

	var req dto.TopUpRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid Request Body",
			Success: false,
		})
		return
	}

	if req.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Amount must be greater than 0",
			Success: false,
		})
		return
	}

	if req.PaymentMethodId <= 0 {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Payment method is required",
			Success: false,
		})
		return
	}

	err := tc.transactionService.TopUp(ctx.Request.Context(), userClaims.Id, req)

	log.Println("error", err)
	if err != nil {

		if err.Error() == "payment method not found" {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Payment method not found",
				Success: false,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Top up failed",
			Success: false,
			Error:   err.Error(),
		})
		return

	}
	ctx.JSON(http.StatusOK, dto.Response[any]{
		Message: "Top up Succes",
		Success: true,
		Data:    nil,
	})
}
