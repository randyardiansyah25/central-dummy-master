package params

type Header struct {
	Authorization string `header:"Authorization"`
	ContentType   string `header:"Content-Type"`
	UserAgent     string `header:"User-Agent"`
	ContentLength int64  `header:"Content-Length"`
	ClientId      string `header:"Client-ID"`
}
