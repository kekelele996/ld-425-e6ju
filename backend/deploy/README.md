# 部署说明

本目录用于存放部署相关文件。

- 根目录 `docker-compose.yml` 为推荐的一键部署方式。
- `backend/Dockerfile` 使用 Go 多阶段构建，产物为静态二进制。
- 如需 Kubernetes 部署，可基于 `backend/Dockerfile` 构建镜像后挂载环境变量与命名卷。
