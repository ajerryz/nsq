# nsqlookupd
在 NSQ 分布式消息系统中，nsqlookupd 是负责服务发现的核心组件，扮演着 “目录服务” 的角色。它主要维护 nsqd 节点的元数据（如节点地址、主题信息等），帮助生产者和消费者高效地找到目标 nsqd 节点，从而实现集群的动态扩展和容错。


# 一、nsqlookupd 快速了解
## 1.1 nsqlookupd 的核心定位
nsqlookupd 是 NSQ 集群的 “导航系统”，解决了以下关键问题：
- 当 nsqd 节点动态增减（如扩容、下线）时，消费者如何自动发现新的可用节点？
- 生产者如何知道哪些 nsqd 节点正在处理目标主题（Topic）的消息？

简单来说，nsqlookupd 为整个集群提供了节点发现与元数据管理能力，让其他组件（nsqd、消费者）无需硬编码节点地址，即可动态协同工作。


## 1.2 nsqlookupd 的核心功能
1. 维护 nsqd 节点的元数据

nsqlookupd 会持续跟踪集群中所有活跃的 nsqd 节点，并存储它们的关键信息：
- nsqd 节点的 TCP 地址（用于消费者连接）和 HTTP 地址（用于管理）；
- 该 nsqd 节点上创建的所有主题（Topic） 及关联的通道（Channel）；
- 节点的在线状态（通过心跳机制判断）。

2. 处理 nsqd 节点的注册与心跳
- nsqd 节点启动后，会主动向配置的 nsqlookupd 节点发送注册请求，上报自身地址和当前管理的 Topic 信息。
- 之后，nsqd 会定期（默认每 30 秒）向 nsqlookupd 发送心跳包（含最新的 Topic 信息），证明自己处于活跃状态。
- 如果 nsqlookupd 长时间（默认 60 秒）未收到某个 nsqd 的心跳，会将其标记为 “离线”，并从元数据中移除，避免消费者连接无效节点。

3. 响应消费者的查询请求
- 当消费者启动并希望订阅某个 Topic 时，会向 nsqlookupd 发送查询请求（如 “哪些 nsqd 节点在处理 order_pay 这个 Topic？”）。
- nsqlookupd 会根据自身维护的元数据，返回所有正在处理该 Topic 的活跃 nsqd 节点列表，消费者再根据列表连接对应的 nsqd 节点并订阅消息。

4. 支持生产者的负载均衡（间接）
- 虽然生产者发送消息时通常直接指定 nsqd 节点，但部分场景下（如多节点部署），生产者也可通过查询 nsqlookupd 获取处理目标 Topic 的 nsqd 列表，再通过自定义策略（如轮询）选择节点发送消息，实现负载均衡。

## 1.3 nsqlookupd 的工作机制
1. 无状态设计

nsqlookupd 是无状态的，每个实例独立工作，不与其他 nsqlookupd 节点通信或共享数据。这意味着：
- 可部署多个 nsqlookupd 实例（推荐 3+），通过冗余提高可用性；
- 客户端（nsqd、消费者）只需连接其中一个或多个实例，无需关心集群一致性（因为元数据最终会通过 nsqd 的多节点注册同步）。

2. 数据存储方式

nsqlookupd 仅在内存中维护元数据（不持久化），结构大致如下：
```text
{
  "nsqd节点1": {
    "tcp_address": "192.168.1.10:4150",
    "http_address": "192.168.1.10:4151",
    "topics": ["order_pay", "user_register"],
    "last_heartbeat": 1620000000  // 最后一次心跳时间戳
  },
  "nsqd节点2": { ... }
}
```
当 nsqlookupd 重启时，元数据会丢失，但 nsqd 节点会在重连后重新注册，因此无需持久化。

3. 与其他组件的交互流程

- nsqd → nsqlookupd：
  1. nsqd 启动时，向配置的所有 nsqlookupd 发送 IDENTIFY 命令（注册自身地址和 Topic 信息）；
  2. 之后每 30 秒发送 PING 命令（心跳），更新 Topic 信息和在线状态；
  3. nsqd 退出前，发送 UNREGISTER 命令（可选，主动注销）。
- 消费者 → nsqlookupd：
  1. 消费者启动时，向 nsqlookupd 发送 HTTP 请求（如 http://lookupd:4161/lookup?topic=order_pay）；
  2. nsqlookupd 返回处理 order_pay 的活跃 nsqd 节点列表；
  3. 消费者根据列表，通过 TCP 连接 nsqd 并发送 SUB 命令订阅消息。

## 1.4 nsqlookupd 的配置与启动
nsqlookupd 可通过命令行或配置文件调整行为，核心参数如下：
- `--tcp-address`：TCP 服务地址（默认 0.0.0.0:4160，用于接收 nsqd 的注册和心跳）；
- `--http-address`：HTTP 服务地址（默认 0.0.0.0:4161，用于接收消费者的查询请求）；
- `--broadcast-address`：对外暴露的地址（若 nsqlookupd 部署在容器内，需指定宿主机地址，方便 nsqd / 消费者访问）；
- `--inactive-producer-timeout`：nsqd 节点超时时间（默认 60 秒，超过此时间未收到心跳则标记为离线）。

启动示例:
```shell
# 启动一个 nsqlookupd 实例，指定 TCP 和 HTTP 端口
nsqlookupd --tcp-address=0.0.0.0:4160 --http-address=0.0.0.0:4161
```
高可用部署:
```shell
# 推荐同时启动多个 nsqlookupd 实例（如 3 个），nsqd 和消费者连接时指定所有实例地址，避免单点故障：
# 注意：nsqlookupd 是无状态，内存存储的。
# 高可用关键点在于部署多个nsqlookupd,然后nsqd向多个nsqlookupd注册
# nsqd 关联多个 nsqlookupd
nsqd --lookupd-tcp-address=192.168.1.1:4160 --lookupd-tcp-address=192.168.1.2:4160 ...
```

## 1.5 nsqlookupd 的特点与局限性
特点
- 轻量高效：内存存储元数据，无磁盘 IO，响应速度快；
- 高可用：无状态设计，支持多实例部署，单个实例故障不影响集群；
- 动态适配：自动感知 nsqd 节点的上下线，无需人工干预；
- 简化配置：客户端无需硬编码 nsqd 地址，降低维护成本。

局限性
- 元数据最终一致性：由于多 nsqlookupd 实例独立工作，短期内可能存在元数据不一致（如某 nsqd 刚注册到实例 A，但未同步到实例 B），但会通过 nsqd 心跳逐渐一致；
- 无持久化：重启后元数据丢失，需依赖 nsqd 重新注册（通常几秒内完成，影响较小）；
- 仅支持 NSQ 生态：专为 NSQ 设计，不适合作为通用服务发现工具。

## 1.6 总结
nsqlookupd 是 NSQ 集群的 “分布式目录”，通过维护 nsqd 节点的元数据、处理注册 / 心跳 / 查询，实现了节点的动态发现与协同。其无状态设计和轻量特性，让 NSQ 集群能够轻松扩展到大规模节点，同时保持高可用性和低延迟。理解 nsqlookupd 的工作机制，是掌握 NSQ 分布式消息传递流程的关键。