# 配置说明

配置读取逻辑位于 `config/config.go`。服务会优先读取 `.env`，如果没有 `.env`，则使用系统环境变量。

## 环境变量

| 变量 | 必填 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `GRSAI_API_KEY` | 是 | 无 | 请求 GRS AI 上游服务时使用的 Bearer Token |
| `GRSAI_BASE_URL` | 否 | `https://grsai.dakka.com.cn` | GRS AI 上游服务地址 |
| `SERVER_PORT` | 否 | `8080` | 当前 NewAPI 服务监听端口 |
| `PROXY_API_KEY` | 否 | 空 | 设置后，客户端请求当前服务也必须携带 `Authorization: Bearer ${PROXY_API_KEY}` |
| `RUSTFS_ENDPOINT` | 否 | 空 | RustFS/S3 兼容服务地址；设置后图片结果会转存到对象存储 |
| `RUSTFS_ACCESS_KEY` | 否 | 空 | RustFS Access Key |
| `RUSTFS_SECRET_KEY` | 否 | 空 | RustFS Secret Key |
| `RUSTFS_BUCKET` | 否 | `grsai-images` | 图片对象存储 bucket |
| `RUSTFS_REGION` | 否 | `us-east-1` | S3 签名 region，RustFS 通常用默认值即可 |
| `RUSTFS_PREFIX` | 否 | `grsai` | 图片对象 key 前缀 |
| `RUSTFS_PUBLIC_BASE_URL` | 否 | 空 | 返回给客户端的公开访问基础地址；为空时使用 `RUSTFS_ENDPOINT` |
| `RUSTFS_MAX_IMAGE_BYTES` | 否 | `52428800` | 单张图片最大下载/上传大小 |

## 示例

```dotenv
GRSAI_API_KEY=sk-your-key
GRSAI_BASE_URL=https://grsai.dakka.com.cn
SERVER_PORT=8080
PROXY_API_KEY=local-dev-token

# 可选：把 GRS AI 生成的图片 URL 转存到 RustFS/S3 兼容存储
RUSTFS_ENDPOINT=https://fs.sendi.wang
RUSTFS_ACCESS_KEY=admin
RUSTFS_SECRET_KEY=your-rustfs-secret
RUSTFS_BUCKET=grsai-images
RUSTFS_REGION=us-east-1
RUSTFS_PREFIX=grsai
RUSTFS_PUBLIC_BASE_URL=https://fs.sendi.wang
```

## 启动时校验

如果 `GRSAI_API_KEY` 为空，服务会直接退出：

```text
GRSAI_API_KEY environment variable is required
```

这可以避免服务在没有上游凭据时被误暴露。
