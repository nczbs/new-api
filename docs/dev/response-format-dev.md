# 渠道统一响应格式开发设计

## 背景

部分上游渠道在客户端请求非流式响应时，仍返回错误的响应头，例如：

```http
Content-Type: text/event-stream
```

但此时响应体实际应按非流式 JSON 处理，下游客户端也应收到 JSON 类型响应。当前 relay 链路中存在两类相关行为：

1. 非流式响应写回时会复制上游响应头。
2. 部分 handler 会根据上游 `Content-Type: text/event-stream` 将 `info.IsStream` 自动升级为流式。

因此，当上游只是不规范地返回了错误 `Content-Type` 时，系统可能把客户端非流式请求误判为流式请求，并向客户端透传错误响应头。

## 目标

新增渠道级“统一响应格式”能力，用于修正当前渠道的非流式聊天响应格式。

第一版目标：

- 配置为渠道级设置，只影响当前渠道，不影响所有渠道。
- 默认关闭，避免改变现有渠道行为。
- 仅覆盖 `/v1/chat/completions` 与 `/v1/responses`。
- 仅覆盖上游 `StatusCode == http.StatusOK` 的目标接口非流式响应。
- 开启后以客户端请求语义为准，非流式请求不得被上游错误 `Content-Type` 自动升级为流式。
- 开启后目标接口 `200 OK` 非流式响应最终返回：

```http
Content-Type: application/json; charset=utf-8
```

## 非目标

第一版不处理以下内容：

- 不覆盖 `/v1/messages`。
- 不覆盖 Gemini native，例如 `/v1beta/models/*:generateContent`。
- 不覆盖 `/v1/responses/compact`。
- 不覆盖旧 `/v1/completions`。
- 不覆盖 embeddings、images、rerank、audio、realtime、MJ、Suno/task 等入口。
- 不处理非 200 错误响应。
- 不做 SSE body 到 JSON body 的聚合转换。
- 不展示、不执行自定义响应规则列表。
- 第一版仅覆盖默认前端 `web/default` 的渠道创建/编辑 UI，不同步 `web/classic`。

如果上游在客户端非流式请求下实际返回 SSE body，而不仅是响应头错误，第一版仍按非流式 JSON 解析；解析失败时沿用现有错误处理路径。

## 配置设计

统一响应格式配置存储在渠道的 `setting` JSON 中，而不是全局配置，也不新增数据库列。

注意：这里的 `setting.response_format` 是渠道级设置，不等同于 OpenAI 请求体中的 `response_format` 参数。

建议配置结构：

```json
{
  "response_format": {
    "enabled": true,
    "mode": "client_stream",
    "rules": []
  }
}
```

字段说明：

| 字段                        | 类型      | 第一版含义                                               |
|---------------------------|---------|-----------------------------------------------------|
| `response_format.enabled` | boolean | 是否启用统一响应格式                                          |
| `response_format.mode`    | string  | 响应模式判定策略。第一版唯一有效值为 `client_stream`，表示以客户端原始请求是否流式为准 |
| `response_format.rules`   | array   | 第一版预留为空数组，不执行                                       |

`response_format.mode` 说明：

- `mode` 是统一响应格式能力的策略枚举字段，不是模型输出格式，不等同于 OpenAI 请求体里的 `response_format`。
- 第一版只支持 `client_stream`。
- `client_stream` 表示响应模式以客户端原始请求中的 `stream` 语义为准：
    - 客户端未请求流式响应时，上游错误返回 `Content-Type: text/event-stream` 也不得把当前请求自动升级为流式。
    - 客户端请求流式响应时，继续保持现有流式处理逻辑。
- 第一版不支持按上游 `Content-Type`、固定 JSON、规则列表等其它策略决定响应模式。

兼容策略：

- 缺失 `response_format` 时视为关闭。
- `response_format.enabled` 为 `false` 时视为关闭。
- `response_format.enabled` 为 `true` 且 `response_format.mode` 缺失或为 `client_stream` 时生效。
- `response_format.enabled` 为 `true` 但 `response_format.mode` 为未知值时，保守视为关闭，避免错误配置改变线上行为。
- `response_format.rules` 第一版不解析执行，仅保留后续扩展位置。
- 前端开启保存时应保留已有 `response_format.rules`，不得因为开启开关清空用户手写或未来版本写入的规则。
- 前端关闭保存时必须写入完整默认对象，明确表达该能力已关闭：

```json
{
  "response_format": {
    "enabled": false,
    "mode": "client_stream",
    "rules": []
  }
}
```

## 后端设计

### ChannelSettings

在 `dto.ChannelSettings` 中新增嵌套结构，用于解析 `setting` 中的 `response_format`。

建议结构命名：

```go
type ChannelResponseFormatSettings struct {
Enabled bool                        `json:"enabled,omitempty"`
Mode    string                      `json:"mode,omitempty"`
Rules   []ChannelResponseFormatRule `json:"rules,omitempty"`
}

type ChannelResponseFormatRule struct {
// 第一版仅预留结构，暂不定义执行语义。
}
```

