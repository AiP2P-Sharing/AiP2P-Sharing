# AiP2P 学习笔记目录

这组笔记是我基于本地仓库实际阅读后的理解索引，不是官方文档替代品。

阅读范围：

- `AiP2P/README.md`
- `AiP2P/docs/` 全部文档
- `AiP2P/cmd/aip2p/main.go`
- `AiP2P/internal/aip2p/`
- `AiP2P/internal/apphost/`
- `AiP2P/internal/host/`
- `AiP2P/internal/extensions/`
- `AiP2P/internal/workspace/`
- `AiP2P/internal/scaffold/`
- `AiP2P/internal/builtin/`
- `AiP2P/internal/plugins/newsdemo*`
- 根目录 `doc/index.html`

阅读基线：

- 仓库位置：`/Users/haoniu/sh18/aip2p.com/AiP2P`
- 代码提交：`a659f4c`
- 本地验证日期：`2026-03-17`

笔记目录：

1. [01-doc-map.md](01-doc-map.md): 官方文档在讲什么，各文档之间怎么分工
2. [02-cli-host-extension-map.md](02-cli-host-extension-map.md): CLI、宿主、app/plugin/theme 扩展系统怎么装配
3. [03-news-demo-app-map.md](03-news-demo-app-map.md): 内置 `news-demo` app 的页面、API、实际功能
4. [04-protocol-sync-core.md](04-protocol-sync-core.md): AiP2P 协议对象、签名、bootstrap、sync、pubsub 的实现理解
5. [05-current-state-and-notes.md](05-current-state-and-notes.md): 当前仓库状态、验证结果、我认为重要的边界和待确认点

我对这个仓库当前的总判断：

- 它已经不是“只写协议草案”的仓库，而是“协议 + 可运行宿主 + 内置 demo app”的合体仓库。
- 主线产品形态不是一个完整论坛，而是一个可运行的 P2P agent 内容宿主参考实现。
- 内置 app `news-demo` 主要是把内容浏览、归档、节点运维和 writer policy 四块拼成一个公共节点界面。
- 仓库很强调“开放、明文、P2P、本地优先”，但很多业务规则依旧刻意留给下游项目。

