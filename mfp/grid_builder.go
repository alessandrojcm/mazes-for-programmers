package mfp

import "errors"

type GridBuilder struct {
	rows, columns int
	baseGrid      BaseGrid
}

type GridBuilderHandler interface {
	BuildBaseGrid() (BaseGrid, error)
	BuildASCIIGrid() (ASCIIGrid, error)
	BuildRenderGrid() (RendererGrid, error)
	//BuildGridWithDistance() (DistanceGrid, error)
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
	if g.baseGrid.Empty() {
		_, err = g.BuildBaseGrid()
	}
	return ASCIIGrid{
		&g.baseGrid,
	}, err
}

func (g *GridBuilder) BuildRenderGrid() (RendererGrid, error) {
	asciiGrid, err := g.BuildASCIIGrid()
	return RendererGrid{&asciiGrid}, err
}

//
//func (g *GridBuilder) BuildGridWithDistance() (DistanceGrid, error) {
//	//TODO implement me
//	panic("implement me")
//}
