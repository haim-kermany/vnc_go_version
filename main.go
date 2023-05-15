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

type Row struct {
	hight  int
	y      int
	inUse  bool
	matrix *Matrix
}

func NewRow(matrix *Matrix, hight int) *Row {
	r := &Row{}
	r.hight = hight
	r.matrix = matrix
	r.inUse = false
	return r
}

type Col struct {
	width  int
	x      int
	inUse  bool
	matrix *Matrix
}

func NewCol(matrix *Matrix, width int) *Col {
	c := &Col{}
	c.width = width
	c.matrix = matrix
	c.inUse = false
	return c
}

type Position struct {
	firstRow *Row
	lastRow  *Row
	firstCol *Col
	lastCol  *Col
	x_offset int
	y_offset int
}

func NewPosition(firstRow *Row, lastRow *Row, firstCol *Col, lastCol *Col) *Position {
	return &Position{firstRow, lastRow, firstCol, lastCol, 0, 0}
}

type PositionIndexes struct {
	firstRowIndex int
	lastRowIndex  int
	firstColIndex int
	lastColIndex  int
}

func (i *PositionIndexes) updateIndexes(new_i *PositionIndexes) {
	if i.firstRowIndex > new_i.firstRowIndex {
		i.firstRowIndex = new_i.firstRowIndex
	}
	if i.firstColIndex > new_i.firstColIndex {
		i.firstColIndex = new_i.firstColIndex
	}
	if i.lastRowIndex < new_i.lastRowIndex {
		i.lastRowIndex = new_i.lastRowIndex
	}
	if i.lastColIndex < new_i.lastColIndex {
		i.lastColIndex = new_i.lastColIndex
	}

}

type Matrix struct {
	rows []*Row
	cols []*Col
}

func NewMatrix(nrows int, ncols int) *Matrix {
	m := Matrix{}
	m.rows = make([]*Row, nrows)
	m.cols = make([]*Col, ncols)
	for i := range m.rows {
		m.rows[i] = NewRow(&m, borderDistance)
	}
	for i := range m.cols {
		m.cols[i] = NewCol(&m, borderDistance)
	}
	return &m
}

func (matrix *Matrix) getPosiotonIndexes(position *Position) PositionIndexes {
	positionIndex := PositionIndexes{}
	for i, row := range matrix.rows {
		if row == position.firstRow {
			positionIndex.firstRowIndex = i
		}
		if row == position.lastRow {
			positionIndex.lastRowIndex = i
		}
	}
	for i, col := range matrix.cols {
		if col == position.firstCol {
			positionIndex.firstColIndex = i
		}
		if col == position.lastCol {
			positionIndex.lastColIndex = i
		}
	}
	return positionIndex
}

func (matrix *Matrix) getPosioton(positionIndexes PositionIndexes) *Position {
	return NewPosition(
		matrix.rows[positionIndexes.firstRowIndex],
		matrix.rows[positionIndexes.lastRowIndex],
		matrix.cols[positionIndexes.firstColIndex],
		matrix.cols[positionIndexes.lastColIndex],
	)
}

func (matrix *Matrix) removeUnused() {
	rows := []*Row{}
	cols := []*Col{}
	for _, row := range matrix.rows {
		if row.inUse {
			rows = append(rows, row)
		}
	}
	for _, col := range matrix.cols {
		if col.inUse {
			cols = append(cols, col)
		}
	}
	matrix.rows = rows
	matrix.cols = cols
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
	m = AddBorders(network, m)
	setIconsPositions(network, m)
	SetLayersLocation(m)
	resolveDrawioInfo(network)
}

