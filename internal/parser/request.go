package parser

type RequestHeaders struct {
	Method string
	HttpVersion string
	Path string
	Headers map[string]string
}

type Request struct {
	Headers RequestHeaders
	Body string
}
