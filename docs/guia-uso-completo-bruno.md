# Guia Completo: Fluxo de Uso da API EngLog com Bruno

> "The best way to learn is by doing." - Richard Branson 🚀

Este guia detalha o fluxo completo para incluir entradas no log usando a coleção Bruno, desde o registro/login até o logout, com todas as referências específicas dos requests.

## 📋 Pré-requisitos

1. **Ambiente em Execução**:
   ```bash
   # Iniciar ambiente de desenvolvimento
   make dev-up

   # Verificar se a API está funcionando
   make dev-api
   ```

2. **Bruno Client**: Tenha o Bruno instalado e configurado
3. **Coleção Importada**: Importe a coleção da pasta `bruno-collection/`

## 🎯 Fluxo Completo: Passo a Passo

### **Passo 1: Verificar Saúde da API** 🏥

**Request**: `Health/Health Check.bru`
- **Método**: `GET`
- **URL**: `{{base_url}}/health`
- **Autenticação**: Nenhuma
- **Objetivo**: Verificar se a API está funcionando

**Resposta Esperada**:
```json
{
  "status": "ok",
  "timestamp": "2025-07-31T10:30:00Z"
}
```

---

### **Passo 2: Registrar Novo Usuário** 👤

**Request**: `Auth/Register User.bru`
- **Método**: `POST`
- **URL**: `{{base_url}}/v1/auth/register`
- **Autenticação**: Nenhuma

**Body de Exemplo**:
```json
{
  "email": "engineer@example.com",
  "password": "SecurePass123!",
  "first_name": "Maria",
  "last_name": "Silva",
  "timezone": "America/Sao_Paulo"
}
```

**Resposta Esperada**:
```json
{
  "user": {
    "id": "uuid-do-usuario",
    "email": "engineer@example.com",
    "first_name": "Maria",
    "last_name": "Silva",
    "timezone": "America/Sao_Paulo"
  },
  "tokens": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ..."
  }
}
```

**Script Automático**: O Bruno automaticamente salva os tokens nas variáveis de ambiente.

---

### **Passo 3: Login (Alternativa ao Registro)** 🔐

**Request**: `Auth/Login User.bru`
- **Método**: `POST`
- **URL**: `{{base_url}}/v1/auth/login`
- **Autenticação**: Nenhuma

**Body de Exemplo**:
```json
{
  "email": "engineer@example.com",
  "password": "SecurePass123!"
}
```

**Resposta Esperada**: Similar ao registro, mas sem dados de criação.

---

### **Passo 4: Verificar Perfil do Usuário** ✅

**Request**: `Auth/Get Me.bru`
- **Método**: `GET`
- **URL**: `{{base_url}}/v1/auth/me`
- **Autenticação**: Bearer Token (`{{access_token}}`)

**Resposta Esperada**:
```json
{
  "id": "uuid-do-usuario",
  "email": "engineer@example.com",
  "first_name": "Maria",
  "last_name": "Silva",
  "timezone": "America/Sao_Paulo",
  "created_at": "2025-07-31T08:00:00Z",
  "last_login": "2025-07-31T10:30:00Z"
}
```

---

### **Passo 5: Criar Projeto (Opcional)** 📁

**Request**: `Projects/Create Project.bru`
- **Método**: `POST`
- **URL**: `{{base_url}}/v1/projects`
- **Autenticação**: Bearer Token (`{{access_token}}`)

**Body de Exemplo**:
```json
{
  "name": "Sistema de Autenticação",
  "description": "Desenvolvimento do módulo de autenticação JWT",
  "status": "active",
  "color": "#3B82F6"
}
```

**Resposta Esperada**:
```json
{
  "id": "project-uuid",
  "name": "Sistema de Autenticação",
  "description": "Desenvolvimento do módulo de autenticação JWT",
  "status": "active",
  "color": "#3B82F6",
  "user_id": "user-uuid",
  "created_at": "2025-07-31T10:35:00Z"
}
```

⚠️ **Importante**: Anote o `project_id` para usar nas entradas de log.

---

### **Passo 6: Criar Tags (Opcional)** 🏷️

