package repository

import (
	"gorm.io/gorm"
	"tesla01/bisa_patungan/internal/model"
)

type CampaignRepository interface {
	FindAll() ([]model.Campaign, error)
	FindByUserID(userID int) ([]model.Campaign, error)
	FindByID(id int) (model.Campaign, error)
	Save(campaign model.Campaign) (model.Campaign, error)
	Update(campaign model.Campaign) (model.Campaign, error)
	CreateImage(campaignImage model.CampaignImage) (model.CampaignImage, error)
	MarkAllImagesAsNonPrimary(campaignID int) (bool, error)
}

type CampaignRepositoryImpl struct {
	db *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) *CampaignRepositoryImpl {
	return &CampaignRepositoryImpl{db: db}
}

func (r *CampaignRepositoryImpl) FindAll() ([]model.Campaign, error) {
	var campaigns []model.Campaign

	err := r.db.Find(&campaigns).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *CampaignRepositoryImpl) FindByUserID(userID int) ([]model.Campaign, error) {
	var campaigns []model.Campaign

	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *CampaignRepositoryImpl) FindByID(id int) (model.Campaign, error) {
	var campaign model.Campaign

	err := r.db.Preload("User").Preload("CampaignImages").Where("id = ?", id).Find(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil

}

func (r *CampaignRepositoryImpl) Save(campaign model.Campaign) (model.Campaign, error) {
	err := r.db.Create(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *CampaignRepositoryImpl) Update(campaign model.Campaign) (model.Campaign, error) {
	err := r.db.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *CampaignRepositoryImpl) CreateImage(campaignImage model.CampaignImage) (model.CampaignImage, error) {
	err := r.db.Create(&campaignImage).Error

	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}

func (r *CampaignRepositoryImpl) MarkAllImagesAsNonPrimary(campaignID int) (bool, error) {
	err := r.db.Model(&model.CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
