package chanrpc

type CallFunc func(interface{}) interface{}
type Callback func(interface{}, error)

type CallInfo struct {
	f       CallFunc
	arg     interface{}
	cb      Callback
	chanRet chan *RetInfo
}

type RetInfo struct {
	ret interface{}
	err error
	cb  Callback
}

