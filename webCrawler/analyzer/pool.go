package analyzer

import (
	"webCrawler/middleware"
	"reflect"
	"errors"
	"fmt"
)

type AnalyzerPool interface {
	Take() (Analyzer,error)
	Return(ar Analyzer) error
	Total() uint32
	Used() uint32
}

type myAnalyzerPool struct {
	pool middleware.Pool
	etype reflect.Type
}

type GenmyAnalyzer func() Analyzer

func NewmyAnalyzerPool(total uint32, gen GenmyAnalyzer) (AnalyzerPool, error) {
	etype:=reflect.TypeOf(gen())
	genEntity:= func() middleware.Entity{
		return gen()
	}
	pool,err:=middleware.NewPool(total,etype,genEntity)
	if err != nil {
		return nil,err
	}
	anPool:=&myAnalyzerPool{pool:pool,etype:etype}
	return anPool,nil
}

func(m *myAnalyzerPool) Take() (Analyzer,error) {
	entity,err:=m.pool.Take()
	if err != nil {
		return nil,err
	}
	v,ok:=entity.(Analyzer)
	if!ok {
		errMsg:=fmt.Sprintf("the type of analyze entity is %v\n",m.etype)
		panic(errors.New(errMsg))
	}
	return v,nil
}

func(m *myAnalyzerPool) Return(ar Analyzer) error {
	return m.pool.Return(ar)
}
func(m *myAnalyzerPool) Total() uint32 {
	return m.pool.Total()
}
func(m *myAnalyzerPool) Used() uint32 {
	return m.pool.Used()
}
