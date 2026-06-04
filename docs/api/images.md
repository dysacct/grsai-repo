# 图片生成与编辑

图片接口对外保持 OpenAI 风格，对内请求 GRS AI 的绘图接口。绘图结果来自上游 SSE 任务的最终成功事件。

## 文生图

### `POST /v1/images/generations`

```bash
curl http://localhost:8080/v1/images/generations \
  -H 'Content-Type: application/json' \
  -d '{
    "model": "gpt-image-2",
    "prompt": "一张干净的产品文档封面，带有 API 网关和光线轨迹",
    "size": "1024x1024",
    "quality": "high"
  }'
```

### 参数

| 字段 | 类型 | 必填 | 默认值 | 说明 |
| --- | --- | --- | --- | --- |
| `prompt` | string | 是 | 无 | 图片提示词 |
| `model` | string | 否 | `gpt-image-2` | 图片模型 |
| `n` | number | 否 | `1` | 当前会发起一次上游绘图任务 |
| `size` | string | 否 | `1024x1024` | OpenAI 尺寸或 GRS AI 原生尺寸 |
| `quality` | string | 否 | 无 | 透传给上游 |
| `response_format` | string | 否 | `url` | 为 `b64_json` 时返回 base64 |

## 图生图

### `POST /v1/images/edits`

支持两种输入方式。

### JSON

```bash
curl http://localhost:8080/v1/images/edits \
  -H 'Content-Type: application/json' \
  -d '{
    "model": "nano-banana-pro",
    "image": "https://example.com/source.png",
    "prompt": "保留主体，把背景换成简洁的科技展台",
    "size": "1024x1024"
  }'
```

### Multipart

```bash
curl http://localhost:8080/v1/images/edits \
  -F model=nano-banana-pro \
  -F prompt='保留主体，把背景换成简洁的科技展台' \
  -F size=1024x1024 \
  -F image=@source.png
```

## 尺寸映射

| 输入尺寸 | 上游 size | aspect_ratio |
| --- | --- | --- |
| `256x256` | `1K` | `1:1` |
| `512x512` | `1K` | `1:1` |
| `1024x1024` | `1K` | `1:1` |
| `1024x1792` | `2K` | `9:16` |
| `1792x1024` | `2K` | `16:9` |
| `1K` / `2K` / `4K` | 原样透传 | 空 |

未知尺寸会原样透传给上游。
