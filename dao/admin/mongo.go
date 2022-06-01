package admin

import (
	"context"

	"github.com/chenjie199234/config/ecode"
	"github.com/chenjie199234/config/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (d *Dao) initmongo() error {
	_, e := d.mongo.Database("admin").Collection("node").InsertOne(context.Background(), &model.Node{
		NodeId:       []uint32{0},
		NodeName:     "root",
		NodeData:     "",
		CurNodeIndex: 0,
		Children:     []uint32{},
	})
	if e != nil && !mongo.IsDuplicateKeyError(e) {
		return e
	}
	return nil
}
func (d *Dao) MongoGetUser(ctx context.Context, userid primitive.ObjectID) (*model.User, error) {
	user := &model.User{}
	if e := d.mongo.Database("admin").Collection("user").FindOne(ctx, bson.M{"_id": userid}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, e
	}
	return user, nil
}

//if nodeids is not empty or nil,only the node in the required nodeids will return
func (d *Dao) MongoGetUserNodes(ctx context.Context, userid primitive.ObjectID, nodeids [][]uint32) (*model.UserNodes, error) {
	filter := bson.M{"user_id": userid}
	if len(nodeids) > 0 {
		filter["node_id"] = bson.M{"$in": nodeids}
	}
	cursor, e := d.mongo.Database("admin").Collection("usernode").Find(ctx, filter)
	if e != nil {
		return nil, e
	}
	defer cursor.Close(ctx)
	result := &model.UserNodes{
		R: make([][]uint32, 0),
		W: make([][]uint32, 0),
		X: make([][]uint32, 0),
	}
	for cursor.Next(ctx) {
		tmp := &model.UserNode{}
		if e := cursor.Decode(tmp); e != nil {
			return nil, e
		}
		if tmp.R == 1 {
			result.R = append(result.R, tmp.NodeId)
		}
		if tmp.W == 1 {
			result.W = append(result.W, tmp.NodeId)
		}
		if tmp.X == 1 {
			result.X = append(result.X, tmp.NodeId)
		}
	}
	return result, cursor.Err()
}

