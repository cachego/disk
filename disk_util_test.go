package cache_test

import (
	"testing"

	cache "github.com/cachego/disk"
)

func TestDeleteFile(t *testing.T) {
	err := cache.CoverFile(".cache", "test", []byte("test"))
	if err != nil {
		t.Errorf("expected value to be nil, got '%s'", err)
	}
	err = cache.DeleteFile(".cache", "test")
	if err != nil {
		t.Errorf("expected value to be nil, got '%s'", err)
	}
}

func TestGetFileData(t *testing.T) {
	err := cache.CoverFile(".cache", "test", []byte("data1"))
	if err != nil {
		t.Errorf("expected value to be nil, got '%s'", err)
	}
	data, err := cache.GetFileData(".cache", "test")
	if err != nil {
		t.Errorf("expected value to be nil, got '%s'", err)
	}
	if string(data) != "data1" {
		t.Errorf("expected value to be 'test', got '%s'", data)
	}
	err = cache.DeleteFile(".cache", "test")
	if err != nil {
		t.Errorf("expected value to be nil, got '%s'", err)
	}
}
