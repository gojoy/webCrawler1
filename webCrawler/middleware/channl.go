package middleware

import (
	"webCrawler/base"
	"errors"
	"sync"
	"fmt"
)

var statusNameMap=map[ChannelManagerStatus]string{
	CHANNEL_STATUS_UNINITIALIZED:"uninitialized",
	CHANNEL_STATUS_INITIALZED:"initialized",
	CHANNEL_STATUS_CLOSED:"closed",
}

var defaultlen uint=18

type ChannelManagerStatus uint8

const (
	CHANNEL_STATUS_UNINITIALIZED  ChannelManagerStatus=0
	CHANNEL_STATUS_INITIALZED	ChannelManagerStatus=1
	CHANNEL_STATUS_CLOSED	ChannelManagerStatus=2
)
type ChannelManager interface {
	Init(len base.ChannelArgs,rest bool) bool

	Close() bool

	ReqChan() (chan  base.Request,error)

	RespChan() (chan base.Response,error)

	ItemChan() (chan base.Item,error)

	ErrorChan() (chan error,error)

	Status() ChannelManagerStatus

	Summary() string
}



type myChanManager struct {
	//ChannelLen uint
	reqCh	chan base.Request
	respCh chan base.Response
	itemCh chan base.Item
	errorCh chan error
	status  ChannelManagerStatus
	rwmutex sync.RWMutex
}

func NewmyChanManager(chanLen base.ChannelArgs) ChannelManager  {
	chanmnr:=&myChanManager{}
	chanmnr.Init(chanLen,true)
	return chanmnr
}

func(m *myChanManager) Init(chanLen base.ChannelArgs,reset bool) bool {
	if err:=chanLen.Check();err!=nil {
		panic(errors.New(fmt.Sprintf("the channen len is invlid! check err %s\n",err)))
	}
	if !reset&&m.status==CHANNEL_STATUS_INITIALZED {
		return false
	}
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()
	//m.ChannelLen=chanLen
	m.reqCh=make(chan base.Request,chanLen.ReqChanLen())
	m.respCh=make(chan base.Response,chanLen.RespChanLen())
	m.itemCh=make(chan base.Item,chanLen.ItemChanLen())
	m.errorCh=make(chan error,chanLen.ErrorChanLen())
	m.status=CHANNEL_STATUS_INITIALZED
	return true
}

func(m *myChanManager) Close() bool  {
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()
	if m.status!=CHANNEL_STATUS_INITIALZED {
		return false
	}
	close(m.itemCh)
	close(m.errorCh)
	close(m.respCh)
	close(m.reqCh)
	m.status=CHANNEL_STATUS_CLOSED
	return true
}

func(m *myChanManager) CheckStatus() error  {
	if m.status==CHANNEL_STATUS_INITIALZED  {
		return nil
	}
	statusName,ok:=statusNameMap[m.status]
	if !ok {
		statusName=fmt.Sprintf("%d",m.status)
	}
	errMsg:=fmt.Sprintf("the status means:%s!\n",statusName)
	return errors.New(errMsg)
}

func(m *myChanManager) ReqChan() ( chan base.Request, error)  {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	if err:=m.CheckStatus();err!=nil {
		return nil,err
	}
	return m.reqCh,nil
}

func(m *myChanManager) RespChan() (chan base.Response,error) {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	if err:=m.CheckStatus();err!=nil{
		return nil,err
	}
	return m.respCh,nil
}

func(m *myChanManager) ItemChan() (chan base.Item,error)  {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	if err := m.CheckStatus(); err != nil {
		return nil,err
	}
	return m.itemCh,nil
}

func(m *myChanManager) ErrorChan() (chan error,error) {
	m.rwmutex.RUnlock()
	defer m.rwmutex.RUnlock()
	if err := m.CheckStatus(); err != nil {
		return nil,err
	}
	return m.errorCh,nil
}

func(m *myChanManager) Summary() string  {
	var s string
	s=fmt.Sprintf("the len req is %d,and status is %v\n",len(m.reqCh),statusNameMap[m.status])
	return s
}

func(m *myChanManager) ChanLen() uint  {
	return uint(len(m.reqCh))
}

func(m *myChanManager) Status() ChannelManagerStatus {
	return m.status
}

