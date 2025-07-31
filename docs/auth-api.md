# Authentication API Documentation

> "Security is not a product, but a process." 🔐

Este documento descreve como usar a API de autenticação JWT implementada no sistema EngLog.

## Visão Geral

A API de autenticação fornece um sistema completo de gerenciamento de usuários com:
- Estratégia de tokens duplos (access + refresh)
- Rotação automática de refresh tokens
- Sistema de denylist para tokens revogados
- Middleware de proteção para rotas
- Hash seguro de senhas com bcrypt

## Endpoints Disponíveis

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints Públicos

#### 1. Registro de Usuário
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "timezone": "America/Sao_Paulo"  // opcional, padrão: UTC
}
```

**Resposta de Sucesso (201):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "timezone": "America/Sao_Paulo",
    "preferences": null,
    "created_at": "2025-07-25T23:30:00Z"
  },
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

#### 2. Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Resposta de Sucesso (200):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "timezone": "America/Sao_Paulo",
    "preferences": null,
    "created_at": "2025-07-25T23:30:00Z"
  },
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

#### 3. Renovação de Token
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Resposta de Sucesso (200):**
```json
{
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

#### 4. Logout
```http
POST /auth/logout
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Resposta de Sucesso (200):**
```json
{
  "message": "Successfully logged out"
}
```

### Endpoints Protegidos

#### 5. Obter Perfil do Usuário
```http
GET /me
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Resposta de Sucesso (200):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "timezone": "America/Sao_Paulo",
    "preferences": null,
    "created_at": "2025-07-25T23:30:00Z"
  }
}
```

## Autenticação

### Header de Autorização
Todas as rotas protegidas requerem o header de autorização:
```
Authorization: Bearer <access_token>
```

### Estrutura dos Tokens

#### Access Token
- **Duração**: 15 minutos
- **Uso**: Autenticação de requests para rotas protegidas
- **Tipo**: JWT com claims customizados

#### Refresh Token
- **Duração**: 30 dias
- **Uso**: Renovação de access tokens
- **Segurança**: Inclui JTI (JWT ID) para controle de denylist

## Códigos de Erro

### Erros de Autenticação (401)
```json
{
  "error": "Invalid token",
  "details": "token is expired"
}
```

### Erros de Validação (400)
```json
{
  "error": "Invalid request format",
  "details": "email: required field is missing"
}
```

### Conflitos (409)
```json
{
  "error": "User already exists"
}
```

### Erros do Servidor (500)
```json
{
  "error": "Failed to create access token"
}
```

## Fluxo Recomendado

### 1. Autenticação Inicial
1. Registrar ou fazer login
2. Armazenar ambos os tokens de forma segura
3. Usar access token para requests autenticados

### 2. Renovação de Tokens
1. Quando access token expirar (401), usar refresh token
2. Armazenar novos tokens
3. O refresh token antigo é automaticamente invalidado

### 3. Logout
1. Enviar refresh token para endpoint de logout
2. Token é adicionado à denylist
3. Remover tokens do armazenamento local

## Exemplo de Implementação (JavaScript)

```javascript
class AuthService {
  constructor(baseURL) {
    this.baseURL = baseURL;
    this.accessToken = localStorage.getItem('access_token');
    this.refreshToken = localStorage.getItem('refresh_token');
  }

  async login(email, password) {
    const response = await fetch(`${this.baseURL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });

    if (response.ok) {
      const data = await response.json();
      this.setTokens(data.tokens);
      return data;
    }
    throw new Error('Login failed');
  }

  async refreshTokens() {
    if (!this.refreshToken) throw new Error('No refresh token');

    const response = await fetch(`${this.baseURL}/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: this.refreshToken })
    });

    if (response.ok) {
      const data = await response.json();
      this.setTokens(data.tokens);
      return data;
    }
    throw new Error('Token refresh failed');
  }

  async authenticatedRequest(url, options = {}) {
    let response = await fetch(url, {
      ...options,
      headers: {
        ...options.headers,
        'Authorization': `Bearer ${this.accessToken}`
      }
    });

    if (response.status === 401) {
      await this.refreshTokens();
      response = await fetch(url, {
        ...options,
        headers: {
          ...options.headers,
          'Authorization': `Bearer ${this.accessToken}`
        }
      });
    }

    return response;
  }

  setTokens(tokens) {
    this.accessToken = tokens.access_token;
    this.refreshToken = tokens.refresh_token;
    localStorage.setItem('access_token', this.accessToken);
    localStorage.setItem('refresh_token', this.refreshToken);
  }

  async logout() {
    if (this.refreshToken) {
      await fetch(`${this.baseURL}/auth/logout`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: this.refreshToken })
      });
    }

    this.clearTokens();
  }

  clearTokens() {
    this.accessToken = null;
    this.refreshToken = null;
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
  }
}
```

## Variáveis de Ambiente

Configure as seguintes variáveis para personalizar o comportamento:

```bash
# Chave secreta para assinatura JWT (OBRIGATÓRIO em produção)
JWT_SECRET_KEY=your-very-secure-secret-key-here

# Tempo de vida dos tokens (opcional)
JWT_ACCESS_TOKEN_TTL=15m
JWT_REFRESH_TOKEN_TTL=720h  # 30 dias

# Intervalo de limpeza de tokens expirados (opcional)
JWT_DENYLIST_CLEANUP_INTERVAL=24h
```

## Considerações de Segurança

### Produção
- ✅ Use uma chave secreta forte (mínimo 32 caracteres)
- ✅ Configure HTTPS obrigatório
- ✅ Implemente rate limiting
- ✅ Monitor tentativas de login
- ✅ Log eventos de segurança

### Desenvolvimento
- ⚠️ Nunca commite chaves secretas
- ⚠️ Use variáveis de ambiente
- ⚠️ Teste fluxos de expiração
- ⚠️ Valide entrada de dados

## Monitoramento

O sistema registra automaticamente:
- Tentativas de login bem-sucedidas/falhadas
- Operações de refresh token
- Limpeza de tokens expirados
- Erros de autenticação

## Troubleshooting

### Token Inválido
- Verifique se o token não expirou
- Confirme que está usando o header correto
- Verifique se o token não foi revogado

### Refresh Token Falhou
- Token pode ter expirado (30 dias)
- Token pode ter sido usado em rotação anterior
- Token pode ter sido revogado por logout

### Performance
- Monitor frequência de refresh
- Considere ajustar TTL dos tokens
- Verifique logs de cleanup
