package helper

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	"net/http"
)

var (
	serverError       = errors.New(consts.ErrServer)
	unauthorizedError = errors.New(consts.ErrUnauthorized)
	errWrongInputJson = errors.New(consts.ErrWrongInputJson)
	errNotFound       = errors.New(consts.ErrNotFound)
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data":   data,
		"errors": nil,
	})
}

func UnauthorizedErrorResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"data":   nil,
		"errors": unauthorizedError.Error(),
	})
}

func InternalServerErrorResponse(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"data":   nil,
		"errors": serverError.Error(),
	})
}

func JsonErrorResponse(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"data":   nil,
		"errors": errWrongInputJson.Error(),
	})
}

func BadRequestErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"data":   nil,
		"errors": err,
	})
}

func NotFoundErrorResponse(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"data":   nil,
		"errors": errNotFound,
	})
}

func ValidateErrorResponse(c *gin.Context, err error) {
	data := make(map[string]interface{})

	json.Unmarshal([]byte(err.Error()), &data)

	c.JSON(http.StatusBadRequest, gin.H{
		"data":   nil,
		"errors": data,
	})
}