**Request**: `Tags/Create Tag.bru`
- **Método**: `POST`
- **URL**: `{{base_url}}/v1/tags`
- **Autenticação**: Bearer Token (`{{access_token}}`)

**Body de Exemplo**:
```json
{
  "name": "backend",
  "color": "#10B981"
}
```

**Resposta Esperada**:
```json
{
  "id": "tag-uuid",
  "name": "backend",
  "color": "#10B981",
  "user_id": "user-uuid",
  "usage_count": 0
}
```

---

### **Passo 7: Criar Primeira Entrada de Log** 📝

**Request**: `Log Entries/Create Log Entry.bru`
- **Método**: `POST`
- **URL**: `{{base_url}}/v1/logs`
- **Autenticação**: Bearer Token (`{{access_token}}`)

**Body Atualizado para EngLog**:
```json
{
  "title": "Implementação de autenticação JWT",
  "description": "Desenvolvimento completo do sistema de autenticação usando JWT com refresh tokens e middleware de proteção",
  "type": "development",
  "value_rating": "high",
  "impact_level": "team",
  "project_id": "project-uuid-aqui",
  "start_time": "2025-07-31T09:00:00Z",
  "end_time": "2025-07-31T11:30:00Z",
  "tags": ["jwt", "authentication", "security", "backend"]
}
```

**Campos Obrigatórios**:
- `title`: Título da atividade
- `type`: Tipo da atividade (development, meeting, code_review, debugging, etc.)
- `value_rating`: Valor percebido (low, medium, high, critical)
- `impact_level`: Nível de impacto (personal, team, department, company)
- `start_time`: Horário de início (ISO 8601)
- `end_time`: Horário de fim (ISO 8601)

**Campos Opcionais**:
- `description`: Descrição detalhada
- `project_id`: ID do projeto (UUID)
- `tags`: Array de strings com tags

**Resposta Esperada**:
```json
{
  "id": "log-entry-uuid",
  "title": "Implementação de autenticação JWT",
  "description": "Desenvolvimento completo...",
  "type": "development",
  "value_rating": "high",
  "impact_level": "team",
  "duration_minutes": 150,
  "start_time": "2025-07-31T09:00:00Z",
  "end_time": "2025-07-31T11:30:00Z",
  "project_id": "project-uuid",
  "user_id": "user-uuid",
  "created_at": "2025-07-31T11:30:00Z",
  "updated_at": "2025-07-31T11:30:00Z"
}
```

---

### **Passo 8: Criar Múltiplas Entradas** 📦

**Request**: `Log Entries/Bulk Create Log Entries.bru`
- **Método**: `POST`
- **URL**: `{{base_url}}/v1/logs/bulk`
- **Autenticação**: Bearer Token (`{{access_token}}`)

**Body de Exemplo**:
```json
{
  "entries": [
    {
      "title": "Code review do PR #123",
      "description": "Revisão detalhada do pull request com melhorias no sistema de logs",
      "type": "code_review",
      "value_rating": "medium",
      "impact_level": "team",
      "start_time": "2025-07-31T14:00:00Z",
      "end_time": "2025-07-31T14:45:00Z",
      "project_id": "project-uuid-aqui",
      "tags": ["code-review", "logs", "quality"]
    },
    {
      "title": "Reunião de planning",
      "description": "Planejamento das funcionalidades da sprint 12",
      "type": "meeting",
      "value_rating": "high",
      "impact_level": "team",
      "start_time": "2025-07-31T15:00:00Z",
      "end_time": "2025-07-31T16:00:00Z",
      "tags": ["planning", "sprint", "team"]
    },
    {
      "title": "Debug do problema de performance",
      "description": "Investigação e correção de queries lentas no endpoint de analytics",
      "type": "debugging",
      "value_rating": "critical",
      "impact_level": "company",
      "start_time": "2025-07-31T16:30:00Z",
      "end_time": "2025-07-31T18:00:00Z",
      "project_id": "project-uuid-aqui",
      "tags": ["performance", "database", "optimization"]
    }
  ]
}
```

---

### **Passo 9: Listar Entradas Criadas** 📋

