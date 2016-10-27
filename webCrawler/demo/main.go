package main

import (
	"fmt"
	"webCrawler/base"
	"webCrawler/scheduler"
	"errors"
	"time"
	//"net/url"
	//"io"
	//"github.com/PuerkitoBio/goquery"
	"net/http"
	//"strings"
	//"webCrawler/analyzer"
	"webCrawler/itempipeline"
	"log"
	"webCrawler/analyzer"
)
/*
type Genp func(resp *http.Response, reth uint32) ([]base.Data, []error)

func GenParase() Genp {
	return parseForATag
}

func parseForATag(httpResp *http.Response, respDepth uint32) ([]base.Data, []error) {
	// TODO 支持更多的HTTP响应状态
	if httpResp.StatusCode != 200 {
		err := errors.New(
			fmt.Sprintf("Unsupported status code %d. (httpResponse=%v)", httpResp))
		return nil, []error{err}
	}
	var reqUrl *url.URL = httpResp.Request.URL
	var httpRespBody io.ReadCloser = httpResp.Body
	defer func() {
		if httpRespBody != nil {
			httpRespBody.Close()
		}
	}()
	dataList := make([]base.Data, 0)
	errs := make([]error, 0)
	// 开始解析
	doc, err := goquery.NewDocumentFromReader(httpRespBody)
	if err != nil {
		errs = append(errs, err)
		return dataList, errs
	}
	// 查找“A”标签并提取链接地址
	doc.Find("a").Each(func(index int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		// 前期过滤
		if !exists || href == "" || href == "#" || href == "/" {
			return
		}
		href = strings.TrimSpace(href)
		lowerHref := strings.ToLower(href)
		// 暂不支持对Javascript代码的解析。
		if href != "" && !strings.HasPrefix(lowerHref, "javascript") {
			aUrl, err := url.Parse(href)
			if err != nil {
				errs = append(errs, err)
				return
			}
			if !aUrl.IsAbs() {
				aUrl = reqUrl.ResolveReference(aUrl)
			}
			httpReq, err := http.NewRequest("GET", aUrl.String(), nil)
			if err != nil {
				errs = append(errs, err)
			} else {
				req := base.NewRequest(httpReq, respDepth)
				dataList = append(dataList, req)
			}
		}
		text := strings.TrimSpace(sel.Text())
		if text != "" {
			imap := make(map[string]interface{})
			imap["parent_url"] = reqUrl
			imap["a.text"] = text
			imap["a.index"] = index
			item := base.Item(imap)
			dataList = append(dataList, &item)
		}
	})
	return dataList, errs
}
*/
var parseForATag analyzer.ParseResponse=base.GenParase()

func getResponseParsers() []analyzer.ParseResponse {
	parsers := []analyzer.ParseResponse{
		parseForATag,
	}
	return parsers
}


// 获得条目处理器的序列。
func getItemProcessors() []itempipeline.ProcessItem {
	itemProcessors := []itempipeline.ProcessItem{
		processItem,
	}
	return itemProcessors
}

// 生成HTTP客户端。
func genHttpClient() *http.Client {
	return &http.Client{}
}


func processItem(item base.Item) (result base.Item, err error) {
	if item == nil {
		return nil, errors.New("Invalid item!")
	}
	// 生成结果
	result = make(map[string]interface{})
	for k, v := range item {
		result[k] = v
		//fmt.Printf("get item is %v,value is %v\n",k,v)
	}
	if _, ok := result["number"]; !ok {
		result["number"] = len(result)
	}
	time.Sleep(10 * time.Millisecond)
	return result, nil
}


func main() {
	channelArgs := base.NewChannelArgs(10, 10, 10, 10)
	poolBaseArgs := base.NewPoolBaseArgs(4, 4)
	crawlDepth := uint32(3)
	httpClientGenerator := genHttpClient
	respParsers := getResponseParsers()
	itemProcessors := getItemProcessors()
	startUrl := "http://www.qq.com"
	firstHttpReq, err := http.NewRequest("GET", startUrl, nil)
	if err != nil {
		log.Panicln(err)
		return
	}
	scheduler:=scheduler.NewmyScheduler()
	scheduler.Start(
		channelArgs,
		poolBaseArgs,
		crawlDepth,
		httpClientGenerator,
		respParsers,
		itemProcessors,
		firstHttpReq)

time.Sleep(5*time.Second)

	fmt.Printf("this is main\n")
}
