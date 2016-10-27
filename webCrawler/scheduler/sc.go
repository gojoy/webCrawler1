package scheduler

import (
	dl "webCrawler/downloader"
	mdw "webCrawler/middleware"
	 "webCrawler/base"
	anlz "webCrawler/analyzer"
	ipl "webCrawler/itempipeline"
	"net/http"
	"fmt"
	"log"
	"errors"
	"sync/atomic"
	"strings"
	"time"
)
const (
	DOWNLOADER_CODE   = "downloader"
	ANALYZER_CODE     = "analyzer"
	ITEMPIPELINE_CODE = "item_pipeline"
	SCHEDULER_CODE    = "scheduler"
)
var logr *log.Logger=base.GetLogger()

type Scheduler interface {
	Start(chanlen base.ChannelArgs,poolSize base.PoolBaseArgs,
	crawDepth uint32,httpClientGenerator GenHttpClient,
	respParasers []anlz.ParseResponse,
	itemProcessor []ipl.ProcessItem,
	firstRequest *http.Request) (err error)

	Stop()	bool

	Runing()	bool

	ErrorChan() <-chan error

	Idle()	bool

	Summary(pre string)	SchedSummary
}

type GenHttpClient func() *http.Client



type myScheduler struct {
	//poolSize uint32
	poolAgrs	base.PoolBaseArgs
	chanArgs	base.ChannelArgs
	//channelLen uint32
	crawlDepth uint32
	primaryDomain	string
	chanman mdw.ChannelManager
	stopSign mdw.StopSign
	dlpool dl.PageDownloadPool
	analyzerPool anlz.AnalyzerPool
	itemPipeline ipl.ItemPipeLine
	//运行标记 0未运行 1 已运行 2 已停止
	running uint32
	urlMap	map[string]bool
	reqcache requestCache
}

func NewmyScheduler() Scheduler {
	return &myScheduler{}
}


func(m *myScheduler) Start(chanlen base.ChannelArgs,poolSize base.PoolBaseArgs,
crawDepth uint32,httpClientGenerator GenHttpClient,
respParasers []anlz.ParseResponse,
itemProcessor []ipl.ProcessItem,
firstRequest *http.Request) (err error) {
	defer func() {
		if p:=recover();p!=nil{
			errMsg:=fmt.Sprintf("error in scheduler is %s\n",p)
			log.Fatal(errMsg)
			err=errors.New(errMsg)
		}
	}()
	if atomic.LoadUint32(&m.running)==1{
		return errors.New("the scheduler is running")
	}
	atomic.StoreUint32(&m.running,1)
	if err := chanlen.Check(); err != nil {
		return err
	}
	m.chanArgs=chanlen
	if err := poolSize.Check(); err != nil {
		return err
	}
	m.poolAgrs=poolSize
	m.crawlDepth=crawDepth
	m.chanman=genChannelManager(chanlen)

	if httpClientGenerator==nil {
		return errors.New("httpGen is nil\n")
	}
	dlp,err:=genDownloaderPool(m.poolAgrs.PageDownloaderPoolSize(),httpClientGenerator)
	if err!=nil{
		return err
	}
	m.dlpool=dlp
	anlzpool,err:=generateAnalyzerPool(m.poolAgrs.AnalyzerPoolSize())
	if err != nil {
		return err
	}
	m.analyzerPool=anlzpool

	if itemProcessor==nil{
		return errors.New("itempro is nil")
	}
	for i,v:=range itemProcessor{
		if v==nil{
			return errors.New(fmt.Sprintf("the %d num itemprocess is nil\n",i))
		}
	}
	itp:=genItemPipeLine(itemProcessor)
	m.itemPipeline=itp
	if m.stopSign==nil{
		m.stopSign=mdw.NewmyStopSign()
	}else {
		m.stopSign.Reset()
	}
	m.urlMap=make(map[string]bool)
	pd,err:=getPrimaryDomain(firstRequest.Host)
	if err!=nil{
		return err
	}
	m.primaryDomain=pd
	m.reqcache=NewrequestCache()

	//缓存
	cacherequest:=base.NewRequest(firstRequest,0)
	m.reqcache.put(cacherequest)
	m.startDownload()
	m.activeAnalyzer(respParasers)
	m.openItemPipeLine()
	m.schedule(100*time.Millisecond)
	return
}

