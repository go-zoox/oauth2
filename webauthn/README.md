# WebAuthn OAuth2 Provider

这个包为 [go-zoox/oauth2](https://github.com/go-zoox/oauth2) 库提供 WebAuthn（无密码认证）支持。

## 🚀 特性

- 🔐 **无密码认证** - 使用生物识别、硬件密钥或 PIN 码登录
- 🛡️ **更高安全性** - 基于公钥密码学，防钓鱼攻击
- 📱 **多平台支持** - 支持指纹、面容识别、Windows Hello、Touch ID 等
- 🔑 **硬件密钥** - 支持 YubiKey、SoloKey 等 FIDO2 设备
- ⚡ **快速登录** - 比传统密码更快的认证体验
- 🎯 **标准兼容** - 完全符合 W3C WebAuthn 标准

## 📋 前提条件

1. Go 1.19+ 
2. 支持 WebAuthn 的浏览器（Chrome 67+、Firefox 60+、Safari 14+、Edge 18+）
3. HTTPS 环境（生产环境必须）

## 🔧 安装

```bash
go get github.com/go-zoox/oauth2/webauthn
go get github.com/go-webauthn/webauthn
```

## 🏗️ 配置

### 1. 基本配置

```go
import (
    "github.com/go-zoox/oauth2/webauthn"
)

// 创建用户和会话存储
userStore := webauthn.NewSimpleUserStore()
sessionStore := webauthn.NewSimpleSessionStore()

// 创建 WebAuthn 客户端
client, err := webauthn.New(&webauthn.WebAuthnConfig{
    // 基本 OAuth2 配置
    ClientID:      "your-app-id",
    ClientSecret:  "your-app-secret", 
    RedirectURI:   "https://yourdomain.com/auth/callback",
    Scope:         "webauthn",

    // WebAuthn 特定配置
    RPDisplayName: "您的应用名称",           // 显示给用户的应用名称
    RPID:          "yourdomain.com",      // 您的域名
    RPOrigins:     []string{"https://yourdomain.com"}, // 允许的源
    
    // 存储接口
    UserStore:     userStore,
    SessionStore:  sessionStore,
    
    // 可选: 超时设置（毫秒）
    Timeout:       60000, // 60 秒
})
```

### 2. 环境变量配置

```bash
# WebAuthn 配置
export WEBAUTHN_RP_DISPLAY_NAME="您的应用名称"
export WEBAUTHN_RP_ID="yourdomain.com"
export WEBAUTHN_RP_ORIGIN="https://yourdomain.com"

# 可选: 端口设置
export PORT=8080
```

## 💡 使用方法

### 基本示例

```go
package main

import (
    "log"
    "net/http"
    
    "github.com/go-zoox/oauth2"
    "github.com/go-zoox/oauth2/webauthn"
)

func main() {
    // 初始化存储
    userStore := webauthn.NewSimpleUserStore()
    sessionStore := webauthn.NewSimpleSessionStore()

    // 创建 WebAuthn 客户端
    client, err := webauthn.New(&webauthn.WebAuthnConfig{
        ClientID:      "demo-app",
        ClientSecret:  "demo-secret",
        RedirectURI:   "http://localhost:8080/auth/callback",
        RPDisplayName: "WebAuthn Demo",
        RPID:          "localhost",
        RPOrigins:     []string{"http://localhost:8080"},
        UserStore:     userStore,
        SessionStore:  sessionStore,
    })
    if err != nil {
        log.Fatal("创建 WebAuthn 客户端失败:", err)
    }

    // 注册路由
    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/login", loginHandler)
    
    log.Println("服务器启动在 http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 注册新用户

```go
func registerHandler(w http.ResponseWriter, r *http.Request) {
    // 获取 WebAuthn 客户端
    webauthnClient := client.(*webauthn.client)
    
    // 开始注册流程
    options, sessionID, err := webauthnClient.BeginRegistration(
        "user123",           // 用户 ID
        "user@example.com",  // 用户名
        "用户显示名称",        // 显示名称
    )
    if err != nil {
        http.Error(w, "注册初始化失败", http.StatusInternalServerError)
        return
    }
    
    // 将选项发送给前端
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(options)
    
    // 存储会话 ID（用于后续验证）
    // 在实际应用中，您可能会使用 HTTP 会话或其他机制
}
```

### 完成注册

```go
func finishRegistrationHandler(w http.ResponseWriter, r *http.Request) {
    // 解析前端发送的凭据响应
    var credentialResponse []byte
    // ... 从请求中解析凭据数据
    
    // 获取 WebAuthn 客户端
    webauthnClient := client.(*webauthn.client)
    
    // 完成注册
    err := webauthnClient.FinishRegistration(
        "user123",           // 用户 ID
        sessionID,           // 会话 ID
        credentialResponse,  // 凭据响应
    )
    if err != nil {
        http.Error(w, "注册失败", http.StatusUnauthorized)
        return
    }
    
    // 注册成功
    w.WriteHeader(http.StatusOK)
}
```

### 用户登录

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
    // 获取 WebAuthn 客户端
    webauthnClient := client.(*webauthn.client)
    
    // 开始登录流程
    options, sessionID, err := webauthnClient.BeginLogin("user123")
    if err != nil {
        http.Error(w, "登录初始化失败", http.StatusInternalServerError)
        return
    }
    
    // 将选项发送给前端
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(options)
}
```

## 🌐 前端集成

### JavaScript 示例

```javascript
// 注册新用户
async function register(username, displayName) {
    try {
        // 1. 开始注册
        const response = await fetch('/register/begin', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, displayName })
        });
        
        const options = await response.json();
        
        // 2. 转换数据格式
        options.publicKey.challenge = base64urlDecode(options.publicKey.challenge);
        options.publicKey.user.id = base64urlDecode(options.publicKey.user.id);
        
        // 3. 创建凭据
        const credential = await navigator.credentials.create(options);
        
        // 4. 完成注册
        await fetch('/register/finish', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                username,
                credential: {
                    id: credential.id,
                    rawId: base64urlEncode(credential.rawId),
                    type: credential.type,
                    response: {
                        attestationObject: base64urlEncode(credential.response.attestationObject),
                        clientDataJSON: base64urlEncode(credential.response.clientDataJSON)
                    }
                }
            })
        });
        
        alert('注册成功！');
    } catch (error) {
        console.error('注册失败:', error);
    }
}

// 用户登录
async function login(username) {
    try {
        // 1. 开始登录
        const response = await fetch('/login/begin', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username })
        });
        
        const options = await response.json();
        
        // 2. 转换数据格式
        options.publicKey.challenge = base64urlDecode(options.publicKey.challenge);
        if (options.publicKey.allowCredentials) {
            options.publicKey.allowCredentials.forEach(cred => {
                cred.id = base64urlDecode(cred.id);
            });
        }
        
        // 3. 获取断言
        const assertion = await navigator.credentials.get(options);
        
        // 4. 完成登录
        await fetch('/login/finish', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                username,
                assertion: {
                    id: assertion.id,
                    rawId: base64urlEncode(assertion.rawId),
                    type: assertion.type,
                    response: {
                        authenticatorData: base64urlEncode(assertion.response.authenticatorData),
                        clientDataJSON: base64urlEncode(assertion.response.clientDataJSON),
                        signature: base64urlEncode(assertion.response.signature),
                        userHandle: assertion.response.userHandle ? base64urlEncode(assertion.response.userHandle) : null
                    }
                }
            })
        });
        
        alert('登录成功！');
    } catch (error) {
        console.error('登录失败:', error);
    }
}

// 工具函数
function base64urlDecode(str) {
    return Uint8Array.from(atob(str.replace(/-/g, '+').replace(/_/g, '/')), c => c.charCodeAt(0));
}

function base64urlEncode(buffer) {
    return btoa(String.fromCharCode(...new Uint8Array(buffer)))
        .replace(/\+/g, '-')
        .replace(/\//g, '_')
        .replace(/=/g, '');
}
```

## 🗃️ 存储接口

### 自定义用户存储

```go
type CustomUserStore struct {
    // 您的数据库连接或其他存储
    db *sql.DB
}

func (s *CustomUserStore) GetUser(userID string) (webauthn.WebAuthnUser, error) {
    // 从数据库获取用户
    // 返回实现了 WebAuthnUser 接口的用户对象
}

func (s *CustomUserStore) CreateUser(userID, username, displayName string) (webauthn.WebAuthnUser, error) {
    // 在数据库中创建新用户
}

func (s *CustomUserStore) UpdateUser(user webauthn.WebAuthnUser) error {
    // 更新数据库中的用户信息
}
```

### 自定义会话存储

```go
type CustomSessionStore struct {
    // Redis 或其他会话存储
    redis *redis.Client
}

func (s *CustomSessionStore) StoreSession(sessionID string, data *webauthn.SessionData) error {
    // 存储会话数据到 Redis
}

func (s *CustomSessionStore) GetSession(sessionID string) (*webauthn.SessionData, error) {
    // 从 Redis 获取会话数据
}

func (s *CustomSessionStore) DeleteSession(sessionID string) error {
    // 删除 Redis 中的会话数据
}
```

## 🔒 安全考虑

### 1. HTTPS 要求
```bash
# 生产环境必须使用 HTTPS
# 开发环境可以使用 localhost
```

### 2. 域名配置
```go
// 确保 RPID 和 RPOrigins 配置正确
RPIDclient, err := webauthn.New(&webauthn.WebAuthnConfig{
    RPID:      "yourdomain.com",  // 不包含协议和端口
    RPOrigins: []string{
        "https://yourdomain.com",     // 生产环境
        "https://www.yourdomain.com", // www 子域名
    },
})
```

### 3. 会话管理
```go
// 使用安全的会话存储
// 设置适当的会话超时
// 实施会话固定保护
```

### 4. 用户验证
```go
// 实施适当的用户验证要求
// 考虑使用用户验证首选项
```

## 🎯 完整示例

查看 [example/webauthn](../example/webauthn/) 目录获取完整的工作示例，包括：

- 完整的 HTML 界面（中文）
- JavaScript WebAuthn 客户端代码
- 注册和登录流程
- 错误处理
- 用户体验优化

### 运行示例

```bash
# 1. 进入示例目录
cd example/webauthn

# 2. 设置环境变量（可选）
export WEBAUTHN_RP_DISPLAY_NAME="WebAuthn 演示"
export WEBAUTHN_RP_ID="localhost"
export WEBAUTHN_RP_ORIGIN="http://localhost:8080"

# 3. 运行示例
go run main.go

# 4. 访问 http://localhost:8080
```

## 🌍 浏览器支持

| 浏览器 | 版本要求 | 支持的验证器 |
|--------|----------|-------------|
| Chrome | 67+ | 平台验证器、USB 安全密钥、BLE |
| Firefox | 60+ | 平台验证器、USB 安全密钥 |
| Safari | 14+ | Touch ID、Face ID、USB 安全密钥 |
| Edge | 18+ | Windows Hello、USB 安全密钥 |

## 🔧 故障排除

### 常见问题

1. **"此操作已超时"**
   - 检查超时设置
   - 确保用户在限定时间内完成操作

2. **"不支持的设备"**
   - 确认浏览器支持 WebAuthn
   - 检查设备是否有可用的验证器

3. **"域名不匹配"**
   - 检查 RPID 和 RPOrigins 配置
   - 确保在正确的域名下运行

4. **HTTPS 错误**
   - 生产环境必须使用 HTTPS
   - 开发环境可以使用 localhost

### 调试技巧

```javascript
// 启用详细日志
console.log('WebAuthn options:', options);
console.log('Credential created:', credential);

// 检查浏览器支持
if (!window.PublicKeyCredential) {
    console.error('WebAuthn 不受支持');
}
```

## 📚 扩展阅读

- [W3C WebAuthn 规范](https://www.w3.org/TR/webauthn/)
- [FIDO Alliance](https://fidoalliance.org/)
- [WebAuthn.io 演示](https://webauthn.io/)
- [Can I Use WebAuthn](https://caniuse.com/webauthn)

## 🤝 贡献

欢迎贡献代码！请查看主项目的贡献指南。

## 📄 许可证

本包是 [go-zoox/oauth2](https://github.com/go-zoox/oauth2) 库的一部分，遵循相同的许可条款。