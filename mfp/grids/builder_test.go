package grids_test

import (
	"github.com/stretchr/testify/assert"
	"mazes-for-programmers/mfp/grids"
	"testing"
)

func TestGridBuilder_BuildBaseGrid(t *testing.T) {
	t.Parallel()
	builder := grids.NewBuilder(4, 4)

	grid, err := builder.BuildBaseGrid()

	assert.ErrorIs(t, err, nil)
	assert.Equal(t, grid.Size(), 16)
}

func TestGridBuilder_BuildASCIIGrid(t *testing.T) {
	t.Parallel()
	builder := grids.NewBuilder(4, 4)

	grid, err := builder.BuildASCIIGrid()

	assert.ErrorIs(t, err, nil)
	assert.Equal(t, grid.Size(), 16)
}
