package main

import (
	"os"
	"text/template"
)

type DrawioElementInterface interface {
	GetId() uint
	SetId()
}

type DrawioElement struct {
	deType string
	id     uint
}

func (di DrawioElement) GetId() uint {
	return di.id
}

var id uint = 100

func (di *DrawioElement) SetId() {
	di.id = id
	id += 5
}

// //////////////////////////////////////////////////////////////////////////////
type DrawioLocatedElement struct {
	DrawioElement
	parentId uint
	x        int
	y        int
	name     string
}

func (di DrawioLocatedElement) GetX() int {
	return di.x
}
func (di DrawioLocatedElement) GetY() int {
	return di.y
}
func (di DrawioLocatedElement) GetParentId() uint {
	return di.parentId
}
func (di DrawioLocatedElement) GetName() string {
	return di.name
}

// //////////////////////////////////////////////////////////////////////////////
type DrawioSquereElement struct {
	DrawioLocatedElement
	width int
	hight int
}

func (di DrawioSquereElement) GetHight() int {
	return di.hight
}
func (di DrawioSquereElement) GetWidth() int {
	return di.width
}

// //////////////////////////////////////////////////////////////////////////////////////
type DrawioSubnetElement struct {
	DrawioSquereElement
	ip  string
	key string
}

func (di DrawioSubnetElement) GetIP() string {
	return di.ip
}
func (di DrawioSubnetElement) GetKey() string {
	return di.key
}

// ///////////////////////////////////////////////////////////////////////////////////////////////
type DrawioIconElement struct {
	DrawioLocatedElement
}

// ///////////////////////////////////////////////////////////////////////////////////////////////
type DrawioNetworkInterfaceElement struct {
	DrawioIconElement
	floating_ip string
	svi         string
}

func (di DrawioNetworkInterfaceElement) GetFip() string {
	return di.floating_ip
}

func (di DrawioNetworkInterfaceElement) GetSvi() string {
	return di.svi
}

type DrawioConnectElementPoint struct {
	x int
	y int
}

func (di DrawioConnectElementPoint) GetX() int {
	return di.x
}
func (di DrawioConnectElementPoint) GetY() int {
	return di.y
}

// //////////////////////////////////////////////////////////////////////////////
type DrawioConnectElement struct {
	DrawioElement
	srcId  uint
	dstId  uint
	label  string
	points []DrawioConnectElementPoint
}

func (di DrawioConnectElement) GetSrcId() uint {
	return di.srcId
}
func (di DrawioConnectElement) GetDstId() uint {
	return di.dstId
}
func (di DrawioConnectElement) GetLabel() string {
	return di.label
}

func (di DrawioConnectElement) GetPoints() []DrawioConnectElementPoint {
	return di.points
}

func NewDrawioConnectElement(etype string, id uint, srcid uint, dstid uint, label string, points []DrawioConnectElementPoint) *DrawioConnectElement {
	return &DrawioConnectElement{DrawioElement{etype, id}, srcid, dstid, label, points}
}

///////////////////////////////////////
////
///
//
//
//
//
////
///
//
//
//
//
////
///
//
//
//
//
////
///
//
//
//
//
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type TreeNodeInterface interface {
	GetParent() TreeNodeInterface
	AddIconTreeNode(icon TreeNodeInterface)
	GetIconTreeNodes() []TreeNodeInterface
	GetChildren() ([]TreeNodeInterface, []TreeNodeInterface)
	GetType() string
	GetDrawioElementInterface() DrawioElementInterface
	SetPosition(position *Position)
	GetPosition() *Position
	SetDrawioInfo()
}

///////////////////////////////////////////////////////////////////////////

type TreeNode struct {
	DrawioElementInterface
	parent   TreeNodeInterface
	elements []TreeNodeInterface
	position *Position
}

func (tn *TreeNode) GetDrawioElementInterface() DrawioElementInterface {
	return tn.DrawioElementInterface
}

func (tn *TreeNode) GetParent() TreeNodeInterface {
	return tn.parent
}
func (tn *TreeNode) AddIconTreeNode(icon TreeNodeInterface) {
	tn.elements = append(tn.elements, icon)
}
func (tn *TreeNode) GetIconTreeNodes() []TreeNodeInterface {
	return tn.elements
}

