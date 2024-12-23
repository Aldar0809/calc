package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aldar0809/calc/api"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    api.RequestBody
		expectedStatus int
		expectedResult float64
		expectedError  string
	}{
		{
			name:           "Valid expression",
			requestBody:    api.RequestBody{Expression: "2+2*2"},
			expectedStatus: http.StatusOK,
			expectedResult: 6,
		},
		{
			name:           "Invalid expression",
			requestBody:    api.RequestBody{Expression: "2+a"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "Expression is not valid",
		},
		{
			name:           "Division by zero",
			requestBody:    api.RequestBody{Expression: "2/0"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "Expression is not valid",
		},
		{
			name:           "Empty expression",
			requestBody:    api.RequestBody{Expression: ""},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "Expression is not valid",
		},
		{
			name:           "Valid expression with parentheses",
			requestBody:    api.RequestBody{Expression: "(2+2)*2"},
			expectedStatus: http.StatusOK,
			expectedResult: 8,
		},
		{
			name:           "Invalid expression with mismatched parentheses",
			requestBody:    api.RequestBody{Expression: "(2+2*2"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedError:  "Expression is not valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBodyJSON, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(requestBodyJSON))
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			api.CalculateHandler(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, recorder.Code)
			}

			var response api.ResponseBody
			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			if err != nil {
				t.Errorf("Failed to unmarshal response body: %v", err)
			}

			if tt.expectedError != "" && response.Error != tt.expectedError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedError, response.Error)
			}

			if tt.expectedResult != 0 && response.Result != tt.expectedResult {
				t.Errorf("Expected result %f, got %f", tt.expectedResult, response.Result)
			}
		})
	}
}
