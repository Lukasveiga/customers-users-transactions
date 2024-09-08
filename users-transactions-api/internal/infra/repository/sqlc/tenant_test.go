package infra

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTenantRepository(t *testing.T) {

	t.Run("[GetTenant] should find tenant by id", func(t *testing.T) {
		tenant, err := testQueries.GetTenant(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotEmpty(t, tenant)
		assert.Equal(t, "Tenant A", tenant.Name)
	})

	t.Run("[GetTenant] should return error when tenant not found", func(t *testing.T) {
		tenant, err := testQueries.GetTenant(context.Background(), 99)

		assert.Empty(t, tenant)
		assert.EqualError(t, err, sql.ErrNoRows.Error())
	})
}
