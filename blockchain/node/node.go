package node

type node struct {
	Address string
	Peers   []string
}

var self node

func GetSelf() node {
	return self
}

func Init(address string, peers []string) node {
	self := node{
		Address: address,
		Peers:   peers,
	}
	return self
}
