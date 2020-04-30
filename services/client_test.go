package services

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/putdotio/go-putio"
	"golang.org/x/oauth2"
)

func TestNewPutioClient(t *testing.T) {
	tk := "token"
	os.Setenv("PUT_IO_TOKEN", tk)
	defer os.Clearenv()
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tk})
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	expectedClient := putio.NewClient(oauthClient)

	tests := []struct {
		name string
		want *putio.Client
	}{
		{
			name: "Same client",
			want: expectedClient,
		},
		{
			name: "Duplicate for singleton",
			want: expectedClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPutioClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPutioClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
