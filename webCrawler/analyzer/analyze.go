package analyzer

import (
	"webCrawler/base"
	"webCrawler/middleware"
	"errors"
	"fmt"
	"net/http"
)


type ParseResponse func(httpresp *http.Response,respdepth uint32) ([]base.Data,[]error)

type Analyzer interface {
	Id() uint32
	Analyze(respParasers []ParseResponse,resp base.Response) ([]base.Data,[]error)
}
var AnalyzeIdGenerator middleware.IdGenertor=middleware.NewIdGenerator()

func GetAnalyzerId() uint32  {
	return AnalyzeIdGenerator.GetUint32()
}

type myAnalyzer struct {
	id uint32
}

func NewmyAnalyzer() Analyzer {
	return &myAnalyzer{id:GetAnalyzerId()}
}
func(m *myAnalyzer) Id() uint32 {
	return m.id
}
//依次对传入的 []ParseResponse分析函数进行处理，返回处理数据切片和错误列表切片
//传入参数 分析函数列表和需要分析的响应数据
func(m *myAnalyzer) Analyze(respParasers []ParseResponse, resp base.Response) ([]base.Data, []error) {
	//判断传入函数若为空 报错
	if respParasers==nil {
		return nil,[]error{errors.New("paraser is nil\n")}
	}
	//判断需要处理的响应数据是否为空 报错
	if !resp.Valid() {
		return nil,[]error{errors.New("response is nil\n")}
	}
	//记录URL
	myurl:=resp.HttpResponse().Request.URL
	info:=fmt.Sprintf("log respone url is %v\n",myurl)
	fmt.Println(info)
	respdepth:=resp.Depth()
	//返回数据列表
	dataList:=make([]base.Data,0)
	//返回错误列表
	errorList:=make([]error,0)
	//将处理函数切片分开，每个为单独函数
	httpResp := resp.HttpResponse()

	for i,respParaser:=range respParasers {
		if respParaser==nil {
			errMsg:=fmt.Sprintf("the num %d respParaser is nil!\n",i)
			errorList=append(errorList,errors.New(errMsg))
			continue
		}
		//处理函数respParaser 输入响应返回resp和深度respdepth
		//返回[]base.Data和[]error
		pDataList,pErrorList:=respParaser(httpResp,respdepth)
		if pDataList!=nil {
			for _,pdata:=range pDataList {
				dataList=appendDataList(dataList,pdata,respdepth)
			}
		}
		if pErrorList!=nil {
			for _,perror:=range pErrorList{
				errorList=appendErrorList(errorList,perror)
			}
		}
	}
	return dataList,errorList

}

func appendDataList(datalist []base.Data,data base.Data,respdepth uint32) []base.Data{
	if data==nil {
		return datalist
	}
	 req,ok:=data.(*base.Request)
	if !ok {
		return append(datalist,data)
	}
	newpath:=respdepth+1
	if req.Depth()!=newpath {
		req=base.NewRequest(req.HttpReq(),newpath)
	}
	return append(datalist,req)
}

func appendErrorList(errorlist []error,err error) []error  {
	if err == nil {
		return errorlist
	}
	return append(errorlist,err)
}