func NewTreeNode(drawioElementInterface DrawioElementInterface, parent TreeNodeInterface) TreeNode {
	drawioElementInterface.SetId()
	return TreeNode{drawioElementInterface, parent, []TreeNodeInterface{}, &Position{}}

}
func (tn *TreeNode) GetChildren() ([]TreeNodeInterface, []TreeNodeInterface) {
	return []TreeNodeInterface{}, []TreeNodeInterface{}
}

func (tn *TreeNode) SetPosition(position *Position) {
	tn.position = position
}
func (tn *TreeNode) GetPosition() *Position {
	return tn.position
}

///////////////////////////////////////////////////////////////////////

type NetworkTreeNode struct {
	TreeNode
	vpcs []TreeNodeInterface
}

func NewNetworkTreeNode() *NetworkTreeNode {
	return &NetworkTreeNode{NewTreeNode(&DrawioSquereElement{}, nil), []TreeNodeInterface{}}
}

func (tn *NetworkTreeNode) GetChildren() ([]TreeNodeInterface, []TreeNodeInterface) {
	return tn.vpcs, tn.elements
}
func (tn *NetworkTreeNode) GetType() string { return "pub" }

func (tn *NetworkTreeNode) SetDrawioInfo() {
	di := tn.DrawioElementInterface.(*DrawioSquereElement)
	position := tn.GetPosition()
	di.parentId = 0
	di.name = "network name"
	di.width = position.lastCol.width + position.lastCol.x - position.firstCol.x
	di.hight = position.lastRow.hight + position.lastRow.y - position.firstRow.y
	di.x = int(borderDistance)
	di.y = int(borderDistance)
}

// ////////////////////////////////////////////////////////////////////////////////////////
type VpcTreeNode struct {
	TreeNode
	zones []TreeNodeInterface
	sgs   []TreeNodeInterface
}

func NewVpcTreeNode(parent *NetworkTreeNode) *VpcTreeNode {
	vpc := VpcTreeNode{NewTreeNode(&DrawioSquereElement{}, parent), []TreeNodeInterface{}, []TreeNodeInterface{}}
	parent.vpcs = append(parent.vpcs, &vpc)
	return &vpc
}
func (tn *VpcTreeNode) GetChildren() ([]TreeNodeInterface, []TreeNodeInterface) {
	return tn.zones, append(tn.elements, tn.sgs...)
}
func (tn *VpcTreeNode) GetType() string { return "vpc" }

func (tn *VpcTreeNode) SetDrawioInfo() {
	di := tn.DrawioElementInterface.(*DrawioSquereElement)
	position := tn.GetPosition()
	parentPosition := tn.GetParent().GetPosition()
	di.parentId = tn.GetParent().GetDrawioElementInterface().GetId()
	di.name = "vpc name"
	di.width = position.lastCol.width + position.lastCol.x - position.firstCol.x
	di.hight = position.lastRow.hight + position.lastRow.y - position.firstRow.y
	di.x = position.firstCol.x - parentPosition.firstCol.x
	di.y = position.firstRow.y - parentPosition.firstRow.y
}

///////////////////////////////////////////////////////////////////////

type ZoneTreeNode struct {
	TreeNode
	subnets []TreeNodeInterface
}

func NewZoneTreeNode(parent *VpcTreeNode) *ZoneTreeNode {
	zone := ZoneTreeNode{NewTreeNode(&DrawioSquereElement{}, parent), []TreeNodeInterface{}}
	parent.zones = append(parent.zones, &zone)
	return &zone
}
func (tn *ZoneTreeNode) GetChildren() ([]TreeNodeInterface, []TreeNodeInterface) {
	return tn.subnets, tn.elements
}
func (tn *ZoneTreeNode) GetType() string { return "zone" }

func (tn *ZoneTreeNode) SetDrawioInfo() {
	di := tn.DrawioElementInterface.(*DrawioSquereElement)
	position := tn.GetPosition()
	parentPosition := tn.GetParent().GetPosition()
	di.parentId = tn.GetParent().GetDrawioElementInterface().GetId()
	di.name = "zone name"
	di.width = position.lastCol.width + position.lastCol.x - position.firstCol.x
	di.hight = position.lastRow.hight + position.lastRow.y - position.firstRow.y
	di.x = position.firstCol.x - parentPosition.firstCol.x
	di.y = position.firstRow.y - parentPosition.firstRow.y
}

