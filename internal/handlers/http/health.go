package httphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func(h *Handler) PingStatus(c *gin.Context){
	c.String(http.StatusOK, "ok")
}
