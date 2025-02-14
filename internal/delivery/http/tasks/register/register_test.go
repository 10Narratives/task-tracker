package register_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/register"
	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/register/mocks"
	"github.com/10Narratives/task-tracker/pkg/logging/slogdiscard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(m *mocks.TaskRegistrar)
		expectedStatus int
		expectedResp   register.Response
	}{
		{
			name:        "valid request",
			requestBody: `{"date":"20250205","title":"Test Task","comment":"This is a test","repeat":"d 7"}`,
			mockSetup: func(m *mocks.TaskRegistrar) {
				m.On("Register", mock.Anything, "20250205", "Test Task", "This is a test", "d 7").
					Return(int64(1), nil)
			},
			expectedStatus: http.StatusOK,
			expectedResp:   register.Response{ID: "1"},
		},
		{
			name:        "empty request body",
			requestBody: ``,
			mockSetup: func(m *mocks.TaskRegistrar) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   register.Response{Err: "empty request"},
		},
		{
			name:        "invalid JSON format",
			requestBody: `{invalid_json}`,
			mockSetup: func(m *mocks.TaskRegistrar) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   register.Response{Err: "failed to decode request body"},
		},
		{
			name:        "validation error - missing fields",
			requestBody: `{"date":"","title":"","repeat":""}`,
			mockSetup: func(m *mocks.TaskRegistrar) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   register.Response{Err: "field Date is required, field Title is required"},
		},
		{
			name:        "validation error - wrong date format",
			requestBody: `{"date":"2025-02-05","title":"Test task","repeat":""}`,
			mockSetup: func(m *mocks.TaskRegistrar) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   register.Response{Err: "field Date must be in YYYYMMDD date format"},
		},
		{
			name:        "validation error - wrong repeat format",
			requestBody: `{"date":"20250205","title":"Test task","repeat":"d 500"}`,
			mockSetup: func(m *mocks.TaskRegistrar) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   register.Response{Err: "field Repeat must satisfy expected patterns"},
		},
		{
			name:        "task registration fails",
			requestBody: `{"date":"20250205","title":"Test Task","comment":"This is a test","repeat":"d 7"}`,
			mockSetup: func(m *mocks.TaskRegistrar) {
				m.On("Register", mock.Anything, "20250205", "Test Task", "This is a test", "d 7").
					Return(int64(0), errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   register.Response{Err: "failed to add task"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRegistrar := new(mocks.TaskRegistrar)
			tc.mockSetup(mockRegistrar)

			handler := register.New(slogdiscard.NewDiscardLogger(), mockRegistrar)

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedStatus, recorder.Code)

			var actualResp register.Response
			_ = json.Unmarshal(recorder.Body.Bytes(), &actualResp)

			assert.Equal(t, tc.expectedResp, actualResp)
			mockRegistrar.AssertExpectations(t)
		})
	}
}
