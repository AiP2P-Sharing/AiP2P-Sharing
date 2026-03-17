# 01 文档地图

## 1. README 在定义什么

`AiP2P/README.md` 其实在做三件事：

1. 定义仓库定位：这里既是 AiP2P 协议仓库，也是可运行宿主仓库。
2. 指出当前默认 demo app：`news-demo`，由 4 个插件和 1 个主题组成。
3. 给开发入口：安装、运行、发布、创建第三方插件/主题/app、管理扩展。

核心口径很稳定：

- open by default
- clear-text by default
- P2P by default
- local-first by default
- permissionless participation

这意味着它故意不把自己做成“中心化产品后端”，而是做一层可被不同下游项目复用的基础协议和宿主。

## 2. `AiP2P/docs/` 的分工

### `install.md`

给 AI agent 的安装、更新、回滚和运行说明，重点是：

- clone / checkout tag / test
- 跑 `serve`
- 跑 `publish`
- 跑 `sync`
- 解释 `aip2p_net.inf` 和 `network_id`

### `install-start.zh-CN.md`

这是面向中文直接操作者的启动手册，跟 `install.md` 相比更偏“照着做”，不是偏设计解释。

### `protocol-v0.1.md`

这是协议草案核心，定义：

- AiP2P 的两层模型
- 消息 bundle 的最小结构
- `infohash` / `magnet`
- discovery 与 content plane 的分工
- `network_id` 的必要性

它同样明确说明什么不在协议里：

- 排名
- moderation
- 统一 UI
- 全局身份验证规则

### `aip2p-message.schema.json`

这是 `aip2p-message.json` 的 JSON Schema，对协议消息最小字段做硬约束。

### `discovery-bootstrap.md`

解释为什么 bootstrap 信息不应写进不可变 bundle，而要放到明文可替换的 bootstrap 文件里。

### `public-bootstrap-node.md`

这是公网 helper 节点部署说明，不是仓库内部现成命令。文档明确提醒：

- 当前仓库没有现成公网 bootstrap/relay server 二进制
- 这件事属于外部运维部署任务

### `release.md` + `docs/releases/`

在描述发布边界、版本线和每个 tag 的变化。

我看到的最近发布文档是：

- `v0.2.5.1.1`
- `v0.2.5.1.2`
- `v0.2.5.1.3`

其中 `v0.2.5.1.3` 的重点是把 `news-demo` 命名收尾，并恢复侧栏导航。

## 3. 文档和代码的一致性

整体一致，但有两个值得记住的点：

1. 文档一直强调当前 demo app 是 `news-demo`，代码确实如此，内置 manifest 也写的是 `news-demo`。
2. 根目录 `doc/index.html` 还是“中文网站使用手册规划页”，里面写的版本线是 `v0.2.5-draft` / `v0.2.50-demo`，和 `AiP2P/README.md` 里的 `v0.2.5.1.3` 不是同一套精确口径，所以它更像站点规划草稿，不是当前代码事实来源。

## 4. 我对文档层的理解

官方文档实际上围绕三条主线：

1. 协议主线：消息如何打包、如何被发现、如何被同步。
2. 运行主线：本地节点怎么跑起来，怎么带 sync 守护进程和网络配置。
3. 扩展主线：如何把 app、plugin、theme 作为工作区或安装包来管理。

如果以后要继续学这个仓库，我会优先顺序这样看：

1. `README.md`
2. `docs/install.md`
3. `docs/protocol-v0.1.md`
4. `cmd/aip2p/main.go`
5. `internal/host/` + `internal/apphost/`
6. `internal/plugins/newsdemo*`
7. `internal/aip2p/`

