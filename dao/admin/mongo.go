package admin

import (
	"context"
	"errors"
	"strconv"

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
func (d *Dao) MongoGetUserPermission(ctx context.Context, userid primitive.ObjectID, nodeid []uint32) (canread, canwrite, admin bool, e error) {
	noderoute := make([][]uint32, 0, len(nodeid))
	for i := range nodeid {
		noderoute = append(noderoute, nodeid[:i+1])
	}
	var usernodes *model.UserNodes
	usernodes, e = d.MongoGetUserNodes(ctx, userid, noderoute)
	if e != nil {
		return
	}
	canread, canwrite, admin = usernodes.CheckNode(nodeid)
	return
}

var errNoChange = errors.New("no change")

//if admin is true,canread and canwrite will be ignore
//if admin is false and canread is false too,means delete this user from this node
//if admin is false and canwrite is true,then canread must be tree too
func (d *Dao) MongoUpdateUserPermission(ctx context.Context, operateUserid, targetUserid primitive.ObjectID, nodeid []uint32, admin, canread, canwrite bool) (e error) {
	if admin {
		//ignore
		canread = false
		canwrite = false
	} else if !canread && canwrite {
		e = ecode.ErrReq
		return
	}
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
			if e == errNoChange {
				e = nil
			}
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	if !admin && !canread {
		//delete
		//check admin
		var x bool
		if _, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid); e != nil || !x {
			if e == nil {
				e = ecode.ErrPermission
			}
			return
		}
		//all check success,delete database
		filter := bson.M{"user_id": targetUserid}
		for i, v := range nodeid {
			filter["node_id."+strconv.Itoa(i)] = v
		}
		if _, e = d.mongo.Database("admin").Collection("usernode").DeleteMany(sctx, filter); e != nil {
			return
		}
		//if node's length <= 1 means delete this user from the root
		//this is same as delete this user completely
		//so need to delete the account
		if len(nodeid) <= 1 {
			_, e = d.mongo.Database("admin").Collection("user").DeleteOne(sctx, bson.M{"_id": targetUserid})
		}
		return
	}
	//update
	//check target user exist
	var num int64
	num, e = d.mongo.Database("admin").Collection("user").CountDocuments(sctx, bson.M{"_id": targetUserid})
	if e != nil {
		return
	}
	if num == 0 {
		e = ecode.ErrReq
		return
	}
	var r, x bool
	//get target user permission on parent path
	if r, _, x, e = d.MongoGetUserPermission(sctx, targetUserid, nodeid[:len(nodeid)-1]); e != nil {
		return
	}
	if x {
		e = errNoChange
		return
	}
	if !r {
		e = ecode.ErrPNodeReadPermission
		return
	}
	if admin {
		//check operate user parent path admin
		_, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid[:len(nodeid)-1])
	} else {
		//check operate user current path admin
		_, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid)
	}
	if e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success
	filter := bson.M{"user_id": targetUserid, "node_id": nodeid}
	updater := bson.M{"$set": bson.M{"r": canread, "w": canwrite, "x": admin}}
	var result *mongo.UpdateResult
	if result, e = d.mongo.Database("admin").Collection("usernode").UpdateOne(sctx, filter, updater, options.Update().SetUpsert(true)); e != nil {
		return
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		e = errNoChange
		return
	}
	if admin {
		filter = bson.M{"user_id": targetUserid}
		for i, v := range nodeid {
			filter["node_id."+strconv.Itoa(i)] = v
		}
		filter["node_id."+strconv.Itoa(len(nodeid))] = bson.M{"$exists": true}
		_, e = d.mongo.Database("admin").Collection("usernode").DeleteMany(sctx, filter)
	}
	return
}

//if nodeids are not empty or nil,only the node in the required nodeids will return
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
		if tmp.R {
			result.R = append(result.R, tmp.NodeId)
		}
		if tmp.W {
			result.W = append(result.W, tmp.NodeId)
		}
		if tmp.X {
			result.X = append(result.X, tmp.NodeId)
		}
	}
	return result, cursor.Err()
}

//if userids are not empty or nil,only the user in the required userids will return
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
		if tmp.R {
			result.R = append(result.R, tmp.UserId)
		}
		if tmp.W {
			result.W = append(result.W, tmp.UserId)
		}
		if tmp.X {
			result.X = append(result.X, tmp.UserId)
		}
	}
	return result, cursor.Err()
}

