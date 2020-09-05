package chanrpc

// CallFunc 调用函数原型
type CallFunc func(interface{}) interface{}

// Callback 回调函数原型
type Callback func(interface{}, error)

// CallInfo 调用信息
type CallInfo struct {
	f       CallFunc
	arg     interface{}
	cb      Callback
	chanRet chan *RetInfo
}

// 返回信息
type RetInfo struct {
	ret interface{}
	err error
	cb  Callback
}