//if userids is not empty or nil,only the user in the required userids will return
func (d *Dao) MongoGetNodeUsers(ctx context.Context, nodeid []uint32, userids []primitive.ObjectID) (*model.NodeUsers, error) {
	filter := bson.M{"node_id": nodeid}
	if len(userids) > 0 {
		filter["user_id"] = bson.M{"$in": userids}
	}
	cursor, e := d.mongo.Database("admin").Collection("usernode").Find(ctx, filter)
	if e != nil {
		return nil, e
	}
	defer cursor.Close(ctx)
	result := &model.NodeUsers{
		R: make([]primitive.ObjectID, 0),
		W: make([]primitive.ObjectID, 0),
		X: make([]primitive.ObjectID, 0),
	}
	for cursor.Next(ctx) {
		tmp := &model.UserNode{}
		if e := cursor.Decode(tmp); e != nil {
			return nil, e
		}
		if tmp.R == 1 {
			result.R = append(result.R, tmp.UserId)
		}
		if tmp.W == 1 {
			result.W = append(result.W, tmp.UserId)
		}
		if tmp.X == 1 {
			result.X = append(result.X, tmp.UserId)
		}
	}
	return result, cursor.Err()
}
func (d *Dao) MongoAddNode(ctx context.Context, userid primitive.ObjectID, pnodeid []uint32, name, data string) (e error) {
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	//check parent exist
	parent := &model.Node{}
	e = d.mongo.Database("admin").Collection("node").FindOne(sctx, bson.M{"node_id": pnodeid}).Decode(parent)
	if e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrReq
		}
		return
	}
	//check admin
	var usernodes *model.UserNodes
	noderoute := make([][]uint32, 0, len(pnodeid))
	for i := range pnodeid {
		noderoute = append(noderoute, pnodeid[:i+1])
	}
	usernodes, e = d.MongoGetUserNodes(sctx, userid, noderoute)
	if e != nil {
		return
	}
	if _, _, x := usernodes.CheckNode(pnodeid); !x {
		e = ecode.ErrAuth
		return
	}
	//all check success,modify database
	if _, e = d.mongo.Database("admin").Collection("node").InsertOne(sctx, &model.Node{
		NodeId:       append(pnodeid, parent.CurNodeIndex+1),
		NodeName:     name,
		NodeData:     data,
		CurNodeIndex: 0,
		Children:     []uint32{},
	}); e != nil {
		return
	}
	updater := bson.M{
		"$set": bson.M{
			"cur_node_index": parent.CurNodeIndex + 1,
		},
		"$push": bson.M{
			"children": parent.CurNodeIndex + 1,
		},
	}
	if _, e = d.mongo.Database("admin").Collection("node").UpdateOne(sctx, bson.M{"node_id": pnodeid}, updater); e != nil {
		return
	}
	return
}
func (d *Dao) MongoUpdateNode(ctx context.Context, userid primitive.ObjectID, nodeid []uint32, name, data string) (e error) {
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	//check admin
	var usernodes *model.UserNodes
	noderoute := make([][]uint32, 0, len(nodeid))
	for i := range nodeid {
		noderoute = append(noderoute, nodeid[:i+1])
	}
	usernodes, e = d.MongoGetUserNodes(sctx, userid, noderoute)
	if e != nil {
		return
	}
	if _, _, x := usernodes.CheckNode(nodeid); !x {
		e = ecode.ErrAuth
		return
	}
	//all check success,modify database
	_, e = d.mongo.Database("admin").Collection("node").UpdateOne(sctx, bson.M{"node_id": nodeid}, bson.M{"$set": bson.M{"node_name": name, "node_data": data}})
	return
}
func (d *Dao) MongoListNode(ctx context.Context, userid primitive.ObjectID, pnodeid []uint32) (nodes []*model.Node, e error) {
	//check canread or admin
	var usernodes *model.UserNodes
	noderoute := make([][]uint32, 0, len(pnodeid))
	for i := range pnodeid {
		noderoute = append(noderoute, pnodeid[:i+1])
	}
	usernodes, e = d.MongoGetUserNodes(ctx, userid, noderoute)
	if e != nil {
		return
	}
	if r, _, x := usernodes.CheckNode(pnodeid); !r && !x {
		e = ecode.ErrAuth
		return
	}
	//all check success,search database
	parent := &model.Node{}
	if e = d.mongo.Database("admin").Collection("node").FindOne(ctx, bson.M{"node_id": pnodeid}).Decode(parent); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrReq
		}
		return
	}
	if len(parent.Children) == 0 {
		return nil, nil
	}
	nodeids := make([][]uint32, 0, len(parent.Children))
	for _, v := range parent.Children {
		nodeids = append(nodeids, append(parent.NodeId, v))
	}
	var cursor *mongo.Cursor
	cursor, e = d.mongo.Database("admin").Collection("node").Find(ctx, bson.M{"node_id": bson.M{"$in": nodeids}})
	if e != nil {
		return
	}
	defer cursor.Close(ctx)
	nodes = make([]*model.Node, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		tmp := &model.Node{}
		if e = cursor.Decode(tmp); e != nil {
			return
		}
		nodes = append(nodes, tmp)
	}
	e = cursor.Err()
	return
}
func (d *Dao) MongoDeleteNode(ctx context.Context, userid primitive.ObjectID, nodeid []uint32) (e error) {
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	//check admin
	var usernodes *model.UserNodes
	noderoute := make([][]uint32, 0, len(nodeid))
	for i := range nodeid {
		noderoute = append(noderoute, nodeid[:i+1])
	}
	usernodes, e = d.MongoGetUserNodes(sctx, userid, noderoute)
	if e != nil {
		return
	}
	if _, _, x := usernodes.CheckNode(nodeid); !x {
		e = ecode.ErrAuth
		return
	}
	//all check success,delete database
	return d.delnode(sctx, [][]uint32{nodeid})
}
func (d *Dao) delnode(ctx context.Context, nodeids [][]uint32) error {
	cursor, e := d.mongo.Database("admin").Collection("node").Find(ctx, bson.M{"node_id": bson.M{"$in": nodeids}})
	if e != nil {
		return e
	}
	nodes := make([]*model.Node, 0, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		tmp := &model.Node{}
		if e = cursor.Decode(tmp); e != nil {
			cursor.Close(ctx)
			return e
		}
		nodes = append(nodes, tmp)
	}
	cursor.Close(ctx)
	if cursor.Err() != nil {
		return cursor.Err()
	}
	if len(nodes) == 0 {
		return nil
	}
	if _, e = d.mongo.Database("admin").Collection("node").DeleteMany(ctx, bson.M{"node_id": bson.M{"$in": nodeids}}); e != nil {
		return e
	}
	if _, e = d.mongo.Database("admin").Collection("usernode").DeleteMany(ctx, bson.M{"node_id": bson.M{"$in": nodeids}}); e != nil {
		return e
	}
	childrennum := 0
	for _, node := range nodes {
		childrennum += len(node.Children)
	}
	newnodeids := make([][]uint32, 0, childrennum)
	for _, node := range nodes {
		for _, child := range node.Children {
			newnodeids = append(newnodeids, append(node.NodeId, child))
		}
	}
	if len(newnodeids) == 0 {
		return nil
	}
	return d.delnode(ctx, newnodeids)
}
