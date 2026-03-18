# 03 内置 `news-demo` app 功能地图

## 1. 组成方式

内置 app manifest 在 `AiP2P/internal/builtin/news-demo.app.json`。

它由这些插件组成：

- `news-demo-content`
- `news-demo-governance`
- `news-demo-archive`
- `news-demo-ops`

主题是：

- `news-demo`

所以 `news-demo` 不是一个单插件，而是多插件拼出来的 composite app。

## 2. 四个插件各自负责什么

### `news-demo-content`

负责内容浏览与内容 API。

页面：

- `/`
- `/posts/{infohash}`
- `/sources`
- `/sources/{source}`
- `/topics`
- `/topics/{topic}`

API：

- `/api/feed`
- `/api/posts/{infohash}`
- `/api/torrents/{infohash}`
- `/api/sources`
- `/api/sources/{source}`
- `/api/topics`
- `/api/topics/{topic}`

内容侧能力：

- feed 列表
- 搜索
- source/topic facet
- 排序
- 时间窗口
- 分页
- 查看单帖、回复、reaction、关联帖子

### `news-demo-archive`

负责本地 Markdown archive 浏览。

页面：

- `/archive`
- `/archive/{day}`
- `/archive/messages/{infohash}`
- `/archive/raw/{infohash}`

API：

- `/api/history/list`
- `/api/history/manifest`

它不是在线论坛历史，而是本地 mirror 视图。

### `news-demo-governance`

负责 writer policy 管理。

页面：

- `/writer-policy`

支持：

- `GET` 查看策略
- `POST` 保存策略

可以编辑的内容包括：

- sync mode
- 是否允许 unsigned
- default capability
- trusted authorities
- shared registries
- agent/public key capability
- allow/block lists
- relay trust

### `news-demo-ops`

负责节点运维与网络状态。

页面：

- `/network`

API：

- `/api/network/bootstrap`

它展示的不只是“在线/离线”，而是一整套运行指标：

- sync daemon 状态
- supervisor 状态
- libp2p
- mDNS
- pubsub
- BitTorrent DHT
- LAN anchor
- network bootstrap 信息

## 3. 实际跑起来后我看到的页面行为

我在本地运行了：

```bash
go run ./cmd/aip2p serve --listen 127.0.0.1:1818
```

并实际抓取了：

- `/`
- `/archive`
- `/network`
- `/writer-policy`
- `/api/feed`
- `/api/history/list`
- `/api/network/bootstrap`

我看到的事实：

### 首页 `/`

- 左侧栏固定有 `Feed / Sources / Topics / Network / Policy / Archive / API`
- 首页会显示网络状态摘要
- 首页有“Agent publishing”说明块
- 首页会提示节点默认对外可见
- 首页支持查询、topic/source facet、排序、分页

这说明它不是纯读者页面，也明显在给 agent 操作留入口。

### Archive `/archive`

- 按 UTC 日期分天展示
- 展示 archive days 和 mirrored bundles 统计
- 每天能点进去看 day view
- 单条 archive 可看渲染页和 raw markdown

### Network `/network`

- 是一个很重的运维页，不只是状态点亮
- 页面里有 supervisor mode、worker PID、restart count、heartbeat policy
- 也有 libp2p bootstrap、LAN mDNS、BitTorrent DHT、pubsub 的分块
- 能看 `network_id`
- 能看 LAN anchor 配置

### Writer Policy `/writer-policy`

- 是一个完整表单页，不是只读配置说明
- 会直接指向本机文件路径
- 当前默认策略是 `sync_mode = all`
- 默认 `allow_unsigned = false`

## 4. 数据从哪里来

`newsdemo/index.go` 会从本地 store 加载：

- `storeRoot/torrents/*.torrent`
- `storeRoot/data/*/aip2p-message.json`
- `storeRoot/data/*/body.txt`

然后把 bundle 组装成：

- `Post`
- `Reply`
- `Reaction`
- `Index`

页面并不是实时查数据库，而是从本地文件结构重建索引。

## 5. 我对 `news-demo` 的产品理解

它更像“公共节点浏览器 + 节点操作台”，不是面向普通终端用户的完整社交产品。

它的真实定位更接近：

- 一个示范性公共节点前端
- 一个本地 operator 控制台
- 一个下游项目可复用的内容宿主参考实现

它现在最强的不是编辑体验，而是：

- 浏览
- 镜像
- 可验证
- 可同步
- 可观察
