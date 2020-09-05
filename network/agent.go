package network

// Agent 网络代理
type Agent interface {
	Run()
	OnClose()
}
