package main

type Layer struct {
	thickness int
	distance  int
	matrix    *Matrix
	index     int
}

func NewLayer(matrix *Matrix, index int) *Layer {
	return &Layer{0, 0, matrix, index}
}
func (l *Layer) SetThickness(thickness int) { l.thickness = thickness }

type Location struct {
	firstRow *Layer
	lastRow  *Layer
	firstCol *Layer
	lastCol  *Layer
	x_offset int
	y_offset int
}

func NewLocation(firstRow *Layer, lastRow *Layer, firstCol *Layer, lastCol *Layer) *Location {
	return &Location{firstRow, lastRow, firstCol, lastCol, 0, 0}
}
func NewCellLocation(firstRow *Layer, firstCol *Layer) *Location {
	return &Location{firstRow, firstRow, firstCol, firstCol, 0, 0}
}
func (l *Location) NextRow() *Layer {
	return l.lastRow.matrix.rows[l.lastRow.index+1]
}
func (l *Location) NextCol() *Layer {
	return l.lastCol.matrix.cols[l.lastCol.index+1]
}
func (l *Location) PrevRow() *Layer {
	return l.firstRow.matrix.rows[l.firstRow.index-1]
}
func (l *Location) PrevCol() *Layer {
	return l.firstCol.matrix.cols[l.firstCol.index-1]
}

////////////////////////////////////////////////////

type Matrix struct {
	rows []*Layer
	cols []*Layer
}

func NewMatrix(nrows int, ncols int) *Matrix {
	m := Matrix{}
	m.rows = make([]*Layer, nrows)
	m.cols = make([]*Layer, ncols)
	for i := range m.rows {
		m.rows[i] = NewLayer(&m, i)
	}
	for i := range m.cols {
		m.cols[i] = NewLayer(&m, i)
	}
	return &m
}

func (matrix *Matrix) removeUnused() {
	rows := []*Layer{}
	cols := []*Layer{}
	for _, row := range matrix.rows {
		if row.thickness > 0 {
			row.index = len(rows)
			rows = append(rows, row)
		}
	}
	for _, col := range matrix.cols {
		if col.thickness > 0 {
			col.index = len(cols)
			cols = append(cols, col)
		}
	}
	matrix.rows = rows
	matrix.cols = cols
}

func (matrix *Matrix) Expend(new_i func(int) int) {
	new_rows := make([]*Layer, new_i(len(matrix.rows)))
	new_cols := make([]*Layer, new_i(len(matrix.cols)))
	for i := range new_rows {
		new_rows[i] = NewLayer(matrix, i)
	}
	for i := range new_cols {
		new_cols[i] = NewLayer(matrix, i)
	}

	for i, row := range matrix.rows {
		new_rows[new_i(i)] = row
		row.index = new_i(i)
	}
	for i, col := range matrix.cols {
		new_cols[new_i(i)] = col
		col.index = new_i(i)
	}
	matrix.rows = new_rows
	matrix.cols = new_cols

}
func (m *Matrix) SetLayersDistance() {
	var y int = 0
	for _, row := range m.rows {
		row.distance = y
		y = y + row.thickness
	}
	var x int = 0
	for _, col := range m.cols {
		col.distance = x
		x = x + col.thickness
	}
}
