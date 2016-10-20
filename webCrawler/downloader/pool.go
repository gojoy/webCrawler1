package downloader

import (
	"webCrawler/middleware"
	"reflect"
	"fmt"
	"errors"
)

type PageDownloadPool interface {
	Take() (PageDownloader,error)
	Return(dr PageDownloader) error
	Total() uint32
	Used() uint32
}

type myDownloaderPool struct {
	pool middleware.Pool
	etype reflect.Type
}

type GenPageDownloader func() PageDownloader

func NewmyDownloaderPool(total uint32,gen GenPageDownloader) (PageDownloadPool,error) {
	etype:=reflect.TypeOf(gen())
	genEntity:= func() middleware.Entity{
		return gen()
	}
	pool,err:=middleware.NewPool(total,etype,genEntity)
	if err != nil {
		return nil,err
	}
	depool:=&myDownloaderPool{pool:pool,etype:etype}
	return depool,nil
}
func(m myDownloaderPool) Take() (PageDownloader,error) {
	entity,err:=m.pool.Take()
	if err != nil {
		return nil,err
	}
	dl,ok:=entity.(PageDownloader)
	if !ok {
		errMsg:=fmt.Sprintf("the type of entity is not %s\n",m.etype)
		panic(errors.New(errMsg))
	}
	return dl,nil
}
func(m myDownloaderPool) Return(dr PageDownloader) error {
	return  m.pool.Return(dr)
}
func(m *myDownloaderPool) Total() uint32 {
	return m.pool.Total()
}
func(m *myDownloaderPool) Used() uint32 {
	return m.pool.Used()
}


