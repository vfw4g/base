package bean

import (
	"github.com/vfw4g/base/errors"
	"reflect"
)

type Instance interface{}

type BeanFactory func() Instance

type beanMeta struct {
	Instance    Instance
	Bfactory    BeanFactory
	Singleton   bool
	HasRegister bool
}

func (bm *beanMeta) refresh() {
	new := bm.Bfactory()
	if new == nil {
		new = bm.Instance
		return
	}
	vnew := reflect.ValueOf(bm.Bfactory()).Elem()
	vold := reflect.ValueOf(bm.Instance).Elem()
	vold.Set(vnew)
}

var beans = make(map[string]*beanMeta)

func GetBeanByNameDelay(name string, temp interface{}) {
	t := reflect.TypeOf(temp)
	if t.Kind() != reflect.Ptr {
		panic(errors.New("param temp must be the ptr type"))
	}
	if bm, ok := beans[name]; !ok {
		//beanMeta not exist
		bm := &beanMeta{
			Instance: temp,
		}
		beans[name] = bm
		return
	} else {
		//beanMeta has exist
		if bm.HasRegister {
			//has register
			if bm.Singleton {
				//is singleton
				if bm.Instance != nil {
					temp = bm.Instance
				}
			} else {
				//not singleton
				temp = bm.Bfactory()
			}
		} else {
			//has not register
			return
		}
	}
}

func GetBeanByName(name string) Instance {
	if bm, ok := beans[name]; !ok {
		//not register
		errors.Printf(errors.New("bean has not register"))
		return nil
	} else if bm.Singleton {
		//is singleton
		if bm.Instance == nil {
			//hasn't init bean
			instance := bm.Bfactory()
			bm.Instance = instance
		}
		return bm.Instance
	} else {
		//is not singleton
		return bm.Bfactory()
	}
}

//注册一个Singleton Bean，如果已经注册，则panic一个错误
//建议name用全路径引用，如：github.com/vfw4g/base/logger.Logger
func RegisterSingletonBean(name string, bf BeanFactory) {
	if HasRegister(name) {
		//has register
		panic(errors.New("bean has register"))
	} else if BeanExist(name) {
		//not register but delay
		bm := beans[name]
		bm.Bfactory = bf
		bm.refresh()
		bm.HasRegister = true
		bm.Singleton = true
	} else {
		//not register and not delay
		beans[name] = &beanMeta{
			Instance:    bf(),
			Bfactory:    bf,
			Singleton:   true,
			HasRegister: true,
		}
	}
}

//注册一个Singleton Bean，如果已经注册，则panic一个错误
//建议name用全路径引用，如：github.com/vfw4g/base/logger.Logger
func GetRegisterSingletonBean(name string, bf BeanFactory) Instance {
	RegisterSingletonBean(name, bf)
	return GetBeanByName(name)
}

//注册一个Singleton Bean，当且仅当其未注册时才注册，不会抛出任何错误
//建议name用全路径引用，如：github.com/vfw4g/base/logger.Logger
func RegisterSingletonBeanNX(name string, bf BeanFactory) {
	if !BeanExist(name) {
		beans[name] = &beanMeta{
			Instance:  bf(),
			Bfactory:  bf,
			Singleton: true,
		}
	}
	//else do nothing...
}

//注册一个Prototype Bean，如果已经注册，则panic一个错误
//建议name用全路径引用，如：github.com/vfw4g/base/logger.Logger
func RegisterPrototypeBean(name string, bf BeanFactory) {
	if BeanExist(name) {
		//has not register
		panic(errors.New("bean has register"))
	}
	beans[name] = &beanMeta{
		Bfactory: bf,
	}
}

//返回Bean的注册结果，true则已注册
//建议name用全路径引用，如：github.com/vfw4g/base/logger.Logger
func BeanExist(name string) bool {
	if _, hasExist := beans[name]; hasExist {
		return true
	}
	return false
}

func HasRegister(name string) bool {
	if BeanExist(name) && beans[name].HasRegister {
		return true
	}
	return false
}
