# 配置说明

配置读取逻辑位于 `config/config.go`。服务会优先读取 `.env`，如果没有 `.env`，则使用系统环境变量。

## 环境变量

| 变量 | 必填 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `GRSAI_API_KEY` | 是 | 无 | 请求 GRS AI 上游服务时使用的 Bearer Token |
| `GRSAI_BASE_URL` | 否 | `https://grsai.dakka.com.cn` | GRS AI 上游服务地址 |
| `SERVER_PORT` | 否 | `8080` | 当前 NewAPI 服务监听端口 |

## 示例

```dotenv
GRSAI_API_KEY=sk-your-key
GRSAI_BASE_URL=https://grsai.dakka.com.cn
SERVER_PORT=8080
```

## 启动时校验

如果 `GRSAI_API_KEY` 为空，服务会直接退出：

```text
GRSAI_API_KEY environment variable is required
```

这可以避免服务在没有上游凭据时被误暴露。
