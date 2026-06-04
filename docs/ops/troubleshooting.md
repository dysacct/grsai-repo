# 故障排查

## 启动后立即退出

检查 `GRSAI_API_KEY` 是否为空：

```bash
grep GRSAI_API_KEY .env
```

服务启动时会强制校验这个变量。

## 请求聊天接口返回 500

常见原因：

| 原因 | 检查方式 |
| --- | --- |
| 上游 Key 无效 | 确认 `GRSAI_API_KEY` 是否正确 |
| 上游地址不可达 | 检查 `GRSAI_BASE_URL` |
| 模型不可用 | 先请求 `/v1/models` 查看声明模型，再确认上游支持情况 |

## 流式请求没有实时输出

确认客户端使用了流式读取方式。`curl` 需要加 `-N`：

```bash
curl -N http://localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{"model":"gpt-5.5","stream":true,"messages":[{"role":"user","content":"hello"}]}'
```

如果前面还有 Nginx，检查是否关闭了代理缓冲：

```nginx
proxy_buffering off;
```

## 图片接口响应为空

当 `response_format` 为 `b64_json` 时，服务会尝试从图片 URL 下载内容再转 base64。如果服务器无法访问图片 URL，`b64_json` 可能为空。生产环境更推荐直接返回 URL。

## 文档站构建失败

先确认 Node 版本：

```bash
node -v
```

VitePress 当前要求 Node.js 20+。如果版本符合，再清理缓存后重装：

```bash
rm -rf node_modules package-lock.json docs/.vitepress/cache
npm install
npm run docs:build
```
