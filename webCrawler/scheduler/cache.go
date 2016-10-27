package scheduler

import (
	"webCrawler/base"
	"sync"
	"fmt"
)

type requestCache interface {
	put(req *base.Request) bool

	get() *base.Request

	capacity() int

	length() int

	close()

	summary() string
}

type requestCacheSlice struct {
	cache []*base.Request
	mu sync.RWMutex
	//状态 0运行，1停止
	status int
}
func NewrequestCache() requestCache {
	rc:=&requestCacheSlice{
		cache:make([]*base.Request,0),
		status:0,
	}
	return rc
}

func(m *requestCacheSlice) put(req *base.Request) bool  {
	if req==nil{
		return false
	}
	if m.status == 1 {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache=append(m.cache,req)
	return true
}

func(m *requestCacheSlice) get()*base.Request  {
	if m.status == 1 {
		return nil
	}
	if m.length()==0 {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	rq:=m.cache[0]
	m.cache=m.cache[1:]
	return rq
}
func(m *requestCacheSlice) close()  {
	if m.status==1 {
		return
	}
	m.status=1
}

func(m *requestCacheSlice) length()int  {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.cache)
}

func(m *requestCacheSlice) capacity() int  {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return cap(m.cache)
}

func(m *requestCacheSlice) summary() string {
	return fmt.Sprintf("status is %s,and len is %v,cap is %d\n",statusmap[m.status],m.length(),m.capacity())
}

var statusmap =map[int]string{
	0:"runing",
	1:"close",
}