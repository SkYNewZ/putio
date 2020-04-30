package handlers

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

func Test_parseDirectoryID(t *testing.T) {
	dummyRequest := http.Request{}
	mux.SetURLVars(&dummyRequest, map[string]string{
		"folder": "hello",
	})

	tests := []struct {
		name    string
		r       *http.Request
		want    int64
		wantErr bool
	}{
		{
			name: "Invalid folder",
			r: func() *http.Request {
				dummyRequest := http.Request{}
				return mux.SetURLVars(&dummyRequest, map[string]string{
					"folder": "hello",
				})
			}(),
			wantErr: true,
			want:    0,
		},
		{
			name: "Valid folder",
			r: func() *http.Request {
				dummyRequest := http.Request{}
				return mux.SetURLVars(&dummyRequest, map[string]string{
					"folder": "1234",
				})
			}(),
			wantErr: false,
			want:    1234,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDirectoryID(tt.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDirectoryID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseDirectoryID() = %v, want %v", got, tt.want)
			}
		})
	}
}
