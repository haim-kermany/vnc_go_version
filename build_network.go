package main

func create_network() TreeNodeInterface {
	network := NewNetworkTreeNode()
	NewInternetTreeNode(network)
	NewInternetTreeNode(network)
	NewInternetTreeNode(network)
	NewInternetTreeNode(network)

	vpc1 := NewVpcTreeNode(network)
	zone11 := NewZoneTreeNode(vpc1)

	NewGetWayTreeNode(zone11)
	NewInternetSeviceTreeNode(vpc1)

	sg11 := NewSGTreeNode(vpc1)
	sg12 := NewSGTreeNode(vpc1)

	subnet111 := NewSubnetTreeNode(zone11)

	ni1 := NewNITreeNode(subnet111, sg11)
	ni2 := NewNITreeNode(subnet111, sg12)
	vsi1 := NewVsiTreeNode(zone11)
	vsi1.AddNI(ni1)
	vsi1.AddNI(ni2)
	NewVsiConnectorTreeNode(network, vsi1, ni1)
	NewVsiConnectorTreeNode(network, vsi1, ni2)

	zone12 := NewZoneTreeNode(vpc1)
	NewGetWayTreeNode(zone12)
	subnet112 := NewSubnetTreeNode(zone11)
	subnet121 := NewSubnetTreeNode(zone12)
	ni4 := NewNITreeNode(subnet112, sg12)
	ni4.GetDrawioElementInterface().(*DrawioIconElement).svi = "svi1"
	ni5 := NewNITreeNode(subnet121, sg11)
	ni5.GetDrawioElementInterface().(*DrawioIconElement).svi = "svi2"
	ni5.GetDrawioElementInterface().(*DrawioIconElement).floating_ip = "fip"

	vpc2 := NewVpcTreeNode(network)
	zone21 := NewZoneTreeNode(vpc2)
	sg21 := NewSGTreeNode(vpc2)

	subnet211 := NewSubnetTreeNode(zone21)

	vsi2 := NewVsiTreeNode(zone21)
	ni6 := NewNITreeNode(subnet211, sg21)
	ni7 := NewNITreeNode(subnet211, sg21)
	ni8 := NewNITreeNode(subnet211, sg21)
	vsi2.AddNI(ni6)
	vsi2.AddNI(ni7)
	vsi2.AddNI(ni8)
	NewVsiConnectorTreeNode(network, vsi2, ni6)
	NewVsiConnectorTreeNode(network, vsi2, ni7)
	NewVsiConnectorTreeNode(network, vsi2, ni8)

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
	vsi3.AddNI(ni10)
	vsi3.AddNI(ni11)
	vsi3.AddNI(ni12)
	NewVsiConnectorTreeNode(network, vsi3, ni10)
	NewVsiConnectorTreeNode(network, vsi3, ni13)
	NewVsiConnectorTreeNode(network, vsi3, ni14)

	vsi4 := NewVsiTreeNode(zone22)
	vsi4.AddNI(ni13)
	vsi4.AddNI(ni14)
	NewVsiConnectorTreeNode(network, vsi4, ni11)
	NewVsiConnectorTreeNode(network, vsi4, ni12)

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
	NewGetWayTreeNode(zone21)
	NewGetWayTreeNode(zone22)

	NewInternetSeviceTreeNode(vpc2)

	return network
}
