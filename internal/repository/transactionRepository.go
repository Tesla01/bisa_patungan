package repository

import (
	"gorm.io/gorm"
	"tesla01/bisa_patungan/internal/model"
)

type TransactionRepository interface {
	GetByCampaignID(campaignID int) ([]model.Transaction, error)
	GetByUserID(userID int) ([]model.Transaction, error)
	GetByID(transactionID int) (model.Transaction, error)
	Save(transaction model.Transaction) (model.Transaction, error)
	Update(transaction model.Transaction) (model.Transaction, error)
}

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{db}
}

func (r *TransactionRepositoryImpl) GetByCampaignID(campaignID int) ([]model.Transaction, error) {
	var transaction []model.Transaction

	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *TransactionRepositoryImpl) GetByUserID(userID int) ([]model.Transaction, error) {
	var transactions []model.Transaction

	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *TransactionRepositoryImpl) Save(transaction model.Transaction) (model.Transaction, error) {
	err := r.db.Create(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *TransactionRepositoryImpl) Update(transaction model.Transaction) (model.Transaction, error) {
	err := r.db.Save(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *TransactionRepositoryImpl) GetByID(transactionID int) (model.Transaction, error) {
	var transaction model.Transaction

	err := r.db.Where("id = ?", transactionID).Find(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
