package batchhll

import (
	"context"

	"github.com/MasterOfBinary/gobatch/batch"
	"github.com/go-redis/redis/v8"
)

func New(rdb *redis.Client, key string) *HLL {
	return &HLL{
		rdb: rdb,
		key: key,
	}
}

// HLL is a very simple batch HyperLogLog Processor.
type HLL struct {
	rdb *redis.Client
	key string
}

func (h *HLL) Process(ctx context.Context, ps *batch.PipelineStage) {
	defer ps.Close()

	allinput := make([]interface{}, 0, 10)

	for in := range ps.Input {
		item := in.Get()
		allinput = append(allinput, item)
	}

	// It seems unlikely that the result would be very useful if the caller doesn't
	// know how many were processed in the batch so just ignore the response.
	err := h.rdb.PFAdd(ctx, h.key, allinput...).Err()

	if err != nil {
		ps.Errors <- err
	}
}
