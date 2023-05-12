package service

import (
	"errors"
	"strconv"
	"tesla01/bisa_patungan/internal/model"
	repository2 "tesla01/bisa_patungan/internal/repository"
)

type TransactionService interface {
	GetTransactionByCampaignID(input model.GetCampaignTransactionsInput) ([]model.Transaction, error)
	GetTransactionByUserID(userID int) ([]model.Transaction, error)
	CreateTransaction(input model.CreateTransactionInput) (model.Transaction, error)
	ProcessPayment(input model.TransactionNotificationInput) error
}

type TransactionServiceImpl struct {
	repository         repository2.TransactionRepository
	campaignRepository repository2.CampaignRepository
	paymentService     PaymentService
}

func NewTransactionService(repository repository2.TransactionRepository, campaignRepository repository2.CampaignRepository, paymentService PaymentService) *TransactionServiceImpl {
	return &TransactionServiceImpl{repository, campaignRepository, paymentService}
}

func (s *TransactionServiceImpl) GetTransactionByCampaignID(input model.GetCampaignTransactionsInput) ([]model.Transaction, error) {

	currentCampaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []model.Transaction{}, err
	}

	if currentCampaign.UserID != input.User.ID {
		return []model.Transaction{}, errors.New("not an owner of the campaign")
	}

	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *TransactionServiceImpl) GetTransactionByUserID(userID int) ([]model.Transaction, error) {
	transaction, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (s *TransactionServiceImpl) CreateTransaction(input model.CreateTransactionInput) (model.Transaction, error) {
	transaction := model.Transaction{
		CampaignID: input.CampaignID,
		UserID:     input.User.ID,
		Amount:     input.Amount,
		Status:     "pending",
	}

	newTransaction, err := s.repository.Save(transaction)

	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := model.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)

	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL

	newTransaction, err = s.repository.Update(newTransaction)

	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

func (s *TransactionServiceImpl) ProcessPayment(input model.TransactionNotificationInput) error {
	transaction_id, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.repository.GetByID(transaction_id)

	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.repository.Update(transaction)

	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(transaction.CampaignID)

	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil
}
