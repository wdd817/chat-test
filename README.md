- 模块拆分
- word filter
    - 多过滤器
    - 过滤器管理
        - 想到的最简单的方式是, 配置指定N个过滤器
        - id 取模来选取过滤器, 防止消息错位
- chat server
- chat session
- chat room
    - 消息统计
    - 消息历史


- 模块
    - 生命周期管理
    - 模块间通信
        - chanrpc
            - 注册消息处理
            - 消息传递 
            - 应该只有异步调用


- 简易chan RPC
    - 所有请求与返回的go struct定义
    - 建立公用req -> resp的映射 (待定)
    - client call server
    - server response


- 用户流程
    - 连接到服务器
    - 输入名字
    - 加入聊天室
    - 获取最近的50条消息
    - 说话
        - 收到消息
        - 敏感词过滤
        - 投递给聊天服务器
    - 命令
        - 收到消息
        - 判断命令
        - 投递给聊天服务器