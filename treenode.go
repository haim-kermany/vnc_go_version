package main

type DrawioElementInterface interface {
	GetId() uint
	GetId1() uint
	GetId2() uint
	SetId()
}

type DrawioElement struct {
	parentId uint
	deType   string
	id       uint
}

func (di *DrawioElement) GetId() uint {
	return di.id
}
func (di *DrawioElement) GetId1() uint {
	return di.id + 1
}
func (di *DrawioElement) GetId2() uint {
	return di.id + 2
}
func (di *DrawioElement) GetParentId() uint {
	return di.parentId
}

var id uint = 100

func (di *DrawioElement) SetId() {
	di.id = id
	id += 5
}

// //////////////////////////////////////////////////////////////////////////////
type DrawioLocatedElement struct {
	DrawioElement
	x    int
	y    int
	name string
}

func (di *DrawioLocatedElement) GetX() int {
	return di.x
}
func (di *DrawioLocatedElement) GetY() int {
	return di.y
}
func (di *DrawioLocatedElement) GetName() string {
	return di.name
}

// //////////////////////////////////////////////////////////////////////////////
type DrawioSquereElement struct {
	DrawioLocatedElement
	width int
	hight int
	ip    string
	key   string
}

func (di *DrawioSquereElement) GetHight() int {
	return di.hight
}
func (di *DrawioSquereElement) GetWidth() int {
	return di.width
}
func (di *DrawioSquereElement) GetIP() string {
	return di.ip
}
func (di *DrawioSquereElement) GetKey() string {
	return di.key
}

// ///////////////////////////////////////////////////////////////////////////////////////////////
type DrawioIconElement struct {
	DrawioLocatedElement
	floating_ip string
	svi         string
}

func (di *DrawioIconElement) GetFip() string {
	return di.floating_ip
}

func (di *DrawioIconElement) GetSvi() string {
	return di.svi
}

// /////////////////////////////////////////////////////////////////////////
type DrawioConnectElementPoint struct {
	x int
	y int
}

