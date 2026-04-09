package stdres

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/example/location-demo/pkg/stderr"
	"github.com/gin-gonic/gin"
)

func JsonOk(c *gin.Context, data interface{}) {
	response(c, "application/json", http.StatusOK, data)
}

func BadRequest(c *gin.Context, data interface{}) {
	response(c, "application/json", http.StatusBadRequest, data)
}

func ServerError(c *gin.Context, data interface{}) {
	response(c, "application/json", http.StatusInternalServerError, data)
}

func UnauthorizeError(c *gin.Context, data interface{}) {
	c.Abort()
	response(c, "application/json", http.StatusUnauthorized, data)
}

func ErrRes(c *gin.Context, err error) {
	if err == nil {
		return
	}
	var e stderr.Error
	if !errors.As(err, &e) {
		response(c, "application/json", http.StatusInternalServerError, err)
		return
	}
	response(c, "application/json", e.HttpCode(), e)
}

func Image(c *gin.Context, contentType string, data []byte) {
	c.Data(http.StatusOK, contentType, data)
}

func CSV(c *gin.Context, fileName string, data [][]string) {
	// Header response
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", fileName))
	c.Header("Cache-Control", "no-cache")

	// Write UTF-8 BOM
	if _, err := c.Writer.Write([]byte("\xEF\xBB\xBF")); err != nil {
		log.Printf("unable to write UTF-8 BOM. Error: %v", err)
		response(c, "application/json", http.StatusInternalServerError, err)
		return
	}

	// Create a CSV writer
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write data to CSV
	if err := writer.WriteAll(data); err != nil {
		log.Println("unable to write data to CSV. Error:", err)
		response(c, "application/json", http.StatusInternalServerError, err)
		return
	}
}

func response(c *gin.Context, contentType string, httpCode int, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("unable to marshal. Error: %s", err)
		return
	}
	c.Data(httpCode, contentType, dataBytes)
}
