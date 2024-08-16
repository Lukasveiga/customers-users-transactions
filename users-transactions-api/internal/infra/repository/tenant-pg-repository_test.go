package infra

import (
	"testing"

	"github.com/Lukasveiga/customers-users-transaction/config"
	"github.com/stretchr/testify/assert"
)

func TestPostgreTenantRepository(t *testing.T) {
	dbClient := config.SetupPgTestcontainers()

	pgTenantRepository := NewPgTenantRepository(dbClient)

	defer dbClient.Close()

	t.Run("should successfully connect to pg container", func(t *testing.T) {
		assert.NotNil(t, &dbClient)
	})

	t.Run("[FindById] should return tenant by id", func(t *testing.T) {
		tenantId := int32(1)

		tenant, err := pgTenantRepository.FindById(tenantId)

		assert.NoError(t, err)

		assert.Equal(t, tenantId, tenant.Id)
		assert.Equal(t, "Tenant A", tenant.Name)
	})

	t.Run("[FindById] should return nil when tenant is not found by id", func(t *testing.T) {
		tenantId := int32(15)

		tenant, err := pgTenantRepository.FindById(tenantId)

		assert.NoError(t, err)

		assert.Nil(t, tenant)
	})

}