`dto.ChannelSettings` 增加：

```go
ResponseFormat ChannelResponseFormatSettings `json:"response_format,omitempty"`
```

### 客户端请求语义

第一版必须新增不可变字段保存客户端原始流式语义，不能把“客户端请求语义”等同于当前 `info.IsStream`。

建议在 `relay/common.RelayInfo` 增加字段：

```go
ClientRequestedStream bool
```

赋值规则：

- 在 `genBaseRelayInfo` 中解析请求后立即赋值，值来源与当前初始化 `info.IsStream` 的 `request.IsStream(c)` 一致。
- 初始化时 `info.IsStream` 与 `info.ClientRequestedStream` 可以相同。
- `ClientRequestedStream` 表示客户端原始请求是否要求流式响应，后续 handler、adapter、上游响应头、内部请求格式转换都不得修改该字段。
- `info.IsStream` 仍表示 relay 执行过程中的最终处理模式，可以被现有逻辑或新 helper 修改。
- `response_format` 的所有“客户端是否流式”判断必须读取 `info.ClientRequestedStream`，不得读取可能已经被上游响应头污染的当前
  `info.IsStream`。

### 流式判定

当前部分代码直接使用类似逻辑：

```go
info.IsStream = info.IsStream || strings.HasPrefix(httpResp.Header.Get("Content-Type"), "text/event-stream")
```

第一版应新增公共 helper，统一判断上游 `Content-Type` 是否允许影响 `info.IsStream`。

建议 helper 形态：

```go
func ApplyUpstreamContentTypeStreamDetection(info *relaycommon.RelayInfo, contentType string)
```

helper 内部应解析上游 `Content-Type` 的 media type，例如使用 `mime.ParseMediaType`，避免大小写、空格、`; charset=utf-8`
等格式差异导致误判。

行为规则：

- 未开启 `response_format.enabled` 时，保持现有逻辑。
- 当前 relay mode 不是 `RelayModeChatCompletions` 或 `RelayModeResponses` 时，保持现有逻辑。
- `info.ClientRequestedStream` 为 `true` 时，保持流式逻辑。
- 开启配置且当前请求是非流式目标接口时，忽略上游 `text/event-stream` 对 `info.IsStream` 的自动升级。

`/v1/responses` 误判点说明：

- 直接 `/v1/responses` 入口的 `ResponsesHelper` 当前没有发现基于上游 `Content-Type` 自动升级 `info.IsStream`
  的逻辑；其流式语义来自请求体 `stream` 字段。
- 已核实的 `RelayModeResponses` 误判点存在于 `relay/chat_completions_via_responses.go`：该内部转换路径会临时把
  `info.RelayMode` 设置为 `RelayModeResponses`，请求上游 Responses API 后再根据上游 `Content-Type: text/event-stream` 修改
  `info.IsStream`。
- helper 接入必须覆盖这条 chat-completions-via-responses 内部转换路径，同时不得改变直接 `/v1/responses` 的既有客户端
  stream 判断语义。

实现约束：

- endpoint 覆盖范围必须基于 `info.RelayMode` 判断，而不是直接对请求路径做前缀匹配。
- `/v1/responses` 只对应 `RelayModeResponses`；`/v1/responses/compact` 对应独立的 `RelayModeResponsesCompact`，必须保持排除。
- 不要使用 `strings.HasPrefix(path, "/v1/responses")` 判断是否命中本能力，否则会把 `/v1/responses/compact` 误纳入覆盖范围。

该 helper 应替代所有现有上游 `Content-Type` 自动升级 `info.IsStream` 的直接写法，避免遗漏和行为分叉。替换后，非目标 relay
mode 必须由 helper 保持原行为。

### 响应头写回

目标接口 `200 OK` 非流式响应写回时增加响应头规范化。

生效条件：

- 当前渠道开启 `response_format.enabled`。
- 当前 relay mode 是 `/v1/chat/completions` 对应的 `RelayModeChatCompletions` 或 `/v1/responses` 对应的
  `RelayModeResponses`。
- `info.ClientRequestedStream` 为 `false`。
- 当前请求最终按非流式响应处理，即写回前 `info.IsStream == false`。
- 上游响应状态为 `http.StatusOK`。

生效行为：

- 覆盖下游响应 `Content-Type` 为 `application/json; charset=utf-8`。
- 继续保留现有其他可透传响应头逻辑。
- `Content-Length` 仍由写回逻辑按最终 body 重新计算。
- 请求 ID 捕获逻辑保持不变。

实现约束：

- 响应头规范化应接入统一写回路径，不要在各 adaptor 或 handler 中分散手写 `c.Writer.Header().Set("Content-Type", ...)`。
- 当前非流式写回会先复制上游响应头，再设置 `Content-Length` 并调用 `WriteHeader`；因此 `Content-Type` 覆盖必须发生在复制上游
  header 之后、`WriteHeader` 之前。
