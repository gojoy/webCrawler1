package middleware

import (
	"sync"
	"fmt"
)

type StopSign interface {
	//发出停止信号
	Sign() bool

	//判断信号是否发出
	Signed() bool

	//重置停止信号 相当于收回停止信号
	Reset()

	//处理停止信号 code代表信号处理方的代号
	Deal(code string)

	DealCount(code string) uint32

	DealTotal() uint32

	Summary() string
}

type myStopSign struct {
	signed bool
	dealCount map[string]uint32
	rwx sync.RWMutex
}

func NewmyStopSign() StopSign {
	dealcount:=make(map[string]uint32)
	return &myStopSign{
		dealCount:dealcount,
	}
}
func(m *myStopSign) Sign() bool {
	m.rwx.Lock()
	defer m.rwx.Unlock()
	if m.signed {
		return false
	}
	m.signed=true
	return true
}
func(m *myStopSign) Signed() bool {
	return m.signed
}

func(m *myStopSign) Deal(code string)  {
	m.rwx.Lock()
	defer m.rwx.Unlock()
	if !m.signed {
		return
	}
	if _, ok := m.dealCount[code]; !ok {
		m.dealCount[code]=1
	}else {
		m.dealCount[code]+=1
	}
}
func(m *myStopSign) Reset()  {
	m.rwx.Lock()
	defer m.rwx.Unlock()
	m.signed= false
	m.dealCount=make(map[string]uint32)
}
func(m *myStopSign) DealCount(code string) uint32  {
	m.rwx.RLock()
	defer m.rwx.RUnlock()
	return m.dealCount[code]
}
func(m *myStopSign) DealTotal() uint32 {
	m.rwx.RLock()
	defer m.rwx.RUnlock()
	var n uint32=0
	for _,v:=range m.dealCount{
		n=n+v
	}
	return n
}
func(m *myStopSign) Summary() string {
	m.rwx.RLock()
	defer m.rwx.RUnlock()
	s:=fmt.Sprintf("the stop status is %v\n",m.signed)
	return s
}