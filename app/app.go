package app

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"sync/atomic"
	"syscall"

	"chat-test/chanrpc"
	"chat-test/log"
	"chat-test/module"
	"chat-test/util"
)

// 节点全局状态
const (
	StateNone = iota // 未开始或已停止
	StateInit        // 正在初始化中
	StateRun         // 正在运行中
	StateStop        // 正在停止中
)

// 单例
var singleton = NewApp()

// mod 模块
type mod struct {
	mi       module.Module
	closeSig chan bool
	wg       sync.WaitGroup
}

// Instance 默认单例
func Instance() *App {
	return singleton
}

// NewApp 创建App
func NewApp() *App {
	app := &App{
		closeSig: make(chan os.Signal, 1),
		state:    StateNone,
	}
	app.wg.Add(1)
	return app
}

// App .
// 有两种启停方式:
// 	1. Start -> Stop: 手动启动和停止app，比较干净，通常用于测试代码
//  2. Run -> Terminate: 基于Start/Stop封装，自动监听OS Signal或通过Terminate来终止，通常用于真正的节点启动流程
type App struct {
	mods     []*mod
	state    int32
	closeSig chan os.Signal
	wg       sync.WaitGroup
}

// SetState 设置状态
func (app *App) setState(s int32) {
	atomic.StoreInt32(&app.state, s)
}

// GetState 获取状态
func (app *App) GetState() int32 {
	return atomic.LoadInt32(&app.state)
}

// Start 非阻塞启动app，需要在当前goroutine调用Stop来停止app
func (app *App) Start(mods ...module.Module) {
	// 单个app不能启动两次
	if app.GetState() != StateNone {
		log.Fatal("app mods cannot start twice")
	}
	if len(mods) == 0 {
		return
	}
	// 注册module 并增加开关
	log.Info("app starting up")
	// register
	for _, mi := range mods {
		m := new(mod)
		m.mi = mi
		m.closeSig = make(chan bool, 1)
		app.mods = append(app.mods, m)
	}
	app.setState(StateInit)
	// 模块初始化
	for _, m := range app.mods {
		mi := m.mi
		if err := mi.OnInit(); err != nil {
			log.Fatal("module %v init error %v", reflect.TypeOf(mi), err)
		}
	}
	// 模块启动
	for _, m := range app.mods {
		m.wg.Add(1)
		go run(m)
	}
	app.setState(StateRun)
}

// Stop 停止App
func (app *App) Stop() {
	if app.GetState() == StateStop {
		return
	}
	app.setState(StateStop)
	// 先进后出
	for i := len(app.mods) - 1; i >= 0; i-- {
		m := app.mods[i]
		m.closeSig <- true
		m.wg.Wait()
		destroy(m)
	}
	app.wg.Done()
	app.setState(StateNone)
}

// Stats 简单查看各个 Service Chan 中的请求个数
func (app *App) Stats() string {
	var ret string
	for _, m := range app.mods {
		ret += fmt.Sprintf("chan: %v, len: %v \r\n", m.mi.Name(), len(m.mi.ChanRPC().ChanCall))
	}
	return ret
}

// GetChanRPC 获取指定名字模块的消息投递通道
func (app *App) GetChanRPC(name string) *chanrpc.Server {
	for _, m := range app.mods {
		if m.mi.Name() == name {
			return m.mi.ChanRPC()
		}
	}
	return nil
}

func run(m *mod) {
	m.mi.Run(m.closeSig)
	m.wg.Done()
}

func destroy(m *mod) {
	defer util.PrintPanicStack()
	m.mi.OnDestroy()
}

// Run 阻塞启动app
func (app *App) Run(mods ...module.Module) {
	app.Start(mods...)
	for {
		signal.Notify(app.closeSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		sig := <-app.closeSig
		log.Info("server closing down (signal: %v)", sig)
		if sig == syscall.SIGHUP {
			continue
		}
		break
	}
	app.Stop()
}

// Terminate 用于模拟信号，终止Run，并等待app停止完成
func (app *App) Terminate() {
	if app.GetState() != StateRun {
		return
	}
	app.closeSig <- syscall.SIGUSR1
	app.wg.Wait()
}
