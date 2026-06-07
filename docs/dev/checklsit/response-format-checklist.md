# 渠道统一响应格式 Development Steps

## 1. 后端配置契约

- [x] 在 `dto.ChannelSettings` 中新增 `ResponseFormat` 字段，对应 `setting.response_format`
- [x] 新增 `ChannelResponseFormatSettings`，包含 `enabled`、`mode`、`rules`
- [x] `rules` 使用 ``Rules []map[string]any `json:"rules,omitempty"` ``，不得使用空 struct 预留规则
- [x] 确保 `response_format.rules` 解析后再序列化时能 round-trip 保留每条规则中的未知字段
- [x] 新增 `client_stream` mode 常量或等价判断逻辑
- [x] 实现启用判定：缺失配置关闭、`enabled=false` 关闭、`enabled=true && mode==""` 生效、`enabled=true && mode=="client_stream"` 生效、未知 mode 关闭
- [x] Test: 覆盖 `ChannelResponseFormatSettings` 启用、关闭、未知 mode 的判定分支
- [x] Test: 覆盖 `response_format.rules` 未知字段的反序列化和重新序列化不丢失

## 2. Relay 原始客户端 Stream 语义

- [x] 在 `relay/common.RelayInfo` 中新增 `ClientRequestedStream bool`
- [x] 在 `genBaseRelayInfo` 中由 `request.IsStream(c)` 同时初始化 `IsStream` 与 `ClientRequestedStream`
- [x] 确保后续 handler、adapter、内部转换逻辑只修改 `IsStream`，不修改 `ClientRequestedStream`
- [x] 检查 `GenRelayInfoResponses`、`GenRelayInfoOpenAI`、`GenRelayInfoClaude`、`GenRelayInfoGemini`、`GenRelayInfoImage` 等入口继承该字段
- [x] Test: 构造 relay info，确认 `ClientRequestedStream` 与初始化请求语义一致，且 helper 修改 `IsStream` 时不会改变它

## 3. 上游 Content-Type 流式判定 Helper

- [x] 新增公共 helper：`ApplyUpstreamContentTypeStreamDetection(info, contentType)`
- [x] 使用 `mime.ParseMediaType` 解析 `Content-Type`
- [x] 支持大小写、前后空格、`; charset=utf-8` 等 `text/event-stream` 变体
- [x] helper 只收敛已有的上游 `Content-Type` -> `info.IsStream` 自动升级逻辑，不给无检测点入口新增检测
- [x] 默认关闭时保持现有行为：上游 `text/event-stream` 可升级 `info.IsStream`
- [x] 非目标 relay mode 保持现有行为：上游 `text/event-stream` 可升级 `info.IsStream`
- [x] 目标 relay mode 仅限 `RelayModeChatCompletions` 与 `RelayModeResponses`
- [x] 排除 `RelayModeResponsesCompact`、`RelayModeCompletions`、`RelayModeGemini`、images、messages 等非目标入口
- [x] 开启 `response_format` 且客户端非流式时，忽略上游错误 `text/event-stream` 对 `info.IsStream` 的升级
- [x] 开启 `response_format` 且客户端流式时，保持现有流式升级逻辑
- [x] 不在直接 `/v1/responses` 的 `ResponsesHelper` 中新增上游 `Content-Type` -> `info.IsStream` 自动检测
- [x] Test: 覆盖目标 mode、非目标 mode、ResponsesCompact、Gemini native、images、messages、completions
- [x] Test: 覆盖 `text/event-stream` media type 的大小写、空格、charset 参数和非法 media type

## 4. 替换现有流式误判点

- [ ] 将 `relay/compatible_handler.go` 中直接基于 `Content-Type` 修改 `info.IsStream` 的逻辑替换为 helper
- [ ] 将 `relay/claude_handler.go` 中直接基于 `Content-Type` 修改 `info.IsStream` 的逻辑替换为 helper
- [ ] 将 `relay/gemini_handler.go` 中直接基于 `Content-Type` 修改 `info.IsStream` 的逻辑替换为 helper
- [ ] 将 `relay/image_handler.go` 中直接基于 `Content-Type` 修改 `info.IsStream` 的逻辑替换为 helper，确保非目标 images 仍保持现有行为
- [ ] 将 `relay/chat_completions_via_responses.go` 中内部 Responses 转换路径的直接判定替换为 helper
- [ ] 不修改直接 `/v1/responses` 为新增基于上游 `Content-Type` 的 stream auto-detect
- [ ] 移除替换后不再需要的 `strings` import
- [ ] Test: `/v1/chat/completions` 非流式请求开启配置后，遇到上游 `Content-Type: text/event-stream` 不升级为流式
- [ ] Test: `chat_completions_via_responses` 内部转换路径开启配置后，遇到上游 `Content-Type: text/event-stream` 不升级为流式
- [ ] Test: `/v1/messages`、Gemini native、images、旧 `/v1/completions` 不受新配置影响

