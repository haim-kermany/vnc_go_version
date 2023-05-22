package main

func create_network2() TreeNodeInterface {
	network := NewNetworkTreeNode()
	vpc1 := NewVpcTreeNode(network)
	zone11 := NewZoneTreeNode(vpc1)
	subnet111 := NewSubnetTreeNode(zone11)
	subnet112 := NewSubnetTreeNode(zone11)
	sg11 := NewSGTreeNode(vpc1)

	NewNITreeNode(subnet111, sg11)
	NewNITreeNode(subnet112, sg11)
	return network
}

func create_network() TreeNodeInterface {
	network := NewNetworkTreeNode()
	i1 := NewInternetTreeNode(network)
	i2 := NewInternetTreeNode(network)
	i3 := NewInternetTreeNode(network)
	i4 := NewInternetTreeNode(network)

	vpc1 := NewVpcTreeNode(network)
	zone11 := NewZoneTreeNode(vpc1)

	gw11 := NewGetWayTreeNode(zone11)
	is1 := NewInternetSeviceTreeNode(vpc1)

	sg11 := NewSGTreeNode(vpc1)
	sg12 := NewSGTreeNode(vpc1)

	subnet111 := NewSubnetTreeNode(zone11)

	ni1 := NewNITreeNode(subnet111, sg11)
	ni2 := NewNITreeNode(subnet111, sg12)
	vsi1 := NewVsiTreeNode(zone11)
	NewVsiLineTreeNode(network, vsi1, ni1)
	NewVsiLineTreeNode(network, vsi1, ni2)

	zone12 := NewZoneTreeNode(vpc1)
	gw12 := NewGetWayTreeNode(zone12)
	subnet112 := NewSubnetTreeNode(zone11)
	subnet121 := NewSubnetTreeNode(zone12)
	ni4 := NewNITreeNode(subnet112, sg12)
	ni4.GetDrawioElement().svi = "svi1"
	ni5 := NewNITreeNode(subnet121, sg11)
	ni5.GetDrawioElement().svi = "svi2"
	ni5.GetDrawioElement().floating_ip = "fip"
	ni5b := NewNITreeNode(subnet121, sg11)
	ni5b.GetDrawioElement().svi = "svi2"
	ni5b.GetDrawioElement().floating_ip = "fip2"

	vpc2 := NewVpcTreeNode(network)
	zone21 := NewZoneTreeNode(vpc2)
	sg21 := NewSGTreeNode(vpc2)

	subnet211 := NewSubnetTreeNode(zone21)

	vsi2 := NewVsiTreeNode(zone21)
	ni6 := NewNITreeNode(subnet211, sg21)
	ni7 := NewNITreeNode(subnet211, sg21)
	ni8 := NewNITreeNode(subnet211, sg21)
	NewVsiLineTreeNode(network, vsi2, ni6)
	NewVsiLineTreeNode(network, vsi2, ni7)
	NewVsiLineTreeNode(network, vsi2, ni8)

	zone22 := NewZoneTreeNode(vpc2)
	zone23 := NewZoneTreeNode(vpc2)
	subnet221 := NewSubnetTreeNode(zone22)
	subnet222 := NewSubnetTreeNode(zone22)
	subnet231 := NewSubnetTreeNode(zone23)
	sg22 := NewSGTreeNode(vpc2)

	ni10 := NewNITreeNode(subnet221, sg22)
	ni11 := NewNITreeNode(subnet222, sg22)
	ni12 := NewNITreeNode(subnet222, sg22)
	ni13 := NewNITreeNode(subnet222, sg22)
	ni14 := NewNITreeNode(subnet222, sg22)

	vsi3 := NewVsiTreeNode(zone22)
	NewVsiLineTreeNode(network, vsi3, ni10)
	NewVsiLineTreeNode(network, vsi3, ni13)
	NewVsiLineTreeNode(network, vsi3, ni14)

	vsi4 := NewVsiTreeNode(zone22)

	NewVsiLineTreeNode(network, vsi4, ni11)
	NewVsiLineTreeNode(network, vsi4, ni12)

	ni20 := NewNITreeNode(subnet231, sg22)
	ni21 := NewNITreeNode(subnet231, nil)
	ni22 := NewNITreeNode(subnet231, sg22)
	ni23 := NewNITreeNode(subnet231, sg22)
	ni24 := NewNITreeNode(subnet231, nil)
	ni25 := NewNITreeNode(subnet231, sg22)
	ni26 := NewNITreeNode(subnet231, nil)
	ni27 := NewNITreeNode(subnet231, nil)
	ni28 := NewNITreeNode(subnet231, sg22)
	ni29 := NewNITreeNode(subnet231, sg22)

	gw21 := NewGetWayTreeNode(zone21)
	gw22 := NewGetWayTreeNode(zone22)

	is2 := NewInternetSeviceTreeNode(vpc2)

	c1 := NewConnectivityLineTreeNode(network, ni4, i4, false)
	c2a := NewConnectivityLineTreeNode(network, ni5, i2, false)
	c2b := NewConnectivityLineTreeNode(network, ni5, i2, true)
	c2c := NewConnectivityLineTreeNode(network, i2, ni5, true)

	c3 := NewConnectivityLineTreeNode(network, ni8, i3, true)
	c4 := NewConnectivityLineTreeNode(network, ni11, i1, true)
	c5 := NewConnectivityLineTreeNode(network, ni12, i1, true)

	c6 := NewConnectivityLineTreeNode(network, ni5b, i3, false)
	c7 := NewConnectivityLineTreeNode(network, ni5b, i4, false)

	c1.SetPassage(gw11, false)
	c2a.SetPassage(ni5, false)
	c2b.SetPassage(ni5, false)
	c2c.SetPassage(ni5, true)
	c3.SetPassage(gw21, false)
	c4.SetPassage(gw22, false)
	c5.SetPassage(gw22, false)
	c6.SetPassage(gw12, false)
	c7.SetPassage(ni5b, false)

	NewConnectivityLineTreeNode(network, ni10, is2, true)
	NewConnectivityLineTreeNode(network, ni1, is1, true)

	NewConnectivityLineTreeNode(network, ni8, ni14, true)

	NewConnectivityLineTreeNode(network, ni20, ni22, true)
	NewConnectivityLineTreeNode(network, ni21, ni24, true)
	NewConnectivityLineTreeNode(network, ni23, ni27, true)
	NewConnectivityLineTreeNode(network, ni25, ni29, true)
	NewConnectivityLineTreeNode(network, ni26, ni27, true)
	NewConnectivityLineTreeNode(network, ni28, ni26, true)
	NewConnectivityLineTreeNode(network, ni22, ni28, true)

	return network
}
