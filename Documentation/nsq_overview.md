# nsq
NSQ 是一个开源的实时分布式消息平台，设计目标是处理高吞吐量、低延迟的消息传递场景，具备简单、可扩展、容错性强等特点。它广泛应用于日志收集、实时分析、服务解耦等场景。

# 一、NSQ 核心架构
NSQ 的架构由三个核心组件构成，各组件职责明确，协同工作实现分布式消息传递：
- nsqd：消息节点（核心组件），负责接收、存储、投递消息，以及与消费者直接交互。
- nsqlookupd：服务发现节点，维护 nsqd 节点的元数据（如节点地址、主题信息），帮助消费者找到存储目标消息的 nsqd 节点。
- nsqadmin：Web 管理界面，用于监控集群状态（如消息量、节点健康度）、管理主题 / 通道等。

# 默认端口情况
```text

nsqadmin( --http-address  default http:4171)
   --lookupd-http-address
   --nsqd-http-address 
    
nsqlookupd                                      nsqd    (--http-address  default:tcp:4150        --http-address default:http:4151)
--tcp-address    tpc:4160                   <----       --lookupd-tcp-address    
--http-address   http:4161
```

# 二、其他工具
- `nsq_stat`:令行工具，用于查看 NSQ 集群的统计信息。
- `nsq_tail`:命令行工具，用于实时查看指定主题 / 通道的消息内容（类似 tail -f）。
- `nsq_to_file`:消息导出工具，用于将指定主题的消息持久化到本地文件。
- `nsq_to_http`:消息转发工具，用于将消息发送到 HTTP 服务端点。
- `nsq_to_nsq`:跨集群消息同步工具，用于将一个 NSQ 集群的消息转发到另一个 NSQ 集群。
- `to_nsq`:消息生产工具，用于向 NSQ 主题发送消息（从标准输入或文件读取）。
- `nsqadmin`: NSQ 的Web 管理界面，用于监控和管理整个 NSQ 集群。
- `nsqd`:NSQ 的核心消息代理节点，负责接收、存储和转发消息。
- `nsqlookupd`:NSQ 的服务发现组件，负责管理集群拓扑信息。