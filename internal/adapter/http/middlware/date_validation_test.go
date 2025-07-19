package middlware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/LiquidCats/rater/internal/adapter/http/middlware"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDateValidationMiddleware_Handle(t *testing.T) {
	tests := []struct {
		name           string
		dateParam      string
		expectedStatus int
		expectAborted  bool
		setupContext   func(*gin.Context)
	}{
		{
			name:           "empty date parameter should continue",
			dateParam:      "",
			expectedStatus: 200, // Handler should be called
			expectAborted:  false,
		},
		{
			name:           "invalid date format should abort with validation error",
			dateParam:      "invalid-date",
			expectedStatus: http.StatusUnprocessableEntity,
			expectAborted:  true,
		},
		{
			name:           "valid date in past should continue",
			dateParam:      "2020-01-01T12:00:00",
			expectedStatus: 200,
			expectAborted:  false,
		},
		{
			name:           "valid date in future should abort with validation error",
			dateParam:      "2030-01-01T12:00:00",
			expectedStatus: http.StatusUnprocessableEntity,
			expectAborted:  true,
		},
		{
			name:           "date exactly now should continue",
			dateParam:      time.Now().Format(entity.DefaultFormat),
			expectedStatus: 200,
			expectAborted:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			gin.SetMode(gin.TestMode)
			router := gin.New()

			dateValidationMiddleware := middlware.NewDateValidation()
			// Add middleware and a test handler
			handlerCalled := false
			router.GET("/test/*date", dateValidationMiddleware.Handle, func(c *gin.Context) {
				handlerCalled = true
				c.Status(http.StatusOK)
			})

			// Create request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test/"+tt.dateParam, nil)

			// Execute
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code, "HTTP status mismatch")

			if !tt.expectAborted {
				assert.True(t, handlerCalled, "Handler should be called when middleware continues")
				return
			}

			assert.False(t, handlerCalled, "Handler should not be called when middleware aborts")

			// Verify error response structure for aborted requests
			assert.Equal(t, http.StatusUnprocessableEntity, w.Code, "HTTP status mismatch")

			var response gin.H
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Failed to unmarshal error response")

			// Check response structure
			assert.Equal(t, "fail", response["status"])
			assert.Contains(t, response, "data")

			// Check data contains date field with error messages
			data, ok := response["data"].(map[string]interface{})
			require.True(t, ok, "Data should be a map")
			assert.Contains(t, data, "date")

			// Verify date field contains error messages
			dateErrors, ok := data["date"].([]interface{})
			require.True(t, ok, "Date errors should be an array")
			assert.NotEmpty(t, dateErrors, "Should have at least one error message")
		})
	}
}

func TestDateValidationMiddleware_Handle_FutureDateTime(t *testing.T) {
	// Test with a specific future date to ensure the future validation works
	gin.SetMode(gin.TestMode)
	router := gin.New()

	handlerCalled := false
	router.Use(middlware.NewDateValidation().Handle)
	router.GET("/test/:date", func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	// Create a date that's definitely in the future
	futureDate := time.Now().Add(24 * time.Hour).Format(entity.DefaultFormat)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test/"+futureDate, nil)
	router.ServeHTTP(w, req)

	assert.False(t, handlerCalled, "Handler should not be called for future dates")
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Check response structure
	assert.Equal(t, "fail", response["status"])
	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)

	// Check date field contains "must be in the past" message
	dateErrors, ok := data["date"].([]interface{})
	require.True(t, ok)
	assert.NotEmpty(t, dateErrors)

	// Check that one of the error messages contains "must be in the past"
	foundMessage := false
	for _, errMsg := range dateErrors {
		if str, ok := errMsg.(string); ok && strings.Contains(str, "must be in the past") {
			foundMessage = true
			break
		}
	}
	assert.True(t, foundMessage, "Should contain 'must be in the past' error message")
}

func TestDateValidationMiddleware_Handle_ParseError(t *testing.T) {
	// Test that parse errors are properly handled and include error details
	gin.SetMode(gin.TestMode)
	router := gin.New()

	handlerCalled := false
	router.Use(middlware.NewDateValidation().Handle)
	router.GET("/test/:date", func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	invalidDate := "not-a-date"
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test/"+invalidDate, nil)
	router.ServeHTTP(w, req)

	assert.False(t, handlerCalled, "Handler should not be called for invalid dates")
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Check response structure
	assert.Equal(t, "fail", response["status"])
	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)

	// Check date field contains format error message
	dateErrors, ok := data["date"].([]interface{})
	require.True(t, ok)
	assert.NotEmpty(t, dateErrors)

	// Check that one of the error messages contains format information
	foundFormatMessage := false
	for _, errMsg := range dateErrors {
		if str, ok := errMsg.(string); ok && strings.Contains(str, "datetime must be in format:") {
			foundFormatMessage = true
			break
		}
	}
	assert.True(t, foundFormatMessage, "Should contain format error message")
}

func TestDateValidationMiddleware_ResponseStructure(t *testing.T) {
	// Test to verify the exact response structure matches the DTO pattern
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(middlware.NewDateValidation().Handle)
	router.GET("/test/:date", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Test with invalid date to trigger validation error
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test/invalid-date", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	// Parse response and verify exact structure
	var response gin.H
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify top-level structure
	assert.Equal(t, "fail", response["status"])
	assert.Contains(t, response, "data")

	// Verify data structure
	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok, "Data should be a map")
	assert.Contains(t, data, "date")

	// Verify date field contains array of strings
	dateErrors, ok := data["date"].([]interface{})
	require.True(t, ok, "Date field should be an array")
	require.NotEmpty(t, dateErrors, "Should have error messages")

	// Verify all error messages are strings
	for i, errMsg := range dateErrors {
		_, ok := errMsg.(string)
		assert.True(t, ok, "Error message at index %d should be a string", i)
	}
}