///////////////////////////////////////////////////////////////////////

type SGTreeNode struct {
	TreeNode
}

func NewSGTreeNode(parent *VpcTreeNode) *SGTreeNode {
	sg := SGTreeNode{NewTreeNode(&DrawioSquereElement{}, parent)}
	parent.sgs = append(parent.sgs, &sg)
	return &sg
}
func (tn *SGTreeNode) GetType() string { return "sg" }

func (tn *SGTreeNode) SetDrawioInfo() {
	di := tn.DrawioElementInterface.(*DrawioSquereElement)
	di.parentId = tn.GetParent().GetDrawioElementInterface().GetId()
	di.name = "sg name"
	di.width = 400
	di.hight = 400
	di.x = 40
	di.y = 40
}

/////////////////////////////////////////////////////////////////////////

type SubnetTreeNode struct {
	TreeNode
}

func NewSubnetTreeNode(parent *ZoneTreeNode) *SubnetTreeNode {
	subnet := SubnetTreeNode{NewTreeNode(&DrawioSubnetElement{}, parent)}
	parent.subnets = append(parent.subnets, &subnet)
	return &subnet
}

func (tn *SubnetTreeNode) GetChildren() ([]TreeNodeInterface, []TreeNodeInterface) {
	return []TreeNodeInterface{}, tn.elements
}
func (tn *SubnetTreeNode) GetType() string { return "subnet" }

func (tn *SubnetTreeNode) SetDrawioInfo() {
	di := tn.DrawioElementInterface.(*DrawioSubnetElement)
	position := tn.GetPosition()
	parentPosition := tn.GetParent().GetPosition()
	di.parentId = tn.GetParent().GetDrawioElementInterface().GetId()
	di.name = "subnet name"
	di.width = position.lastCol.width + position.lastCol.x - position.firstCol.x
	di.hight = position.lastRow.hight + position.lastRow.y - position.firstRow.y
	di.x = position.firstCol.x - parentPosition.firstCol.x
	di.y = position.firstRow.y - parentPosition.firstRow.y
}

/////////////////////////////////////////////////////////////////////////

type IconTreeNode struct {
	TreeNode
	sg TreeNodeInterface
}

// ///////////////////////////////////////////
type NITreeNode struct {
	IconTreeNode
}

func NewNITreeNode(parent TreeNodeInterface, sg *SGTreeNode) *NITreeNode {
	ni := NITreeNode{IconTreeNode{NewTreeNode(&DrawioNetworkInterfaceElement{}, parent), sg}}
	parent.AddIconTreeNode(&ni)
	if sg != nil {
		sg.AddIconTreeNode(&ni)
	}
	return &ni
}
func (tn *NITreeNode) GetType() string { return "ni" }

func (tn *NITreeNode) SetDrawioInfo() {
	di := tn.DrawioElementInterface.(*DrawioNetworkInterfaceElement)
	position := tn.GetPosition()
	parentPosition := tn.GetParent().GetPosition()
	di.parentId = tn.GetParent().GetDrawioElementInterface().GetId()
	di.name = "icon name"
	di.x = position.firstCol.x - parentPosition.firstCol.x + position.firstCol.width/2 - iconSize/2 + position.x_offset
	di.y = position.firstRow.y - parentPosition.firstRow.y + position.firstRow.hight/2 - iconSize/2 + position.y_offset
}

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

