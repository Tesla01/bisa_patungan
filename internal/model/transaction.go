package model

import (
	"time"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	User       User
	Campaign   Campaign
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
