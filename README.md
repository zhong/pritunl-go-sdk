# Pritunl Go SDK

基于 Pritunl Web 后端 API 封装的 Go SDK，支持开源版通过管理员 API Token 进行认证。

## 特性

- 纯标准库实现，无第三方依赖
- 自动完成 `Auth-Token` / `Auth-Timestamp` / `Auth-Nonce` / `Auth-Signature` HMAC-SHA256 签名
- 支持跳过 TLS 验证（方便自签名证书环境）
- 封装了常用接口：Status、Organization、User、Server、Key

## 安装

```bash
go get github.com/example/pritunl-go-sdk
```

## 快速开始

```go
package main

import (
    "fmt"
    "log"

    pritunl "github.com/example/pritunl-go-sdk"
)

func main() {
    client := pritunl.NewClient(
        "https://your-pritunl-server",
        "YOUR_API_TOKEN",
        "YOUR_API_SECRET",
        true, // 跳过 TLS 验证
    )

    status, err := client.GetStatus()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Pritunl version: %s\n", status.ServerVersion)
}
```

完整示例见 [`example/main.go`](example/main.go)。

## 运行示例

```bash
cd example
go run main.go \
  -base https://your-pritunl-server \
  -token YOUR_API_TOKEN \
  -secret YOUR_API_SECRET \
  -insecure=true
```

或者使用环境变量：

```bash
export PRITUNL_API_TOKEN=your_token
export PRITUNL_API_SECRET=your_secret
go run main.go -base https://your-pritunl-server
```

## API 覆盖

### Status

- `GetStatus() (*Status, error)`

### Organization

- `ListOrganizations() ([]Organization, error)`
- `GetOrganization(id string) (*Organization, error)`
- `CreateOrganization(name string) (*Organization, error)`
- `UpdateOrganization(id, name string) (*Organization, error)`
- `DeleteOrganization(id string) error`
- `FindOrganizationByName(name string) (*Organization, error)`

### User

- `ListUsers(orgID string) ([]User, error)`
- `GetUser(orgID, userID string) (*User, error)`
- `CreateUser(orgID string, req CreateUserRequest) (*User, error)`
- `UpdateUser(orgID, userID string, req CreateUserRequest) (*User, error)`
- `DeleteUser(orgID, userID string) error`
- `FindUserByName(orgID, name string) (*User, error)`

### Server

- `ListServers() ([]Server, error)`
- `GetServer(id string) (*Server, error)`
- `CreateServer(req CreateServerRequest) (*Server, error)`
- `UpdateServer(id string, req CreateServerRequest) (*Server, error)`
- `DeleteServer(id string) error`
- `ServerOperation(id, operation string) error`
- `AttachOrganization(serverID, orgID string) error`
- `DetachOrganization(serverID, orgID string) error`
- `AddRoute(serverID string, route ServerRoute) error`
- `DeleteRoute(serverID, network string) error`

### Key / Profile

- `GetKeyLink(orgID, userID string) (*KeyLink, error)`
- `DownloadKeyArchive(orgID, userID, format string) ([]byte, error)`
- `DownloadKeyConfig(orgID, userID, serverID string) ([]byte, error)`
- `DownloadLinkedKeyArchive(keyID, format string) ([]byte, error)`

## 签名说明

Pritunl API 使用 HMAC-SHA256 签名，签名原文为：

```
API_TOKEN & TIMESTAMP & NONCE & METHOD & PATH
```

其中 `PATH` 不包含查询参数。SDK 已在每次请求时自动构造该签名。

## 开源版限制

- 此 SDK 只能调用开源版已开放的后端接口
- Enterprise 专属功能（如 Site-to-site Links、SSO 等）仍需要有效 License
- 管理员 API Token 默认在 Web UI 中不可见，需要通过数据库脚本开启，参考仓库中的 `enable_pritunl_api_token_nopymongo.py`
