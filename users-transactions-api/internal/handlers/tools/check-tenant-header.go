package tools

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CheckTenantHeader(c *gin.Context) (int32, bool) {
	tenantId, err := strconv.ParseInt(c.GetHeader("tenant-id"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant-id"})
		return -1, false
	}

	return int32(tenantId), true
}
