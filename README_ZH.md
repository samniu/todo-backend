# Todo应用后端

一款使用Go语言开发的Todo应用后端服务，提供RESTful API和WebSocket实时通知功能。

*[English](README.md)*

## 概述

该后端服务为Todo应用提供了强大的API支持，具有用户认证、任务管理和实时通知功能。它设计用于支持在线和离线操作，当网络连接恢复时能够无缝同步数据。

## 功能特性

- 🔒 **用户认证**: 使用JWT实现安全的注册和登录系统
- ✅ **任务管理**: 完整的CRUD操作支持任务创建、查询、更新和删除
- 🔔 **实时通知**: WebSocket支持即时更新
- 🔄 **离线同步支持**: 特殊的API设计用于处理客户端离线操作

## 快速开始

### 前提条件

- Go 1.18+
- PostgreSQL/MySQL (可配置)

### 安装

1. 克隆仓库

```bash
git https://github.com/samniu/todo-backend.git
cd todo-backend
```

2. 设置环境变量

创建一个`.env`文件:

```
DB_CONNECTION=postgres://username:password@localhost:5432/todo_db
JWT_SECRET=your_jwt_secret_key
PORT=8080
```

3. 运行服务器

```bash
go run cmd/main.go
```

服务器将在`http://localhost:8080`启动。

## 许可证

此项目采用 [MIT 许可证](LICENSE)。