**Request**: `Log Entries/Get Log Entries.bru`
- **Método**: `GET`
- **URL**: `{{base_url}}/v1/logs?limit=10&offset=0`
- **Autenticação**: Bearer Token (`{{access_token}}`)

**Parâmetros de Query Disponíveis**:
- `limit`: Limite de resultados (padrão: 10)
- `offset`: Offset para paginação (padrão: 0)
- `type`: Filtrar por tipo de atividade
- `project_id`: Filtrar por projeto
- `value_rating`: Filtrar por valor (low, medium, high, critical)
- `impact_level`: Filtrar por impacto (personal, team, department, company)
- `start_date`: Data de início (YYYY-MM-DD)
- `end_date`: Data de fim (YYYY-MM-DD)
- `search`: Busca em título e descrição

**Exemplo com Filtros**:
```
{{base_url}}/v1/logs?type=development&value_rating=high&limit=5&start_date=2025-07-31
```

---

### **Passo 10: Obter Entrada Específica** 🔍

**Request**: `Log Entries/Get Log Entry.bru`
- **Método**: `GET`
- **URL**: `{{base_url}}/v1/logs/{id}`
- **Autenticação**: Bearer Token (`{{access_token}}`)

Substitua `{id}` pelo UUID da entrada que deseja consultar.

---

### **Passo 11: Atualizar Entrada** ✏️

**Request**: `Log Entries/Update Log Entry.bru`
- **Método**: `PUT`
- **URL**: `{{base_url}}/v1/logs/{id}`
- **Autenticação**: Bearer Token (`{{access_token}}`)

**Body de Exemplo**:
```json
{
  "title": "Implementação de autenticação JWT - CONCLUÍDA",
  "description": "Desenvolvimento completo do sistema de autenticação usando JWT com refresh tokens, middleware de proteção e testes unitários",
  "type": "development",
  "value_rating": "critical",
  "impact_level": "company",
  "project_id": "project-uuid-aqui",
  "start_time": "2025-07-31T09:00:00Z",
  "end_time": "2025-07-31T12:00:00Z",
  "tags": ["jwt", "authentication", "security", "backend", "completed"]
}
```

---

### **Passo 12: Ver Analytics** 📊

**Request**: `Analytics/Get Productivity Metrics.bru`
- **Método**: `GET`
- **URL**: `{{base_url}}/v1/analytics/productivity`
- **Autenticação**: Bearer Token (`{{access_token}}`)

**Parâmetros Opcionais**:
- `start_date`: Data inicial (YYYY-MM-DD)
- `end_date`: Data final (YYYY-MM-DD)

**Resposta Esperada**:
```json
{
  "total_entries": 4,
  "total_minutes": 360,
  "average_duration": 90.0,
  "entries_by_type": {
    "development": 2,
    "code_review": 1,
    "meeting": 1,
    "debugging": 1
  },
  "entries_by_value": {
    "critical": 2,
    "high": 1,
    "medium": 1,
    "low": 0
  },
  "entries_by_impact": {
    "company": 2,
    "team": 2,
    "department": 0,
    "personal": 0
  },
  "period": {
    "start_date": "2025-07-31",
    "end_date": "2025-07-31"
  }
}
```

---

### **Passo 13: Ver Resumo de Atividades** 📈

**Request**: `Analytics/Get Activity Summary.bru`
- **Método**: `GET`
- **URL**: `{{base_url}}/v1/analytics/summary`
- **Autenticação**: Bearer Token (`{{access_token}}`)

---

### **Passo 14: Excluir Entrada (Opcional)** 🗑️

**Request**: `Log Entries/Delete Log Entry.bru`
- **Método**: `DELETE`
- **URL**: `{{base_url}}/v1/logs/{id}`
- **Autenticação**: Bearer Token (`{{access_token}}`)

---

### **Passo 15: Renovar Token (Se Necessário)** 🔄

**Request**: `Auth/Refresh Token.bru`
- **Método**: `POST`
- **URL**: `{{base_url}}/v1/auth/refresh`
- **Autenticação**: Nenhuma

**Body**:
```json
{
  "refresh_token": "{{refresh_token}}"
}
```

---

### **Passo 16: Logout** 👋

