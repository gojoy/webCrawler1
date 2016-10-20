package middleware

import (
	"sync"
	"math"
)

type IdGenertor interface {
	GetUint32() uint32
}

type myIdGenertor struct {
	sn uint32
	rwx sync.RWMutex
	ended bool
}

func NewIdGenerator()	IdGenertor  {
	return &myIdGenertor{}
}
func(m *myIdGenertor) GetUint32() uint32 {
	m.rwx.Lock()
	defer m.rwx.Unlock()
	if m.ended {
		defer func() {
			m.sn=0
			m.ended=false
		}()
		return m.sn
	}
	id:=m.sn
	if id < math.MaxUint32 {
		m.sn++
	}else {
		m.ended=true
	}
	return id
}