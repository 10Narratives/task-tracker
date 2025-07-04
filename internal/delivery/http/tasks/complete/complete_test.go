package complete_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/complete"
	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/complete/mocks"
	"github.com/10Narratives/task-tracker/internal/lib/logging/handlers/slogdiscard"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCompleteHandler(t *testing.T) {
	tests := []struct {
		name       string
		mockSetup  func(m *mocks.TaskCompleter)
		id         string
		wantStatus int
		wantResp   complete.Response
	}{
		{
			name: "successful complete",
			mockSetup: func(m *mocks.TaskCompleter) {
				m.On("Complete", mock.Anything, int64(100)).Return(nil)
			},
			id:         "100",
			wantStatus: http.StatusOK,
			wantResp:   complete.Response{},
		},
		{
			name: "unsuccessful complete - invalid id",
			mockSetup: func(m *mocks.TaskCompleter) {
			},
			id:         "invalid",
			wantStatus: http.StatusBadRequest,
			wantResp:   complete.Response{Err: "gotten invalid id"},
		},
		{
			name: "unsuccessful complete - empty id",
			mockSetup: func(m *mocks.TaskCompleter) {
			},
			id:         "",
			wantStatus: http.StatusBadRequest,
			wantResp:   complete.Response{Err: "gotten invalid id"},
		},
		{
			name: "unsuccessful complete - database error",
			mockSetup: func(m *mocks.TaskCompleter) {
				m.On("Complete", mock.Anything, int64(100)).Return(errors.New("database error"))
			},
			id:         "100",
			wantStatus: http.StatusInternalServerError,
			wantResp:   complete.Response{Err: "failed to complete task"},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mock := mocks.NewTaskCompleter(t)
			tc.mockSetup(mock)

			handler := complete.New(slogdiscard.NewDiscardLogger(), mock)

			url := "/api/tasks/done"
			if tc.id != "" {
				url += "?id=" + tc.id
			}

			req := httptest.NewRequest(http.MethodPost, url, nil)
			rec := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Post(`/api/tasks/done`, handler)
			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			var actualResp complete.Response
			_ = json.Unmarshal(rec.Body.Bytes(), &actualResp)

			assert.Equal(t, tc.wantResp, actualResp)
			mock.AssertExpectations(t)
		})
	}
}