func (di *DrawioConnectElementPoint) GetX() int {
	return di.x
}
func (di *DrawioConnectElementPoint) GetY() int {
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

func (di *DrawioConnectElement) GetSrcId() uint {
	return di.srcId
}
func (di *DrawioConnectElement) GetDstId() uint {
	return di.dstId
}
func (di *DrawioConnectElement) GetLabel() string {
	return di.label
}

func (di *DrawioConnectElement) GetPoints() []DrawioConnectElementPoint {
	return di.points
}
func (di *DrawioConnectElement) AddPoint(x int, y int) {
	di.points = append(di.points, DrawioConnectElementPoint{x, y})
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
	SetParent(TreeNodeInterface)
	AddIconTreeNode(icon TreeNodeInterface)
	GetIconTreeNodes() []TreeNodeInterface
	GetChildren() ([]TreeNodeInterface, []TreeNodeInterface)
	GetType() string
	IsLine() bool
	GetName() string
	GetDrawioElementInterface() DrawioElementInterface
	SetPosition(position *Position)
	GetPosition() *Position
	SetDrawioInfo()
	AllocPassenger() int
}

///////////////////////////////////////////////////////////////////////////

type TreeNode struct {
	DrawioElementInterface
	parent   TreeNodeInterface
	elements []TreeNodeInterface
	position *Position
}

func (tn *TreeNode) AllocPassenger() int {
	return 0
}
func (tn *TreeNode) GetName() string { return "name?" }

func (tn *TreeNode) GetDrawioElementInterface() DrawioElementInterface {
	return tn.DrawioElementInterface
}

func (tn *TreeNode) GetParent() TreeNodeInterface {
	return tn.parent
}
func (tn *TreeNode) SetParent(p TreeNodeInterface) {
	tn.parent = p
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
func (tn *TreeNode) IsLine() bool {
	return false
}

func (tn *TreeNode) GetType() string {
	return "No type"
}

///////////////////////////////////////////////////////////////////////

type SquereTreeNode struct {
	TreeNode
}

func NewSquereTreeNode(parent TreeNodeInterface) SquereTreeNode {
	return SquereTreeNode{NewTreeNode(&DrawioSquereElement{}, parent)}
}

func (tn *SquereTreeNode) SetDrawioInfo() {
	di := tn.DrawioElementInterface.(*DrawioSquereElement)
	position := tn.GetPosition()
	parentPosition := position
	if tn.GetParent() != nil {
		parentPosition = tn.GetParent().GetPosition()
		di.parentId = tn.GetParent().GetDrawioElementInterface().GetId()
	}
	di.name = tn.GetName()
	di.width = position.lastCol.width + position.lastCol.x - position.firstCol.x
	di.hight = position.lastRow.hight + position.lastRow.y - position.firstRow.y
	di.x = position.firstCol.x - parentPosition.firstCol.x
	di.y = position.firstRow.y - parentPosition.firstRow.y
}

//////////////////////////////////////////////////////////////////////////////

type NetworkTreeNode struct {
	SquereTreeNode
	vpcs        []TreeNodeInterface
	connections []TreeNodeInterface
}

func NewNetworkTreeNode() *NetworkTreeNode {
	return &NetworkTreeNode{NewSquereTreeNode(nil), []TreeNodeInterface{}, []TreeNodeInterface{}}
}

func (tn *NetworkTreeNode) GetChildren() ([]TreeNodeInterface, []TreeNodeInterface) {
	return tn.vpcs, append(tn.elements, tn.connections...)
}
func (tn *NetworkTreeNode) GetType() string { return "pub" }

// ////////////////////////////////////////////////////////////////////////////////////////
type VpcTreeNode struct {
	SquereTreeNode
	zones []TreeNodeInterface
	sgs   []TreeNodeInterface
}

func NewVpcTreeNode(parent *NetworkTreeNode) *VpcTreeNode {
	vpc := VpcTreeNode{NewSquereTreeNode(parent), []TreeNodeInterface{}, []TreeNodeInterface{}}
	parent.vpcs = append(parent.vpcs, &vpc)
	return &vpc
}
func (tn *VpcTreeNode) GetChildren() ([]TreeNodeInterface, []TreeNodeInterface) {
	return tn.zones, append(tn.elements, tn.sgs...)
}
func (tn *VpcTreeNode) GetType() string { return "vpc" }

///////////////////////////////////////////////////////////////////////

type ZoneTreeNode struct {
	SquereTreeNode
	subnets []TreeNodeInterface
}

func NewZoneTreeNode(parent *VpcTreeNode) *ZoneTreeNode {
	zone := ZoneTreeNode{NewSquereTreeNode(parent), []TreeNodeInterface{}}
	parent.zones = append(parent.zones, &zone)
	return &zone
}
func (tn *ZoneTreeNode) GetChildren() ([]TreeNodeInterface, []TreeNodeInterface) {
	return tn.subnets, tn.elements
}
func (tn *ZoneTreeNode) GetType() string { return "zone" }

///////////////////////////////////////////////////////////////////////

type SGTreeNode struct {
	SquereTreeNode
}

func NewSGTreeNode(parent *VpcTreeNode) *SGTreeNode {
	sg := SGTreeNode{NewSquereTreeNode(parent)}
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
	SquereTreeNode
}

func NewSubnetTreeNode(parent *ZoneTreeNode) *SubnetTreeNode {
	subnet := SubnetTreeNode{NewSquereTreeNode(parent)}
	parent.subnets = append(parent.subnets, &subnet)
	return &subnet
}

func (tn *SubnetTreeNode) GetChildren() ([]TreeNodeInterface, []TreeNodeInterface) {
	return []TreeNodeInterface{}, tn.elements
}
func (tn *SubnetTreeNode) GetType() string { return "subnet" }

/////////////////////////////////////////////////////////////////////////

type IconTreeNode struct {
	TreeNode
	nPassengers int
}

func NewIconTreeNode(parent TreeNodeInterface) IconTreeNode {
	return IconTreeNode{NewTreeNode(&DrawioIconElement{}, parent), 0}
}

func (tn *IconTreeNode) AllocPassenger() int {
	n := tn.nPassengers
	tn.nPassengers++
	return []int{0, 10, -10, 20, -20, 15, -15, 25, -25}[n]
}
func (tn *IconTreeNode) SetDrawioInfo() {
	di := tn.DrawioElementInterface.(*DrawioIconElement)
	position := tn.GetPosition()
	parentPosition := tn.GetParent().GetPosition()
	di.parentId = tn.GetParent().GetDrawioElementInterface().GetId()
	di.name = "icon name"
	di.x = position.firstCol.x - parentPosition.firstCol.x + position.firstCol.width/2 - iconSize/2 + position.x_offset
	di.y = position.firstRow.y - parentPosition.firstRow.y + position.firstRow.hight/2 - iconSize/2 + position.y_offset
}

// ///////////////////////////////////////////
type NITreeNode struct {
	IconTreeNode
	sg     TreeNodeInterface
	hasVsi bool
}

func NewNITreeNode(parent TreeNodeInterface, sg *SGTreeNode) *NITreeNode {
	ni := NITreeNode{NewIconTreeNode(parent), sg, true}
	parent.AddIconTreeNode(&ni)
	if sg != nil {
		sg.AddIconTreeNode(&ni)
	}
	return &ni
}
func (tn *NITreeNode) GetType() string { return "ni" }
func (tn *NITreeNode) HasVsi() bool    { return tn.hasVsi }

// ///////////////////////////////////////////
type GetWayTreeNode struct {
	IconTreeNode
}

func NewGetWayTreeNode(parent TreeNodeInterface) *GetWayTreeNode {
	gw := GetWayTreeNode{NewIconTreeNode(parent)}
	parent.AddIconTreeNode(&gw)
	return &gw
}
func (tn *GetWayTreeNode) GetType() string { return "gateway" }

// ///////////////////////////////////////////
type VsiTreeNode struct {
	IconTreeNode
	nis []TreeNodeInterface
}

func NewVsiTreeNode(parent TreeNodeInterface) *VsiTreeNode {
	vsi := VsiTreeNode{NewIconTreeNode(parent), []TreeNodeInterface{}}
	parent.AddIconTreeNode(&vsi)
	return &vsi
}
func (tn *VsiTreeNode) AddNI(ni TreeNodeInterface) {
	tn.nis = append(tn.nis, ni)
}
func (tn *VsiTreeNode) GetVsiSubnets() *map[TreeNodeInterface]bool {
	vsiSubents := map[TreeNodeInterface]bool{}
	for _, ni := range tn.nis {
		vsiSubents[ni.GetParent()] = true
	}
	return &vsiSubents
}

func (tn *VsiTreeNode) GetType() string { return "vsi" }

// ///////////////////////////////////////////
type InternetTreeNode struct {
	IconTreeNode
}

func NewInternetTreeNode(parent TreeNodeInterface) *InternetTreeNode {
	inter := InternetTreeNode{NewIconTreeNode(parent)}
	parent.AddIconTreeNode(&inter)
	return &inter
}
func (tn *InternetTreeNode) GetType() string { return "internet" }

// ////////////////////////////////////////////////////////////////
type InternetSeviceTreeNode struct {
	IconTreeNode
}

func NewInternetSeviceTreeNode(parent TreeNodeInterface) *InternetSeviceTreeNode {
	inter := InternetSeviceTreeNode{NewIconTreeNode(parent)}
	parent.AddIconTreeNode(&inter)
	return &inter
}
func (tn *InternetSeviceTreeNode) GetType() string { return "internet_service" }

// ////////////////////////////////////////////////////////////////
type LineTreeNode struct {
	TreeNode
	src TreeNodeInterface
	dst TreeNodeInterface
}

func (tn *LineTreeNode) IsLine() bool {
	return true
}

func (tn *LineTreeNode) SetDrawioInfo() {
	di := tn.DrawioElementInterface.(*DrawioConnectElement)
	di.srcId = tn.src.GetDrawioElementInterface().GetId()
	di.dstId = tn.dst.GetDrawioElementInterface().GetId()
	if tn.GetParent().GetType() == "ni" {
		di.parentId = tn.GetParent().GetDrawioElementInterface().GetId2()
	} else {
		di.parentId = tn.GetParent().GetDrawioElementInterface().GetId()
	}
}

// ////////////////////////////////////////////////////////////////
type VsiLineTreeNode struct {
	LineTreeNode
}

func NewVsiLineTreeNode(network TreeNodeInterface, vsi TreeNodeInterface, ni TreeNodeInterface) *VsiLineTreeNode {
	vsi.(*VsiTreeNode).AddNI(ni)
	conn := VsiLineTreeNode{LineTreeNode{NewTreeNode(&DrawioConnectElement{}, network), vsi, ni}}
	network.(*NetworkTreeNode).connections = append(network.(*NetworkTreeNode).connections, &conn)
	return &conn
}
func (tn *VsiLineTreeNode) GetType() string { return "vsi_connector" }

// ////////////////////////////////////////////////////////////////
type ConnectivityTreeNode struct {
	LineTreeNode
	directed bool
}

func NewConnectivityLineTreeNode(network TreeNodeInterface, src TreeNodeInterface, dst TreeNodeInterface, directed bool) *ConnectivityTreeNode {
	conn := ConnectivityTreeNode{LineTreeNode{NewTreeNode(&DrawioConnectElement{}, network), src, dst}, directed}
	network.(*NetworkTreeNode).connections = append(network.(*NetworkTreeNode).connections, &conn)
	return &conn
}
func (tn *ConnectivityTreeNode) SetPassage(passage TreeNodeInterface, reverse bool) {
	tn.SetParent(passage)
	passengerNumber := passage.AllocPassenger()
	if !reverse {
		tn.DrawioElementInterface.(*DrawioConnectElement).AddPoint(iconSize, iconSize/2+passengerNumber)
		tn.DrawioElementInterface.(*DrawioConnectElement).AddPoint(0, iconSize/2+passengerNumber)
	} else {
		tn.DrawioElementInterface.(*DrawioConnectElement).AddPoint(0, iconSize/2+passengerNumber)
		tn.DrawioElementInterface.(*DrawioConnectElement).AddPoint(iconSize, iconSize/2+passengerNumber)
	}

}

func (tn *ConnectivityTreeNode) GetType() string {
	if tn.directed {
		return "diredge"
	} else {
		return "undiredge"
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func GetElements(tn TreeNodeInterface) []TreeNodeInterface {
	subtrees, leafs := tn.GetChildren()
	ret := append(leafs, tn)
	for _, element := range subtrees {
		sub := GetElements(element)
		ret = append(ret, sub...)
	}
	return ret
}
