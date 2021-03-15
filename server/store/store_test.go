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
