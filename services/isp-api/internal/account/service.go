package account

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (Account, error)
	GetByAccountNo(ctx context.Context, accountNo string) (Account, error)
	Create(ctx context.Context, account Account) error
	Update(ctx context.Context, account Account) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, id string) (Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByAccountNo(ctx context.Context, accountNo string) (Account, error) {
	return s.repo.GetByAccountNo(ctx, accountNo)
}

func (s *Service) Create(ctx context.Context, account Account) error {
	if err := Validate(account); err != nil {
		return err
	}
	return s.repo.Create(ctx, account)
}

func (s *Service) Update(ctx context.Context, account Account) error {
	if err := Validate(account); err != nil {
		return err
	}
	return s.repo.Update(ctx, account)
}