## 5. 响应头规范化 Context Flag

- [ ] 在 `constant/context_key.go` 中新增响应头规范化 context key
- [ ] 在 relay 层新增 helper，根据 `RelayInfo` 和上游状态码计算是否需要规范化响应头
- [ ] 使用已确认的 A 方案：relay 层写入 context flag，统一写回路径只读取 flag 执行 header 覆盖
- [ ] 生效条件限制为：配置开启、`RelayModeChatCompletions` 或 `RelayModeResponses`、`ClientRequestedStream=false`、`info.IsStream=false`、上游 `StatusCode == http.StatusOK`
- [ ] 直接 `/v1/responses` 的 `200 OK` 非流式写回需要设置响应头规范化 flag
- [ ] `chat_completions_via_responses` 内部转换路径的 `200 OK` 非流式写回需要设置响应头规范化 flag
- [ ] 非 200 错误响应不得设置规范化 flag
- [ ] 客户端流式请求不得设置规范化 flag
- [ ] ResponsesCompact 不得设置规范化 flag
- [ ] Test: 覆盖 context flag 的启用、关闭、非 200、客户端流式、ResponsesCompact 分支

## 6. 统一写回路径接入

- [ ] 在 `service.IOCopyBytesGracefully` 中读取 context flag
- [ ] 在复制上游 header 之后、设置 `Content-Length` 和 `WriteHeader` 之前覆盖 `Content-Type`
- [ ] 覆盖值固定为 `application/json; charset=utf-8`
- [ ] 保留其他上游 header 复制逻辑
- [ ] 保留 `Content-Length` 按最终 body 重新计算逻辑
- [ ] 保留 `X-Oneapi-Request-Id` 捕获和本地 request id 保护逻辑
- [ ] 不在 adapter 或 handler 中分散手写 `Content-Type` 覆盖
- [ ] Test: 开启配置后，目标接口 `200 OK` 非流式响应最终返回 `application/json; charset=utf-8`
- [ ] Test: `Content-Length` 仍等于最终 body 长度
- [ ] Test: 请求 ID 捕获逻辑不变

## 7. Responses 直接入口与内部转换校验

- [ ] 检查直接 `/v1/responses` 入口，确认流式语义仍来自请求体 `stream`
- [ ] 直接 `/v1/responses` 只新增 `200 OK` 非流式响应头规范化 flag，不新增 stream auto-detect
- [ ] 检查 `/v1/responses/compact` 对应 `RelayModeResponsesCompact`，确认不被 response_format 覆盖
- [ ] 在 `chat_completions_via_responses` 临时切换到 `RelayModeResponses` 时设置或清理规范化 flag
- [ ] 在 `chat_completions_via_responses` 临时切换到 `RelayModeResponses` 时替换已有 stream auto-detect
- [ ] 确保内部转换结束后恢复原 `RelayMode` 和 `RequestURLPath`
- [ ] Test: 直接 `/v1/responses` 的 stream 判断不因上游响应头改变客户端语义
- [ ] Test: 直接 `/v1/responses` 开启配置后，`200 OK` 非流式响应最终返回 `application/json; charset=utf-8`
- [ ] Test: `chat_completions_via_responses` 开启配置后，`200 OK` 非流式响应最终返回 `application/json; charset=utf-8`
- [ ] Test: `/v1/responses/compact` 不升级、不规范化、不改变现有响应行为

## 8. 前端表单类型与默认值

- [ ] 在 `web/default/src/features/channels/types.ts` 中为 `ChannelSettings` 增加 `response_format`
- [ ] 定义前端 `ChannelResponseFormatSettings` 类型，包含 `enabled`、`mode`、`rules`
- [ ] 在 `channel-form.ts` 的 Zod schema 中新增统一响应格式表单字段
- [ ] 在默认表单值中设置统一响应格式开关为关闭
- [ ] 编辑旧渠道时，`setting.response_format` 缺失应回填为关闭
- [ ] 编辑已有渠道时，正确读取 `response_format.enabled`
- [ ] 编辑已有渠道时，`response_format.rules` 应作为对象数组保留，不能被类型收窄为固定空对象
- [ ] Test: 新建渠道默认关闭
- [ ] Test: 编辑旧渠道缺少 `response_format` 时默认关闭

