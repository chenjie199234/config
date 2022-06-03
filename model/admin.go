package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Node struct {
	NodeId       []uint32 `bson:"node_id"` //self's id
	NodeName     string   `bson:"node_name"`
	NodeData     string   `bson:"node_data"`
	CurNodeIndex uint32   `bson:"cur_node_index"` //auto increment,this is for child's node_id
}
type User struct {
	ID primitive.ObjectID `bson:"_id"` //user's id
	//TODO add login way
}
type UserNode struct {
	UserId primitive.ObjectID `bson:"user_id"`
	NodeId []uint32           `bson:"node_id"`
	R      uint32             `bson:"r"`
	W      uint32             `bson:"w"`
	X      uint32             `bson:"x"`
}

type UserNodes struct {
	R [][]uint32
	W [][]uint32
	X [][]uint32
}
type NodeUsers struct {
	R []primitive.ObjectID
	W []primitive.ObjectID
	X []primitive.ObjectID
}

func (u *UserNodes) CheckNode(nodeid []uint32) (canread, canwrite, admin bool) {
	for i := 0; i < len(nodeid); i++ {
		tmp := nodeid[:i+1]
		//check admin first
		has := false
		for _, x := range u.X {
			if len(x) != len(tmp) {
				continue
			}
			same := true
			for j := range x {
				if x[j] != tmp[j] {
					same = false
					break
				}
			}
			if same {
				has = true
				break
			}
		}
		if has {
			return true, true, true
		}
		//check can read
		has = false
		for _, r := range u.R {
			if len(r) != len(tmp) {
				continue
			}
			same := true
			for j := range r {
				if r[j] != tmp[j] {
					same = false
					break
				}
			}
			if same {
				has = true
				break
			}
		}
		if !has {
			return false, false, false
		}
	}
	canread = true
	admin = false
	for _, w := range u.W {
		if len(w) != len(nodeid) {
			continue
		}
		same := true
		for j := range w {
			if w[j] != nodeid[j] {
				same = false
				break
			}
		}
		if same {
			canwrite = true
			break
		}
	}
	return
}
func (n *NodeUsers) CheckUser(userid string) (canread, canwrite, admin bool) {
	for _, x := range n.X {
		if x.Hex() == userid {
			return true, true, true
		}
	}
	for _, r := range n.R {
		if r.Hex() == userid {
			canread = true
			break
		}
	}
	if !canread {
		return false, false, false
	}
	for _, w := range n.W {
		if w.Hex() == userid {
			canwrite = true
			break
		}
	}
	return
}