// ///////////////////////////////////////////////////////////////
func LayoutIcons(network TreeNodeInterface) *Matrix {
	m := NewMatrix(100, 100)
	indexes := PositionIndexes{}
	for _, vpc := range network.(*NetworkTreeNode).vpcs {
		indexes.firstRowIndex = 0
		indexes.lastRowIndex = 0
		for _, zone := range vpc.(*VpcTreeNode).zones {
			indexes.firstRowIndex = 0
			indexes.lastRowIndex = 0
			for _, subnet := range zone.(*ZoneTreeNode).subnets {
				rowFull := false
				for i, icon := range subnet.GetIconTreeNodes() {
					icon.SetPosition(m.getPosioton(indexes))
					if rowFull || i == len(subnet.GetIconTreeNodes())-1 {
						indexes.firstRowIndex++
						indexes.lastRowIndex++
					}
					rowFull = !rowFull
				}
			}
			indexes.firstColIndex++
			indexes.lastColIndex++
		}
	}
	return m
}

// ///////////////////////////////////////////////////////////////////////////////////////
func AddBorders(network TreeNodeInterface, m *Matrix) *Matrix {
	new_m := NewMatrix(6*len(m.rows)+1, 6*len(m.cols)+1)
	for i, row := range m.rows {
		row.inUse = true
		row.hight = subnetHight
		new_m.rows[i*6+3] = row
	}
	for i, col := range m.cols {
		col.inUse = true
		col.width = subnetWidth
		new_m.cols[i*6+3] = col
	}

	for _, vpc := range network.(*NetworkTreeNode).vpcs {
		vpcIndexes := PositionIndexes{len(m.rows), 0, len(m.cols), 0}
		for _, zone := range vpc.(*VpcTreeNode).zones {
			zoneIndexes := PositionIndexes{len(m.rows), 0, len(m.cols), 0}
			for _, subnet := range zone.(*ZoneTreeNode).subnets {
				subnetIndexes := PositionIndexes{len(m.rows), 0, len(m.cols), 0}
				for _, icon := range subnet.GetIconTreeNodes() {
					iconIndexes := m.getPosiotonIndexes(icon.GetPosition())
					subnetIndexes.updateIndexes(&iconIndexes)
				}
				subnet.SetPosition(m.getPosioton(subnetIndexes))

				vpcIndexes.updateIndexes(&subnetIndexes)
				zoneIndexes.updateIndexes(&subnetIndexes)
				newSubnetIndexes := PositionIndexes{
					subnetIndexes.firstRowIndex*6 + 3,
					subnetIndexes.lastRowIndex*6 + 3,
					subnetIndexes.firstColIndex*6 + 3,
					subnetIndexes.lastColIndex*6 + 3}
				new_m.rows[newSubnetIndexes.firstRowIndex-1].inUse = true
				new_m.rows[newSubnetIndexes.lastRowIndex+1].inUse = true
				new_m.cols[newSubnetIndexes.firstColIndex-1].inUse = true
				new_m.cols[newSubnetIndexes.lastColIndex+1].inUse = true

			}
			newZoneIndexes := PositionIndexes{
				zoneIndexes.firstRowIndex*6 + 2,
				zoneIndexes.lastRowIndex*6 + 4,
				zoneIndexes.firstColIndex*6 + 2,
				zoneIndexes.lastColIndex*6 + 4}

			new_m.rows[newZoneIndexes.firstRowIndex-1].inUse = true
			new_m.cols[newZoneIndexes.firstColIndex-1].inUse = true

			new_m.rows[newZoneIndexes.firstRowIndex].inUse = true
			new_m.rows[newZoneIndexes.lastRowIndex].inUse = true
			new_m.cols[newZoneIndexes.firstColIndex].inUse = true
			new_m.cols[newZoneIndexes.lastColIndex].inUse = true

			zone.SetPosition(new_m.getPosioton(newZoneIndexes))
		}
		newVpcIndexes := PositionIndexes{
			vpcIndexes.firstRowIndex*6 + 1,
			vpcIndexes.lastRowIndex*6 + 5,
			vpcIndexes.firstColIndex*6 + 1,
			vpcIndexes.lastColIndex*6 + 5}

		new_m.rows[newVpcIndexes.firstRowIndex-1].inUse = true
		new_m.cols[newVpcIndexes.firstColIndex-1].inUse = true

		new_m.rows[newVpcIndexes.firstRowIndex].inUse = true
		new_m.rows[newVpcIndexes.lastRowIndex].inUse = true
		new_m.cols[newVpcIndexes.firstColIndex].inUse = true
		new_m.cols[newVpcIndexes.lastColIndex].inUse = true

		vpc.SetPosition(new_m.getPosioton(newVpcIndexes))
	}
	new_m.rows[len(new_m.rows)-1].inUse = true
	new_m.cols[len(new_m.cols)-1].inUse = true
	new_m.removeUnused()
	network.SetPosition(NewPosition(new_m.rows[0], new_m.rows[len(m.rows)-1], new_m.cols[0], new_m.cols[len(m.cols)-1]))
	return new_m
}

