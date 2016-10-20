package middleware

import (
	"reflect"
	"fmt"
	"errors"
	"sync"
)

type Entity interface {
	Id() uint32
}

type Pool interface {
	Take() (Entity,error)
	Return(er Entity) error
	Total() uint32
	Used() uint32
}

type myPool struct {
	total uint32
	etype reflect.Type
	genEntity	func()	Entity
	container chan Entity
	continerId	map[uint32]bool
	rwx sync.RWMutex
}

func NewPool(tol uint32,entityType reflect.Type,genEntity func() Entity) (Pool,error)  {
	if tol==0 {
		errMsg:=fmt.Sprintf("the pool cannot initialized cas total is %d\n",tol)
		return nil,errors.New(errMsg)
	}
	size:=int(tol)
	coner:=make(chan Entity,size)
	conerId:=make(map[uint32]bool)
	for i := 0; i < size; i++ {
		newEntity:=genEntity()
		if entityType!=reflect.TypeOf(newEntity) {
			errMsg:=fmt.Sprintf("the type of func genE is not %s\n",entityType)
			return nil,errors.New(errMsg)
		}
		conerId[newEntity.Id()]=true
		coner<-newEntity
	}
	pool:=&myPool{
		total:tol,
		etype:entityType,
		genEntity:genEntity,
		container:coner,
		continerId:conerId,
	}
	return pool,nil
}

func(mp *myPool) Take() (Entity,error)  {
	entity,ok:=<-mp.container
	if !ok {
		return nil,errors.New("the inner container is invalid\n")
	}
	mp.rwx.Lock()
	defer mp.rwx.Unlock()
	mp.continerId[entity.Id()]=false
	return entity,nil
}

func(m *myPool) Return(er Entity) error  {
	if er==nil {
		return errors.New("cannot return nil\n")
	}
	if m.etype != reflect.TypeOf(er) {
		errMsg:=fmt.Sprintf("the runing is not %s\n",m.etype)
		return errors.New(errMsg)
	}

	m.rwx.Lock()
	 v,ok:=m.continerId[er.Id()]
	if ok&&!v {
		m.continerId[er.Id()]=true
	}
	m.rwx.Unlock()
	if !ok{
		errMsg:=fmt.Sprintf("the is %d is illeged\n ",er.Id())
		return errors.New(errMsg)
	}
	if v {
		errMsg:=fmt.Sprintf("the id %d has return\n",er.Id())
		return errors.New(errMsg)
	}

	m.container<-er

	return nil
}

func(m *myPool) Total() uint32  {
	return m.total
}

func(m *myPool) Used() uint32 {
	var num uint32=0
	for _,v:=range m.continerId{
		if v {
			num++
		}
	}
	return num
}