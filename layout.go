package main

const (
	borderDistance int = 40
	subnetWidth    int = 8 * 40
	subnetHight    int = 6 * 40
	iconSize       int = 60
	iconSpace      int = 4 * 40
)

func layout(network TreeNodeInterface) {
	m := LayoutSubnetIcons(network)
	m.Expend(func(i int) int { return 6*i + 3 })
	SetSquersLocations(network)
	m.removeUnused()
	setIconsLocations(network, m)
	m.SetLayersDistance()
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
func LayoutSubnetIcons(network TreeNodeInterface) *Matrix {
	m := NewMatrix(100, 100)
	colIndex := 0
	for _, vpc := range network.(*NetworkTreeNode).vpcs {
		for _, zone := range vpc.(*VpcTreeNode).zones {
			rowIndex := 0
			for _, subnet := range zone.(*ZoneTreeNode).subnets {
				rowFull := false
				for i, icon := range subnet.GetIconTreeNodes() {
					icon.SetLocation(NewCellLocation(m.rows[rowIndex], m.cols[colIndex]))
					m.rows[rowIndex].SetThickness(subnetHight)
					m.cols[colIndex].SetThickness(subnetWidth)
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
	return m
}

/////////////////////////////////////////////////////////////

func ExpendLocation(tn TreeNodeInterface) {
	l := tn.GetLocation()
	l = NewLocation(l.PrevRow(), l.NextRow(), l.PrevCol(), l.NextCol())
	tn.SetLocation(l)
	l.firstRow.SetThickness(borderDistance)
	l.lastRow.SetThickness(borderDistance)
	l.firstCol.SetThickness(borderDistance)
	l.lastCol.SetThickness(borderDistance)

}

// ///////////////////////////////////////////////////////////////////////////////////////

func SetSquersLocations(network TreeNodeInterface) {

	for _, vpc := range network.(*NetworkTreeNode).vpcs {
		for _, zone := range vpc.(*VpcTreeNode).zones {
			for _, subnet := range zone.(*ZoneTreeNode).subnets {
				MergeLocations(subnet)
				subnet.GetLocation().PrevRow().SetThickness(borderDistance)
				subnet.GetLocation().PrevCol().SetThickness(borderDistance)
			}
			MergeLocations(zone)
			ExpendLocation(zone)
			zone.GetLocation().PrevRow().SetThickness(borderDistance)
			zone.GetLocation().PrevCol().SetThickness(borderDistance)
		}
		MergeLocations(vpc)
		ExpendLocation(vpc)
		vpc.GetLocation().PrevRow().SetThickness(borderDistance)
		vpc.GetLocation().PrevCol().SetThickness(borderDistance)
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
			vpc.GetLocation().firstRow.SetThickness(iconSpace)
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
					zone.GetLocation().firstRow.SetThickness(iconSpace)
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
		network.GetLocation().firstCol.SetThickness(iconSpace)
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
			MergeLocations(sg)
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
