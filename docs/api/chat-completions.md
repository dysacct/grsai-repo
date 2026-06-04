# 聊天补全

## `POST /v1/chat/completions`

请求会转发到 GRS AI 的 `/v1/chat/completions`。

## 请求体

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `model` | string | 是 | 模型 ID |
| `messages` | array | 是 | 对话消息 |
| `stream` | boolean | 否 | 为 `true` 时启用 SSE 流式响应 |
| `max_tokens` | number | 否 | 最大输出 token |
| `temperature` | number | 否 | 采样温度 |
| `top_p` | number | 否 | nucleus sampling 参数 |
| `tools` | any | 否 | 工具调用参数，原样透传 |
| `tool_choice` | any | 否 | 工具选择参数，原样透传 |

## 普通请求
Hyh3202276686@@@
```bash
curl http://localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{
    "model": "gpt-5.5",
    "messages": [
      { "role": "user", "content": "你好，介绍一下你自己" }
    ]
  }'
```

## 流式请求

```bash
curl -N http://localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{
    "model": "gpt-5.5",
    "stream": true,
    "messages": [
      { "role": "user", "content": "流式输出三条部署建议" }
    ]
  }'
```

## 多模态消息

`content` 可以是字符串，也可以是 OpenAI 风格的内容数组：

```json
{
  "model": "gemini-3.1-pro",
  "messages": [
    {
      "role": "user",
      "content": [
        { "type": "text", "text": "请描述这张图" },
        {
          "type": "image_url",
          "image_url": {
            "url": "https://example.com/image.png"
          }
        }
      ]
    }
  ]
}
```