func create_network() TreeNodeInterface {
	network := NewNetworkTreeNode()
	vpc1 := NewVpcTreeNode(network)
	vpc2 := NewVpcTreeNode(network)
	zone11 := NewZoneTreeNode(vpc1)
	zone12 := NewZoneTreeNode(vpc1)
	zone21 := NewZoneTreeNode(vpc2)
	zone22 := NewZoneTreeNode(vpc2)
	zone23 := NewZoneTreeNode(vpc2)

	sg11 := NewSGTreeNode(vpc1)
	sg12 := NewSGTreeNode(vpc1)
	sg21 := NewSGTreeNode(vpc2)
	sg22 := NewSGTreeNode(vpc2)

	subnet111 := NewSubnetTreeNode(zone11)
	subnet112 := NewSubnetTreeNode(zone11)

	subnet121 := NewSubnetTreeNode(zone12)

	subnet211 := NewSubnetTreeNode(zone21)

	subnet221 := NewSubnetTreeNode(zone22)
	subnet222 := NewSubnetTreeNode(zone22)
	subnet231 := NewSubnetTreeNode(zone23)

	NewNITreeNode(subnet111, sg11)
	NewNITreeNode(subnet111, sg12)
	NewNITreeNode(subnet112, sg12)
	NewNITreeNode(subnet121, sg11)

	NewNITreeNode(subnet211, sg21)
	NewNITreeNode(subnet211, sg21)
	NewNITreeNode(subnet211, sg21)

	NewNITreeNode(subnet221, sg22)
	NewNITreeNode(subnet222, sg22)
	NewNITreeNode(subnet222, sg22)
	NewNITreeNode(subnet222, sg22)
	NewNITreeNode(subnet222, sg22)

	NewNITreeNode(subnet231, sg22)
	NewNITreeNode(subnet231, sg22)
	NewNITreeNode(subnet231, sg22)
	NewNITreeNode(subnet231, sg22)
	NewNITreeNode(subnet231, sg22)
	NewNITreeNode(subnet231, sg22)
	NewNITreeNode(subnet231, sg22)
	NewNITreeNode(subnet231, sg22)
	NewNITreeNode(subnet231, sg22)
	NewNITreeNode(subnet231, sg22)
	NewNITreeNode(zone11, nil)
	NewNITreeNode(zone12, nil)
	NewNITreeNode(zone21, nil)
	NewNITreeNode(zone22, nil)

	NewNITreeNode(vpc1, nil)
	NewNITreeNode(vpc2, nil)
	NewNITreeNode(network, nil)
	NewNITreeNode(network, nil)
	NewNITreeNode(network, nil)
	NewNITreeNode(network, nil)
	return network
}

func layout(network TreeNodeInterface) {
	m := LayoutSubnets(network)
	AddBorders(network, m)
	SetLayersSize(network)
	SetLayersLocation(m)
	setIconsPositions(network, m)
	resolveDrawioInfo(network)
}

