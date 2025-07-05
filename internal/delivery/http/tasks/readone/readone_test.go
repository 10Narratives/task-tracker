package readone_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/readone"
	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/readone/mocks"
	"github.com/10Narratives/task-tracker/internal/lib/logging/handlers/slogdiscard"
	"github.com/10Narratives/task-tracker/internal/models"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReadoneHandler(t *testing.T) {
	tests := []struct {
		name       string
		mockSetup  func(m *mocks.TaskReader)
		id         string
		wantStatus int
		wantResp   readone.Response
	}{
		{
			name: "successful reading",
			mockSetup: func(m *mocks.TaskReader) {
				m.
					On("Task", mock.Anything, int64(100)).
					Return(models.Task{ID: 100, Date: "20250402", Title: "title", Comment: "comment", Repeat: "d 7"}, nil)
			},
			id:         "100",
			wantStatus: http.StatusOK,
			wantResp:   readone.Response{ID: "100", Date: "20250402", Title: "title", Comment: "comment", Repeat: "d 7"},
		},
		{
			name: "unsuccessful reading - invalid id",
			mockSetup: func(m *mocks.TaskReader) {
			},
			id:         "invalid",
			wantStatus: http.StatusBadRequest,
			wantResp:   readone.Response{Err: "gotten invalid id"},
		},
		{
			name: "unsuccessful reading - database error",
			mockSetup: func(m *mocks.TaskReader) {
				m.
					On("Task", mock.Anything, int64(100)).
					Return(models.Task{}, errors.New("database error"))
			},
			id:         "100",
			wantStatus: http.StatusInternalServerError,
			wantResp:   readone.Response{Err: "failed to find task by id"},
		},
		{
			name: "unsuccessful reading - not found",
			mockSetup: func(m *mocks.TaskReader) {
				m.
					On("Task", mock.Anything, int64(100)).
					Return(models.Task{}, nil)
			},
			id:         "100",
			wantStatus: http.StatusNotFound,
			wantResp:   readone.Response{Err: "task not found"},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mock := mocks.NewTaskReader(t)
			tc.mockSetup(mock)

			handler := readone.New(slogdiscard.NewDiscardLogger(), mock)

			url := "/api/tasks/done"
			if tc.id != "" {
				url += "?id=" + tc.id
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Get(`/api/tasks/done`, handler)
			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			var actualResp readone.Response
			_ = json.Unmarshal(rec.Body.Bytes(), &actualResp)

			assert.Equal(t, tc.wantResp, actualResp)
			mock.AssertExpectations(t)
		})
	}
}
