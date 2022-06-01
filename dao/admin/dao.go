package admin

import (
	csql "database/sql"

	credis "github.com/chenjie199234/Corelib/redis"
	cmongo "go.mongodb.org/mongo-driver/mongo"
)

//Dao this is a data operation layer to operate admin service's data
type Dao struct {
	sql   *csql.DB
	redis *credis.Pool
	mongo *cmongo.Client
}

//NewDao Dao is only a data operation layer
//don't write business logic in this package
//business logic should be written in service package
func NewDao(sql *csql.DB, redis *credis.Pool, mongo *cmongo.Client) *Dao {
	d := &Dao{
		sql:   sql,
		redis: redis,
		mongo: mongo,
	}
	if e := d.initmongo(); e != nil {
		panic("[dao.admin] init error: " + e.Error())
	}
	return d
}
