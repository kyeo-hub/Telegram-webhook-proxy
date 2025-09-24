# Telegram Bot API 反向代理

这是一个简单的 Telegram Bot API 反向代理服务器，用于代理对 https://api.telegram.org 的请求。它支持 API 密钥验证以防止滥用，并通过环境变量保护 Telegram Token。

## 功能

- 反向代理 Telegram Bot API
- API 密钥验证防止滥用
- 使用环境变量保护 Telegram Token
- 可自定义端口
- 支持Docker部署

## 使用方法

### 1. 设置环境变量

创建 `.env` 文件或直接设置环境变量：

```bash
export TELEGRAM_TOKEN="your_telegram_bot_token"
export API_KEYS="key1,key2,key3"
export PORT=8080
```

或者创建 `.env` 文件：

```env
TELEGRAM_TOKEN=your_telegram_bot_token
API_KEYS=key1,key2,key3
PORT=8080
```

### 2. 编译

```bash
go build -o telegram-proxy
```

或者使用Makefile：

```bash
make build
```

### 3. 运行

```bash
./telegram-proxy
```

或者使用环境变量直接运行：

```bash
TELEGRAM_TOKEN="your_telegram_bot_token" API_KEYS="key1,key2,key3" ./telegram-proxy
```

## Docker部署

### 构建镜像

```bash
docker build -t telegram-proxy .
```

### 运行容器

```bash
docker run -d \
  --name telegram-proxy \
  -p 8080:8080 \
  -e TELEGRAM_TOKEN="your_telegram_bot_token" \
  -e API_KEYS="key1,key2,key3" \
  telegram-proxy
```

### 使用docker-compose

创建 `docker-compose.yml` 文件：

```yaml
version: '3.8'
services:
  telegram-proxy:
    build: .
    ports:
      - "8080:8080"
    environment:
      - TELEGRAM_TOKEN=your_telegram_bot_token
      - API_KEYS=key1,key2,key3
    restart: unless-stopped
```

然后运行：

```bash
docker-compose up -d
```

## 使用代理

代理服务器运行后，你可以通过以下方式访问 Telegram Bot API：

```
# 原始请求
https://api.telegram.org/bot<TOKEN>/METHOD

# 通过代理请求 (需要添加 X-API-Key 头)
http://localhost:8080/bot<TOKEN>/METHOD
```

或者在 URL 中传递 API Key：

```
http://localhost:8080/bot<TOKEN>/METHOD?api_key=your_api_key
```

### 示例

获取 bot 信息：

```bash
curl -H "X-API-Key: key1" http://localhost:8080/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getMe
```

或者：

```bash
curl "http://localhost:8080/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getMe?api_key=key1"
```

发送消息示例：

```bash
curl -X POST "http://localhost:8080/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/sendMessage?api_key=key1" \
  -d "chat_id=123456789&text=Hello World"
```

## API 密钥

为了防止滥用，所有请求都需要提供有效的 API 密钥。可以通过以下两种方式之一提供：

1. 在请求头中添加 `X-API-Key: your_api_key`
2. 在查询参数中添加 `api_key=your_api_key`

## 配置项

| 环境变量 | 必需 | 描述 |
|---------|------|------|
| TELEGRAM_TOKEN | 是 | 你的 Telegram Bot Token |
| API_KEYS | 是 | API 密钥列表，用逗号分隔 |
| PORT | 否 | 服务器端口，默认为 8080 |

## 健康检查

访问 `/health` 端点可以检查服务器是否正常运行：

```bash
curl http://localhost:8080/health
```