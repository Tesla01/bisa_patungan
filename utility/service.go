package utility

type Service interface {
	CheckHealth() (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CheckHealth() (bool, error) {
	return true, nil
}
