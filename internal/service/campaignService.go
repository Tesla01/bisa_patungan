package service

import (
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"tesla01/bisa_patungan/internal/model"
	repository2 "tesla01/bisa_patungan/internal/repository"
)

type CampaignService interface {
	GetCampaigns(userID int) ([]model.Campaign, error)
	GetCampaignByID(input model.GetCampaignDetailInput) (model.Campaign, error)
	CreateCampaign(input model.CreateCampaignInput) (model.Campaign, error)
	UpdateCampaign(inputID model.GetCampaignDetailInput, inputData model.CreateCampaignInput) (model.Campaign, error)
	SaveCampaignImage(input model.CreateCampaignImageInput, fileLocation string) (model.CampaignImage, error)
}

type CampaignServiceImpl struct {
	repository repository2.CampaignRepository
}

func NewCampaignService(repository repository2.CampaignRepository) *CampaignServiceImpl {
	return &CampaignServiceImpl{repository}
}

func (s *CampaignServiceImpl) GetCampaigns(userID int) ([]model.Campaign, error) {
	if userID != 0 {
		campaigns, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}

	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (s *CampaignServiceImpl) GetCampaignByID(input model.GetCampaignDetailInput) (model.Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (s *CampaignServiceImpl) CreateCampaign(input model.CreateCampaignInput) (model.Campaign, error) {
	campaign := model.Campaign{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		Perks:            input.Perks,
		GoalAmount:       input.GoalAmount,
		UserID:           input.User.ID,
	}

	campaignSlug := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(campaignSlug)

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *CampaignServiceImpl) UpdateCampaign(inputID model.GetCampaignDetailInput, inputData model.CreateCampaignInput) (model.Campaign, error) {
	campaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserID != inputData.User.ID {
		fmt.Printf("Campaign User %d, Input User %d \n", campaign.UserID, inputData.User.ID)
		return campaign, errors.New("not an owner of the campaign")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return campaign, err
	}

	return updatedCampaign, nil
}

func (s *CampaignServiceImpl) SaveCampaignImage(input model.CreateCampaignImageInput, fileLocation string) (model.CampaignImage, error) {
	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return model.CampaignImage{}, err
	}

	if campaign.UserID != input.User.ID {
		return model.CampaignImage{}, errors.New("not an owner of the campaign")
	}

	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return model.CampaignImage{}, err
		}
	}

	campaignImage := model.CampaignImage{
		CampaignID: input.CampaignID,
		IsPrimary:  isPrimary,
		FileName:   fileLocation,
	}

	newCampaignImage, err := s.repository.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}
