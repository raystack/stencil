package search

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNamespaceNotFound = errors.New("namespace not found")
)

type InMemoryStore struct {
	indexMap map[string]map[string]map[Schema]struct{}
	*sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		indexMap: make(map[string]map[string]map[Schema]struct{}),
	}
}

func (m *InMemoryStore) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	m.RLock()
	defer m.RUnlock()
	namespaceMap, ok := m.indexMap[req.Namespace]
	if !ok {
		return nil, ErrNamespaceNotFound
	}

	fieldMap, ok := namespaceMap[req.Field]
	if !ok {
		return &SearchResponse{}, nil
	}

	schemas := make([]*Schema, 0)
	for schema, _ := range fieldMap {
		schemas = append(schemas, &schema)
	}

	return &SearchResponse{Schemas: schemas}, nil
}

func (m *InMemoryStore) Index(ctx context.Context, req *IndexRequest) error {
	m.Lock()
	defer m.Unlock()
	namespaceMap, ok := m.indexMap[req.Namespace]
	if !ok {
		namespaceMap = make(map[string]map[Schema]struct{})
	}

	for _, field := range req.Fields {
		fieldMap, ok := namespaceMap[field]
		if !ok {
			fieldMap = make(map[Schema]struct{})
		}
		fieldMap[Schema{
			Namespace: req.Namespace,
			Version:   req.Version,
			Message:   req.Message,
			Name:      req.Name,
			Latest:    req.Latest,
			Package:   req.Package,
		}] = struct{}{}
		namespaceMap[field] = fieldMap
	}

	m.indexMap[req.Namespace] = namespaceMap
	return nil
}