// ////////////////////////////////////////////////////////////////////////////////////////

// /////////////////////////////////////////////////////////////////////////////////////////////////////////
func setIconsPositions(network TreeNodeInterface, m *Matrix) {
	for _, vpc := range network.(*NetworkTreeNode).vpcs {
		icons := vpc.GetIconTreeNodes()
		if len(vpc.GetIconTreeNodes()) > 0 {
			vpc.GetPosition().firstRow.hight = iconSpace
		}
		iconIndex := 0
		indexes := m.getPosiotonIndexes(vpc.GetPosition())
		for ci := indexes.firstColIndex; ci <= indexes.lastColIndex && iconIndex < len(icons); ci++ {
			if m.cols[ci].width >= iconSpace || (ci > indexes.firstColIndex && m.cols[ci-1].width >= iconSpace) {
				icons[iconIndex].SetPosition(NewPosition(m.rows[indexes.firstRowIndex], m.rows[indexes.firstRowIndex], m.cols[ci], m.cols[ci]))
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
						icon.SetPosition(NewPosition(subnet.GetPosition().firstRow, subnet.GetPosition().firstRow, subnet.GetPosition().firstCol, subnet.GetPosition().firstCol))
						icon.GetPosition().y_offset = iconSize
					} else {
						positionIndex := m.getPosiotonIndexes(icon.(*VsiTreeNode).nis[0].GetParent().GetPosition())
						positionIndex.firstRowIndex += 1
						position := m.getPosioton(positionIndex)
						position.lastRow = position.firstRow
						position.lastCol = position.firstCol
						position.x_offset = subnetWidth/2 - iconSize/2
						icon.SetPosition(position)
					}

				} else if icon.IsGateway() {
					col := zone.(*ZoneTreeNode).subnets[0].GetPosition().firstCol
					row := zone.GetPosition().firstRow
					zone.GetPosition().firstRow.hight = iconSpace
					icon.SetPosition(NewPosition(row, row, col, col))
					icon.GetPosition().x_offset -= subnetWidth/2 - iconSize/2

				}
			}
			for _, subnet := range zone.(*ZoneTreeNode).subnets {
				icons := subnet.GetIconTreeNodes()
				for _, icon1 := range icons {
					for _, icon2 := range icons {
						if icon1 != icon2 {
							if icon1.GetPosition().firstRow == icon2.GetPosition().firstRow && icon1.GetPosition().firstCol == icon2.GetPosition().firstCol {
								icon1.GetPosition().x_offset = iconSize
								icon2.GetPosition().x_offset = -iconSize
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
		network.GetPosition().firstCol.width = iconSpace
	}
	for ri, row := range m.rows {
		if iconIndex >= len(icons) {
			break
		}
		if row.hight >= iconSpace || (ri > 0 && m.rows[ri-1].hight >= iconSpace) {
			icons[iconIndex].SetPosition(NewPosition(row, row, m.cols[0], m.cols[0]))
			iconIndex++
		}
	}
}

func setSGPositions(network TreeNodeInterface) {
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
		row.y = y
		y = y + row.hight
	}
	var x int = 0
	for _, col := range m.cols {
		col.x = x
		x = x + col.width
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
