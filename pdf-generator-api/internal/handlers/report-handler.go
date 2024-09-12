package handlers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/Lukasveiga/customers-users-transactions/internal/usecases"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	transactionReport *usecases.TransactionReport
}

func NewReportHandler(transactionReport *usecases.TransactionReport) *ReportHandler {
	return &ReportHandler{transactionReport: transactionReport}
}

func (h *ReportHandler) SendReport(c *gin.Context) {
	tenantId, err := strconv.ParseInt(c.Param("tenantId"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant-id"})
		return
	}

	accountId, err := strconv.ParseInt(c.Param("accountId"), 0, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account id"})
		return
	}

	path, err := h.transactionReport.GeneratePdfReport(usecases.GenerateInputParams{
		TenantId:  int32(tenantId),
		AccountId: int32(accountId),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	_, err = os.Stat(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "File not found",
		})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename=transactions-report.pdf")
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", "0")

	// Send the file as a response
	c.File(path)
}