func LayoutSubnets(network TreeNodeInterface) *Matrix {
	m := NewMatrix(4, 5)
	network.(*NetworkTreeNode).vpcs[0].(*VpcTreeNode).zones[0].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[0].SetPosition(NewPosition(m.rows[0], m.rows[0], m.cols[0], m.cols[0]))
	network.(*NetworkTreeNode).vpcs[0].(*VpcTreeNode).zones[0].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[1].SetPosition(NewPosition(m.rows[0], m.rows[0], m.cols[0], m.cols[0]))

	network.(*NetworkTreeNode).vpcs[0].(*VpcTreeNode).zones[0].(*ZoneTreeNode).subnets[1].(*SubnetTreeNode).elements[0].SetPosition(NewPosition(m.rows[1], m.rows[1], m.cols[0], m.cols[0]))

	network.(*NetworkTreeNode).vpcs[0].(*VpcTreeNode).zones[1].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[0].SetPosition(NewPosition(m.rows[0], m.rows[0], m.cols[1], m.cols[1]))

	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[0].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[0].SetPosition(NewPosition(m.rows[0], m.rows[0], m.cols[2], m.cols[2]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[0].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[1].SetPosition(NewPosition(m.rows[0], m.rows[0], m.cols[2], m.cols[2]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[0].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[2].SetPosition(NewPosition(m.rows[1], m.rows[1], m.cols[2], m.cols[2]))

	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[1].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[0].SetPosition(NewPosition(m.rows[0], m.rows[0], m.cols[3], m.cols[3]))

	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[1].(*ZoneTreeNode).subnets[1].(*SubnetTreeNode).elements[0].SetPosition(NewPosition(m.rows[1], m.rows[1], m.cols[3], m.cols[3]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[1].(*ZoneTreeNode).subnets[1].(*SubnetTreeNode).elements[1].SetPosition(NewPosition(m.rows[1], m.rows[1], m.cols[3], m.cols[3]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[1].(*ZoneTreeNode).subnets[1].(*SubnetTreeNode).elements[2].SetPosition(NewPosition(m.rows[1], m.rows[1], m.cols[4], m.cols[4]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[1].(*ZoneTreeNode).subnets[1].(*SubnetTreeNode).elements[3].SetPosition(NewPosition(m.rows[1], m.rows[1], m.cols[4], m.cols[4]))

	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[2].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[0].SetPosition(NewPosition(m.rows[2], m.rows[2], m.cols[2], m.cols[2]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[2].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[1].SetPosition(NewPosition(m.rows[2], m.rows[2], m.cols[2], m.cols[2]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[2].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[2].SetPosition(NewPosition(m.rows[2], m.rows[2], m.cols[3], m.cols[3]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[2].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[3].SetPosition(NewPosition(m.rows[2], m.rows[2], m.cols[4], m.cols[4]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[2].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[4].SetPosition(NewPosition(m.rows[3], m.rows[3], m.cols[2], m.cols[2]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[2].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[5].SetPosition(NewPosition(m.rows[3], m.rows[3], m.cols[2], m.cols[2]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[2].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[6].SetPosition(NewPosition(m.rows[3], m.rows[3], m.cols[3], m.cols[3]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[2].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[7].SetPosition(NewPosition(m.rows[3], m.rows[3], m.cols[3], m.cols[3]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[2].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[8].SetPosition(NewPosition(m.rows[3], m.rows[3], m.cols[4], m.cols[4]))
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[2].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).elements[9].SetPosition(NewPosition(m.rows[3], m.rows[3], m.cols[4], m.cols[4]))
	return m
}

func AddBorders(network TreeNodeInterface, m *Matrix) {
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

	m.rows = []*Row{}
	m.cols = []*Col{}
	for _, row := range new_m.rows {
		if row.inUse {
			m.rows = append(m.rows, row)
		}
	}
	for _, col := range new_m.cols {
		if col.inUse {
			m.cols = append(m.cols, col)
		}
	}
	network.SetPosition(NewPosition(m.rows[0], m.rows[len(m.rows)-1], m.cols[0], m.cols[len(m.cols)-1]))
}

func SetLayersSize(network TreeNodeInterface) {
	if len(network.GetIconTreeNodes()) > 0 {
		network.GetPosition().firstCol.width = iconSpace
	}
	for _, vpc := range network.(*NetworkTreeNode).vpcs {
		if len(vpc.GetIconTreeNodes()) > 0 {
			vpc.GetPosition().firstRow.hight = iconSpace
		}
		for _, zone := range vpc.(*VpcTreeNode).zones {
			if len(zone.GetIconTreeNodes()) > 0 {
				zone.GetPosition().firstCol.width = iconSpace
			}
		}
	}
}

func setIconsPositions(network TreeNodeInterface, m *Matrix) {
	icons := network.GetIconTreeNodes()
	iconIndex := 0
	for ri, row := range m.rows {
		if iconIndex >= len(icons) {
			break
		}
		if row.hight >= iconSpace || (ri > 0 && m.rows[ri-1].hight >= iconSpace) {
			icons[iconIndex].SetPosition(NewPosition(row, row, m.cols[0], m.cols[0]))
			iconIndex++
		}
	}
	for _, vpc := range network.(*NetworkTreeNode).vpcs {
		icons := vpc.GetIconTreeNodes()
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
			iconIndex := 0
			indexes := m.getPosiotonIndexes(zone.GetPosition())
			for ri := indexes.firstRowIndex; ri <= indexes.lastRowIndex && iconIndex < len(icons); ri++ {
				if m.rows[ri].hight >= iconSpace || (ri > indexes.firstRowIndex && m.rows[ri-1].hight >= iconSpace) {
					icons[iconIndex].SetPosition(NewPosition(m.rows[ri], m.rows[ri], m.cols[indexes.firstColIndex], m.cols[indexes.firstColIndex]))
					iconIndex++
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

func GetElements(tn TreeNodeInterface) []TreeNodeInterface {
	subtrees, leafs := tn.GetChildren()
	ret := append(leafs, tn)
	for _, element := range subtrees {
		sub := GetElements(element)
		ret = append(ret, sub...)
	}
	return ret
}

func main() {
	network := create_network()
	layout(network)

	elements := GetElements(network)

	templateFuncMap := template.FuncMap{
		"add": func(x uint, y uint) uint {
			return x + y
		},
	}
	tmpl, err := template.New("drawio_template.txt").Funcs(templateFuncMap).ParseFiles("drawio_template.txt")
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
