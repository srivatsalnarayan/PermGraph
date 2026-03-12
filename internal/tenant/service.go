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

	// start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer tx.Rollback(ctx)

	var tenantID int64

	// 1️⃣ insert tenant
	err = tx.QueryRow(ctx,
		`INSERT INTO zanzibar_core.tenants (tenant_name)
		 VALUES ($1)
		 RETURNING tenant_id`,
		name,
	).Scan(&tenantID)

	if err != nil {
		return 0, err
	}

	// 2️⃣ initialize revision counter
	_, err = tx.Exec(ctx,
		`INSERT INTO zanzibar_core.revision_counter (tenant_id, current_revision)
		 VALUES ($1, 1)`,
		tenantID,
	)
	if err != nil {
		return 0, err
	}

	// 3️⃣ insert revision history
	_, err = tx.Exec(ctx,
		`INSERT INTO zanzibar_core.tenant_revision
		 VALUES ($1, 1, now(), NULL)`,
		tenantID,
	)
	if err != nil {
		return 0, err
	}

	// build schema name
	schemaName := fmt.Sprintf("tenant_%d", tenantID)

	// 4️⃣ create tenant schema
	_, err = tx.Exec(ctx,
		fmt.Sprintf(`CREATE SCHEMA %s`, schemaName),
	)
	if err != nil {
		return 0, err
	}

	// 5️⃣ create auth_tuple table
	_, err = tx.Exec(ctx,
		fmt.Sprintf(`
		CREATE TABLE %s.auth_tuple (
			tuple_id BIGSERIAL PRIMARY KEY,
			object_type TEXT,
			object_id TEXT,
			relation TEXT,
			subject_type TEXT,
			subject_id TEXT,
			valid_from_rev BIGINT,
			valid_to_rev BIGINT
		)`, schemaName),
	)
	if err != nil {
		return 0, err
	}

	// 6️⃣ create authorization_model table
	_, err = tx.Exec(ctx,
		fmt.Sprintf(`
		CREATE TABLE %s.authorization_model (
			model_id BIGSERIAL PRIMARY KEY,
			model_json JSONB,
			valid_from_rev BIGINT,
			valid_to_rev BIGINT
		)`, schemaName),
	)
	if err != nil {
		return 0, err
	}

	// commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	fmt.Println("Tenant fully initialized:", tenantID)

	return tenantID, nil
}