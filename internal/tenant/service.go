package tenant

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) CreateTenant(ctx context.Context, name string) (int64, error) {

	var tenantID int64

	err := s.db.QueryRow(ctx,
		`INSERT INTO zanzibar_core.tenants (tenant_name)
		 VALUES ($1)
		 RETURNING tenant_id`,
		name,
	).Scan(&tenantID)

	if err != nil {
		return 0, err
	}

	fmt.Println("Tenant created:", tenantID)

	return tenantID, nil
}