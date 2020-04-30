package vfwrepository

type Entity interface {
	GetId() interface{}
}

type Repository interface {
	//返回指定ID的实体
	Get(cond Entity) (bool, error)
	//保存实体对象
	Save(entity Entity) (affected int64, err error)
	//保存实体对象
	Update(entity Entity) (affected int64, err error)
	//物理删除实体对象
	Del(entity Entity) (affected int64, err error)
	//逻辑删除实体对象
	MarkDel(entity Entity) (affected int64, err error)
	//判断对象是否存在
	Exist(cond Entity) (isExist bool, err error)
	//SQL查询
	FindSqlWithLimit(results interface{}, sql string, limit, start int, params ...interface{}) (count int64, err error)

	FindAll(results interface{}, desc ...string) (count int64, err error)

	Find(results interface{}, desc ...string) (err error)

	FindLimit(results interface{}, limit, start int, desc ...string) (count int64, err error)
}
