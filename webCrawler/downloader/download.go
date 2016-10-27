package downloader

import (
	"webCrawler/base"
	"webCrawler/middleware"
	"net/http"
)

type PageDownloader interface {
	Id() uint32
	Download(req base.Request) (*base.Response,error)
}

var downloaderIdGenertor middleware.IdGenertor=middleware.NewIdGenerator()

func GetDownloaderId()	uint32  {
	return downloaderIdGenertor.GetUint32()
}

type myPageDownloader struct {
	id uint32
	httpClient http.Client
}

func NewPageDownloader(client *http.Client) PageDownloader  {
	if client == nil {
		client=&http.Client{}
	}
	id:=GetDownloaderId()
	return &myPageDownloader{
		id:id,
		httpClient:*client,
	}
}
func(m *myPageDownloader) Id() uint32 {
	return m.id
}
func(m *myPageDownloader) Download(req base.Request) (*base.Response,error) {
	httpreq:=req.HttpReq()
	httpresp,err:=m.httpClient.Do(httpreq)
	if err!=nil {
		return nil,err
	}
	return base.NewResponse(httpresp,req.Depth()),nil
}
