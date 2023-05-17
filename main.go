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

	type Sizes struct {
		IconSize   int
		FipXOffset int
		FipYOffset int
		VSIXOffset int
		VSIYOffset int
		VSISize    int
	}
	type DrawioInfo struct {
		Sizes
		Nodes []TreeNodeInterface
	}
	err = tmpl.Execute(fo, DrawioInfo{Sizes{iconSize, -iconSize, 30, 30, -10, 40}, elements})
	if err != nil {
		panic(err)
	}

}
