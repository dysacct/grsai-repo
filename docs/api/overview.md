# 接口总览

GRSAI NewAPI 暴露一组 OpenAI 风格接口，路由定义位于 `router/router.go`。

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/v1/models` | 获取服务声明支持的模型 |
| `POST` | `/v1/chat/completions` | 聊天补全，支持普通响应和 SSE 流式响应 |
| `POST` | `/v1/images/generations` | 文生图 |
| `POST` | `/v1/images/edits` | 图生图或图片编辑 |

## 请求头

如果设置了 `PROXY_API_KEY`，客户端请求当前服务时必须携带：

```http
Authorization: Bearer ${PROXY_API_KEY}
```

未设置 `PROXY_API_KEY` 时，当前服务不校验客户端 `Authorization`。服务请求上游 GRS AI 时会自动添加：

```http
Authorization: Bearer ${GRSAI_API_KEY}
Content-Type: application/json
```

## 错误格式

错误响应保持 OpenAI 风格的外层结构：

```json
{
  "error": {
    "message": "Missing required parameter: 'prompt'.",
    "type": "invalid_request_error",
    "param": "prompt"
  }
}
```

## 响应类型

| 场景 | Content-Type |
| --- | --- |
| 普通 JSON | `application/json` |
| 流式聊天 | `text/event-stream` |
