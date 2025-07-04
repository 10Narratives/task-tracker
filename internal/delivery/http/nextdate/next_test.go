package next_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	next "github.com/10Narratives/task-tracker/internal/delivery/http/nextdate"
	"github.com/10Narratives/task-tracker/internal/lib/logging/handlers/slogdiscard"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestNextDateHandler(t *testing.T) {
	tests := []struct {
		name       string
		now        string
		date       string
		repeat     string
		wantStatus int
		wantResp   next.Response
	}{
		{
			name:       "d 10",
			now:        "20250210",
			date:       "20250210",
			repeat:     "d 10",
			wantStatus: http.StatusOK,
			wantResp:   next.Response{NextDate: "20250220"},
		},
		{
			name:       "invalid dateformat",
			now:        "2025-02-10",
			date:       "2025-02-10",
			repeat:     "d 10",
			wantStatus: http.StatusBadRequest,
			wantResp:   next.Response{Err: "field Now must be in YYYYMMDD date format, field Date must be in YYYYMMDD date format"},
		},
		{
			name:       "invalid repeat",
			now:        "20250210",
			date:       "20250210",
			repeat:     "d 1000",
			wantStatus: http.StatusBadRequest,
			wantResp:   next.Response{Err: "field Repeat must satisfy expected patterns"},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := next.New(slogdiscard.NewDiscardLogger())
			pattern := "/api/nextdate"
			url := pattern + "?now=" + tc.now + "&date=" + tc.date + "&repeat=" + url.QueryEscape(tc.repeat)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Get(pattern, handler)
			r.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			var actualResp next.Response
			_ = json.Unmarshal(rec.Body.Bytes(), &actualResp)

			assert.Equal(t, tc.wantResp, actualResp)
		})
	}
}