**Request**: `Auth/Logout User.bru`
- **Método**: `POST`
- **URL**: `{{base_url}}/v1/auth/logout`
- **Autenticação**: Nenhuma

**Body**:
```json
{
  "refresh_token": "{{refresh_token}}"
}
```

**Script Automático**: O Bruno limpa automaticamente os tokens das variáveis de ambiente.

---

## 🎯 Fluxos de Trabalho Recomendados

### **Fluxo Diário Básico**

1. **Login** → `Auth/Login User.bru`
2. **Criar entradas ao longo do dia** → `Log Entries/Create Log Entry.bru`
3. **Ver progresso** → `Analytics/Get Productivity Metrics.bru`
4. **Logout** → `Auth/Logout User.bru`

### **Fluxo de Desenvolvimento Intensivo**

1. **Login** → `Auth/Login User.bru`
2. **Criar projeto** → `Projects/Create Project.bru`
3. **Criar tags relevantes** → `Tags/Create Tag.bru`
4. **Criar múltiplas entradas** → `Log Entries/Bulk Create Log Entries.bru`
5. **Monitorar analytics** → `Analytics/Get Activity Summary.bru`
6. **Atualizar entradas** → `Log Entries/Update Log Entry.bru`
7. **Logout** → `Auth/Logout User.bru`

### **Fluxo de Revisão/Auditoria**

1. **Login** → `Auth/Login User.bru`
2. **Listar entradas com filtros** → `Log Entries/Get Log Entries.bru`
3. **Ver detalhes específicos** → `Log Entries/Get Log Entry.bru`
4. **Analisar métricas** → `Analytics/Get Productivity Metrics.bru`
5. **Logout** → `Auth/Logout User.bru`

---

## 🔧 Configurações do Bruno

### **Variáveis de Ambiente**

No arquivo `environments/Local.bru`:
```javascript
vars {
  base_url: http://localhost:8080
  access_token:
  refresh_token:
}
```

### **Scripts Automáticos**

Os requests de autenticação incluem scripts que:
- Salvam tokens automaticamente após login/registro
- Limpam tokens após logout
- Facilitam o fluxo contínuo de testes

---

## 🐛 Troubleshooting

### **Token Expirado**
- Execute `Auth/Refresh Token.bru` para renovar
- Se falhar, faça novo login com `Auth/Login User.bru`

### **Erro 401 Unauthorized**
- Verifique se `{{access_token}}` está preenchido
- Confirme se fez login recentemente

### **Erro 400 Bad Request**
- Verifique o formato dos dados no body
- Confirme que campos obrigatórios estão presentes
- Valide formato das datas (ISO 8601)

### **Erro de Conexão**
- Confirme que a API está rodando: `make dev-api`
- Verifique se o `base_url` está correto no ambiente

---

## 📚 Referências Rápidas

### **Tipos de Atividade Válidos**
- `development` - Desenvolvimento de código
- `meeting` - Reuniões e calls
- `code_review` - Revisões de código
- `debugging` - Correção de bugs
- `documentation` - Documentação
- `testing` - Testes
- `deployment` - Deploy e DevOps
- `research` - Pesquisa e estudos
- `planning` - Planejamento
- `learning` - Aprendizado
- `maintenance` - Manutenção
- `support` - Suporte
- `other` - Outros

### **Níveis de Valor**
- `low` - Atividades rotineiras
- `medium` - Trabalho padrão
- `high` - Trabalho importante
- `critical` - Atividades críticas

### **Níveis de Impacto**
- `personal` - Desenvolvimento individual
- `team` - Impacto no time
- `department` - Impacto departamental
- `company` - Impacto organizacional

---

## ✅ Checklist de Verificação

- [ ] API está rodando (`make dev-api`)
- [ ] Coleção Bruno importada
- [ ] Ambiente Local configurado
- [ ] Health check passou
- [ ] Usuário registrado/logado
- [ ] Tokens salvos automaticamente
- [ ] Projeto criado (se necessário)
- [ ] Tags criadas (se necessário)
- [ ] Entradas de log criadas
- [ ] Analytics verificadas
- [ ] Logout executado

Este guia fornece um fluxo completo e prático para usar a API EngLog com todas as referências específicas da coleção Bruno!
