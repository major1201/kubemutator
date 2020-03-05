package session

import "github.com/gin-gonic/gin"

const (
	// RequestIDSessionKey request id session context key
	RequestIDSessionKey = "REQUEST_ID"
)

// RequestIDFromContext returns ID in this context, or 0 if not contained.
func RequestIDFromContext(c *gin.Context) RequestID {
	if id, ok := c.Value(RequestIDSessionKey).(RequestID); ok {
		return id
	}
	return 0
}
