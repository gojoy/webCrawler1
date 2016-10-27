package base

import (
	"strings"
	"net/http"
	"errors"
	"fmt"
	"net/url"
	"io"
	"github.com/PuerkitoBio/goquery"
	//"webCrawler/analyzer"
)

//type Genp func(resp *http.Response, reth uint32) ([] Data, []error)

func GenParase() func(resp *http.Response, reth uint32) ([] Data, []error) {
	return parseForATag
}

func parseForATag(httpResp *http.Response, respDepth uint32) ([] Data, []error) {
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
	dataList := make([] Data, 0)
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
				req :=  NewRequest(httpReq, respDepth)
				dataList = append(dataList, req)
			}
		}
		text := strings.TrimSpace(sel.Text())
		if text != "" {
			imap := make(map[string]interface{})
			imap["parent_url"] = reqUrl
			imap["a.text"] = text
			imap["a.index"] = index
			item :=  Item(imap)
			// TODO
			//进行输出
			//fmt.Printf("value is %v\n",imap)
			dataList = append(dataList, &item)
		}
	})
	return dataList, errs
}
