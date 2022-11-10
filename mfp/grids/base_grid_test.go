package grids_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"mazes-for-programmers/mfp"
	"mazes-for-programmers/mfp/grids"
	"testing"
)

func TestNewGrid(t *testing.T) {
	t.Parallel()
	builder := grids.NewBuilder(
		4,
		4,
	)
	grid, err := builder.BuildBaseGrid()
	assert.ErrorIs(t, err, nil)

	assert.Equal(t, 16, grid.Size())

	for column := range grid.EachRow() {
		assert.Equal(t, 4, len(column), fmt.Sprint(column))
	}
}

func TestRowGenerator(t *testing.T) {
	t.Parallel()
	builder := grids.NewBuilder(
		4,
		4,
	)
	grid, err := builder.BuildBaseGrid()
	assert.ErrorIs(t, err, nil)
	for column := range grid.EachRow() {
		assert.IsType(t, []*mfp.Cell{}, column)
	}
}

func TestCellGenerator(t *testing.T) {
	t.Parallel()
	builder := grids.NewBuilder(
		4,
		4,
	)
	grid, _ := builder.BuildBaseGrid()
	for cell := range grid.EachCell() {
		assert.IsType(t, &mfp.Cell{}, cell)
	}
}

func TestRandom(t *testing.T) {
	t.Parallel()
	builder := grids.NewBuilder(
		4,
		4,
	)
	grid, err := builder.BuildBaseGrid()
	assert.ErrorIs(t, err, nil)
	_, err = grid.RandomCell()
	assert.ErrorIs(t, err, nil)
}
