# gau-authorization-service

## Introduction | Giới thiệu

**English:**  
This repository provides an authorization service written in Go, designed to manage authentication, token refresh, and user permissions. It is suitable for microservices architectures and can be deployed using Docker or Kubernetes.

**Tiếng Việt:**  
Repo này cung cấp dịch vụ xác thực và phân quyền viết bằng Go, dùng để quản lý xác thực, làm mới token và quyền người dùng. Phù hợp với kiến trúc microservices và có thể triển khai bằng Docker hoặc Kubernetes.

---

## Directory Structure | Cấu trúc thư mục

```
tnqbao-gau-authorization-service/
├── Dockerfile
├── entrypoint.sh
├── go.mod
├── go.sum
├── init-env.sh
├── main.go
├── .dockerignore
├── .env.example
├── config/
│   ├── env_config.go
│   └── main.go
├── controller/
│   ├── main.go
│   └── refresh_token.go
├── deploy/
│   ├── docker/
│   │   ├── docker-compose.yml
│   │   ├── keep
│   │   └── .env.example
│   └── k8s/
│       ├── production/
│       │   ├── apply.sh
│       │   ├── apply_envsubst.sh
│       │   ├── kustomization.yaml
│       │   ├── unapply.sh
│       │   ├── .env.example
│       │   ├── base/
│       │   │   └── keep
│       │   └── template/
│       │       ├── configmap.yaml
│       │       ├── deployment.yaml
│       │       ├── hpa.yaml
│       │       ├── secret.yaml
│       │       └── service.yaml
│       └── staging/
│           ├── apply.sh
│           ├── apply_envsubst.sh
│           ├── kustomization.yaml
│           ├── unapply.sh
│           ├── .env.example
│           ├── base/
│           │   └── keep
│           └── template/
│               ├── configmap.yaml
│               ├── deployment.yaml
│               ├── hpa.yaml
│               ├── secret.yaml
│               └── service.yaml
├── entity/
│   └── refesh_token.go
├── infra/
│   ├── main.go
│   ├── posgrest.go
│   └── redis.go
├── middlewares/
│   ├── cors.go
│   └── main.go
├── migrations/
│   ├── 000001_create_refresh_tokens.down.sql
│   └── 000001_create_refresh_tokens.up.sql
├── repository/
│   ├── bitmap.go
│   ├── main.go
│   └── token.go
└── routes/
    └── routes.go
```

### 📑 Directory Description | Mô tả thư mục

| Path | Description | Mô tả |
|------|-------------|--------|
| `Dockerfile`, `entrypoint.sh` | Docker image build and startup script | File build và khởi động Docker |
| `go.mod`, `go.sum` | Go module definitions | Định nghĩa module Go |
| `init-env.sh` | Script to generate `.env` from `.env.example` | Tạo file `.env` từ `.env.example` |
| `.env.example` | Sample environment variables | Biến môi trường mẫu |
| `config/` | Environment loading and configuration logic | Logic cấu hình và load môi trường |
| `controller/` | HTTP handlers (e.g., refresh token) | Xử lý HTTP |
| `deploy/docker/` | Docker Compose setup | Cấu hình triển khai với Docker |
| `deploy/k8s/` | Kubernetes manifests and scripts for staging/production | Manifest và script triển khai trên K8s |
| `entity/` | Domain models (e.g., refresh tokens) | Các model dữ liệu chính |
| `infra/` | PostgreSQL, Redis setup and connections | Thiết lập DB và Redis |
| `middlewares/` | CORS and other middleware logic | Middleware |
| `migrations/` | SQL migration files | Các file migration SQL |
| `repository/` | Data access and business logic | Truy cập và xử lý dữ liệu |
| `routes/` | API route definitions | Định nghĩa route |

---

## Deployment | Triển khai

### 🧪 Init Environment | Khởi tạo môi trường

```bash
./init-env.sh
```

> Tạo file `.env` từ `.env.example`, nếu chưa có. Bạn nên kiểm tra lại các giá trị bên trong `.env` trước khi chạy tiếp.

---

### 🐳 Docker

**English:**  
1. Build the Docker image:
   ```bash
   docker build -t gau-authorization-service .
   ```
2. Run with Docker Compose:
   ```bash
   cd deploy/docker
   docker-compose up -d
   ```

**Tiếng Việt:**  
1. Build image Docker:
   ```bash
   docker build -t gau-authorization-service .
   ```
2. Chạy với Docker Compose:
   ```bash
   cd deploy/docker
   docker-compose up -d
   ```

---

### ☸ Kubernetes

**English:**  
1. Edit environment variables in `deploy/k8s/staging/template/configmap.yaml` and `secret.yaml`.
2. Apply manifests:
   ```bash
   cd deploy/k8s/staging
   ./apply.sh
   ```
3. To remove:
   ```bash
   ./unapply.sh
   ```

**Tiếng Việt:**  
1. Chỉnh sửa biến môi trường trong `deploy/k8s/staging/template/configmap.yaml` và `secret.yaml`.
2. Áp dụng manifest:
   ```bash
   cd deploy/k8s/staging
   ./apply.sh
   ```
3. Để xóa:
   ```bash
   ./unapply.sh
   ```

---

## Contact | Liên hệ

For questions or support, please contact the maintainer.  
Để được hỗ trợ, vui lòng liên hệ người quản lý repo.