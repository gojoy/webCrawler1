package itempipeline

import (
	"webCrawler/base"
	"errors"
	"fmt"
	"sync/atomic"
)

type ItemPipeLine interface {
	Send(i base.Item) []error
	FailFast() bool
	SetFailFast(failFast bool)
	Count() []uint64
	ProcessingNum() uint64
	Summary() string
}

type ProcessItem func(item base.Item) (result base.Item,err error)

type myItemPipeLine struct {
	itemProcessors	[]ProcessItem
	fileFast	bool
	sent	uint64
	accept	uint64
	processed	uint64
	processing 	uint64
}

func NewItemPipeLine(iProcessers []ProcessItem) ItemPipeLine {
	if iProcessers==nil {
		panic("the input processitem is nil\n")
	}
	innerProcessItem:=make([]ProcessItem,0)
	for i,v:=range iProcessers{
		if v==nil{
			panic(errors.New(fmt.Sprintf("the num %d processitem is nil\n",i)))
		}
		innerProcessItem=append(innerProcessItem,v)
	}
	return &myItemPipeLine{itemProcessors:innerProcessItem}
}

func(m *myItemPipeLine) Send(i base.Item) []error {
	atomic.AddUint64(&m.processing,1)
	//递减
	defer atomic.AddUint64(&m.processing,^uint64(0))
	atomic.AddUint64(&m.sent,1)
	errs:=make([]error,0)
	if !i.Valid(){
		errs=append(errs,errors.New("item is nil\n"))
		return errs
	}
	atomic.AddUint64(&m.accept,1)
	var innerItem base.Item=i
	for _,v:=range m.itemProcessors{
		processedItem,err:=v(innerItem)
		if err != nil {
			errs=append(errs,err)
			if m.fileFast {
				break
			}
		}
		if processedItem!=nil{
			innerItem=processedItem
		}
	}
	atomic.AddUint64(&m.processed,1)
	return errs
}
func(m *myItemPipeLine) FailFast() bool {
	return m.fileFast
}
func(m *myItemPipeLine) SetFailFast(failFast bool)  {
	m.fileFast=failFast
}
//返回已发送，已接受和已处理的条目个数
func(m *myItemPipeLine) Count() []uint64 {
	counts:=make([]uint64,3)
	counts[0]=atomic.LoadUint64(&m.sent)
	counts[1]=atomic.LoadUint64(&m.accept)
	counts[3]=atomic.LoadUint64(&m.processed)
	return counts
}
func(m *myItemPipeLine) Summary() string {
	Msg:=fmt.Sprintf("now ItemPipeLine send %d,aceept %d,and processed %d\n",m.Count()[0],m.Count()[1],m.Count()[3])
	return Msg
}
func(m *myItemPipeLine) ProcessingNum() uint64 {
	pn:=atomic.LoadUint64(&m.processing)
	return pn
}