func (d *Dao) MongoAddNode(ctx context.Context, operateUserid primitive.ObjectID, pnodeid []uint32, name, data string) (e error) {
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
	var x bool
	if _, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, pnodeid); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success,modify database
	if _, e = d.mongo.Database("admin").Collection("node").InsertOne(sctx, &model.Node{
		NodeId:       append(pnodeid, parent.CurNodeIndex+1),
		NodeName:     name,
		NodeData:     data,
		CurNodeIndex: 0,
	}); e != nil {
		return
	}
	if _, e = d.mongo.Database("admin").Collection("node").UpdateOne(sctx, bson.M{"node_id": pnodeid}, bson.M{"$inc": bson.M{"cur_node_index": 1}}); e != nil {
		return
	}
	return
}
func (d *Dao) MongoUpdateNode(ctx context.Context, operateUserid primitive.ObjectID, nodeid []uint32, name, data string) (e error) {
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
	var x bool
	if _, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success,update database
	_, e = d.mongo.Database("admin").Collection("node").UpdateOne(sctx, bson.M{"node_id": nodeid}, bson.M{"$set": bson.M{"node_name": name, "node_data": data}})
	return
}
func (d *Dao) MongoMoveNode(ctx context.Context, operateUserid primitive.ObjectID, nodeid, pnodeid []uint32) (e error) {
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
	//check self exist
	var self int64
	self, e = d.mongo.Database("admin").Collection("node").CountDocuments(sctx, bson.M{"node_id": nodeid})
	if e != nil {
		return
	}
	if self == 0 {
		e = ecode.ErrReq
		return
	}
	//check parent exist
	parent := &model.Node{}
	if e = d.mongo.Database("admin").Collection("node").FindOne(sctx, bson.M{"node_id": pnodeid}).Decode(parent); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrReq
		}
		return
	}
	//check admin in current path
	var x bool
	if _, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//check admin in new path
	if _, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, pnodeid); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//update the new parent
	if _, e = d.mongo.Database("admin").Collection("node").UpdateOne(sctx, bson.M{"node_id": pnodeid}, bson.M{"$inc": bson.M{"cur_node_index": 1}}); e != nil {
		return
	}
	filter := bson.M{}
	unset := bson.M{}
	for i, v := range nodeid {
		filter["node_id."+strconv.Itoa(i)] = v
		unset["node_id."+strconv.Itoa(i)] = 1
	}
	updater := bson.A{
		bson.M{"$unset": unset},
		bson.M{"$push": bson.M{
			"$each":     append(parent.NodeId, parent.CurNodeIndex+1),
			"$position": 0,
		}},
		bson.M{"$pull": bson.M{"node_id": nil}},
	}
	//update the node
	if _, e = d.mongo.Database("admin").Collection("node").UpdateMany(sctx, filter, updater); e != nil {
		return
	}
	//update the usernode
	_, e = d.mongo.Database("admin").Collection("usernode").UpdateMany(sctx, filter, updater)
	return
}
func (d *Dao) MongoListNode(ctx context.Context, operateUserid primitive.ObjectID, pnodeid []uint32) (nodes []*model.Node, e error) {
	//check canread or admin
	var r, x bool
	if r, _, x, e = d.MongoGetUserPermission(ctx, operateUserid, pnodeid); e != nil || (!x && !r) {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success,query database
	filter := bson.M{"node_id": bson.M{"$size": len(pnodeid) + 1}}
	for i, v := range pnodeid {
		filter["node_id."+strconv.Itoa(i)] = v
	}
	var cursor *mongo.Cursor
	cursor, e = d.mongo.Database("admin").Collection("node").Find(ctx, filter)
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
func (d *Dao) MongoDeleteNode(ctx context.Context, operateUserid primitive.ObjectID, nodeid []uint32) (e error) {
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
	var x bool
	if _, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success,delete database
	delfilter := bson.M{}
	for i, v := range nodeid {
		delfilter["node_id."+strconv.Itoa(i)] = v
	}
	if _, e = d.mongo.Database("admin").Collection("node").DeleteMany(sctx, delfilter); e != nil {
		return
	}
	if _, e = d.mongo.Database("admin").Collection("usernode").DeleteMany(sctx, delfilter); e != nil {
		return
	}
	return
}
