# 05 当前状态、验证结果、我的备注

## 1. 本地验证结果

我实际运行了：

```bash
cd /Users/haoniu/sh18/aip2p.com/AiP2P
go test ./...
```

结果：

- 全量测试通过

我也实际运行了：

```bash
go run ./cmd/aip2p serve --listen 127.0.0.1:1818
```

启动日志显示：

- 当前宿主实际装配的是 `news-demo-content+news-demo-governance+news-demo-archive+news-demo-ops`
- theme 是 `news-demo`

## 2. 我实际观察到的本地运行态

在这台机器上，当前 demo app 不是空数据：

- `/api/feed` 返回了本地已有帖子
- `/api/history/list` 返回了 7 条 manifest entry
- `/api/network/bootstrap` 返回了当前节点的 peer id、dial addrs、bittorrent nodes

页面侧还能看到：

- 本地 writer policy 文件路径在 `~/.aip2p-sharing/`
- network 页面显示 sync heartbeat 是旧的
- 说明当前宿主能展示同步状态，但同步 worker 最近一次状态写回不是刚发生

## 3. 默认运行目录理解

`newsdemo/runtimepaths.go` 说明默认 runtime 根是：

- `~/.aip2p-sharing`

重要文件和目录：

- `~/.aip2p-sharing/aip2p/.aip2p`
- `~/.aip2p-sharing/archive`
- `~/.aip2p-sharing/subscriptions.json`
- `~/.aip2p-sharing/writer_policy.json`
- `~/.aip2p-sharing/aip2p_sharing_net.inf`
- `~/.aip2p-sharing/Trackerlist.inf`
- `~/.aip2p-sharing/identities`
- `~/.aip2p-sharing/delegations`
- `~/.aip2p-sharing/revocations`

所以当前 demo app 的运行是明显“本地节点中心”的，不是无状态 web app。

## 4. 我认为需要记住的边界

### 不是完整论坛产品

仓库多处明确说自己不定义：

- ranking
- moderation
- truth scoring 标准
- 单一 UI

虽然 demo 页面里已经有 `Vote Score`、`Truth`、`Source Quality` 这些排序口径，但这些仍然是 demo 层，而不是协议强约束。

### 是可运行参考实现，不是纸上协议

代码里已经具备：

- 消息打包
- 验签
- BT 下载与 seed
- libp2p bootstrap / rendezvous / pubsub
- 本地 Web 宿主
- app/plugin/theme 扩展系统

所以它已经足够支撑一个小型下游项目原型。

### 根目录 `doc/index.html` 目前是规划页

这个文件是中文手册站点规划，不等于当前实现说明书。

我会把它视为：

- 对外文档网站的未来结构草稿

而把真实事实来源优先放在：

- `AiP2P/README.md`
- `AiP2P/docs/`
- `AiP2P/internal/*`

## 5. 我自己的后续阅读建议

如果下一步要继续深入，我建议按这个顺序：

1. 读 `internal/plugins/newsdemo/server.go` 和 `NewWithThemeAndOptions()`，把 App 对象初始化链再走一遍。
2. 跑一次 `aip2p publish`，从身份文件到 bundle 落盘完整跟一遍。
3. 跑一次 `aip2p sync --once`，把 queue、manifest、announcement 实际走一遍。
4. 做一个最小第三方 app workspace，验证 app/plugin/theme 的替换边界。

## 6. 一句话总结

AiP2P 当前最准确的理解不是“一个论坛网站”，而是“一个面向 AI agent 的明文 P2P 消息协议，加上一套可直接跑起来的本地节点宿主与 demo app 参考实现”。
