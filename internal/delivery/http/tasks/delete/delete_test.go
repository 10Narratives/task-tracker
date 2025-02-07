package delete_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/delete"
	"github.com/10Narratives/task-tracker/internal/delivery/http/tasks/delete/mocks"
	"github.com/10Narratives/task-tracker/pkg/logging/slogdiscard"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteHandler(t *testing.T) {
	tests := []struct {
		name       string
		mockSetup  func(m *mocks.TaskRemover)
		id         string
		wantStatus int
		wantResp   delete.Response
	}{
		{
			name: "successful deletion",
			mockSetup: func(m *mocks.TaskRemover) {
				m.On("Delete", mock.Anything, int64(100)).Return(nil)
			},
			id:         "100",
			wantStatus: http.StatusOK,
			wantResp:   delete.Response{},
		},
		{
			name: "unsuccessful deletion - invalid id",
			mockSetup: func(m *mocks.TaskRemover) {
			},
			id:         "invalid",
			wantStatus: http.StatusBadRequest,
			wantResp:   delete.Response{Err: "gotten invalid id"},
		},
		{
			name: "unsuccessful deletion - empty id",
			mockSetup: func(m *mocks.TaskRemover) {
			},
			id:         "",
			wantStatus: http.StatusBadRequest,
			wantResp:   delete.Response{Err: "gotten invalid id"},
		},
		{
			name: "unsuccessful deletion - database error",
			mockSetup: func(m *mocks.TaskRemover) {
				m.On("Delete", mock.Anything, int64(100)).Return(errors.New("database error"))
			},
			id:         "100",
			wantStatus: http.StatusInternalServerError,
			wantResp:   delete.Response{Err: "failed to delete task"},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mock := mocks.NewTaskRemover(t)
			tc.mockSetup(mock)

			handler := delete.New(slogdiscard.NewDiscardLogger(), mock)

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
			var actualResp delete.Response
			_ = json.Unmarshal(rec.Body.Bytes(), &actualResp)

			assert.Equal(t, tc.wantResp, actualResp)
			mock.AssertExpectations(t)
		})
	}
}
