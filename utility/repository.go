package utility

type Repository interface {
	Check() (bool, error)
}

type repository struct {
}

func NewRepository() *repository {
	return &repository{}
}

func (r *repository) Check() (bool, error) {
	return true, nil
}
