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

// ///////////////////////////////////////////////////////////////////////////////////////////////
type DrawioIconElement struct {
	DrawioLocatedElement
}

func (di DrawioIconElement) GetHight() uint {
	return 60
}
func (di DrawioIconElement) GetWidth() uint {
	return 60
}

func NewDrawioIconElement(etype string, x uint, y uint, w uint, h uint, id uint, pid uint, name string) *DrawioIconElement {
	return &DrawioIconElement{DrawioLocatedElement{DrawioElement{etype, id}, pid, x, y, name}}
}

// //////////////////////////////////////////////////////////////////////////////
type DrawioConnectElement struct {
	DrawioElement
	srcId uint
	dstId uint
	name  string
}

func (di DrawioElement) GetIP() string {
	return "00000000000"
}
func (di DrawioElement) GetKey() string {
	return "noKey"
}

////////////////////////////////////////////////////////////////////////////////

func main() {

	elements := []DrawioElementInterface{
		NewDrawioSquereElement("pub", 40, 40, 1920, 1040, 895, 0, "noname"),
		NewDrawioIconElement("internet", 50.0, 50.0, 60, 60, 900, 895, "161.26.0.0/16"),
		NewDrawioIconElement("internet", 50.0, 210.0, 60, 60, 905, 895, "166.8.0.0/14"),
		NewDrawioIconElement("internet", 50.0, 370.0, 60, 60, 910, 895, "192.0.1.0/24"),
		NewDrawioSquereElement("vpc", 160, 40, 1720, 960, 1230, 895, "noname"),
		NewDrawioSquereElement("zone", 40, 40, 520, 880, 1235, 1230, "us-south-3"),
		NewDrawioSquereElement("subnet", 160, 600, 320, 240, 1240, 1235, "ky-testenv-transit-subnet-3"),
		NewDrawioIconElement("ni_vsi", 130.0, 90.0, 60, 60, 1245, 1240, "transit-2-instance-ky"),
		NewDrawioSquereElement("subnet", 160, 40, 320, 240, 1250, 1235, "ky-testenv-edge-subnet-3"),
		NewDrawioIconElement("ni_vsi", 130.0, 90.0, 60, 60, 1255, 1250, "edge-2-instance-ky"),
		NewDrawioSquereElement("subnet", 160, 320, 320, 240, 1260, 1235, "ky-testenv-private-subnet-3"),
		NewDrawioIconElement("ni_vsi", 130.0, 90.0, 60, 60, 1265, 1260, "private-2-instance-ky"),
		NewDrawioIconElement("gateway", 50.0, 130.0, 60, 60, 1270, 1235, "ky-testenv-gateway-3"),
		NewDrawioSquereElement("zone", 600, 40, 520, 880, 1275, 1230, "us-south-2"),
		NewDrawioSquereElement("subnet", 160, 600, 320, 240, 1280, 1275, "ky-testenv-transit-subnet-2"),
		NewDrawioIconElement("ni_vsi", 130.0, 90.0, 60, 60, 1285, 1280, "transit-1-instance-ky"),
		NewDrawioSquereElement("subnet", 160, 40, 320, 240, 1290, 1275, "ky-testenv-edge-subnet-2"),
		NewDrawioIconElement("ni_vsi", 130.0, 90.0, 60, 60, 1295, 1290, "edge-1-instance-ky"),
		NewDrawioSquereElement("subnet", 160, 320, 320, 240, 1300, 1275, "ky-testenv-private-subnet-2"),
		NewDrawioIconElement("ni_vsi", 130.0, 90.0, 60, 60, 1305, 1300, "private-1-instance-ky"),
		NewDrawioIconElement("gateway", 50.0, 410.0, 60, 60, 1310, 1275, "ky-testenv-gateway-2"),
		NewDrawioSquereElement("zone", 1160, 40, 520, 880, 1315, 1230, "us-south-1"),
		NewDrawioSquereElement("subnet", 160, 320, 320, 240, 1320, 1315, "ky-testenv-private-subnet-1"),
		NewDrawioIconElement("ni_vsi", 130.0, 90.0, 60, 60, 1325, 1320, "private-0-instance-ky"),
		NewDrawioSquereElement("subnet", 160, 40, 320, 240, 1330, 1315, "ky-testenv-edge-subnet-1"),
		NewDrawioIconElement("ni_vsi", 130.0, 90.0, 60, 60, 1335, 1330, "edge-0-instance-ky"),
		NewDrawioSquereElement("subnet", 160, 600, 320, 240, 1340, 1315, "ky-testenv-transit-subnet-1"),
		NewDrawioIconElement("ni_vsi", 130.0, 90.0, 60, 60, 1345, 1340, "transit-0-instance-ky"),
		NewDrawioIconElement("gateway", 50.0, 690.0, 60, 60, 1350, 1315, "ky-testenv-gateway-1"),
		NewDrawioSquereElement("sg", 290.0, 120, 1260.0, 160, 1355, 1230, "ky-testenv-default-sg"),
		NewDrawioSquereElement("sg", 290.0, 400, 1260.0, 160, 1360, 1230, "ky-testenv-default-sg"),
		NewDrawioSquereElement("sg", 290.0, 680, 1260.0, 160, 1365, 1230, "ky-testenv-default-sg"),
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

	err = tmpl.Execute(fo, elements)
	if err != nil {
		panic(err)
	}

}
