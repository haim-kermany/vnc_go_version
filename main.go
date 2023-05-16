package main

import (
	"os"
	"text/template"
)

// //////////////////////////////////////////////////////////////////////////////
//
//
//
//
//
//
//
//
//
//
// //////////////////////////////////////////////////////////////////////////////

const (
	borderDistance int = 40
	subnetWidth    int = 8 * 40
	subnetHight    int = 6 * 40
	iconSize       int = 60
	iconSpace      int = 4 * 40
)

type Layer struct {
	thickness int
	location  int
	inUse     bool
	matrix    *Matrix
	index     int
}

func NewLayer(matrix *Matrix, index int) *Layer {
	return &Layer{borderDistance, 0, false, matrix, index}
}

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
		if row.inUse {
			row.index = len(rows)
			rows = append(rows, row)
		}
	}
	for _, col := range matrix.cols {
		if col.inUse {
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

///////////////////////////////////////////////////////////////////////////////////////////////
//////////
//////////
//////////
//////////
//////////
//////////
//////////
//////////
//////////
///////////////////////////////////////////////////////////////////////////////////////////////

func layout(network TreeNodeInterface) {
	m := LayoutIcons(network)
	m.Expend(func(i int) int { return 6*i + 3 })
	AddBorders(network, m)
	m.removeUnused()
	setIconsLocations(network, m)
	SetLayersLocation(m)
	resolveDrawioInfo(network)
}

func MergeLocations(tn TreeNodeInterface) {
	subtrees, leafs := tn.GetChildren()
	locations := []*Location{}
	for _, c := range append(subtrees, leafs...) {
		locations = append(locations, c.GetLocation())
	}
	var firstRow, lastRow, firstCol, lastCol *Layer = nil, nil, nil, nil
	for _, l := range locations {
		if l == nil {
			continue
		}
		if firstRow == nil || l.firstRow.index < firstRow.index {
			firstRow = l.firstRow
		}
		if lastRow == nil || l.lastRow.index > lastRow.index {
			lastRow = l.lastRow
		}
		if lastCol == nil || l.firstCol.index < firstCol.index {
			firstCol = l.firstCol
		}
		if lastCol == nil || l.lastCol.index > lastCol.index {
			lastCol = l.lastCol
		}
	}
	tn.SetLocation(NewLocation(firstRow, lastRow, firstCol, lastCol))
}

// ///////////////////////////////////////////////////////////////
func LayoutIcons(network TreeNodeInterface) *Matrix {
	m := NewMatrix(100, 100)
	colIndex := 0
	for _, vpc := range network.(*NetworkTreeNode).vpcs {
		for _, zone := range vpc.(*VpcTreeNode).zones {
			rowIndex := 0
			for _, subnet := range zone.(*ZoneTreeNode).subnets {
				rowFull := false
				for i, icon := range subnet.GetIconTreeNodes() {
					icon.SetLocation(NewCellLocation(m.rows[rowIndex], m.cols[colIndex]))
					m.rows[rowIndex].inUse = true
					m.cols[colIndex].inUse = true
					if rowFull || i == len(subnet.GetIconTreeNodes())-1 {
						rowIndex++
					}
					rowFull = !rowFull
				}
			}
			colIndex++
		}
	}
	m.removeUnused()
	for _, row := range m.rows {
		row.thickness = subnetHight
	}
	for _, col := range m.cols {
		col.thickness = subnetWidth
	}
	return m
}

/////////////////////////////////////////////////////////////

func ExpendLocation(tn TreeNodeInterface) {
	l := tn.GetLocation()
	l = NewLocation(l.PrevRow(), l.NextRow(), l.PrevCol(), l.NextCol())
	tn.SetLocation(l)
	l.firstRow.inUse = true
	l.lastRow.inUse = true
	l.firstCol.inUse = true
	l.lastCol.inUse = true

}

// ///////////////////////////////////////////////////////////////////////////////////////

func AddBorders(network TreeNodeInterface, m *Matrix) {

	for _, vpc := range network.(*NetworkTreeNode).vpcs {
		for _, zone := range vpc.(*VpcTreeNode).zones {
			for _, subnet := range zone.(*ZoneTreeNode).subnets {
				MergeLocations(subnet)
				subnet.GetLocation().PrevRow().inUse = true
				subnet.GetLocation().PrevCol().inUse = true
			}
			MergeLocations(zone)
			ExpendLocation(zone)
			zone.GetLocation().PrevRow().inUse = true
			zone.GetLocation().PrevCol().inUse = true
		}
		MergeLocations(vpc)
		ExpendLocation(vpc)
		vpc.GetLocation().PrevRow().inUse = true
		vpc.GetLocation().PrevCol().inUse = true
	}
	MergeLocations(network)
	ExpendLocation(network)
}

// ////////////////////////////////////////////////////////////////////////////////////////

// /////////////////////////////////////////////////////////////////////////////////////////////////////////
func setIconsLocations(network TreeNodeInterface, m *Matrix) {
	for _, vpc := range network.(*NetworkTreeNode).vpcs {
		icons := vpc.GetIconTreeNodes()
		if len(vpc.GetIconTreeNodes()) > 0 {
			vpc.GetLocation().firstRow.thickness = iconSpace
		}
		iconIndex := 0
		firstColIndex := vpc.GetLocation().firstCol.index
		lastColIndex := vpc.GetLocation().lastCol.index
		for ci := firstColIndex; ci <= lastColIndex && iconIndex < len(icons); ci++ {
			if m.cols[ci].thickness >= iconSpace || (ci > firstColIndex && m.cols[ci-1].thickness >= iconSpace) {
				icons[iconIndex].SetLocation(NewCellLocation(vpc.GetLocation().firstRow, m.cols[ci]))
				iconIndex++
			}
		}
		for _, zone := range vpc.(*VpcTreeNode).zones {
			icons := zone.GetIconTreeNodes()
			for _, icon := range icons {
				if icon.IsVSI() {
					vsiSubents := icon.(*VsiTreeNode).GetVsiSubnets()
					if len(*vsiSubents) == 1 {
						subnet := icon.(*VsiTreeNode).nis[0].GetParent()
						icon.SetParent(subnet)
						icon.SetLocation(NewCellLocation(subnet.GetLocation().firstRow, subnet.GetLocation().firstCol))
						icon.GetLocation().y_offset = iconSize
					} else {
						vpcLocation := icon.(*VsiTreeNode).nis[0].GetParent().GetLocation()
						location := NewCellLocation(vpcLocation.firstRow, vpcLocation.firstCol)
						location = NewCellLocation(vpcLocation.NextRow(), vpcLocation.firstCol)
						location.x_offset = subnetWidth/2 - iconSize/2
						icon.SetLocation(location)
					}

				} else if icon.IsGateway() {
					col := zone.(*ZoneTreeNode).subnets[0].GetLocation().firstCol
					row := zone.GetLocation().firstRow
					zone.GetLocation().firstRow.thickness = iconSpace
					icon.SetLocation(NewCellLocation(row, col))
					icon.GetLocation().x_offset -= subnetWidth/2 - iconSize/2

				}
			}
			for _, subnet := range zone.(*ZoneTreeNode).subnets {
				icons := subnet.GetIconTreeNodes()
				for _, icon1 := range icons {
					for _, icon2 := range icons {
						if icon1 != icon2 {
							if icon1.GetLocation().firstRow == icon2.GetLocation().firstRow && icon1.GetLocation().firstCol == icon2.GetLocation().firstCol {
								icon1.GetLocation().x_offset = iconSize
								icon2.GetLocation().x_offset = -iconSize
							}
						}
					}
				}
			}
		}
	}
	icons := network.GetIconTreeNodes()
	iconIndex := 0
	if len(icons) > 0 {
		network.GetLocation().firstCol.thickness = iconSpace
	}
	for ri, row := range m.rows {
		if iconIndex >= len(icons) {
			break
		}
		if row.thickness >= iconSpace || (ri > 0 && m.rows[ri-1].thickness >= iconSpace) {
			icons[iconIndex].SetLocation(NewCellLocation(row, m.cols[0]))
			iconIndex++
		}
	}
}

func setSGLocations(network TreeNodeInterface) {
	for _, vpc := range network.(*NetworkTreeNode).vpcs {
		for _, sg := range vpc.(*VpcTreeNode).sgs {
			sg.SetDrawioInfo()
		}
	}
}

func resolveDrawioInfo(network TreeNodeInterface) {
	network.SetDrawioInfo()
	for _, conn := range network.(*NetworkTreeNode).connections {
		conn.SetDrawioInfo()
	}
	for _, icon := range network.GetIconTreeNodes() {
		icon.SetDrawioInfo()
	}
	for _, vpc := range network.(*NetworkTreeNode).vpcs {
		vpc.SetDrawioInfo()
		for _, icon := range vpc.GetIconTreeNodes() {
			icon.SetDrawioInfo()
		}
		for _, sg := range vpc.(*VpcTreeNode).sgs {
			sg.SetDrawioInfo()
		}
		for _, zone := range vpc.(*VpcTreeNode).zones {
			zone.SetDrawioInfo()
			for _, icon := range zone.GetIconTreeNodes() {
				icon.SetDrawioInfo()
			}
			for _, subnet := range zone.(*ZoneTreeNode).subnets {
				subnet.SetDrawioInfo()
				for _, icon := range subnet.GetIconTreeNodes() {
					icon.SetDrawioInfo()
				}
			}
		}
	}
}

func SetLayersLocation(m *Matrix) {
	var y int = 0
	for _, row := range m.rows {
		row.location = y
		y = y + row.thickness
	}
	var x int = 0
	for _, col := range m.cols {
		col.location = x
		x = x + col.thickness
	}
}

///////////////////////////////////////////////////

//
//
//
//
//
// //////////////////////////////////////////////////////////////////////////////

func main() {
	network := create_network()
	layout(network)

	elements := GetElements(network)

	// templateFuncMap := template.FuncMap{
	// 	"add": func(x uint, y uint) uint {
	// 		return x + y
	// 	},
	// }
	// tmpl, err := template.New("drawio_template.txt").Funcs(templateFuncMap).ParseFiles("drawio_template.txt")
	tmpl, err := template.New("drawio_template.txt").ParseFiles("drawio_template.txt")
	if err != nil {
		panic(err)
	}
	fo, err := os.Create("output.drawio")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	type DrawioInfo struct {
		IconSize uint
		Nodes    []TreeNodeInterface
	}
	err = tmpl.Execute(fo, DrawioInfo{60, elements})
	if err != nil {
		panic(err)
	}

}
