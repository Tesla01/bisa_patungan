package service

import (
	"github.com/veritrans/go-midtrans"
	"strconv"
	"tesla01/bisa_patungan/internal/model"
)

type PaymentService interface {
	GetPaymentURL(transaction model.Transaction, user model.User) (string, error)
}

type PaymentServiceImpl struct {
}

func NewPaymentService() *PaymentServiceImpl {
	return &PaymentServiceImpl{}
}

func (s *PaymentServiceImpl) GetPaymentURL(transaction model.Transaction, user model.User) (string, error) {
	client := midtrans.NewClient()
	//Don't commit server & client key
	client.ServerKey = ""
	client.ClientKey = ""
	client.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: client,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)

	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}
