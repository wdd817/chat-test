### chat practise
#### 项目简介
- 运行环境为Go1.5
- 项目入口就是main.go
- 分三个模块
    - gate模块, 承载连接
    - word filter, 敏感词过滤
    - chat, 聊天室
- 每个模块单独goroutine
- 第三方库
    - github.com/gorilla/websocket, websocket功能
    - github.com/gookit/config 用于读取yaml

#### 测试方法
- go run test/test.go
    - 两个方式
        - 终端
        - 网页端, http://localhost:8888

#### 项目目录说明:
- app 单例, 管理module
- chanrpc 实现了单机下, 基于chan的rpc
- chat 聊天室, 主要是用来存储消息的
- conf 全局配置变量
- gate 承载连接
- log 简易的logger
- module 定义了module接口, 与一个集成了chanrpc的skeleton
- network 网络相关
- proto 内部rpc通信协议
- util 打印堆栈
- wordfilter 过滤模块
