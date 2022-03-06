package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type responseBodyInterceptor struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyInterceptor) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func newResponseBodyInterceptor(c *gin.Context) *responseBodyInterceptor {
	writerInterceptor := &responseBodyInterceptor{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
	c.Writer = writerInterceptor

	return writerInterceptor
}

func getJson(data string) (map[string]interface{}, error) {
	var unmarshalled map[string]interface{}
	err := json.Unmarshal([]byte(data), &unmarshalled)

	return unmarshalled, err
}

type Repository interface {
	Get(key string) (*CachedResponseDTO, error)
	Set(key string, status int, body string) error
}

type Idempotency struct {
	repo Repository
}

func (idempotency *Idempotency) getIdempotencyKey(c *gin.Context) string {
	idempotencyKey := c.Request.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		c.String(400, "Idempotency-Key is required in the headers")
		c.Abort()
	}

	return idempotencyKey
}

func (idempotency *Idempotency) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		idempotencyKey := idempotency.getIdempotencyKey(c)
		cachedResponse, err := idempotency.repo.Get(idempotencyKey)
		if err != nil {
			fmt.Println(err)
			c.String(500, err.Error())
			return
		}
		isCached := len(cachedResponse.Data) > 0

		if isCached {
			jsonResponse, err := getJson(cachedResponse.Data)

			if err != nil {
				fmt.Println(err)
				c.String(500, err.Error())
				return
			}

			c.AbortWithStatusJSON(cachedResponse.Status, jsonResponse)
		} else {
			writerInterceptor := newResponseBodyInterceptor(c)
			c.Next()
			statusCode := writerInterceptor.Status()
			shouldCache := statusCode >= 200 && statusCode < 300
			if shouldCache {
				body := writerInterceptor.body.String()
				err := idempotency.repo.Set(idempotencyKey, statusCode, body)
				if err != nil {
					fmt.Println(err)
					c.String(500, err.Error())
					return
				}
			}
		}
	}
}

func NewIdempotencyMiddleware(repository Repository) *Idempotency {
	return &Idempotency{repo: repository}
}
