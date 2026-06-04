# 模型列表

## `GET /v1/models`

返回当前服务声明支持的模型列表。

```bash
curl http://localhost:8080/v1/models
```

## 响应示例

```json
{
  "object": "list",
  "data": [
    {
      "id": "gpt-5.5",
      "object": "model",
      "created": 1750000000,
      "owned_by": "grsai"
    }
  ]
}
```

## 已声明模型

### 聊天模型

| 模型 |
| --- |
| `gpt-5.5` |
| `gpt-5.4` |
| `gemini-3.1-pro` |
| `gemini-3.1-flash-lite` |
| `gemini-3-flash` |
| `gemini-3-pro` |
| `gemini-2.5-flash` |
| `gemini-2.5-pro` |

### 图片模型

| 模型 |
| --- |
| `gpt-image-2` |
| `gpt-image-2-vip` |
| `nano-banana` |
| `nano-banana-pro` |
| `nano-banana-pro-vt` |
| `nano-banana-2` |
| `nano-banana-fast` |
| `nano-banana-pro-cl` |
| `nano-banana-2-cl` |
| `nano-banana-2-4k-cl` |
| `nano-banana-pro-vip` |
| `nano-banana-pro-4k-vip` |