//开始下载任务
func(m *myScheduler) startDownload()  {
	go func() {
		for  {
			req,ok:=<-m.getReqChan()
			if !ok {
				break
			}
			go m.download(req)
		}
	}()
}
//获取请求通道
func(m *myScheduler) getReqChan() chan base.Request{
	reqChan,err:=m.chanman.ReqChan()
	if err!=nil{
		panic(err)
	}
	return reqChan
}

//获取响应通道
func(m *myScheduler) getRespChan() chan base.Response {
	respChan,err:=m.chanman.RespChan()
	if err != nil {
		panic(err)
	}
	return respChan
}

//获取错误通道
func(m *myScheduler) getErrorChan() chan error{
	errorChan,err:=m.chanman.ErrorChan()
	if err!=nil{
		panic(err)
	}
	return errorChan
}

//获取条目通道
func(m *myScheduler) getItemChan() chan base.Item{
	itemChan,err:=m.chanman.ItemChan()
	if err != nil {
		panic(err)
	}
	return itemChan
}

func(m *myScheduler) download(req base.Request)  {
	defer func() {
		if p := recover(); p != nil {
			errMsg:=fmt.Sprintf("fatal download error %s\n",p)
			log.Fatal(errMsg)
		}
	}()
	downloader,err:=m.dlpool.Take()
	if err!=nil {
		errMsg:=fmt.Sprintf("downloader pool err:%s\n",err)
		m.sendError(errors.New(errMsg),SCHEDULER_CODE)
	}
	defer func() {
		err:=m.dlpool.Return(downloader)
		if err!=nil{
			fmt.Printf("return downloader error %v\n",err)
		}
	}()
	code:=generateCode(DOWNLOADER_CODE,downloader.Id())
	respp,err:=downloader.Download(req)
	if respp!=nil {
		//logr.Println(err)
		m.sendResp(*respp,code)
	}
	if err != nil {
		m.sendError(err,code)
	}
}

//将响应返回给响应通道
func(m *myScheduler) sendResp(resp base.Response,code string) bool  {
	if m.stopSign.Signed() {
		m.stopSign.Deal(code)
		return false
	}
	m.getRespChan()<-resp
	return true
}
func(m *myScheduler) sendItem(item base.Item,code string) bool {
	if m.stopSign.Signed() {
		m.stopSign.Deal(code)
		return false
	}
	m.getItemChan()<-item
	return true
}

func(m *myScheduler) sendError(err error,code string) bool {
	if err==nil{
		return false
	}
	codePrefix:=parseCode(code)[0]
	var errorType base.ErrorType
	switch codePrefix {
	case DOWNLOADER_CODE:
		errorType=base.DOWNLOAD_ERROR
	case ANALYZER_CODE:
		errorType=base.ANALYZER_ERROR
	case ITEMPIPELINE_CODE:
		errorType=base.ITEM_PROCESSER_ERROR
	}
	cError:=base.NewCrawlerError(errorType,err.Error())
	if m.stopSign.Signed() {
		m.stopSign.Deal(code)
		return false
	}
	go func() {
		m.getErrorChan()<-cError
	}()
	return true
}

func(m *myScheduler) activeAnalyzer(respParaser []anlz.ParseResponse)  {
	go func() {
		for   {
			resp,ok:=<-m.getRespChan()
			if !ok {
				break
			}
			go m.analyze(respParaser,resp)
		}
	}()
}
func(m *myScheduler) analyze(respPara []anlz.ParseResponse, resp base.Response)  {
	defer func() {
		if p:=recover();p!=nil {
			errmsg:=fmt.Sprintf("fatal analyzer err %s\n",p)
			log.Fatal(errmsg)
		}
	}()
	analyzer,err:=m.analyzerPool.Take()
	if err!=nil {
		errMsg:=fmt.Sprintf("take analyzer error %s\n",err)
		m.sendError(errors.New(errMsg),SCHEDULER_CODE)
		return
	}
	code:=generateCode(ANALYZER_CODE,analyzer.Id())
	datalist,errlist:=analyzer.Analyze(respPara,resp)
	if datalist!=nil {
		for _,data:=range datalist {
			if data==nil {
				continue
			}
			switch d:=data.(type) {
			case *base.Request:
				m.saveReqToCache(*d,code)
			case *base.Item:
				m.sendItem(*d,code)
			default:
				errMsg:=fmt.Sprintf("the data type %T value is %v unsupport",d,d)
				logr.Printf("no req no item is %v",d)
				m.sendError(errors.New(errMsg),code)
			}
		}
	}
	if errlist!=nil {
		for _,err:=range errlist {
			if err != nil {
				m.sendError(err,code)
			}
		}
	}
}

