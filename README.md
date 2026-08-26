# 内容安全策略与审核决策服务

纯Go实现的内容安全策略服务，支持关键词、正则和标签规则，按优先级返回通过、拒绝、隔离或人工复核结论。检测器、队列、Redis租约与PostgreSQL仓储均通过接口隔离，本地内存适配器可完成策略发布、内容检测和审核任务领取/提交流程，日志不输出内容原文。

## 运行

需要Go 1.22或更高版本：

```bash
go test ./...
go run ./cmd/content-moderation-service
```

默认监听 `:8083`。可通过 `CONTENT_MODERATION_HTTP_ADDR`、`CONTENT_MODERATION_ENVIRONMENT` 和 `CONTENT_MODERATION_SHUTDOWN_SECONDS` 覆盖配置。`configs/config.yaml`提供配置示例；SQL迁移位于`migrations`，Docker文件位于`deploy`。

## API流程

```bash
curl -X POST localhost:8083/v1/moderation/policies -H 'Content-Type: application/json' -d '{"id":"text-v1","name":"Text policy","rules":[{"id":"r1","name":"review phrase","kind":"keyword","pattern":"manual review","risk":"medium","priority":1,"decision":"review"}]}'
curl -X POST localhost:8083/v1/moderation/policies/publish -H 'Content-Type: application/json' -d '{"id":"text-v1"}'
curl -X POST localhost:8083/v1/moderation/check -H 'Content-Type: application/json' -d '{"id":"sample-1","text":"requires manual review","policy_id":"text-v1"}'
curl -X POST localhost:8083/v1/moderation/tasks/claim -H 'Content-Type: application/json' -d '{"reviewer":"alice"}'
```

策略版本发布会刷新规则缓存；人工复核使用有时限的任务租约，检测器超时和错误通过队列重试。服务收到SIGINT/SIGTERM后优雅关闭。

## 目录

`cmd`为启动入口，`internal/domain`为纯领域模型和规则求值，`internal/application`为用例，`internal/adapter`为HTTP和缓存适配器，`internal/infrastructure`为配置、日志、指标和仓储实现，`api`为OpenAPI，`migrations`为数据库迁移，`deploy`和`scripts`提供部署与校验辅助。
