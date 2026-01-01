# nsqadmin
## 1.nsqadmin是什么?
- nsqadmin 是 NSQ 的“只读 + 运维控制台”
- 它 不参与消息收发，不存储状态，完全依赖 nsqd / nsqlookupd 提供的数据

本质是:
- 一个 HTTP Server
- 定期向 `lookupd / nsqd` 拉取信息
- 聚合成一个 可视化管理界面


## 2.nsq三件套的关系
```text
Producer ──▶ nsqd ──▶ Consumer
                ▲
                │
           nsqlookupd
                ▲
                │
            nsqadmin
```
⚠️ nsqadmin 不走消息通道，只走 HTTP API


# nsqadmin 命令参数
```text
Usage of nsqadmin:
  -acl-http-header string
        HTTP header to check for authenticated admin users (default "X-Forwarded-User")
  -admin-user value
        admin user (may be given multiple times; if specified, only these users will be able to perform privileged actions; acl-http-header is used to determine the authenticated user)
  -allow-config-from-cidr string
        A CIDR from which to allow HTTP requests to the /config endpoint (default "127.0.0.1/8")
  -base-path string
        URL base path (default "/")
  -config string
        path to config file
  -dev-static-dir string
        (development use only)
  -graphite-url string
        graphite HTTP address
  -http-address string
        <addr>:<port> to listen on for HTTP clients (default "0.0.0.0:4171")
  -http-client-connect-timeout duration
        timeout for HTTP connect (default 2s)
  -http-client-request-timeout duration
        timeout for HTTP request (default 5s)
  -http-client-tls-cert string
        path to certificate file for the HTTP client
  -http-client-tls-insecure-skip-verify
        configure the HTTP client to skip verification of TLS certificates
  -http-client-tls-key string
        path to key file for the HTTP client
  -http-client-tls-root-ca-file string
        path to CA file for the HTTP client
  -log-level value
        日志级别: debug, info, warn, error, or fatal (default INFO)
  -log-prefix string
        log message prefix (default "[nsqadmin] ")
  -lookupd-http-address value
        nsqlookupd地址，可以多次配置(may be given multiple times)
  -notification-http-endpoint string
        HTTP endpoint (fully qualified) to which POST notifications of admin actions will be sent
  -nsqd-http-address value
        nsqdHTTP地址，可以提供多个 (may be given multiple times)
  -proxy-graphite
        proxy HTTP requests to graphite
  -statsd-counter-format string
        The counter stats key formatting applied by the implementation of statsd. If no formatting is desired, set this to an empty string. (default "stats.counters.%s.count")
  -statsd-gauge-format string
        The gauge stats key formatting applied by the implementation of statsd. If no formatting is desired, set this to an empty string. (default "stats.gauges.%s")
  -statsd-interval duration
        time interval nsqd is configured to push to statsd (must match nsqd) (default 1m0s)
  -statsd-prefix string
        prefix used for keys sent to statsd (%s for host replacement, must match nsqd) (default "nsq.%s")
  -verbose
        [deprecated] has no effect, use --log-level
  -version
        输出版本
```