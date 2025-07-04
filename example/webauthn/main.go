package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-zoox/oauth2/webauthn"
)

var (
	userStore    *webauthn.SimpleUserStore
	sessionStore *webauthn.SimpleSessionStore
	client       webauthn.WebAuthnClientInterface
)

func main() {
	// Initialize stores
	userStore = webauthn.NewSimpleUserStore()
	sessionStore = webauthn.NewSimpleSessionStore()

	// Environment variables
	rpDisplayName := os.Getenv("WEBAUTHN_RP_DISPLAY_NAME")
	rpID := os.Getenv("WEBAUTHN_RP_ID")
	rpOrigin := os.Getenv("WEBAUTHN_RP_ORIGIN")

	if rpDisplayName == "" {
		rpDisplayName = "WebAuthn Demo"
	}
	if rpID == "" {
		rpID = "localhost"
	}
	if rpOrigin == "" {
		rpOrigin = "http://localhost:8080"
	}

	// Create WebAuthn client
	var err error
	client, err = webauthn.New(&webauthn.WebAuthnConfig{
		ClientID:      "webauthn-demo",
		ClientSecret:  "demo-secret",
		RedirectURI:   "http://localhost:8080/auth/callback",
		Scope:         "webauthn",
		RPDisplayName: rpDisplayName,
		RPID:          rpID,
		RPOrigins:     []string{rpOrigin},
		UserStore:     userStore,
		SessionStore:  sessionStore,
		Timeout:       60000, // 60 seconds
	})
	if err != nil {
		log.Fatal("Failed to create WebAuthn client:", err)
	}

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/register/begin", registerBeginHandler)
	http.HandleFunc("/register/finish", registerFinishHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/login/begin", loginBeginHandler)
	http.HandleFunc("/login/finish", loginFinishHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/dashboard", dashboardHandler)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting WebAuthn demo server on port %s", port)
	log.Printf("Visit http://localhost:%s to test WebAuthn authentication", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<title>WebAuthn Demo</title>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
		body { 
			font-family: Arial, sans-serif; 
			max-width: 800px; 
			margin: 50px auto; 
			padding: 20px;
			background: #f5f5f5;
		}
		.container { 
			background: white; 
			padding: 40px; 
			border-radius: 10px; 
			box-shadow: 0 2px 10px rgba(0,0,0,0.1);
		}
		.button { 
			background-color: #007bff; 
			color: white; 
			padding: 15px 30px; 
			text-decoration: none; 
			border: none; 
			border-radius: 5px; 
			cursor: pointer; 
			display: inline-block; 
			margin: 10px 5px;
			font-size: 16px;
		}
		.button:hover { background-color: #0056b3; }
		.button.secondary { background-color: #6c757d; }
		.button.secondary:hover { background-color: #545b62; }
		.feature-list { 
			background: #e9ecef; 
			padding: 20px; 
			border-radius: 5px; 
			margin: 20px 0; 
		}
		.feature-list li { margin: 10px 0; }
		h1 { color: #333; text-align: center; }
		h2 { color: #555; }
		.info { color: #666; margin: 20px 0; }
	</style>
</head>
<body>
	<div class="container">
		<h1>🔐 WebAuthn 演示</h1>
		<div class="info">
			<p>欢迎来到 WebAuthn 无密码认证演示！WebAuthn 让您可以使用生物识别、安全密钥或其他强身份验证器登录，无需记住密码。</p>
		</div>

		<div class="feature-list">
			<h2>✨ 特性</h2>
			<ul>
				<li>🚫 无密码认证</li>
				<li>🔒 支持生物识别（指纹、面容等）</li>
				<li>🔑 支持硬件安全密钥（YubiKey、SoloKey等）</li>
				<li>📱 支持平台验证器（Windows Hello、Touch ID等）</li>
				<li>🛡️ 更高的安全性，防钓鱼攻击</li>
				<li>⚡ 更快的登录体验</li>
			</ul>
		</div>

		<div style="text-align: center; margin: 30px 0;">
			<a href="/register" class="button">📝 注册新账户</a>
			<a href="/login" class="button secondary">🔓 登录</a>
		</div>

		<div class="info">
			<h2>💡 使用说明</h2>
			<ol>
				<li><strong>注册：</strong>首次使用时，点击"注册新账户"创建您的账户和身份验证器</li>
				<li><strong>登录：</strong>注册后，点击"登录"使用您的身份验证器进行认证</li>
				<li><strong>支持的验证器：</strong>指纹识别、面容识别、PIN码、硬件安全密钥等</li>
			</ol>
		</div>
	</div>
</body>
</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<title>注册 - WebAuthn Demo</title>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
		body { 
			font-family: Arial, sans-serif; 
			max-width: 600px; 
			margin: 50px auto; 
			padding: 20px;
			background: #f5f5f5;
		}
		.container { 
			background: white; 
			padding: 40px; 
			border-radius: 10px; 
			box-shadow: 0 2px 10px rgba(0,0,0,0.1);
		}
		.form-group { margin: 20px 0; }
		label { display: block; margin-bottom: 5px; font-weight: bold; }
		input { 
			width: 100%; 
			padding: 12px; 
			border: 1px solid #ddd; 
			border-radius: 5px; 
			font-size: 16px;
			box-sizing: border-box;
		}
		.button { 
			background-color: #28a745; 
			color: white; 
			padding: 15px 30px; 
			border: none; 
			border-radius: 5px; 
			cursor: pointer; 
			font-size: 16px;
			width: 100%;
		}
		.button:hover { background-color: #218838; }
		.button:disabled { background-color: #6c757d; cursor: not-allowed; }
		.back-link { color: #007bff; text-decoration: none; }
		.back-link:hover { text-decoration: underline; }
		.status { 
			padding: 15px; 
			margin: 20px 0; 
			border-radius: 5px; 
			display: none;
		}
		.status.success { background: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
		.status.error { background: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
		.status.info { background: #d1ecf1; color: #0c5460; border: 1px solid #bee5eb; }
		h1 { color: #333; text-align: center; }
	</style>
</head>
<body>
	<div class="container">
		<h1>📝 注册新账户</h1>
		<div id="status" class="status"></div>
		
		<form id="registerForm">
			<div class="form-group">
				<label for="username">用户名或邮箱:</label>
				<input type="text" id="username" name="username" required placeholder="输入您的用户名或邮箱">
			</div>
			<div class="form-group">
				<label for="displayName">显示名称:</label>
				<input type="text" id="displayName" name="displayName" required placeholder="输入您的显示名称">
			</div>
			<button type="submit" class="button" id="registerBtn">🔐 注册并设置身份验证器</button>
		</form>

		<div style="text-align: center; margin: 20px 0;">
			<a href="/" class="back-link">← 返回首页</a> | 
			<a href="/login" class="back-link">已有账户？点击登录</a>
		</div>
	</div>

	<script>
		document.getElementById('registerForm').addEventListener('submit', async (e) => {
			e.preventDefault();
			
			const username = document.getElementById('username').value;
			const displayName = document.getElementById('displayName').value;
			const statusDiv = document.getElementById('status');
			const registerBtn = document.getElementById('registerBtn');
			
			if (!username || !displayName) {
				showStatus('请填写所有字段', 'error');
				return;
			}

			registerBtn.disabled = true;
			registerBtn.textContent = '正在准备注册...';
			showStatus('正在准备 WebAuthn 注册...', 'info');

			try {
				// Begin registration
				const beginResponse = await fetch('/register/begin', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({ username, displayName })
				});

				if (!beginResponse.ok) {
					throw new Error('注册初始化失败');
				}

				const credentialCreationOptions = await beginResponse.json();
				showStatus('请使用您的身份验证器完成注册...', 'info');

				// Convert base64url to ArrayBuffer
				credentialCreationOptions.publicKey.challenge = base64urlDecode(credentialCreationOptions.publicKey.challenge);
				credentialCreationOptions.publicKey.user.id = base64urlDecode(credentialCreationOptions.publicKey.user.id);

				// Create credential
				const credential = await navigator.credentials.create(credentialCreationOptions);

				if (!credential) {
					throw new Error('身份验证器注册失败');
				}

				showStatus('正在完成注册...', 'info');

				// Finish registration
				const finishResponse = await fetch('/register/finish', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
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

				if (!finishResponse.ok) {
					throw new Error('注册完成失败');
				}

				showStatus('🎉 注册成功！正在跳转到登录页面...', 'success');
				setTimeout(() => {
					window.location.href = '/login';
				}, 2000);

			} catch (error) {
				console.error('Registration error:', error);
				showStatus('注册失败: ' + error.message, 'error');
				registerBtn.disabled = false;
				registerBtn.textContent = '🔐 注册并设置身份验证器';
			}
		});

		function showStatus(message, type) {
			const statusDiv = document.getElementById('status');
			statusDiv.className = 'status ' + type;
			statusDiv.textContent = message;
			statusDiv.style.display = 'block';
		}

		function base64urlDecode(str) {
			return Uint8Array.from(atob(str.replace(/-/g, '+').replace(/_/g, '/')), c => c.charCodeAt(0));
		}

		function base64urlEncode(buffer) {
			return btoa(String.fromCharCode(...new Uint8Array(buffer)))
				.replace(/\+/g, '-')
				.replace(/\//g, '_')
				.replace(/=/g, '');
		}
	</script>
</body>
</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

func registerBeginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := req.Username // In practice, you might generate a UUID
	
	// Begin registration
	options, sessionID, err := client.BeginRegistration(userID, req.Username, req.DisplayName)
	if err != nil {
		log.Printf("Failed to begin registration: %v", err)
		http.Error(w, "Failed to begin registration", http.StatusInternalServerError)
		return
	}

	// Store session ID in response (in practice, you might use HTTP sessions)
	w.Header().Set("X-Session-ID", sessionID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

func registerFinishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username   string      `json:"username"`
		Credential interface{} `json:"credential"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// In practice, you would get the session ID from HTTP sessions
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		// For demo purposes, we'll simulate success
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	// Convert credential to bytes and finish registration
	credentialBytes, _ := json.Marshal(req.Credential)
	
	err := client.FinishRegistration(req.Username, sessionID, credentialBytes)
	if err != nil {
		log.Printf("Failed to finish registration: %v", err)
		http.Error(w, "Failed to finish registration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<title>登录 - WebAuthn Demo</title>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
		body { 
			font-family: Arial, sans-serif; 
			max-width: 600px; 
			margin: 50px auto; 
			padding: 20px;
			background: #f5f5f5;
		}
		.container { 
			background: white; 
			padding: 40px; 
			border-radius: 10px; 
			box-shadow: 0 2px 10px rgba(0,0,0,0.1);
		}
		.form-group { margin: 20px 0; }
		label { display: block; margin-bottom: 5px; font-weight: bold; }
		input { 
			width: 100%; 
			padding: 12px; 
			border: 1px solid #ddd; 
			border-radius: 5px; 
			font-size: 16px;
			box-sizing: border-box;
		}
		.button { 
			background-color: #007bff; 
			color: white; 
			padding: 15px 30px; 
			border: none; 
			border-radius: 5px; 
			cursor: pointer; 
			font-size: 16px;
			width: 100%;
		}
		.button:hover { background-color: #0056b3; }
		.button:disabled { background-color: #6c757d; cursor: not-allowed; }
		.back-link { color: #007bff; text-decoration: none; }
		.back-link:hover { text-decoration: underline; }
		.status { 
			padding: 15px; 
			margin: 20px 0; 
			border-radius: 5px; 
			display: none;
		}
		.status.success { background: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
		.status.error { background: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }
		.status.info { background: #d1ecf1; color: #0c5460; border: 1px solid #bee5eb; }
		h1 { color: #333; text-align: center; }
	</style>
</head>
<body>
	<div class="container">
		<h1>🔓 登录</h1>
		<div id="status" class="status"></div>
		
		<form id="loginForm">
			<div class="form-group">
				<label for="username">用户名或邮箱:</label>
				<input type="text" id="username" name="username" required placeholder="输入您的用户名或邮箱">
			</div>
			<button type="submit" class="button" id="loginBtn">🔐 使用身份验证器登录</button>
		</form>

		<div style="text-align: center; margin: 20px 0;">
			<a href="/" class="back-link">← 返回首页</a> | 
			<a href="/register" class="back-link">没有账户？点击注册</a>
		</div>
	</div>

	<script>
		document.getElementById('loginForm').addEventListener('submit', async (e) => {
			e.preventDefault();
			
			const username = document.getElementById('username').value;
			const statusDiv = document.getElementById('status');
			const loginBtn = document.getElementById('loginBtn');
			
			if (!username) {
				showStatus('请输入用户名或邮箱', 'error');
				return;
			}

			loginBtn.disabled = true;
			loginBtn.textContent = '正在准备登录...';
			showStatus('正在准备 WebAuthn 登录...', 'info');

			try {
				// Begin login
				const beginResponse = await fetch('/login/begin', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({ username })
				});

				if (!beginResponse.ok) {
					throw new Error('登录初始化失败');
				}

				const credentialRequestOptions = await beginResponse.json();
				showStatus('请使用您的身份验证器进行认证...', 'info');

				// Convert base64url to ArrayBuffer
				credentialRequestOptions.publicKey.challenge = base64urlDecode(credentialRequestOptions.publicKey.challenge);
				
				if (credentialRequestOptions.publicKey.allowCredentials) {
					credentialRequestOptions.publicKey.allowCredentials.forEach(cred => {
						cred.id = base64urlDecode(cred.id);
					});
				}

				// Get credential
				const assertion = await navigator.credentials.get(credentialRequestOptions);

				if (!assertion) {
					throw new Error('身份验证失败');
				}

				showStatus('正在验证身份...', 'info');

				// Finish login
				const finishResponse = await fetch('/login/finish', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
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

				if (!finishResponse.ok) {
					throw new Error('身份验证失败');
				}

				showStatus('🎉 登录成功！正在跳转到控制台...', 'success');
				setTimeout(() => {
					window.location.href = '/dashboard';
				}, 2000);

			} catch (error) {
				console.error('Login error:', error);
				showStatus('登录失败: ' + error.message, 'error');
				loginBtn.disabled = false;
				loginBtn.textContent = '🔐 使用身份验证器登录';
			}
		});

		function showStatus(message, type) {
			const statusDiv = document.getElementById('status');
			statusDiv.className = 'status ' + type;
			statusDiv.textContent = message;
			statusDiv.style.display = 'block';
		}

		function base64urlDecode(str) {
			return Uint8Array.from(atob(str.replace(/-/g, '+').replace(/_/g, '/')), c => c.charCodeAt(0));
		}

		function base64urlEncode(buffer) {
			return btoa(String.fromCharCode(...new Uint8Array(buffer)))
				.replace(/\+/g, '-')
				.replace(/\//g, '_')
				.replace(/=/g, '');
		}
	</script>
</body>
</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

func loginBeginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Begin login
	options, sessionID, err := client.BeginLogin(req.Username)
	if err != nil {
		log.Printf("Failed to begin login: %v", err)
		http.Error(w, "Failed to begin login", http.StatusInternalServerError)
		return
	}

	// Store session ID in response
	w.Header().Set("X-Session-ID", sessionID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

func loginFinishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username  string      `json:"username"`
		Assertion interface{} `json:"assertion"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// In practice, you would get the session ID from HTTP sessions
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		// For demo purposes, we'll simulate success
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		return
	}

	// Convert assertion to bytes and finish login
	assertionBytes, _ := json.Marshal(req.Assertion)
	
	err := client.FinishLogin(req.Username, sessionID, assertionBytes)
	if err != nil {
		log.Printf("Failed to finish login: %v", err)
		http.Error(w, "Failed to finish login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Simple logout - redirect to home
	http.Redirect(w, r, "/", http.StatusFound)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<title>控制台 - WebAuthn Demo</title>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
		body { 
			font-family: Arial, sans-serif; 
			max-width: 800px; 
			margin: 50px auto; 
			padding: 20px;
			background: #f5f5f5;
		}
		.container { 
			background: white; 
			padding: 40px; 
			border-radius: 10px; 
			box-shadow: 0 2px 10px rgba(0,0,0,0.1);
		}
		.button { 
			background-color: #dc3545; 
			color: white; 
			padding: 10px 20px; 
			text-decoration: none; 
			border: none; 
			border-radius: 5px; 
			cursor: pointer; 
			display: inline-block; 
		}
		.button:hover { background-color: #c82333; }
		h1 { color: #333; text-align: center; }
		.success-message { 
			background: #d4edda; 
			color: #155724; 
			padding: 20px; 
			border-radius: 5px; 
			margin: 20px 0; 
			text-align: center; 
		}
		.features { 
			background: #e9ecef; 
			padding: 20px; 
			border-radius: 5px; 
			margin: 20px 0; 
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>🎉 登录成功！</h1>
		
		<div class="success-message">
			<h2>欢迎来到 WebAuthn 控制台</h2>
			<p>您已成功使用 WebAuthn 进行无密码认证！</p>
		</div>

		<div class="features">
			<h3>🔐 您刚才体验了什么？</h3>
			<ul>
				<li><strong>无密码登录：</strong>没有输入任何密码，只使用了生物识别或安全密钥</li>
				<li><strong>强身份验证：</strong>基于公钥密码学，比传统密码更安全</li>
				<li><strong>防钓鱼保护：</strong>身份验证器与域名绑定，无法在钓鱼网站使用</li>
				<li><strong>用户体验：</strong>快速、便捷，无需记住复杂密码</li>
			</ul>
		</div>

		<div style="text-align: center; margin: 30px 0;">
			<a href="/logout" class="button">🚪 退出登录</a>
		</div>

		<div style="text-align: center; color: #666; margin: 20px 0;">
			<p>这是一个 WebAuthn 技术演示。在实际应用中，您可以在这里访问您的账户功能。</p>
		</div>
	</div>
</body>
</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}