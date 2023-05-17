package main

// /////////////////////////////////////////////////////////////////////////
type DrawioPoint struct {
	x int
	y int
}

func (di *DrawioPoint) GetX() int {
	return di.x
}
func (di *DrawioPoint) GetY() int {
	return di.y
}

// //////////////////////////////////////////////////////////////////////////////

type DrawioElement struct {
	parentId    uint
	deType      string
	id          uint
	x           int
	y           int
	name        string
	width       int
	hight       int
	ip          string
	key         string
	floating_ip string
	svi         string
	srcId       uint
	dstId       uint
	label       string
	points      []DrawioPoint
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

func (di *DrawioElement) GetX() int {
	return di.x
}
func (di *DrawioElement) GetY() int {
	return di.y
}
func (di *DrawioElement) GetName() string {
	return di.name
}

func (di *DrawioElement) GetHight() int {
	return di.hight
}
func (di *DrawioElement) GetWidth() int {
	return di.width
}
func (di *DrawioElement) GetIP() string {
	return di.ip
}
func (di *DrawioElement) GetKey() string {
	return di.key
}

func (di *DrawioElement) GetFip() string {
	return di.floating_ip
}

func (di *DrawioElement) GetSvi() string {
	return di.svi
}

func (di *DrawioElement) GetSrcId() uint {
	return di.srcId
}
func (di *DrawioElement) GetDstId() uint {
	return di.dstId
}
func (di *DrawioElement) GetLabel() string {
	return di.label
}

func (di *DrawioElement) GetPoints() []DrawioPoint {
	return di.points
}
func (di *DrawioElement) AddPoint(x int, y int) {
	di.points = append(di.points, DrawioPoint{x, y})
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
	SetParent(TreeNodeInterface)
	GetParent() TreeNodeInterface
	GetName() string
	AddIconTreeNode(icon TreeNodeInterface)
	GetIconTreeNodes() []TreeNodeInterface

	GetChildren() ([]TreeNodeInterface, []TreeNodeInterface)

	IsLine() bool
	IsNetwork() bool
	IsVPC() bool
	IsZone() bool
	IsSubnet() bool
	IsSG() bool
	IsVSI() bool
	IsNI() bool
	IsGateway() bool
	IsEndpoint() bool
	IsInternet() bool
	IsInternetService() bool
	IsUser() bool
	IsVsiConnector() bool
	IsDirectedEdge() bool
	IsUnDirectedEdge() bool

	SetLocation(location *Location)
	GetLocation() *Location

	GetDrawioElement() *DrawioElement
	SetDrawioInfo()
	AllocPassenger() int
}

///////////////////////////////////////////////////////////////////////////

type TreeNode struct {
	DrawioElement
	parent   TreeNodeInterface
	elements []TreeNodeInterface
	location *Location
}

func (tn *TreeNode) AllocPassenger() int {
	return 0
}
func (tn *TreeNode) GetName() string { return "name?" }

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

func NewTreeNode(parent TreeNodeInterface) TreeNode {
	tn := TreeNode{}
	tn.parent = parent
	tn.DrawioElement.SetId()
	return tn

}
func (tn *TreeNode) GetChildren() ([]TreeNodeInterface, []TreeNodeInterface) {
	return []TreeNodeInterface{}, []TreeNodeInterface{}
}

func (tn *TreeNode) SetLocation(location *Location) {
	tn.location = location
}
func (tn *TreeNode) GetLocation() *Location {
	return tn.location
}

func (tn *TreeNode) IsLine() bool            { return false }
func (tn *TreeNode) IsNetwork() bool         { return false }
func (tn *TreeNode) IsVPC() bool             { return false }
func (tn *TreeNode) IsZone() bool            { return false }
func (tn *TreeNode) IsSubnet() bool          { return false }
func (tn *TreeNode) IsSG() bool              { return false }
func (tn *TreeNode) IsVSI() bool             { return false }
func (tn *TreeNode) IsNI() bool              { return false }
func (tn *TreeNode) IsGateway() bool         { return false }
func (tn *TreeNode) IsEndpoint() bool        { return false }
func (tn *TreeNode) IsInternet() bool        { return false }
func (tn *TreeNode) IsInternetService() bool { return false }
func (tn *TreeNode) IsUser() bool            { return false }
func (tn *TreeNode) IsVsiConnector() bool    { return false }
func (tn *TreeNode) IsDirectedEdge() bool    { return false }
func (tn *TreeNode) IsUnDirectedEdge() bool  { return false }

func (tn *TreeNode) GetDrawioElement() *DrawioElement {
	return &tn.DrawioElement
}

///////////////////////////////////////////////////////////////////////

type SquereTreeNode struct {
	TreeNode
}

func NewSquereTreeNode(parent TreeNodeInterface) SquereTreeNode {
	return SquereTreeNode{NewTreeNode(parent)}
}

func (tn *SquereTreeNode) SetDrawioInfo() {
	location := tn.GetLocation()
	parentLocation := location
	if tn.GetParent() != nil {
		parentLocation = tn.GetParent().GetLocation()
		tn.DrawioElement.parentId = tn.GetParent().GetDrawioElement().GetId()
	}
	tn.DrawioElement.name = tn.GetName()
	tn.DrawioElement.width = location.lastCol.thickness + location.lastCol.distance - location.firstCol.distance
	tn.DrawioElement.hight = location.lastRow.thickness + location.lastRow.distance - location.firstRow.distance
	tn.DrawioElement.x = location.firstCol.distance - parentLocation.firstCol.distance
	tn.DrawioElement.y = location.firstRow.distance - parentLocation.firstRow.distance
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
func (tn *NetworkTreeNode) IsNetwork() bool { return true }

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
func (tn *VpcTreeNode) IsVPC() bool { return true }

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
func (tn *ZoneTreeNode) IsZone() bool { return true }

///////////////////////////////////////////////////////////////////////

type SGTreeNode struct {
	SquereTreeNode
}

func NewSGTreeNode(parent *VpcTreeNode) *SGTreeNode {
	sg := SGTreeNode{NewSquereTreeNode(parent)}
	parent.sgs = append(parent.sgs, &sg)
	return &sg
}
func (tn *SGTreeNode) IsSG() bool { return true }

func (tn *SGTreeNode) SetDrawioInfo() {
	tn.DrawioElement.parentId = tn.GetParent().GetDrawioElement().GetId()
	tn.DrawioElement.name = "sg name"
	tn.DrawioElement.width = 400
	tn.DrawioElement.hight = 400
	tn.DrawioElement.x = 40
	tn.DrawioElement.y = 40
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
func (tn *SubnetTreeNode) IsSubnet() bool { return true }

/////////////////////////////////////////////////////////////////////////

type IconTreeNode struct {
	TreeNode
	nPassengers int
}

func NewIconTreeNode(parent TreeNodeInterface) IconTreeNode {
	return IconTreeNode{NewTreeNode(parent), 0}
}

func (tn *IconTreeNode) AllocPassenger() int {
	n := tn.nPassengers
	tn.nPassengers++
	return []int{0, 10, -10, 20, -20, 15, -15, 25, -25}[n]
}
func (tn *IconTreeNode) SetDrawioInfo() {
	location := tn.GetLocation()
	parentLocation := tn.GetParent().GetLocation()
	tn.DrawioElement.parentId = tn.GetParent().GetDrawioElement().GetId()
	tn.DrawioElement.name = "icon name"
	tn.DrawioElement.x = location.firstCol.distance - parentLocation.firstCol.distance + location.firstCol.thickness/2 - iconSize/2 + location.x_offset
	tn.DrawioElement.y = location.firstRow.distance - parentLocation.firstRow.distance + location.firstRow.thickness/2 - iconSize/2 + location.y_offset
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
func (tn *NITreeNode) HasVsi() bool { return tn.hasVsi }
func (tn *NITreeNode) IsNI() bool   { return true }

// ///////////////////////////////////////////
type GetWayTreeNode struct {
	IconTreeNode
}

func NewGetWayTreeNode(parent TreeNodeInterface) *GetWayTreeNode {
	gw := GetWayTreeNode{NewIconTreeNode(parent)}
	parent.AddIconTreeNode(&gw)
	return &gw
}
func (tn *GetWayTreeNode) IsGateway() bool { return true }

// ///////////////////////////////////////////
type UserTreeNode struct {
	IconTreeNode
}

func NewUserTreeNode(parent TreeNodeInterface) *UserTreeNode {
	user := UserTreeNode{NewIconTreeNode(parent)}
	parent.AddIconTreeNode(&user)
	return &user
}
func (tn *UserTreeNode) IsUser() bool { return true }

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

func (tn *VsiTreeNode) IsVSI() bool { return true }

// ///////////////////////////////////////////
type InternetTreeNode struct {
	IconTreeNode
}

func NewInternetTreeNode(parent TreeNodeInterface) *InternetTreeNode {
	inter := InternetTreeNode{NewIconTreeNode(parent)}
	parent.AddIconTreeNode(&inter)
	return &inter
}
func (tn *InternetTreeNode) IsInternet() bool { return true }

// ////////////////////////////////////////////////////////////////
type InternetSeviceTreeNode struct {
	IconTreeNode
}

func NewInternetSeviceTreeNode(parent TreeNodeInterface) *InternetSeviceTreeNode {
	inter := InternetSeviceTreeNode{NewIconTreeNode(parent)}
	parent.AddIconTreeNode(&inter)
	return &inter
}
func (tn *InternetSeviceTreeNode) IsInternetService() bool { return true }

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
	tn.DrawioElement.srcId = tn.src.GetDrawioElement().GetId()
	tn.DrawioElement.dstId = tn.dst.GetDrawioElement().GetId()
	if tn.GetParent().IsNI() {
		tn.DrawioElement.parentId = tn.GetParent().GetDrawioElement().GetId2()
	} else {
		tn.DrawioElement.parentId = tn.GetParent().GetDrawioElement().GetId()
	}
}

// ////////////////////////////////////////////////////////////////
type VsiLineTreeNode struct {
	LineTreeNode
}

func NewVsiLineTreeNode(network TreeNodeInterface, vsi TreeNodeInterface, ni TreeNodeInterface) *VsiLineTreeNode {
	vsi.(*VsiTreeNode).AddNI(ni)
	conn := VsiLineTreeNode{LineTreeNode{NewTreeNode(network), vsi, ni}}
	network.(*NetworkTreeNode).connections = append(network.(*NetworkTreeNode).connections, &conn)
	return &conn
}
func (tn *VsiLineTreeNode) IsVsiConnector() bool { return true }

// ////////////////////////////////////////////////////////////////
type ConnectivityTreeNode struct {
	LineTreeNode
	directed bool
}

func NewConnectivityLineTreeNode(network TreeNodeInterface, src TreeNodeInterface, dst TreeNodeInterface, directed bool) *ConnectivityTreeNode {
	conn := ConnectivityTreeNode{LineTreeNode{NewTreeNode(network), src, dst}, directed}
	network.(*NetworkTreeNode).connections = append(network.(*NetworkTreeNode).connections, &conn)
	return &conn
}
func (tn *ConnectivityTreeNode) SetPassage(passage TreeNodeInterface, reverse bool) {
	tn.SetParent(passage)
	passengerNumber := passage.AllocPassenger()
	if !reverse {
		tn.DrawioElement.AddPoint(iconSize, iconSize/2+passengerNumber)
		tn.DrawioElement.AddPoint(0, iconSize/2+passengerNumber)
	} else {
		tn.DrawioElement.AddPoint(0, iconSize/2+passengerNumber)
		tn.DrawioElement.AddPoint(iconSize, iconSize/2+passengerNumber)
	}

}

func (tn *ConnectivityTreeNode) IsDirectedEdge() bool   { return tn.directed }
func (tn *ConnectivityTreeNode) IsUnDirectedEdge() bool { return !tn.directed }

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
