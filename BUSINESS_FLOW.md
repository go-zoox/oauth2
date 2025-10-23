# acme.sh 业务流程图

## 整体架构流程

```mermaid
graph TD
    A[用户命令] --> B{命令类型}
    
    B -->|--issue| C[证书申请流程]
    B -->|--renew| D[证书续期流程]
    B -->|--install-cert| E[证书安装流程]
    B -->|--deploy| F[证书部署流程]
    B -->|--install| G[安装 acme.sh]
    
    C --> H[域名验证]
    D --> H
    H --> I{验证方式}
    
    I -->|HTTP-01| J[Webroot/Standalone验证]
    I -->|DNS-01| K[DNS验证]
    I -->|TLS-ALPN-01| L[TLS-ALPN验证]
    
    J --> M[CA验证通过]
    K --> M
    L --> M
    
    M --> N[生成证书]
    N --> O[保存证书文件]
    O --> P[执行钩子脚本]
    
    E --> Q[复制证书到目标位置]
    Q --> R[重载服务]
    
    F --> S[调用部署钩子]
    S --> T[部署到具体服务]
```

## 证书申请详细流程

```mermaid
sequenceDiagram
    participant U as 用户
    participant A as acme.sh
    participant CA as Certificate Authority
    participant DNS as DNS Provider
    participant WEB as Web Server
    
    U->>A: acme.sh --issue -d domain.com
    A->>A: 检查配置和参数
    A->>A: 生成账户密钥(如果不存在)
    A->>CA: 注册账户
    CA-->>A: 账户注册成功
    
    A->>CA: 请求证书订单
    CA-->>A: 返回验证挑战
    
    alt DNS验证模式
        A->>DNS: 添加 TXT 记录
        DNS-->>A: 记录添加成功
        A->>CA: 通知验证就绪
        CA->>DNS: 验证 TXT 记录
        DNS-->>CA: 验证通过
    else HTTP验证模式
        A->>WEB: 创建验证文件
        WEB-->>A: 文件创建成功
        A->>CA: 通知验证就绪
        CA->>WEB: 访问验证文件
        WEB-->>CA: 返回验证内容
    end
    
    CA-->>A: 验证通过
    A->>CA: 请求签发证书
    CA-->>A: 返回证书
    A->>A: 保存证书文件
    A->>A: 执行后置钩子
    A-->>U: 证书申请完成
```

## DNS API 验证流程

```mermaid
graph TD
    A[开始DNS验证] --> B[检查DNS API配置]
    B --> C{API类型}
    
    C -->|Cloudflare| D[使用CF API]
    C -->|阿里云| E[使用阿里云API]
    C -->|其他| F[使用对应API]
    
    D --> G[获取Zone ID]
    E --> G
    F --> G
    
    G --> H[添加TXT记录]
    H --> I[等待DNS传播]
    I --> J[通知CA验证]
    J --> K[CA验证DNS记录]
    K --> L{验证结果}
    
    L -->|成功| M[删除TXT记录]
    L -->|失败| N[报告错误]
    
    M --> O[验证完成]
    N --> P[验证失败]
```

## 证书续期流程

```mermaid
graph TD
    A[Cron定时任务] --> B[扫描所有证书]
    B --> C{证书是否需要续期}
    
    C -->|是| D[读取证书配置]
    C -->|否| E[跳过此证书]
    
    D --> F[执行续期前钩子]
    F --> G[调用issue函数]
    G --> H{续期是否成功}
    
    H -->|成功| I[执行续期后钩子]
    H -->|失败| J[记录错误日志]
    
    I --> K[重载相关服务]
    J --> L[发送通知]
    
    E --> M[检查下一个证书]
    K --> M
    L --> M
    
    M --> N{还有证书?}
    N -->|是| C
    N -->|否| O[续期任务完成]
```

## 部署钩子系统

```mermaid
graph TD
    A[证书生成完成] --> B[检查部署钩子]
    B --> C{钩子类型}
    
    C -->|nginx| D[nginx部署脚本]
    C -->|apache| E[apache部署脚本]
    C -->|docker| F[docker部署脚本]
    C -->|自定义| G[自定义部署脚本]
    
    D --> H[复制证书文件]
    E --> H
    F --> H
    G --> H
    
    H --> I[更新配置文件]
    I --> J[重载/重启服务]
    J --> K[验证部署结果]
    K --> L[部署完成]
```

## 主要数据结构

### 证书配置文件结构
```
~/.acme.sh/domain.com/
├── domain.com.conf          # 证书配置
├── domain.com.key           # 私钥
├── domain.com.cer           # 证书
├── ca.cer                   # CA证书
└── fullchain.cer            # 完整证书链
```

### 配置参数
- `Le_Domain`: 主域名
- `Le_Alt`: 备用域名
- `Le_Webroot`: Web根目录
- `Le_PreHook`: 前置钩子
- `Le_PostHook`: 后置钩子
- `Le_RenewHook`: 续期钩子
- `Le_DeployHook`: 部署钩子

## 关键业务逻辑

1. **ACME协议实现**: 完整实现 RFC 8555 ACME 协议
2. **多CA支持**: 支持 Let's Encrypt, ZeroSSL, SSL.com 等多个CA
3. **自动化续期**: 默认60天自动续期机制
4. **灵活验证**: 支持多种域名验证方式
5. **钩子系统**: 丰富的钩子机制支持自动化部署
6. **跨平台**: 支持多种操作系统和Shell环境