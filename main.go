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
	x        uint
	y        uint
	name     string
}

func (di DrawioLocatedElement) GetX() uint {
	return di.x
}
func (di DrawioLocatedElement) GetY() uint {
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
	width uint
	hight uint
}

func NewDrawioSquereElement(etype string, x uint, y uint, w uint, h uint, id uint, pid uint, name string) *DrawioSquereElement {
	return &DrawioSquereElement{DrawioLocatedElement{DrawioElement{etype, id}, pid, x, y, name}, w, h}
}
func (di DrawioSquereElement) GetHight() uint {
	return di.hight
}
func (di DrawioSquereElement) GetWidth() uint {
	return di.width
}

// //////////////////////////////////////////////////////////////////////////////////////
type DrawioSubnetElement struct {
	DrawioSquereElement
	ip  string
	key string
}

func NewDrawioSubnetElement(etype string, x uint, y uint, w uint, h uint, id uint, pid uint, name string, ip string, key string) *DrawioSubnetElement {
	return &DrawioSubnetElement{DrawioSquereElement{DrawioLocatedElement{DrawioElement{etype, id}, pid, x, y, name}, w, h}, ip, key}
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

func NewDrawioIconElement(etype string, x uint, y uint, w uint, h uint, id uint, pid uint, name string) *DrawioIconElement {
	return &DrawioIconElement{DrawioLocatedElement{DrawioElement{etype, id}, pid, x, y, name}}
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

func NewDrawioNetworkInterfaceElement(etype string, x uint, y uint, w uint, h uint, id uint, pid uint, name string, fip string, svi string) *DrawioNetworkInterfaceElement {
	return &DrawioNetworkInterfaceElement{DrawioIconElement{DrawioLocatedElement{DrawioElement{etype, id}, pid, x, y, name}}, fip, svi}
}

type DrawioConnectElementPoint struct {
	x uint
	y uint
}

func (di DrawioConnectElementPoint) GetX() uint {
	return di.x
}
func (di DrawioConnectElementPoint) GetY() uint {
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
	GetChildren() ([]TreeNodeInterface, []TreeNodeInterface)
	GetType() string
	GetDrawioElementInterface() DrawioElementInterface
	SetPosition(position Position)
	GetPosition() Position
	SetDrawioInfo()
}

///////////////////////////////////////////////////////////////////////////

type TreeNode struct {
	DrawioElementInterface
	parent   TreeNodeInterface
	elements []TreeNodeInterface
	position Position
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

func NewTreeNode(drawioElementInterface DrawioElementInterface, parent TreeNodeInterface) TreeNode {
	drawioElementInterface.SetId()
	return TreeNode{drawioElementInterface, parent, []TreeNodeInterface{}, Position{}}

}
func (tn *TreeNode) GetChildren() ([]TreeNodeInterface, []TreeNodeInterface) {
	return []TreeNodeInterface{}, []TreeNodeInterface{}
}

func (tn *TreeNode) SetPosition(position Position) {
	tn.position = position
}
func (tn *TreeNode) GetPosition() Position {
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
	di.x = borderDistance
	di.y = borderDistance
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
	di.parentId = tn.GetParent().GetDrawioElementInterface().GetId()
	di.name = "icon name"
	di.x = 80
	di.y = 80
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
	borderDistance uint = 40
	subnetWidth    uint = 8 * 40
	subnetHight    uint = 6 * 40
	iconSize       uint = 60
	iconSpace      uint = 4 * 40
)

type Row struct {
	hight  uint
	y      uint
	inUse  bool
	matrix *Matrix
}

func NewRow(matrix *Matrix, hight uint) *Row {
	r := &Row{}
	r.hight = hight
	r.matrix = matrix
	r.inUse = false
	return r
}

type Col struct {
	width  uint
	x      uint
	inUse  bool
	matrix *Matrix
}

func NewCol(matrix *Matrix, width uint) *Col {
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

func (matrix *Matrix) getPosiotonIndexes(position Position) PositionIndexes {
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

func (matrix *Matrix) getPosioton(positionIndexes PositionIndexes) Position {
	return Position{
		matrix.rows[positionIndexes.firstRowIndex],
		matrix.rows[positionIndexes.lastRowIndex],
		matrix.cols[positionIndexes.firstColIndex],
		matrix.cols[positionIndexes.lastColIndex],
	}
}

func layout(network TreeNodeInterface) {

	m := NewMatrix(4, 5)
	network.(*NetworkTreeNode).vpcs[0].(*VpcTreeNode).zones[0].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).SetPosition(Position{m.rows[0], m.rows[0], m.cols[0], m.cols[0]})
	network.(*NetworkTreeNode).vpcs[0].(*VpcTreeNode).zones[0].(*ZoneTreeNode).subnets[1].(*SubnetTreeNode).SetPosition(Position{m.rows[1], m.rows[1], m.cols[0], m.cols[0]})

	network.(*NetworkTreeNode).vpcs[0].(*VpcTreeNode).zones[1].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).SetPosition(Position{m.rows[0], m.rows[0], m.cols[1], m.cols[1]})

	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[0].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).SetPosition(Position{m.rows[0], m.rows[1], m.cols[2], m.cols[2]})

	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[1].(*ZoneTreeNode).subnets[1].(*SubnetTreeNode).SetPosition(Position{m.rows[0], m.rows[0], m.cols[3], m.cols[3]})
	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[1].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).SetPosition(Position{m.rows[1], m.rows[1], m.cols[3], m.cols[4]})

	network.(*NetworkTreeNode).vpcs[1].(*VpcTreeNode).zones[2].(*ZoneTreeNode).subnets[0].(*SubnetTreeNode).SetPosition(Position{m.rows[2], m.rows[3], m.cols[2], m.cols[4]})

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
				subnetIndexes := m.getPosiotonIndexes(subnet.GetPosition())
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
	network.SetPosition(Position{m.rows[0], m.rows[len(m.rows)-1], m.cols[0], m.cols[len(m.cols)-1]})

	var y uint = 0
	for _, row := range m.rows {
		row.y = y
		y = y + row.hight
	}
	var x uint = 0
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

func resolveDrawioInfo(network TreeNodeInterface) []TreeNodeInterface {
	layout(network)
	elements := GetElements(network)
	for _, element := range elements {
		element.SetDrawioInfo()
	}
	return elements
}

//	type DrawioInfo struct {
//		IconSize uint
//		Elements []DrawioElementInterface
//	}
type DrawioInfo struct {
	IconSize uint
	Nodes    []TreeNodeInterface
}

func main() {
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
	NewNITreeNode(subnet221, sg22)
	NewNITreeNode(subnet222, sg22)
	NewNITreeNode(subnet231, sg22)
	NewNITreeNode(zone11, nil)
	NewNITreeNode(zone12, nil)
	NewNITreeNode(zone21, nil)
	NewNITreeNode(zone22, nil)

	NewNITreeNode(vpc1, nil)
	NewNITreeNode(vpc2, nil)
	NewNITreeNode(network, nil)
	NewNITreeNode(network, nil)

	//elements := []TreeNodeInterface{network}
	elements := resolveDrawioInfo(network)

	// elements := []DrawioElementInterface{
	// 	NewDrawioSquereElement("pub", 40, 40, 1920, 1040, 895, 0, "noname"),
	// 	NewDrawioIconElement("internet", 50.0, 50.0, 60, 60, 900, 895, "161.26.0.0/16"),
	// 	NewDrawioIconElement("internet", 50.0, 210.0, 60, 60, 905, 895, "166.8.0.0/14"),
	// 	NewDrawioIconElement("internet", 50.0, 370.0, 60, 60, 910, 895, "192.0.1.0/24"),
	// 	NewDrawioSquereElement("vpc", 160, 40, 1720, 960, 1230, 895, "noname"),
	// 	NewDrawioSquereElement("zone", 40, 40, 520, 880, 1235, 1230, "us-south-3"),
	// 	NewDrawioSubnetElement("subnet", 160, 600, 320, 240, 1240, 1235, "ky-testenv-transit-subnet-3", "0000000", "key"),
	// 	NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1245, 1240, "transit-2-instance-ky", "", ""),
	// 	NewDrawioSubnetElement("subnet", 160, 40, 320, 240, 1250, 1235, "ky-testenv-edge-subnet-3", "0000000", "key"),
	// 	NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1255, 1250, "edge-2-instance-ky", "", "vsi_name"),
	// 	NewDrawioSubnetElement("subnet", 160, 320, 320, 240, 1260, 1235, "ky-testenv-private-subnet-3", "0000000", "key"),
	// 	NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1265, 1260, "private-2-instance-ky", "", "vsi_name"),
	// 	NewDrawioIconElement("gateway", 50.0, 130.0, 60, 60, 1270, 1235, "ky-testenv-gateway-3"),
	// 	NewDrawioSquereElement("zone", 600, 40, 520, 880, 1275, 1230, "us-south-2"),
	// 	NewDrawioSubnetElement("subnet", 160, 600, 320, 240, 1280, 1275, "ky-testenv-transit-subnet-2", "0000000", "key"),
	// 	NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1285, 1280, "transit-1-instance-ky", "fip_name", "vsi_name"),
	// 	NewDrawioSubnetElement("subnet", 160, 40, 320, 240, 1290, 1275, "ky-testenv-edge-subnet-2", "0000000", "key"),
	// 	NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1295, 1290, "edge-1-instance-ky", "", "vsi_name"),
	// 	NewDrawioSubnetElement("subnet", 160, 320, 320, 240, 1300, 1275, "ky-testenv-private-subnet-2", "0000000", "key"),
	// 	NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1305, 1300, "private-1-instance-ky", "", "vsi_name"),
	// 	NewDrawioIconElement("gateway", 50.0, 410.0, 60, 60, 1310, 1275, "ky-testenv-gateway-2"),
	// 	NewDrawioSquereElement("zone", 1160, 40, 520, 880, 1315, 1230, "us-south-1"),
	// 	NewDrawioSubnetElement("subnet", 160, 320, 320, 240, 1320, 1315, "ky-testenv-private-subnet-1", "0000000", "key"),
	// 	NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1325, 1320, "private-0-instance-ky", "fip_name", "vsi_name"),
	// 	NewDrawioSubnetElement("subnet", 160, 40, 320, 240, 1330, 1315, "ky-testenv-edge-subnet-1", "0000000", "key"),
	// 	NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1335, 1330, "edge-0-instance-ky", "", "vsi_name"),
	// 	NewDrawioSubnetElement("subnet", 160, 600, 320, 240, 1340, 1315, "ky-testenv-transit-subnet-1", "0000000", "key"),
	// 	NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1345, 1340, "transit-0-instance-ky", "", "vsi_name"),
	// 	NewDrawioIconElement("gateway", 50.0, 690.0, 60, 60, 1350, 1315, "ky-testenv-gateway-1"),
	// 	NewDrawioSquereElement("sg", 290.0, 120, 1260.0, 160, 1355, 1230, "ky-testenv-default-sg"),
	// 	NewDrawioSquereElement("sg", 290.0, 400, 1260.0, 160, 1360, 1230, "ky-testenv-default-sg"),
	// 	NewDrawioSquereElement("sg", 290.0, 680, 1260.0, 160, 1365, 1230, "ky-testenv-default-sg"),
	// 	NewDrawioConnectElement("diredge", 915, 1325, 1305, "TCP 1-65535", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 920, 1325, 1345, "TCP 1-65535", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 925, 1325, 1245, "TCP 1-65535", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 930, 1345, 1255, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 190.0}}),
	// 	NewDrawioConnectElement("diredge", 935, 1345, 1295, "TCP 443", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 940, 1345, 1325, "TCP 443", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 945, 1345, 1335, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{1590.0, 530.0}}),
	// 	NewDrawioConnectElement("diredge", 950, 1345, 1265, "TCP 443", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 955, 1345, 1285, "TCP 443", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 960, 1265, 1245, "TCP 1-65535", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 965, 1265, 1305, "TCP 1-65535", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 970, 1265, 1345, "TCP 1-65535", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 975, 1285, 1345, "TCP 1-65535", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 980, 1285, 1245, "TCP 1-65535", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 985, 1285, 1305, "TCP 1-65535", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 990, 1305, 1335, "TCP 443", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 995, 1305, 1325, "TCP 443", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 1000, 1305, 1285, "TCP 443", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 1005, 1305, 1295, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{1030.0, 530.0}}),
	// 	NewDrawioConnectElement("diredge", 1010, 1305, 1255, "TCP 443", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 1015, 1305, 1265, "TCP 443", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 1020, 1255, 900, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{350.0, 280.0}, DrawioConnectElementPoint{290.0, 280.0}}),
	// 	NewDrawioConnectElement("diredge", 1025, 1255, 1305, "TCP 1-65535", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 1030, 1255, 1345, "TCP 1-65535", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 190.0}}),
	// 	NewDrawioConnectElement("diredge", 1035, 1255, 910, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{350.0, 276.0}, DrawioConnectElementPoint{290.0, 276.0}}),
	// 	NewDrawioConnectElement("diredge", 1040, 1255, 905, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{350.0, 272.0}, DrawioConnectElementPoint{290.0, 272.0}, DrawioConnectElementPoint{310.0, 190.0}}),
	// 	NewDrawioConnectElement("diredge", 1045, 1255, 1245, "TCP 1-65535", []DrawioConnectElementPoint{DrawioConnectElementPoint{470.0, 530.0}}),
	// 	NewDrawioConnectElement("diredge", 1050, 1335, 1245, "TCP 1-65535", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 190.0}}),
	// 	NewDrawioConnectElement("diredge", 1055, 1335, 1345, "TCP 1-65535", []DrawioConnectElementPoint{DrawioConnectElementPoint{1590.0, 530.0}}),
	// 	NewDrawioConnectElement("diredge", 1060, 1335, 910, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{1470.0, 840.0}, DrawioConnectElementPoint{1410.0, 840.0}}),
	// 	NewDrawioConnectElement("diredge", 1065, 1335, 905, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{1470.0, 836.0}, DrawioConnectElementPoint{1410.0, 836.0}, DrawioConnectElementPoint{870.0, 190.0}}),
	// 	NewDrawioConnectElement("diredge", 1070, 1335, 1305, "TCP 1-65535", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 1075, 1335, 900, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{1470.0, 832.0}, DrawioConnectElementPoint{1410.0, 832.0}}),
	// 	NewDrawioConnectElement("diredge", 1080, 1245, 1295, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{470.0, 530.0}}),
	// 	NewDrawioConnectElement("diredge", 1085, 1245, 1255, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{470.0, 530.0}}),
	// 	NewDrawioConnectElement("diredge", 1090, 1245, 1265, "TCP 443", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 1095, 1245, 1285, "TCP 443", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 1100, 1245, 1325, "TCP 443", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 1105, 1245, 1335, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 190.0}}),
	// 	NewDrawioConnectElement("diredge", 1110, 1295, 1345, "TCP 1-65535", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("diredge", 1115, 1295, 905, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{910.0, 560.0}, DrawioConnectElementPoint{850.0, 560.0}, DrawioConnectElementPoint{590.0, 190.0}}),
	// 	NewDrawioConnectElement("diredge", 1120, 1295, 910, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{910.0, 556.0}, DrawioConnectElementPoint{850.0, 556.0}}),
	// 	NewDrawioConnectElement("diredge", 1125, 1295, 1245, "TCP 1-65535", []DrawioConnectElementPoint{DrawioConnectElementPoint{470.0, 530.0}}),
	// 	NewDrawioConnectElement("diredge", 1130, 1295, 1305, "TCP 1-65535", []DrawioConnectElementPoint{DrawioConnectElementPoint{1030.0, 530.0}}),
	// 	NewDrawioConnectElement("diredge", 1135, 1295, 900, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{910.0, 552.0}, DrawioConnectElementPoint{850.0, 552.0}}),
	// 	NewDrawioConnectElement("undiredge", 1140, 1285, 1265, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{810.0, 470.0}}),
	// 	NewDrawioConnectElement("undiredge", 1145, 1325, 1265, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 470.0}}),
	// 	NewDrawioConnectElement("undiredge", 1150, 1325, 1285, "", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("undiredge", 1155, 1345, 1305, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{1370.0, 750.0}}),
	// 	NewDrawioConnectElement("undiredge", 1160, 1335, 1255, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 190.0}}),
	// 	NewDrawioConnectElement("undiredge", 1165, 1285, 1255, "", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("undiredge", 1170, 1325, 1255, "", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("undiredge", 1175, 1295, 1255, "", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("undiredge", 1180, 1265, 1255, "", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("undiredge", 1185, 1265, 1335, "", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("undiredge", 1190, 1325, 1335, "", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("undiredge", 1195, 1285, 1335, "", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("undiredge", 1200, 1345, 1245, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 750.0}}),
	// 	NewDrawioConnectElement("undiredge", 1205, 1305, 1245, "TCP 443", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("undiredge", 1210, 1285, 1295, "", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("undiredge", 1215, 1325, 1295, "", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("undiredge", 1220, 1265, 1295, "", []DrawioConnectElementPoint{}),
	// 	NewDrawioConnectElement("undiredge", 1225, 1335, 1295, "", []DrawioConnectElementPoint{}),
	//}
	templateFuncMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
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

	err = tmpl.Execute(fo, DrawioInfo{60, elements})
	if err != nil {
		panic(err)
	}
	// err = tmpl.Execute(fo, DrawioInfo{60, elements})
	// if err != nil {
	// 	panic(err)
	// }

}
