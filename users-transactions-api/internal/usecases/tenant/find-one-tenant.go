package usecases

import (
	"context"
	"database/sql"
	"log/slog"

	infra "github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/infra/repository/sqlc"
	"github.com/Lukasveiga/customers-users-transaction/users-transactions-api/internal/shared"
)

type FindOneTenantUseCase struct {
	repo infra.Querier
}

func NewFindOneTenantUseCase(repo infra.Querier) *FindOneTenantUseCase {
	return &FindOneTenantUseCase{
		repo: repo,
	}
}

func (uc *FindOneTenantUseCase) FindOne(tenantId int32) (*infra.Tenant, error) {
	tenant, err := uc.repo.GetTenant(context.Background(), tenantId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &shared.EntityNotFoundError{
				Object: "tenant",
				Id:     tenantId,
			}
		}
		slog.Error(
			"error to find tenant by id",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	return &tenant, nil
}