- 写回层可以通过显式策略参数、带 relay 上下文的 helper，或由 relay 层设置的 context flag
  接收规范化意图；无论采用哪种方式，最终覆盖动作都必须由统一写回路径执行。
- 不要把渠道级 `response_format` 判定逻辑直接扩散到与 relay 无关的通用 HTTP 工具中；如果需要复用，应由 relay 层先根据
  `RelayInfo` 计算策略，再传给写回层执行。

### 与 Force Format 的关系

`force_format` 与 `response_format` 是两个独立能力：

- `force_format` 负责响应体结构标准化。
- `response_format` 第一版负责响应模式判定与响应头规范化。

二者可以同时开启。第一版不因为 `response_format` 开启而额外改写响应体 schema。

## 前端设计

在渠道创建/编辑的高级设置中新增开关。

建议 UI 位置：

- 位于“Channel Extra Settings”区域。
- 与 `Force Format`、`Thinking to Content`、`Pass Through Body` 等通用渠道行为设置同组。

建议文案：

- Label: `Unified Response Format`
- Description: `Normalize non-streaming chat and responses output based on the client stream mode.`

第一版 UI 行为：

- 默认关闭。
- 编辑旧渠道时，如果 `setting.response_format` 缺失，显示为关闭。
- 开启后保存到 `setting.response_format.enabled=true`。
- 保存时始终写入 `mode="client_stream"`。
- 开启保存时，新建渠道首次开启可写入 `rules=[]`；编辑已有渠道时必须保留现有 `rules`。
- 关闭保存时，写入 `response_format.enabled=false`、`mode="client_stream"` 与 `rules=[]`，不得删除整个 `response_format`
  对象。
- 不展示规则列表入口。
- 表单 schema、默认值、编辑回填、提交序列化、高级设置展开判断、错误字段映射都需要接入该字段。

新增前端文案后，需要同步默认前端 i18n 语言。

## 测试计划

### 后端测试

需要覆盖以下场景：

1. 默认关闭时保持现有行为。
2. 开启后，`/v1/chat/completions` 非流式请求遇到上游 `Content-Type: text/event-stream` 时，不升级为流式。
3. 开启后，`chat_completions_via_responses` 内部转换路径遇到上游 `Content-Type: text/event-stream` 时，不升级为流式。
4. 直接 `/v1/responses` 的流式判断仍来自请求体 `stream` 字段，不因上游响应头改变客户端语义。
5. 开启后，目标接口 `200 OK` 非流式响应最终返回 `application/json; charset=utf-8`。
6. 客户端流式请求不受影响，仍走流式响应。
7. `/v1/messages` 不受影响。
8. Gemini native 不受影响。
9. `/v1/responses/compact` 不受影响。
10. `/v1/completions` 不受影响。
11. 非 200 错误响应不进入统一响应格式处理。
12. 上游 body 实际为 SSE 时，不做聚合，沿用非流式解析失败路径。
13. helper 单元测试覆盖 `RelayModeChatCompletions`、`RelayModeResponses`、`RelayModeResponsesCompact`、Gemini
    native、images、messages、completions 等模式。
14. `Content-Type` media type 解析覆盖大小写、空格、带 charset 参数等情况。
15. 响应头测试确认 `Content-Length` 仍按最终 body 重算，请求 ID 捕获逻辑不变。

### 前端测试与检查

需要覆盖以下场景：

1. 新建渠道默认关闭。
2. 编辑旧渠道缺少 `response_format` 时默认关闭。
3. 开启后表单保存的 `setting` JSON 包含 `response_format.enabled=true`。
4. 关闭后后端按关闭处理。
5. 新增 i18n key 后运行 i18n 同步。
6. 修改 TypeScript/TSX 后运行类型检查。

## 后续扩展

后续可在 `response_format.rules` 中追加响应头与响应体字段改写规则。

建议扩展方向：

- 规则以列表形式在 UI 中展示。
- 每条规则支持启用状态。
- 每条规则支持条件 rule。
- 操作对象可区分响应头与响应体 JSON 字段。
- 响应体字段操作可参考现有请求侧 `param_override` 的 JSON path、operation、condition 思路。
- 响应规则执行器应与请求规则执行器分离，避免混淆请求改写与响应改写的生命周期。

后续规则示例仅作为方向，第一版不实现：

```json
{
  "response_format": {
    "enabled": true,
    "mode": "client_stream",
    "rules": [
      {
        "enabled": true,
        "target": "header",
        "operation": "set",
        "path": "X-Custom-Header",
        "value": "custom-value",
        "conditions": []
      },
      {
        "enabled": true,
        "target": "body",
        "operation": "set",
        "path": "metadata.normalized",
        "value": true,
        "conditions": []
      }
    ]
  }
}
```

## 实施顺序

1. 先合入本文档。
2. 再实现后端配置结构与 helper。
3. 接着接入 `/v1/chat/completions` 与 `/v1/responses` 的流式判定和响应头规范化。
4. 再实现前端高级设置开关。
5. 最后补充测试与 i18n 同步。
