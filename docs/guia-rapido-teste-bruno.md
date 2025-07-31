# Guia Rápido de Teste - Bruno Collection

> "The best debugging tool is still a print statement." - Paul Kocher 🛠️

Este é um guia de referência rápida para testar o fluxo completo da API usando a coleção Bruno.

## 🚀 Teste Rápido (5 minutos)

### 1. Health Check
**Endpoint**: `Health/Health Check.bru`
- Verificar se API está rodando
- Não requer autenticação

### 2. Registro/Login
**Endpoint**: `Auth/Register User.bru` OU `Auth/Login User.bru`
- Tokens salvos automaticamente
- Use email único para registro

### 3. Criar Entrada de Log
**Endpoint**: `Log Entries/Create Log Entry.bru`
- Modifique o exemplo conforme necessário
- Campos obrigatórios: title, type, value_rating, impact_level, start_time, end_time

### 4. Listar Entradas
**Endpoint**: `Log Entries/Get Log Entries.bru`
- Verificar se a entrada foi criada
- Testar filtros opcionais

### 5. Analytics
**Endpoint**: `Analytics/Get Productivity Metrics.bru`
- Ver métricas calculadas
- Confirmar dados estão sendo processados

## 📝 Exemplos Prontos para Usar

### Entrada de Log - Desenvolvimento
```json
{
  "title": "Refatoração do módulo de autenticação",
  "description": "Melhoria na arquitetura e adição de testes unitários",
  "type": "development",
  "value_rating": "high",
  "impact_level": "team",
  "start_time": "2025-07-31T14:00:00Z",
  "end_time": "2025-07-31T16:30:00Z",
  "tags": ["refactoring", "testing", "architecture"]
}
```

### Entrada de Log - Reunião
```json
{
  "title": "Daily standup meeting",
  "description": "Sincronização do time sobre progresso dos projetos",
  "type": "meeting",
  "value_rating": "medium",
  "impact_level": "team",
  "start_time": "2025-07-31T09:00:00Z",
  "end_time": "2025-07-31T09:30:00Z",
  "tags": ["standup", "sync", "team"]
}
```

### Entrada de Log - Debug Crítico
```json
{
  "title": "Correção de bug crítico em produção",
  "description": "Investigação e correção de falha de segurança no endpoint de login",
  "type": "debugging",
  "value_rating": "critical",
  "impact_level": "company",
  "start_time": "2025-07-31T20:00:00Z",
  "end_time": "2025-07-31T22:30:00Z",
  "tags": ["hotfix", "security", "production"]
}
```

## 🔍 Filtros Úteis para Teste

### Listar por Tipo
```
GET /v1/logs?type=development&limit=5
```

### Listar por Valor
```
GET /v1/logs?value_rating=high&impact_level=team
```

### Listar por Data
```
GET /v1/logs?start_date=2025-07-31&end_date=2025-07-31
```

### Buscar por Texto
```
GET /v1/logs?search=authentication&limit=10
```

## 🎯 Cenários de Teste Recomendados

### Cenário 1: Dia de Desenvolvimento
1. Criar entrada "development" (2h)
2. Criar entrada "code_review" (30min)
3. Criar entrada "debugging" (1h)
4. Ver analytics do dia

### Cenário 2: Sprint Planning
1. Criar entrada "planning" (1h)
2. Criar entrada "meeting" (30min)
3. Criar projeto relacionado
4. Ver métricas de planejamento

### Cenário 3: Hotfix Production
1. Criar entrada "debugging" com value_rating="critical"
2. Criar entrada "deployment" relacionada
3. Ver impacto nas métricas

## ⚡ Atalhos de Teclado Bruno

- `Ctrl+Enter`: Executar request
- `Ctrl+Shift+Enter`: Executar todos os requests da pasta
- `Ctrl+K`: Buscar requests
- `Ctrl+N`: Novo request

## 🛠️ Troubleshooting Rápido

| Erro | Solução |
|------|---------|
| 401 Unauthorized | Execute Auth/Login User.bru novamente |
| 400 Bad Request | Verifique formato dos campos obrigatórios |
| 404 Not Found | Confirme se o ID existe (use Get List primeiro) |
| 500 Server Error | Verifique logs da API ou faça Health Check |

## 📊 Validação dos Dados

Após criar algumas entradas, verifique:

1. **Duration**: Deve ser calculado automaticamente
2. **Analytics**: Métricas devem refletir as entradas
3. **Tags**: Devem aparecer nas listas de tags
4. **Projects**: Se associados, devem aparecer nos filtros

## 🔄 Reset para Novo Teste

1. Execute `Auth/Logout User.bru`
2. Limpe variáveis se necessário
3. Execute `Auth/Register User.bru` com novo email
4. Comece fluxo novamente

---

**Dica**: Mantenha sempre uma aba com "Get Log Entries" aberta para verificar rapidamente as criações!
