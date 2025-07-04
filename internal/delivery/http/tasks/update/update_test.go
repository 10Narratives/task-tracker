package update_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/update"
	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/update/mocks"
	"github.com/10Narratives/task-tracker/internal/lib/logging/handlers/slogdiscard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(m *mocks.TaskUpdater)
		expectedStatus int
		expectedResp   update.Response
	}{
		{
			name:        "successful update",
			requestBody: `{"id": "100", "date":"20250205","title":"Test Task","comment":"This is a test","repeat":"d 7"}`,
			mockSetup: func(m *mocks.TaskUpdater) {
				m.On("Update", mock.Anything, int64(100), "20250205", "Test Task", "This is a test", "d 7").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedResp:   update.Response{},
		},
		{
			name:        "unsuccessful update - fail to decode body",
			requestBody: `{"id": "100", "date":"20250205","title Test Task","comment":"This is a test","repeat":"d 7"}`,
			mockSetup: func(m *mocks.TaskUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   update.Response{Err: "failed to decode request body"},
		},
		{
			name:        "unsuccessful update - invalid body date",
			requestBody: `{"id": "100", "date":"2025-02-05","title":"Test Task","comment":"This is a test","repeat":"d 7"}`,
			mockSetup: func(m *mocks.TaskUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   update.Response{Err: "field Date must be in YYYYMMDD date format"},
		},
		{
			name:        "unsuccessful update - invalid body title",
			requestBody: `{"id": "100", "date":"20250205","title":"","comment":"This is a test","repeat":"d 7"}`,
			mockSetup: func(m *mocks.TaskUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   update.Response{Err: "field Title is required"},
		},
		{
			name:        "unsuccessful update - invalid body repeat",
			requestBody: `{"id": "100", "date":"20250205","title":"title","comment":"This is a test","repeat":"d 7,5,6,131231"}`,
			mockSetup: func(m *mocks.TaskUpdater) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp:   update.Response{Err: "field Repeat must satisfy expected patterns"},
		},
		{
			name:        "unsuccessful update - database error",
			requestBody: `{"id": "100", "date":"20250205","title":"Test Task","comment":"This is a test","repeat":"d 7"}`,
			mockSetup: func(m *mocks.TaskUpdater) {
				m.On("Update", mock.Anything, int64(100), "20250205", "Test Task", "This is a test", "d 7").Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   update.Response{Err: "failed to update task"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRegistrar := new(mocks.TaskUpdater)
			tc.mockSetup(mockRegistrar)

			handler := update.New(slogdiscard.NewDiscardLogger(), mockRegistrar)

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedStatus, recorder.Code)

			var actualResp update.Response
			_ = json.Unmarshal(recorder.Body.Bytes(), &actualResp)

			assert.Equal(t, tc.expectedResp, actualResp)
			mockRegistrar.AssertExpectations(t)
		})
	}
}
