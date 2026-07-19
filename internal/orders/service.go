package orders

import (
	"context"
	repo "ecom-local/internal/adapters/postgresql/sqlc"
)

type Service interface {
	PlaceOrder(ctx context.Context, id int64) (repo.Product, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) PlaceOrder(ctx context.Context, id int64) (repo.Product, error) {
	return repo.Product{}, nil
}

