package base

type ErrorType string

type CrawlerError interface {
	Error() string
	Type()	ErrorType
}

type myCrawlerError struct {
	errType  ErrorType
	errMsg	string
	fullMsg	string
}

const (
	DOWNLOAD_ERROR	ErrorType ="Download error"
	ANALYZER_ERROR	ErrorType="Analyze error"
	ITEM_PROCESSER_ERROR	ErrorType="Item_processer_error"
)

func NewCrawlerError(errType ErrorType,errMsg string) CrawlerError  {
	return &myCrawlerError{errType:errType,errMsg:errMsg}
}

func(ce *myCrawlerError) Type() ErrorType {
	return ce.errType
}

func(ce *myCrawlerError) Error() string {
	return ce.errMsg
}
