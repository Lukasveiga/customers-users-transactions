package usecases

import (
	"log/slog"

	"github.com/Lukasveiga/customers-users-transaction/internal/domain"
	port "github.com/Lukasveiga/customers-users-transaction/internal/ports/repository"
	"github.com/Lukasveiga/customers-users-transaction/internal/shared"
)

type FindOneTenantUseCase struct {
	repo port.TenantRepository
}

func NewFindOneTenantUseCase(repo port.TenantRepository) *FindOneTenantUseCase {
	return &FindOneTenantUseCase{
		repo: repo,
	}
}

func (uc *FindOneTenantUseCase) FindOne(tenantId int32) (*domain.Tenant, error) {
	tenant, err := uc.repo.FindById(tenantId)

	if err != nil {
		slog.Error(
			"error to find tenant by id",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	if tenant == nil {
		return nil, &shared.EntityNotFoundError{
			Object: "tenant",
			Id:     tenantId,
		}
	}

	return tenant, nil
}
