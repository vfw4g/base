package vfwrepository

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/joomcode/errorx"
	"xorm.io/core"
)

var repos = newXormRepository()

func newXormRepository() *XormRepository {
	connArgs := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", "txvier", "txvier", "localhost", 3306, "txvier")
	engine, err := xorm.NewEngine("mysql", connArgs)
	if err != nil {
		panic(err)
	}
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "t_")
	engine.SetTableMapper(tbMapper)
	//createTableIfNotExist(engine)
	//engine.Sync(new(entity.User))
	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_DEBUG)
	return &XormRepository{
		engine,
	}
}

/*
func NewXormRepository() *XormRepository {
	return repos
}*/

func NewXormRepository(rdsc RDSConfig) (r *XormRepository, err error) {
	url := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", rdsc.UserName(), rdsc.Password(), rdsc.Host(), rdsc.Port(), rdsc.Schema())
	engine, err := xorm.NewEngine(rdsc.DriverName(), url)
	if err != nil {
		return nil, errorx.InitializationFailed.WrapWithNoMessage(err)
	}
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "t_")
	engine.SetTableMapper(tbMapper)
	//createTableIfNotExist(engine)
	//engine.Sync(new(entity.User))
	engine.ShowSQL(rdsc.ShowSQL())
	//todo from rdsc.LogLevel()
	engine.Logger().SetLevel(core.LOG_DEBUG)
	return &XormRepository{
		engine,
	}, nil
}

type XormRepository struct {
	engine *xorm.Engine
}

func (xr *XormRepository) Get(cond Entity) (isExist bool, err error) {
	isExist, err = xr.engine.Get(cond)
	if err != nil {
		return false, errorx.InternalError.WrapWithNoMessage(err)
	}
	return
}

func (xr *XormRepository) Save(entity Entity) (affected int64, err error) {
	affected, err = xr.engine.InsertOne(entity)
	if err != nil {
		return affected, errorx.InternalError.WrapWithNoMessage(err)
	}
	return
}

func (xr *XormRepository) Update(entity Entity) (affected int64, err error) {
	affected, err = xr.engine.Id(entity.GetId()).Update(entity)
	if err != nil {
		return affected, errorx.InternalError.WrapWithNoMessage(err)
	}
	return
}

func (xr *XormRepository) Del(entity Entity) (affected int64, err error) {
	affected, err = xr.engine.Id(entity.GetId()).Unscoped().Delete(entity)
	if err != nil {
		return affected, errorx.InternalError.WrapWithNoMessage(err)
	}
	return
}

func (xr *XormRepository) MarkDel(entity Entity) (affected int64, err error) {
	affected, err = xr.engine.Id(entity.GetId()).Delete(entity)
	if err != nil {
		return affected, errorx.InternalError.WrapWithNoMessage(err)
	}
	return
}

func (xr *XormRepository) Exist(cond Entity) (isExist bool, err error) {
	isExist, err = xr.engine.Exist(cond)
	if err != nil {
		return false, errorx.InternalError.WrapWithNoMessage(err)
	}
	return isExist, nil
}

func (xr *XormRepository) FindSqlWithLimit(results interface{}, sql string, limit, start int, params ...interface{}) (count int64, err error) {
	count, err = xr.engine.SQL(sql, params).Limit(limit, start).FindAndCount(results)
	if err != nil {
		return 0, errorx.InternalError.WrapWithNoMessage(err)
	}
	return
}

func (xr *XormRepository) FindAll(results interface{}, desc ...string) (count int64, err error) {
	count, err = xr.engine.Desc(desc...).FindAndCount(results)
	if err != nil {
		return 0, errorx.InternalError.WrapWithNoMessage(err)
	}
	return
}

func (xr *XormRepository) Find(results interface{}, desc ...string) (err error) {
	err = xr.engine.Desc(desc...).Find(results)
	if err != nil {
		return errorx.InternalError.WrapWithNoMessage(err)
	}
	return
}

func (xr *XormRepository) FindLimit(results interface{}, limit, start int, desc ...string) (count int64, err error) {
	count, err = xr.engine.Limit(limit, start).Desc(desc...).FindAndCount(results)
	if err != nil {
		return 0, errorx.InternalError.WrapWithNoMessage(err)
	}
	return
}

func (xr *XormRepository) T(results interface{}, limit, start int, desc ...string) (count int64, err error) {
	//session := xr.engine.NewSession()
	//session.
	return
}
