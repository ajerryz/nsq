# 安装
- [二进制脚本安装](https://nsq.io/deployment/installing.html)
- Docker
- Mac: `brew install nsq`
- 源码构建
  - `git clone https://github.com/nsqio/nsq`
  - `cd nsq & make`

# 简单测试拓扑
nsqlookupd(1) + nsqd(1) + nsqadmin(1)

1. 启动nsqlookupd: `./nsqlookupd`
   - tcp: 4160 ,nsqd注册/心跳/元数据同步
   - http: 4161,查询 topic/node/nsqd列表(管理与发现)
2. 启动nsqd: `./nsqd --lookupd-tcp-address localhost:4160`
   - tcp: 4150 ,真正的消息收发(Producer/Consumer)
   - http: 4151 , 管理/监控/运维API
3. 启动nsqadmin: `./nsqadmin -lookupd-http-address localhost:4161`
   - http: 4171 (web控制台)
4. 启动nsqadmin直接监控nsqd: `./nsqadmin -nsqd-http-address localhost:4151`
   - http: 4171(web控制台)

运维端口总结:
```text
nsqadmin ──HTTP(4161)──▶ nsqlookupd
nsqadmin ──HTTP(4151)──▶ nsqd
```

# 简单生产部署拓扑