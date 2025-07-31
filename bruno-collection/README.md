# englog API - Bruno Collection

Esta é uma collection completa do Bruno para testar a API englog, baseada na documentação Swagger encontrada na pasta `api/`.

## Estrutura da Collection

A collection está organizada nas seguintes pastas:

### 🔐 Auth
- **Register User**: Registra um novo usuário e retorna tokens de autenticação
- **Login User**: Faz login e obtém tokens de acesso
- **Get Me**: Obtém informações do perfil do usuário autenticado
- **Refresh Token**: Renova os tokens de acesso usando o refresh token
- **Logout User**: Invalida o refresh token e faz logout

### 🏥 Health
- **Health Check**: Verifica se a API está funcionando
- **Readiness Check**: Verifica se a API está pronta para receber requisições

### 📝 Log Entries
- **Create Log Entry**: Cria uma nova entrada de log
- **Get Log Entries**: Lista todas as entradas de log com filtros opcionais
- **Get Log Entry**: Obtém uma entrada específica por ID
- **Update Log Entry**: Atualiza uma entrada existente
- **Delete Log Entry**: Remove uma entrada
- **Bulk Create Log Entries**: Cria múltiplas entradas de uma vez

### 📁 Projects
- **Create Project**: Cria um novo projeto
- **Get Projects**: Lista todos os projetos
- **Get Project**: Obtém um projeto específico por ID
- **Update Project**: Atualiza um projeto existente
- **Delete Project**: Remove um projeto

### 📊 Analytics
- **Get Productivity Metrics**: Obtém métricas de produtividade
- **Get Activity Summary**: Obtém resumo de atividades

### 🏷️ Tags
- **Create Tag**: Cria uma nova tag
- **Get Tags**: Lista todas as tags
- **Get Popular Tags**: Obtém as tags mais populares
- **Get Recently Used Tags**: Obtém tags recentemente utilizadas
- **Search Tags**: Busca tags por nome
- **Get User Tag Usage**: Obtém estatísticas de uso de tags
- **Get Tag**: Obtém uma tag específica por ID
- **Update Tag**: Atualiza uma tag existente
- **Delete Tag**: Remove uma tag

### 👤 Users
- **Get Profile**: Obtém o perfil do usuário
- **Update Profile**: Atualiza informações do perfil
- **Change Password**: Altera a senha do usuário
- **Delete Account**: Remove a conta do usuário

## Como Usar

### 1. Configuração do Ambiente

O arquivo `environments/Local.bru` contém as variáveis de ambiente:

```
vars {
  base_url: http://localhost:8080
  access_token:
  refresh_token:
}
```

Ajuste a `base_url` conforme necessário para o seu ambiente.

### 2. Fluxo de Autenticação

1. **Registre um usuário** usando "Register User" ou **faça login** com "Login User"
2. Os tokens serão automaticamente salvos nas variáveis de ambiente através dos scripts pós-resposta
3. As demais requisições autenticadas usarão automaticamente o `access_token`

### 3. Variáveis Dinâmicas

Algumas requisições usam variáveis para IDs:
- `{{log_entry_id}}`: ID de uma entrada de log
- `{{project_id}}`: ID de um projeto
- `{{tag_id}}`: ID de uma tag

Você pode definir essas variáveis:
- Manualmente no ambiente Local.bru
- Capturando das respostas de criação através de scripts
- Copiando IDs das listagens

### 4. Scripts Automáticos

As requisições de autenticação incluem scripts que automaticamente:
- Salvam tokens nas variáveis de ambiente após login/registro
- Limpam tokens após logout

### 5. Testes

Cada requisição inclui testes básicos que verificam:
- Status codes esperados
- Presença de campos obrigatórios nas respostas
- Estrutura básica dos dados retornados

## Exemplos de Dados

### Registro de Usuário
```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "first_name": "John",
  "last_name": "Doe",
  "timezone": "UTC"
}
```

### Entrada de Log
```json
{
  "title": "Fixed authentication bug",
  "description": "Resolved issue with JWT token validation",
  "category": "bug_fix",
  "tags": ["authentication", "jwt", "security"],
  "time_spent_minutes": 120,
  "difficulty": 3,
  "mood": "satisfied"
}
```

### Projeto
```json
{
  "name": "englog API",
  "description": "Engineering log management API",
  "repository_url": "https://github.com/garnizeh/englog",
  "status": "active",
  "tags": ["go", "api", "postgresql"]
}
```

## Dicas

1. **Ordem das Requisições**: Comece sempre com autenticação
2. **IDs de Teste**: Use os IDs de exemplo fornecidos ou capture-os das respostas
3. **Filtros**: Muitas listagens suportam filtros via query parameters
4. **Paginação**: Use `limit` e `offset` para paginar resultados
5. **Tokens**: Os tokens são automaticamente gerenciados pelos scripts

## Troubleshooting

- **401 Unauthorized**: Verifique se o token está válido ou renove-o
- **404 Not Found**: Verifique se os IDs usados existem
- **400 Bad Request**: Verifique a estrutura do JSON enviado
- **Conexão recusada**: Verifique se a API está rodando na URL configurada

---

Esta collection cobre todos os endpoints documentados na API englog e está pronta para uso em desenvolvimento e testes.
