package services

import (
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
)

func TestNewCache(t *testing.T) {
	expectedClient := &Cache{cache.New(10*time.Second, 1*time.Minute), true}
	tests := []struct {
		name   string
		enable bool
		want   *Cache
	}{
		{
			name:   "Valid client",
			enable: true,
			want:   expectedClient,
		},
		{
			name:   "Duplicate for singleton",
			enable: true,
			want:   expectedClient,
		},
		{
			name:   "Disable cache",
			enable: false,
			want:   &Cache{cache.New(10*time.Second, 1*time.Minute), false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCache(tt.enable)
			if got.enable != tt.enable {
				t.Errorf("NewCache() enable = %v want %v", got.enable, tt.enable)
			}
		})
	}
}
