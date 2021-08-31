package api_test

import (
	"context"
	"errors"
	"testing"

	"github.com/odpf/stencil/server/api/v1/pb"
	"github.com/odpf/stencil/server/snapshot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestList(t *testing.T) {
	t.Run("should return list", func(t *testing.T) {
		ctx := context.Background()
		_, _, mockService, v1 := setup()
		st := []*snapshot.Snapshot{
			{
				Namespace: "t",
				Name:      "na",
			},
		}
		req := pb.ListSnapshotRequest{
			Namespace: "t",
		}
		mockService.On("List", mock.Anything, &snapshot.Snapshot{Namespace: "t"}).Return(st, nil)
		res, err := v1.List(ctx, &req)
		assert.Nil(t, err)
		assert.Equal(t, "t", res.Snapshots[0].Namespace)
		assert.Equal(t, "na", res.Snapshots[0].Name)
	})

	t.Run("should return error if getting a list fails", func(t *testing.T) {
		ctx := context.Background()
		_, _, mockService, v1 := setup()
		req := pb.ListSnapshotRequest{
			Namespace: "t",
		}
		err := errors.New("list failed")
		mockService.On("List", mock.Anything, &snapshot.Snapshot{Namespace: "t"}).Return(nil, err)
		res, err := v1.List(ctx, &req)
		assert.NotNil(t, err)
		assert.Equal(t, 0, len(res.Snapshots))
	})
}

func TestUpdateLatestVersion(t *testing.T) {
	t.Run("should update latest tag", func(t *testing.T) {
		ctx := context.Background()
		_, _, mockService, v1 := setup()
		st := &snapshot.Snapshot{
			ID:        1,
			Namespace: "t",
			Name:      "na",
		}
		req := &pb.UpdateLatestRequest{
			Id:     1,
			Latest: true,
		}
		mockService.On("GetSnapshotByID", mock.Anything, int64(1)).Return(st, nil)
		mockService.On("UpdateLatestVersion", mock.Anything, st).Return(nil)
		res, err := v1.UpdateLatest(ctx, req)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), res.Id)
		assert.Equal(t, "t", res.Namespace)
		assert.Equal(t, "na", res.Name)
	})
	t.Run("should return not found err if snapshot not found", func(t *testing.T) {
		ctx := context.Background()
		_, _, mockService, v1 := setup()
		st := &snapshot.Snapshot{
			ID:        1,
			Namespace: "t",
			Name:      "na",
		}
		req := &pb.UpdateLatestRequest{
			Id:     1,
			Latest: true,
		}
		mockService.On("GetSnapshotByID", mock.Anything, int64(1)).Return(st, snapshot.ErrNotFound)
		mockService.On("UpdateLatestVersion", mock.Anything, st).Return(nil)
		_, err := v1.UpdateLatest(ctx, req)
		assert.NotNil(t, err)
		s := status.Convert(err)
		assert.Equal(t, codes.NotFound.String(), s.Code().String())
		assert.Equal(t, "not found", s.Message())
	})
	t.Run("should mark as internal error if get snapshot fails", func(t *testing.T) {
		ctx := context.Background()
		_, _, mockService, v1 := setup()
		st := &snapshot.Snapshot{
			ID:        1,
			Namespace: "t",
			Name:      "na",
		}
		req := &pb.UpdateLatestRequest{
			Id:     1,
			Latest: true,
		}
		err := errors.New("internal")
		mockService.On("GetSnapshotByID", mock.Anything, int64(1)).Return(st, err)
		mockService.On("UpdateLatestVersion", mock.Anything, st).Return(nil)
		_, err = v1.UpdateLatest(ctx, req)
		assert.NotNil(t, err)
		s := status.Convert(err)
		assert.Equal(t, codes.Internal.String(), s.Code().String())
		assert.Equal(t, "internal", s.Message())
	})
	t.Run("should mark as internal error if update snapshot fails", func(t *testing.T) {
		ctx := context.Background()
		_, _, mockService, v1 := setup()
		st := &snapshot.Snapshot{
			ID:        1,
			Namespace: "t",
			Name:      "na",
		}
		req := &pb.UpdateLatestRequest{
			Id:     1,
			Latest: true,
		}
		err := errors.New("internal")
		mockService.On("GetSnapshotByID", mock.Anything, int64(1)).Return(st, nil)
		mockService.On("UpdateLatestVersion", mock.Anything, st).Return(err)
		_, err = v1.UpdateLatest(ctx, req)
		assert.NotNil(t, err)
		s := status.Convert(err)
		assert.Equal(t, codes.Unknown.String(), s.Code().String())
		assert.Equal(t, "internal", s.Message())
	})
}
