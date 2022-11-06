package mfp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGrid(t *testing.T) {
	t.Parallel()
	grid, err := NewGrid(4, 4)
	assert.ErrorIs(t, err, nil)

	assert.Equal(t, 4, grid.rows)
	assert.Equal(t, 4, grid.columns)
	assert.Equal(t, 4, len(grid.cells))

	for column := range grid.EachRow() {
		assert.Equal(t, 4, len(column), fmt.Sprint(column))
	}
}

func TestRowGenerator(t *testing.T) {
	t.Parallel()
	grid, err := NewGrid(4, 4)
	assert.ErrorIs(t, err, nil)
	for column := range grid.EachRow() {
		assert.IsType(t, []*Cell{}, column)
	}
}

func TestCellGenerator(t *testing.T) {
	t.Parallel()
	grid, err := NewGrid(4, 4)
	assert.ErrorIs(t, err, nil)
	for cell := range grid.EachCell() {
		assert.IsType(t, &Cell{}, cell)
	}
}

func TestRandom(t *testing.T) {
	t.Parallel()
	grid, err := NewGrid(4, 4)
	assert.ErrorIs(t, err, nil)
	_, err = grid.RandomCell()
	assert.ErrorIs(t, err, nil)
}
