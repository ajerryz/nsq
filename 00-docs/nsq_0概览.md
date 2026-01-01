# nsq
NSQ 是一个开源的实时分布式消息平台，设计目标是处理高吞吐量、低延迟的消息传递场景，具备简单、可扩展、容错性强等特点。它广泛应用于日志收集、实时分析、服务解耦等场景。

# 一、快速入门
构建一个小型的NSQ集群

1. 一个shell中启动`nsqlookupd`
```shell
./nsqlookupd
```
2. 启动nsqd
```shell
./nsqd --lookupd-tcp-address=127.0.0.1:4160
```
3. 启动nsqadmin
```shell
./nsqadmin --lookupd-http-address=127.0.0.1:4161
```
4. 发布初始消息
```shell
curl -d 'hello world 1' 'http://127.0.0.1:4151/pub?topic=test'
```
5. 将消息到文件
```shell
nsq_to_file --topic=test --output-dir=/tmp --lookupd-http-address=127.0.0.1:4161
```
6. 发送更多消息
```shell
curl -d 'hello world 2' 'http://127.0.0.1:4151/pub?topic=test'
curl -d 'hello world 3' 'http://127.0.0.1:4151/pub?topic=test'
```
7. 打开Web控制台(`http://127.0.0.1:4171`)

# 二、NSQ 核心架构
NSQ 的架构由三个核心组件构成，各组件职责明确，协同工作实现分布式消息传递：
- nsqd：消息节点（核心组件），负责接收、存储、投递消息，以及与消费者直接交互。
- nsqlookupd：服务发现节点，维护 nsqd 节点的元数据（如节点地址、主题信息），帮助消费者找到存储目标消息的 nsqd 节点。
- nsqadmin：Web 管理界面，用于监控集群状态（如消息量、节点健康度）、管理主题 / 通道等。



# 三、功能与保障
NSQ是一个实时分布式消息平台。

特性:
- 支持分布式拓扑，无单点故障
- 水平可扩展(无需代理，可无缝的向集群添加更多节点)
- 低延迟推送消息传递
- 负载均衡和多播式消息路由相结合
- 既擅长流式(高吞吐量)工作负载，也擅长面向任务(低吞吐量)工作负载
- 主要保存在内存中(超过水位线的消息会透明的保存在磁盘上)
- 运行时服务发现,供消费者查找生产者(nsqlookupd)
- 传输层安全协议(TLS)
- 与数据格式无关
- 依赖少(易于部署)，默认配置合理、范围有限。
- 简单的TCP协议，支持任何语言的客户端库
- 用于统计、管理操作和生产者的HTTP(接口)
- 与statsd 集成以实现实时监测
- 强大的集群管理页面(nsqadmin)

## NSQ的保证
