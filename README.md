# GRSAI NewAPI

OpenAI-compatible API adapter for GRS AI chat and image models.

## API service

```bash
GRSAI_API_KEY=sk-your-key go run .
```

Without `.env`, the default server is:

```text
http://localhost:8080
```

With `SERVER_PORT=3456` in `.env`, the server is:

```text
http://localhost:3456
```

## Docker

Remote machines do not need Go installed. Docker builds the Go binary inside the image.

Create `.env` in the project root:

```dotenv
GRSAI_API_KEY=sk-your-grsai-key
GRSAI_BASE_URL=https://grsai.dakka.com.cn
SERVER_PORT=3456

# Optional client-side auth for callers of this proxy.
PROXY_API_KEY=local-dev-token

# Optional RustFS/S3 image transfer.
RUSTFS_ENDPOINT=https://fs.sendi.wang
RUSTFS_PUBLIC_BASE_URL=https://fs.sendi.wang
RUSTFS_ACCESS_KEY=admin
RUSTFS_SECRET_KEY=your-rustfs-secret
RUSTFS_BUCKET=apiLyncr
RUSTFS_REGION=us-east-1
RUSTFS_PREFIX=lyncr
RUSTFS_MAX_IMAGE_BYTES=52428800
RUSTFS_PUBLIC_READ=true
```

Build and start:

```bash
docker compose up -d --build
```

View logs:

```bash
docker compose logs -f
```

Test the service:

```bash
curl http://localhost:3456/v1/models
```

If `PROXY_API_KEY` is set, include it:

```bash
curl http://localhost:3456/v1/models \
  -H "Authorization: Bearer local-dev-token"
```

Stop:

```bash
docker compose down
```

## Linux binary

Build a Linux binary for an x86_64 server:

```bash
sh scripts/build-linux.sh amd64
```

Build a Linux binary for an ARM64 server:

```bash
sh scripts/build-linux.sh arm64
```

Copy the generated binary and `.env` to the server, then run:

```bash
chmod +x grsai-api
nohup ./grsai-api > server.log 2>&1 &
```

## Documentation

```bash
npm install
npm run docs:dev
```

Build static docs:

```bash
npm run docs:build
```
# grsai-repo
