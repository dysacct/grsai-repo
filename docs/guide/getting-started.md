# 快速开始

本项目是一个 Go 服务，对外暴露 OpenAI 风格接口，对内请求 GRS AI。

## 环境要求

| 依赖 | 建议 |
| --- | --- |
| Go | 使用仓库 `go.mod` 指定版本或服务器已安装版本 |
| Node.js | 20+，用于运行 VitePress 文档站 |
| GRS AI Key | 通过 `GRSAI_API_KEY` 注入 |

## 启动 API 服务

在项目根目录创建 `.env`：

```dotenv
GRSAI_API_KEY=你的_grsai_key
GRSAI_BASE_URL=https://grsai.dakka.com.cn
SERVER_PORT=8080
```

启动服务：

```bash
go run .
```

默认监听：

```text
http://localhost:8080
```

## 验证模型接口

```bash
curl http://localhost:8080/v1/models
```

如果返回 `object: "list"` 和模型数组，说明服务已经启动。

## 启动文档站

安装依赖：

```bash
npm install
```

本地开发：

```bash
npm run docs:dev
```

构建静态文档：

```bash
npm run docs:build
```

预览构建结果：

```bash
npm run docs:preview
```
