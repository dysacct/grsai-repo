# 部署

下面是一套偏保守的服务器部署方式，适合直接放在 `/data/transfer-api/grsai-newapi` 这类目录中运行。

## 构建 API 服务

```bash
cd /data/transfer-api/grsai-newapi
go build -o grsai-newapi .
```

## 使用 systemd 托管

创建环境文件：

```bash
sudo install -d -m 0750 /etc/grsai-newapi
sudo tee /etc/grsai-newapi/env >/dev/null <<'EOF'
GRSAI_API_KEY=sk-your-key
GRSAI_BASE_URL=https://grsai.dakka.com.cn
SERVER_PORT=8080
EOF
```

创建服务文件：

```ini
[Unit]
Description=GRSAI NewAPI
After=network-online.target
Wants=network-online.target

[Service]
WorkingDirectory=/data/transfer-api/grsai-newapi
EnvironmentFile=/etc/grsai-newapi/env
ExecStart=/data/transfer-api/grsai-newapi/grsai-newapi
Restart=always
RestartSec=3
User=www-data
Group=www-data

[Install]
WantedBy=multi-user.target
```

保存为：

```text
/etc/systemd/system/grsai-newapi.service
```

启动：

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now grsai-newapi
sudo systemctl status grsai-newapi
```

## 构建文档站

```bash
cd /data/transfer-api/grsai-newapi
npm install
npm run docs:build
```

产物目录：

```text
docs/.vitepress/dist
```

可以交给 Nginx、Caddy 或对象存储托管。

## Nginx 示例

```nginx
server {
    listen 80;
    server_name docs.example.com;

    root /data/transfer-api/grsai-newapi/docs/.vitepress/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```
