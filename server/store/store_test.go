package store_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/odpf/stencil/server/store"
	"github.com/stretchr/testify/assert"
	"gocloud.dev/blob"
)

func setupStorage() *store.Store {
	ctx := context.Background()
	bucket, _ := blob.OpenBucket(ctx, "mem://test-store")
	return &store.Store{
		Bucket: bucket,
	}
}

func seedData(s *store.Store, filename, contents string) {
	reader := bytes.NewReader([]byte(contents))
	ctx := context.Background()
	s.Put(ctx, filename, reader)
}

func verifyData(t *testing.T, s *store.Store, path, data string) {
	ctx := context.Background()
	reader, _ := s.Get(ctx, path)
	result, _ := ioutil.ReadAll(reader)
	assert.Equal(t, []byte(data), result)
}

func TestListFiles(t *testing.T) {
	s := setupStorage()
	seedData(s, "/n/k/v", "file data")
	seedData(s, "/n/k/v2", "file data 2")
	result, _ := s.ListFiles("/n/k/")
	assert.Equal(t, []string{"v", "v2"}, result)
}

func TestListDir(t *testing.T) {
	s := setupStorage()
	seedData(s, "/n/k/v", "file data")
	seedData(s, "/n/k2/v", "file data 2")
	result, _ := s.ListDir("/n/")
	assert.Equal(t, []string{"k", "k2"}, result)
}

func TestGet(t *testing.T) {
	s := setupStorage()
	seedData(s, "/n/k/v", "file data")
	ctx := context.Background()
	reader, _ := s.Get(ctx, "/n/k/v")
	result, _ := ioutil.ReadAll(reader)
	assert.Equal(t, []byte("file data"), result)
}

func TestCopy(t *testing.T) {
	s := setupStorage()
	seedData(s, "/n/k/v1", "file data")
	ctx := context.Background()
	_ = s.Copy(ctx, "/n/k/v1", "/n/k/v2")
	verifyData(t, s, "/n/k/v2", "file data")
}

func TestExists(t *testing.T) {
	t.Run("should return true if file exists", func(t *testing.T) {
		s := setupStorage()
		seedData(s, "/n/k/v", "file data")
		ctx := context.Background()
		ok, err := s.Exists(ctx, "/n/k/v")
		assert.Equal(t, true, ok)
		assert.Nil(t, err)
	})
	t.Run("should return error if file not exists", func(t *testing.T) {
		s := setupStorage()
		ctx := context.Background()
		ok, _ := s.Exists(ctx, "/unknown")
		assert.Equal(t, false, ok)
	})
}
