package grids

import (
	"errors"
	rl "github.com/gen2brain/raylib-go/raylib"
	"mazes-for-programmers/mfp"
)

// GridBuilder -- Holds the underlying properties of the grid
type GridBuilder struct {
	rows, columns int
	baseGrid      BaseGrid
}

// GridBuilderHandler -- this interfaces uses the builder pattern to build all the available types of grids
// structs that implement this interface should hold a baseGrid in memory and use that to build upon.
// This way we gan generate new grids without losing the information of the original grid.
type GridBuilderHandler interface {
	BuildBaseGrid() (BaseGrid, error)
	BuildASCIIGrid() (ASCIIGrid, error)
	BuildGridLineRenderer() (RendererGrid, error)
	BuildGridTiledRenderer() (TiledGrid, error)
	BuildGridWithDistanceRenderer(backgroundColor rl.Color) (DistanceRenderGrid, error)
	BuildGridWithDistance() (DistanceGrid, error)
	BuildAnimatableGrid(backgroundColor rl.Color) (AnimatableGrid, error)
}

func NewBuilder(rows, column int) *GridBuilder {
	return &GridBuilder{
		rows,
		column,
		BaseGrid{},
	}
}

func (g *GridBuilder) BuildBaseGrid() (BaseGrid, error) {
	if g.rows <= 0 || g.columns <= 0 {
		return BaseGrid{}, errors.New("rows and columns must be greater than 0")
	}
	g.baseGrid = BaseGrid{
		rows:    g.rows,
		columns: g.columns,
	}
	err := g.baseGrid.PrepareGrid()
	if err != nil {
		return BaseGrid{}, err
	}
	err = g.baseGrid.ConfigureCells()
	if err != nil {
		return BaseGrid{}, err
	}
	return g.baseGrid, nil
}

func (g *GridBuilder) BuildASCIIGrid() (ASCIIGrid, error) {
	var err error
	// Assure we have a grid
	if g.baseGrid.Empty() {
		_, err = g.BuildBaseGrid()
	}
	return ASCIIGrid{
		&g.baseGrid,
	}, err
}

func (g *GridBuilder) BuildGridLineRenderer() (RendererGrid, error) {
	asciiGrid, err := g.BuildASCIIGrid()
	return RendererGrid{&asciiGrid}, err
}

func (g *GridBuilder) BuildGridWithDistance() (DistanceGrid, error) {
	asciiGrid, err := g.BuildASCIIGrid()
	return DistanceGrid{
		&asciiGrid,
		mfp.Distance{},
	}, err
}

func (g *GridBuilder) BuildGridTiledRenderer() (TiledGrid, error) {
	grid, err := g.BuildGridLineRenderer()
	return TiledGrid{
		&grid,
	}, err
}

func (g *GridBuilder) BuildGridWithDistanceRenderer(backgroundColor rl.Color) (DistanceRenderGrid, error) {
	distanceGrid, err := g.BuildGridWithDistance()
	return DistanceRenderGrid{
		backgroundColor,
		&distanceGrid,
	}, err
}

func (g *GridBuilder) BuildAnimatableGrid(backgroundColor rl.Color) (AnimatableGrid, error) {
	distanceGrid, err := g.BuildGridWithDistanceRenderer(backgroundColor)

	return AnimatableGrid{
		&distanceGrid,
	}, err
}