## 9. 前端 Setting JSON 序列化

- [ ] 修改 `buildSettingJSON`，保存 `response_format.enabled`
- [ ] 保存时始终写入 `response_format.mode = "client_stream"`
- [ ] 开启保存时，新渠道可写入 `rules: []`
- [ ] 开启保存编辑已有渠道时，保留已有 `response_format.rules` 及每条规则中的未知字段
- [ ] 关闭保存时，写入完整默认对象：`enabled=false`、`mode="client_stream"`、`rules=[]`
- [ ] 不删除整个 `response_format` 对象
- [ ] 不展示、不执行 `rules` 列表
- [ ] Test: 开启后保存的 `setting` JSON 包含 `response_format.enabled=true`
- [ ] Test: 关闭后保存的 `setting` JSON 包含完整关闭对象
- [ ] Test: 编辑已有渠道开启保存时保留已有 `rules`，并验证未知字段 round-trip 不丢失

## 10. 前端高级设置 UI

- [ ] 在 `Channel Extra Settings` 区域新增 `Unified Response Format` 开关
- [ ] 将开关放在 `Force Format`、`Thinking to Content`、`Pass Through Body` 等通用渠道行为设置同组
- [ ] 添加说明文案：`Normalize non-streaming chat and responses output based on the client stream mode.`
- [ ] 将该字段接入高级设置展开判断
- [ ] 将该字段接入错误字段映射或字段错误显示逻辑
- [ ] 确保 UI 不展示规则列表入口
- [ ] Test: 打开渠道创建/编辑抽屉，确认开关显示、默认状态正确、保存 payload 正确

## 11. 前端 i18n

- [ ] 在 `en.json` 添加 `Unified Response Format`
- [ ] 在 `en.json` 添加 `Normalize non-streaming chat and responses output based on the client stream mode.`
- [ ] 在 `zh.json` 添加对应简体中文翻译
- [ ] 在 `fr.json` 添加对应法语翻译
- [ ] 在 `ja.json` 添加对应日语翻译
- [ ] 在 `ru.json` 添加对应俄语翻译
- [ ] 在 `vi.json` 添加对应越南语翻译
- [ ] 运行 `bun run i18n:sync`
- [ ] Test: 确认新增 `t()` key 在所有 locale 文件中存在且无缺失

## 12. 后端回归测试

- [ ] 新增或扩展 relay common helper 单元测试
- [ ] 新增或扩展响应写回 header 规范化测试
- [ ] 覆盖默认关闭时保持现有行为
- [ ] 覆盖开启后非流式目标接口不被上游错误 `text/event-stream` 升级
- [ ] 覆盖客户端流式请求仍走流式响应
- [ ] 覆盖非目标入口不受影响
- [ ] 覆盖非 200 错误响应不进入统一响应格式处理
- [ ] 覆盖上游 body 实际为 SSE 时不做聚合转换，沿用非流式解析失败路径
- [ ] 覆盖 `response_format.rules` 使用 `[]map[string]any` 保留未知字段，不使用空 struct
- [ ] 运行聚焦测试：`go test ./relay/common ./service ./relay/...`
- [ ] 按失败范围补充运行更窄或更广的 Go 测试

## 13. 前端验证

- [ ] 运行 `bun run typecheck`
- [ ] 运行渠道表单相关测试或可用的前端测试脚本
- [ ] 检查创建渠道 payload 中 `setting.response_format` 的开启和关闭序列化
- [ ] 检查编辑渠道回填中旧渠道、新配置渠道、有保留 rules 渠道的显示状态
- [ ] 检查有未知字段的 `rules` 在前端回填、切换开关、保存后不被清空或裁剪
- [ ] 检查新增 UI 文案不会溢出或破坏高级设置布局

## 14. 最终检查

- [ ] 运行 `gofmt` 格式化修改过的 Go 文件
- [ ] 检查没有直接新增业务代码中的 `encoding/json` marshal/unmarshal 调用
- [ ] 检查没有新增数据库列或迁移
- [ ] 检查没有修改 `web/classic`
- [ ] 检查没有修改受保护项目标识和组织标识
- [ ] 检查 `git diff`，确认改动范围仅限本功能
- [ ] 汇总已运行的后端测试、前端 typecheck、i18n 同步结果
