package main

import (
	"os"
	"text/template"
)

type DrawioElementInterface interface {
	GetType() string
	GetId() uint
}

type DrawioElement struct {
	deType string
	id     uint
}

func (di DrawioElement) GetType() string {
	return di.deType
}
func (di DrawioElement) GetId() uint {
	return di.id
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

// //////////////////////////////////////////////////////////////////////////////
type DrawioInfo struct {
	IconSize uint
	Elements []DrawioElementInterface
}

func main() {

	elements := []DrawioElementInterface{
		NewDrawioSquereElement("pub", 40, 40, 1920, 1040, 895, 0, "noname"),
		NewDrawioIconElement("internet", 50.0, 50.0, 60, 60, 900, 895, "161.26.0.0/16"),
		NewDrawioIconElement("internet", 50.0, 210.0, 60, 60, 905, 895, "166.8.0.0/14"),
		NewDrawioIconElement("internet", 50.0, 370.0, 60, 60, 910, 895, "192.0.1.0/24"),
		NewDrawioSquereElement("vpc", 160, 40, 1720, 960, 1230, 895, "noname"),
		NewDrawioSquereElement("zone", 40, 40, 520, 880, 1235, 1230, "us-south-3"),
		NewDrawioSubnetElement("subnet", 160, 600, 320, 240, 1240, 1235, "ky-testenv-transit-subnet-3", "0000000", "key"),
		NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1245, 1240, "transit-2-instance-ky", "", ""),
		NewDrawioSubnetElement("subnet", 160, 40, 320, 240, 1250, 1235, "ky-testenv-edge-subnet-3", "0000000", "key"),
		NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1255, 1250, "edge-2-instance-ky", "", "vsi_name"),
		NewDrawioSubnetElement("subnet", 160, 320, 320, 240, 1260, 1235, "ky-testenv-private-subnet-3", "0000000", "key"),
		NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1265, 1260, "private-2-instance-ky", "", "vsi_name"),
		NewDrawioIconElement("gateway", 50.0, 130.0, 60, 60, 1270, 1235, "ky-testenv-gateway-3"),
		NewDrawioSquereElement("zone", 600, 40, 520, 880, 1275, 1230, "us-south-2"),
		NewDrawioSubnetElement("subnet", 160, 600, 320, 240, 1280, 1275, "ky-testenv-transit-subnet-2", "0000000", "key"),
		NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1285, 1280, "transit-1-instance-ky", "fip_name", "vsi_name"),
		NewDrawioSubnetElement("subnet", 160, 40, 320, 240, 1290, 1275, "ky-testenv-edge-subnet-2", "0000000", "key"),
		NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1295, 1290, "edge-1-instance-ky", "", "vsi_name"),
		NewDrawioSubnetElement("subnet", 160, 320, 320, 240, 1300, 1275, "ky-testenv-private-subnet-2", "0000000", "key"),
		NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1305, 1300, "private-1-instance-ky", "", "vsi_name"),
		NewDrawioIconElement("gateway", 50.0, 410.0, 60, 60, 1310, 1275, "ky-testenv-gateway-2"),
		NewDrawioSquereElement("zone", 1160, 40, 520, 880, 1315, 1230, "us-south-1"),
		NewDrawioSubnetElement("subnet", 160, 320, 320, 240, 1320, 1315, "ky-testenv-private-subnet-1", "0000000", "key"),
		NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1325, 1320, "private-0-instance-ky", "fip_name", "vsi_name"),
		NewDrawioSubnetElement("subnet", 160, 40, 320, 240, 1330, 1315, "ky-testenv-edge-subnet-1", "0000000", "key"),
		NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1335, 1330, "edge-0-instance-ky", "", "vsi_name"),
		NewDrawioSubnetElement("subnet", 160, 600, 320, 240, 1340, 1315, "ky-testenv-transit-subnet-1", "0000000", "key"),
		NewDrawioNetworkInterfaceElement("ni", 130.0, 90.0, 60, 60, 1345, 1340, "transit-0-instance-ky", "", "vsi_name"),
		NewDrawioIconElement("gateway", 50.0, 690.0, 60, 60, 1350, 1315, "ky-testenv-gateway-1"),
		NewDrawioSquereElement("sg", 290.0, 120, 1260.0, 160, 1355, 1230, "ky-testenv-default-sg"),
		NewDrawioSquereElement("sg", 290.0, 400, 1260.0, 160, 1360, 1230, "ky-testenv-default-sg"),
		NewDrawioSquereElement("sg", 290.0, 680, 1260.0, 160, 1365, 1230, "ky-testenv-default-sg"),
		NewDrawioConnectElement("diredge", 915, 1325, 1305, "TCP 1-65535", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 920, 1325, 1345, "TCP 1-65535", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 925, 1325, 1245, "TCP 1-65535", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 930, 1345, 1255, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 190.0}}),
		NewDrawioConnectElement("diredge", 935, 1345, 1295, "TCP 443", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 940, 1345, 1325, "TCP 443", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 945, 1345, 1335, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{1590.0, 530.0}}),
		NewDrawioConnectElement("diredge", 950, 1345, 1265, "TCP 443", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 955, 1345, 1285, "TCP 443", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 960, 1265, 1245, "TCP 1-65535", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 965, 1265, 1305, "TCP 1-65535", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 970, 1265, 1345, "TCP 1-65535", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 975, 1285, 1345, "TCP 1-65535", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 980, 1285, 1245, "TCP 1-65535", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 985, 1285, 1305, "TCP 1-65535", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 990, 1305, 1335, "TCP 443", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 995, 1305, 1325, "TCP 443", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 1000, 1305, 1285, "TCP 443", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 1005, 1305, 1295, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{1030.0, 530.0}}),
		NewDrawioConnectElement("diredge", 1010, 1305, 1255, "TCP 443", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 1015, 1305, 1265, "TCP 443", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 1020, 1255, 900, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{350.0, 280.0}, DrawioConnectElementPoint{290.0, 280.0}}),
		NewDrawioConnectElement("diredge", 1025, 1255, 1305, "TCP 1-65535", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 1030, 1255, 1345, "TCP 1-65535", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 190.0}}),
		NewDrawioConnectElement("diredge", 1035, 1255, 910, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{350.0, 276.0}, DrawioConnectElementPoint{290.0, 276.0}}),
		NewDrawioConnectElement("diredge", 1040, 1255, 905, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{350.0, 272.0}, DrawioConnectElementPoint{290.0, 272.0}, DrawioConnectElementPoint{310.0, 190.0}}),
		NewDrawioConnectElement("diredge", 1045, 1255, 1245, "TCP 1-65535", []DrawioConnectElementPoint{DrawioConnectElementPoint{470.0, 530.0}}),
		NewDrawioConnectElement("diredge", 1050, 1335, 1245, "TCP 1-65535", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 190.0}}),
		NewDrawioConnectElement("diredge", 1055, 1335, 1345, "TCP 1-65535", []DrawioConnectElementPoint{DrawioConnectElementPoint{1590.0, 530.0}}),
		NewDrawioConnectElement("diredge", 1060, 1335, 910, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{1470.0, 840.0}, DrawioConnectElementPoint{1410.0, 840.0}}),
		NewDrawioConnectElement("diredge", 1065, 1335, 905, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{1470.0, 836.0}, DrawioConnectElementPoint{1410.0, 836.0}, DrawioConnectElementPoint{870.0, 190.0}}),
		NewDrawioConnectElement("diredge", 1070, 1335, 1305, "TCP 1-65535", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 1075, 1335, 900, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{1470.0, 832.0}, DrawioConnectElementPoint{1410.0, 832.0}}),
		NewDrawioConnectElement("diredge", 1080, 1245, 1295, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{470.0, 530.0}}),
		NewDrawioConnectElement("diredge", 1085, 1245, 1255, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{470.0, 530.0}}),
		NewDrawioConnectElement("diredge", 1090, 1245, 1265, "TCP 443", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 1095, 1245, 1285, "TCP 443", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 1100, 1245, 1325, "TCP 443", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 1105, 1245, 1335, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 190.0}}),
		NewDrawioConnectElement("diredge", 1110, 1295, 1345, "TCP 1-65535", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("diredge", 1115, 1295, 905, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{910.0, 560.0}, DrawioConnectElementPoint{850.0, 560.0}, DrawioConnectElementPoint{590.0, 190.0}}),
		NewDrawioConnectElement("diredge", 1120, 1295, 910, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{910.0, 556.0}, DrawioConnectElementPoint{850.0, 556.0}}),
		NewDrawioConnectElement("diredge", 1125, 1295, 1245, "TCP 1-65535", []DrawioConnectElementPoint{DrawioConnectElementPoint{470.0, 530.0}}),
		NewDrawioConnectElement("diredge", 1130, 1295, 1305, "TCP 1-65535", []DrawioConnectElementPoint{DrawioConnectElementPoint{1030.0, 530.0}}),
		NewDrawioConnectElement("diredge", 1135, 1295, 900, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{910.0, 552.0}, DrawioConnectElementPoint{850.0, 552.0}}),
		NewDrawioConnectElement("undiredge", 1140, 1285, 1265, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{810.0, 470.0}}),
		NewDrawioConnectElement("undiredge", 1145, 1325, 1265, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 470.0}}),
		NewDrawioConnectElement("undiredge", 1150, 1325, 1285, "", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("undiredge", 1155, 1345, 1305, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{1370.0, 750.0}}),
		NewDrawioConnectElement("undiredge", 1160, 1335, 1255, "", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 190.0}}),
		NewDrawioConnectElement("undiredge", 1165, 1285, 1255, "", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("undiredge", 1170, 1325, 1255, "", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("undiredge", 1175, 1295, 1255, "", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("undiredge", 1180, 1265, 1255, "", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("undiredge", 1185, 1265, 1335, "", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("undiredge", 1190, 1325, 1335, "", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("undiredge", 1195, 1285, 1335, "", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("undiredge", 1200, 1345, 1245, "TCP 443", []DrawioConnectElementPoint{DrawioConnectElementPoint{1090.0, 750.0}}),
		NewDrawioConnectElement("undiredge", 1205, 1305, 1245, "TCP 443", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("undiredge", 1210, 1285, 1295, "", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("undiredge", 1215, 1325, 1295, "", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("undiredge", 1220, 1265, 1295, "", []DrawioConnectElementPoint{}),
		NewDrawioConnectElement("undiredge", 1225, 1335, 1295, "", []DrawioConnectElementPoint{}),
	}
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

}
