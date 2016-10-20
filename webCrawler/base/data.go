package base
import "net/http"

type Request struct {
	httpReq *http.Request		//HTTP请求的指针值
	depth  uint32			//请求的深度
}

func NewRequest(hreq *http.Request,dpth uint32) *Request  {
	return &Request{httpReq:hreq,depth:dpth}
}

func(r *Request) HttpReq() *http.Request  {
	return r.httpReq
}

func(r *Request) Depth() uint32  {
	return r.depth
}

type Response struct {
	httpResp http.Response
	depth uint32
}

//响应
func NewResponse(hRespone *http.Response,dpth uint32) *Response  {
	return &Response{httpResp:*hRespone,depth:dpth}
}

func(hr Response) HttpResponse() *http.Response  {
	return &hr.httpResp
}

func(hr Response) Depth() uint32  {
	return hr.depth
}

//条目
type Item map[string]interface{}

type Data interface {
	Valid() bool
}

func(hreq *Request) Valid() bool  {

	//return hreq.httpReq!=nil&&hreq.httpReq.URL!=nil
	return hreq.httpReq.URL!=nil
}

func(hr *Response) Valid() bool {
	//return hr.httpResp!=nil&&hr.httpResp.Body!=nil
	return hr.httpResp.Body!=nil
}

func(i Item) Valid()  bool{
	return i!=nil
}





