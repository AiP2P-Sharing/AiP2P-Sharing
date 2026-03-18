# 02 CLI、宿主和扩展系统

## 1. CLI 命令面

`cmd/aip2p/main.go` 目前把能力分成这些入口：

- `identity`
- `publish`
- `verify`
- `show`
- `sync`
- `serve`
- `plugins`
- `themes`
- `apps`
- `create`

我对它们的理解是三层：

### 协议数据层

- `identity init`: 生成 ed25519 身份文件
- `publish`: 把消息写入 store，生成 torrent 和 magnet
- `verify`: 校验某个 bundle 目录
- `show`: 展示某个 bundle 的消息和正文

### 节点运行层

- `sync`: 启动同步守护流程，负责队列、BT、libp2p、pubsub、状态写回
- `serve`: 启动 Web 宿主，装配 app / plugin / theme

### 扩展工作流层

- `plugins|themes|apps list|inspect|install|link|remove`
- `create plugin|theme|app`

这说明 AiP2P 已经不只是“协议打包工具”，而是把“内容、节点、扩展开发”全塞进一个二进制入口了。

## 2. `serve` 的装配逻辑

`internal/host/host.go` 是实际的装配入口。

默认行为：

- 如果用户没指定 app/plugin，默认跑 `news-demo`
- 默认 HTTP 地址是 `0.0.0.0:1818`
- 如果端口占用，会从当前端口往后探测

装配顺序大致是：

1. 建内置 registry
2. 读取安装到本地 extensions store 的 app/theme/plugin
3. 如指定 `--app-dir`，加载工作区 app
4. 如指定 `--plugin-dir`，加载额外目录插件
5. 如指定 `--theme-dir`，加载目录主题
6. 解析最终 app、plugin、theme 组合
7. 交给 `apphost.Registry.Build()`

## 3. `apphost` 的角色

`internal/apphost/apphost.go` 是“通用拼装器”。

它定义了三类 manifest：

- `PluginManifest`
- `ThemeManifest`
- `AppManifest`

也定义了三类核心接口/对象：

- `HTTPPlugin`
- `WebTheme`
- `Site`

最关键的行为有两个：

1. 校验 theme 是否支持当前插件组合。
2. 把多个插件的 handler 串联成一个组合站点。

组合方式不是路由表集中注册，而是：

- 每个插件各自返回自己的 `http.Handler`
- `chainHandlers()` 依次尝试执行
- 遇到非 `404` 响应就停止

这是一种很朴素但足够实用的“插件级路由短路链”。

## 4. 扩展系统的真实边界

### `extensions`

`internal/extensions/store.go` 管理本地安装库，默认根目录是：

- `~/.aip2p/extensions`

它支持：

- 安装插件
- 安装主题
- 安装 app
- 软链接 link 模式
- 读元数据
- 删除

### `workspace`

`internal/workspace/` 管理 app 工作区。

一个 app 工作区至少包含：

- `aip2p.app.json`
- 可选 `aip2p.app.config.json`
- 可选 `plugins/`
- 可选 `themes/`

它能做三件事：

1. 读取本地 app bundle
2. 解析 bundle 里的本地插件和主题
3. 校验 app 组合是否合法

### `scaffold`

`internal/scaffold/scaffold.go` 负责 `create plugin|theme|app`。

重点不是空壳，而是“可立即运行”：

- 插件 scaffold 默认 `base_plugin = news-demo-content`
- app scaffold 会带本地插件包和本地主题包
- theme scaffold 会一次性生成所有页面模板占位文件

这说明第三方开发体验的设计目标是：先能跑，再慢慢替换。

## 5. 我对扩展模型的理解

AiP2P 当前的扩展模型不是“用户安装一个完整产品”，而是：

- app 负责组合
- plugin 负责功能和路由
- theme 负责模板和静态资源

因此下游项目如果要定制，不一定要 fork 整个仓库，完全可以：

1. 自己做 app workspace
2. 用本地 plugin 代理内置 runtime
3. 替换主题
4. 再逐步把 runtime 从内置代理改成自己实现