func(m *myScheduler) saveReqToCache(req base.Request,code string) bool  {
	logr:=base.GetLogger()
	if req.HttpReq()==nil {
		logr.Panicln("save req err:nil HttpReq")
		return false
	}
	if req.HttpReq().URL==nil {
		logr.Println("save req err:nil url")
		return false
	}
	reqUrl:=req.HttpReq().URL

	if  strings.ToLower(reqUrl.Scheme)!="http"{
		logr.Println("save req! err:url not http")
		return false
	}
	if _,ok:=m.urlMap[reqUrl.String()];ok {
		logr.Printf("ingore the req!url is repeated.requrl=%s\n",reqUrl)
		return false
	}
	if pname,_:=getPrimaryDomain(req.HttpReq().Host);pname!=m.primaryDomain {
		logr.Printf("ingore the req!host %s is not in primaryDoamin\n",pname)
		return false
	}
	if req.Depth()>m.crawlDepth {
		logr.Printf("ingore the req!depth %d is bigger than our depth %d",req.Depth(),m.crawlDepth)
		return false
	}
	if m.stopSign.Signed() {
		m.stopSign.Deal(code)
		return false
	}
	m.reqcache.put(&req)
	m.urlMap[reqUrl.String()]=true
	return true
}

//打开条目处理管道
func(m *myScheduler) openItemPipeLine()  {
	go func() {
		m.itemPipeline.SetFailFast(true)
		code:=ITEMPIPELINE_CODE

		//带有缓冲的通道可以使用range来循环读取
		for item:=range m.getItemChan() {
			go func(item base.Item) {
				defer func() {
					if p := recover(); p != nil {
						errMsg:=fmt.Sprintf("fatal item error:%s\n",p)
						log.Fatal(errMsg)
					}
				}()
				errs:=m.itemPipeline.Send(item)
				if errs != nil {
					for _,err:=range errs{
						m.sendError(err,code)
					}
				}
			}(item)
		}
	}()
}

//调度 将请求缓存中的请求放入请求通道中
func(m *myScheduler) schedule(interval time.Duration)  {
	go func() {
		if m.stopSign.Signed() {
			m.stopSign.Deal(SCHEDULER_CODE)
			return
		}
		for  {
			if m.stopSign.Signed() {
				m.stopSign.Deal(SCHEDULER_CODE)
				return
			}
			remainder:=cap(m.getReqChan())-len(m.getReqChan())
			var temp *base.Request
			for remainder>0{
				temp=m.reqcache.get()
				if temp == nil {
					break
				}
				m.getReqChan()<-*temp
				remainder--
			}
			time.Sleep(interval)
		}
	}()
}

//停止调度器
func(m *myScheduler) Stop() bool  {
	if atomic.LoadUint32(&m.running) != 1 {
		return false
	}
	m.stopSign.Sign()
	m.chanman.Close()
	m.reqcache.close()
	atomic.StoreUint32(&m.running,2)
	return true
}

//判断是否在运行，running字段为1则在运行，否则停止
func(m *myScheduler) Runing() bool {
	return atomic.LoadUint32(&m.running)==1
}

func(m *myScheduler) ErrorChan() <-chan error{
	if m.chanman.Status() != mdw.CHANNEL_STATUS_INITIALZED {
		return nil
	}
	return m.getErrorChan()
}

//判断调度器是否空虚 空虚则返回true
func(m *myScheduler) Idle() bool {
	idleDlPool:=m.dlpool.Used()==0
	idleAnlPool:=m.analyzerPool.Used()==0
	idelItemLine:=m.itemPipeline.ProcessingNum()==0
	if idelItemLine&&idleAnlPool&&idleDlPool {
		return true
	}
	return false
}

func(m *myScheduler) Summary(pre string) SchedSummary {
	return NewSchedSummary(m,pre)
}

