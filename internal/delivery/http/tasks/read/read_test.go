package read_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/read"
	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/read/mocks"
	"github.com/10Narratives/task-tracker/internal/models"
	"github.com/10Narratives/task-tracker/pkg/logging/slogdiscard"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReadHandler(t *testing.T) {
	tests := []struct {
		name           string
		search         string
		mockSetup      func(m *mocks.TaskReader)
		expectedStatus int
		expectedResp   read.Response
	}{
		{
			name:   "No search parameter",
			search: "",
			mockSetup: func(m *mocks.TaskReader) {
				m.On("Tasks", mock.Anything, "").Return([]models.Task{
					{ID: 1, Title: "Task 1", Date: "20250207", Comment: "The 1 task"},
					{ID: 2, Title: "Task 2", Date: "20250208", Comment: "The 2 task"},
					{ID: 3, Title: "Task 3", Date: "20250209", Comment: "The 3 task"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResp: read.Response{Tasks: []models.Task{
				{ID: 1, Title: "Task 1", Date: "20250207", Comment: "The 1 task"},
				{ID: 2, Title: "Task 2", Date: "20250208", Comment: "The 2 task"},
				{ID: 3, Title: "Task 3", Date: "20250209", Comment: "The 3 task"},
			}},
		},
		{
			name:   "Search by date",
			search: "20250205",
			mockSetup: func(m *mocks.TaskReader) {
				m.On("Tasks", mock.Anything, "").Return([]models.Task{
					{ID: 1, Title: "Task 1", Date: "20250205", Comment: "The 1 task"},
					{ID: 2, Title: "Task 2", Date: "20250205", Comment: "The 2 task"},
					{ID: 3, Title: "Task 3", Date: "20250205", Comment: "The 3 task"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResp: read.Response{Tasks: []models.Task{
				{ID: 1, Title: "Task 1", Date: "20250205", Comment: "The 1 task"},
				{ID: 2, Title: "Task 2", Date: "20250205", Comment: "The 2 task"},
				{ID: 3, Title: "Task 3", Date: "20250205", Comment: "The 3 task"},
			}},
		},
		{
			name:   "Search by payload",
			search: "ask",
			mockSetup: func(m *mocks.TaskReader) {
				m.On("Tasks", mock.Anything, "").Return([]models.Task{
					{ID: 1, Title: "Task 1", Date: "20250205", Comment: "The 1 task"},
					{ID: 2, Title: "Task 2", Date: "20250205", Comment: "The 2 task"},
					{ID: 3, Title: "Task 3", Date: "20250205", Comment: "The 3 task"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResp: read.Response{Tasks: []models.Task{
				{ID: 1, Title: "Task 1", Date: "20250205", Comment: "The 1 task"},
				{ID: 2, Title: "Task 2", Date: "20250205", Comment: "The 2 task"},
				{ID: 3, Title: "Task 3", Date: "20250205", Comment: "The 3 task"},
			}},
		},
		{
			name:   "Database error",
			search: "",
			mockSetup: func(m *mocks.TaskReader) {
				m.On("Tasks", mock.Anything, "").Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   read.Response{Err: "failed to read tasks"},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mock := mocks.NewTaskReader(t)
			tc.mockSetup(mock)

			handler := read.New(slogdiscard.NewDiscardLogger(), mock)

			url := `/api/tasks`
			if tc.search != "" {
				url += `?search=` + tc.search
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Get(`/api/tasks`, handler)
			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			var actualResp read.Response
			_ = json.Unmarshal(rec.Body.Bytes(), &actualResp)

			assert.Equal(t, tc.expectedResp, actualResp)
			mock.AssertExpectations(t)
		})
	}
}
