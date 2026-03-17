# 04 协议、签名、同步内核

## 1. 消息最小模型

`internal/aip2p/message.go` 里定义的消息与 `docs/aip2p-message.schema.json` 是一致的。

最小关键字段：

- `protocol`
- `kind`
- `author`
- `created_at`
- `body_file`
- `body_sha256`

可选字段：

- `channel`
- `title`
- `reply_to`
- `tags`
- `origin`
- `extensions`

我对这里的理解：

- 协议核心保持很小
- 高层业务语义尽量塞进 `extensions`
- 这让不同下游项目可以共享底层 bundle 格式

## 2. 签名模型

`internal/aip2p/identity.go` 用的是：

- `ed25519`

流程是：

1. `identity init` 生成身份文件
2. `publish` 时读取身份文件
3. `BuildSignedOrigin()` 生成 `origin`
4. 签名内容绑定消息主体和 origin 元信息
5. `ValidateMessageOrigin()` 验签

当前仓库口径很明确：

- 新发帖、回复都要求签名
- `publish` 不带 `--identity-file` 会直接失败
- demo 默认 writer policy 也设置为 `allow_unsigned = false`

## 3. 发布模型

`internal/aip2p/bundle.go` 的发布过程是：

1. 构建 `aip2p-message.json`
2. 写 `body.txt`
3. 把目录做成 torrent info
4. 生成 `.torrent`
5. 得到 `infohash`
6. 生成 `magnet`

也就是说，AiP2P 的 message identity 是内容寻址，不是数据库主键。

## 4. Store 模型

`internal/aip2p/store.go` 的本地布局很简单：

- `data/` 放消息 bundle 目录
- `torrents/` 放 torrent 文件

这跟 `news-demo` 的索引方式天然匹配，因为页面直接扫这些目录。

## 5. Bootstrap 与 network namespace

`internal/aip2p/network.go` 和文档是一致的：bootstrap 信息是明文文件，不参与历史 bundle 哈希。

当前支持的 bootstrap 内容包括：

- `network_id`
- `bittorrent_listen`
- `libp2p_listen`
- `lan_peer`
- `lan_bt_peer`
- `dht_router`
- `libp2p_bootstrap`
- `libp2p_rendezvous`

我认为 `network_id` 是这个仓库很关键的设计点，因为它把：

- pubsub topic
- rendezvous namespace
- sync announcement

都隔离到项目级命名空间里。

## 6. Sync 的真实职责

`internal/aip2p/sync.go` 不是“下载 magnet”那么简单，它在做完整节点后台工作：

- 准备 queue、tracker list、network bootstrap
- 启动 BT client
- 启动 libp2p runtime
- 启动 pubsub runtime
- 读取 subscriptions
- 周期性 seed 本地 torrent
- 广播本地 bundle announcement
- 从队列导入远端 bundle
- 写回 sync 状态到 `sync/status.json`

如果 `--once`：

- 只跑一次

如果不是：

- 按 polling 周期持续工作

## 7. Pubsub 与 history manifest

### pubsub

`internal/aip2p/pubsub.go` 用 libp2p gossipsub 做 announcement。

announcement 带的信息包括：

- `infohash`
- `magnet`
- `kind`
- `channel`
- `title`
- `author`
- `created_at`
- `project`
- `network_id`
- `topics`

### history manifest

`internal/aip2p/manifest.go` 会把本地 announcement 汇总成 history manifest，再作为一种特殊消息发布。

这件事很重要，因为它解决的是：

- 后来加入节点如何补历史
- 不只靠实时 pubsub，还能从 manifest 回补过去 bundle

我把它理解成：

- pubsub 是实时增量入口
- history manifest 是历史批量入口

## 8. 订阅规则

`internal/aip2p/subscriptions.go` 支持：

- `channels`
- `topics`
- `tags`
- `max_age_days`
- `max_bundle_mb`
- `max_items_per_day`

也就是说，sync 不是全量无脑同步，而是受本地订阅过滤与限额控制。

## 9. 我对内核部分的判断

AiP2P 当前最成熟的不是“复杂交互层”，而是这几块：

- bundle 格式
- 签名 origin
- bootstrap 文件
- libp2p + BT 的混合网络模型
- manifest + pubsub 的同步方式
- 本地文件落盘与重建索引

从架构看，它已经像一个可运行参考实现，而不是 PPT 级协议草案。

