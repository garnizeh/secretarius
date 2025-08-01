# EngLog TODO List

> "The best way to get started is to quit talking and begin doing." - Walt Disney 🚀

## High Priority Tasks

### 🎯 Task 0125: Sistema de Configuração Hierárquica de Modelos LLM
**Priority**: HIGH
**Estimated Effort**: 4-6 weeks
**Due Date**: Q3 2025
**Document**: [LLM Config Hierarchy Proposal](docs/llm-config-hierarchy-proposal.md)

**Description**:
Implementar sistema hierárquico de configuração de modelos LLM para substituir o modelo hardcoded atual ("llama3.2:3b") por um sistema flexível com precedência: User > TaskType > System.

**Subtasks**:

#### 📋 **Fase 1: Foundation (1-2 semanas)**
- [ ] **Task 0125.1**: Atualizar protocolo gRPC (`proto/worker.proto`)
  - [ ] Adicionar mensagem `LLMConfig` com campos: model, timeout_seconds, max_retries, fallback_model, parameters
  - [ ] Adicionar `LLMConfig llm_config` na `TaskRequest`
  - [ ] Adicionar `repeated string supported_models` e `LLMConfig default_llm_config` na `RegisterWorkerRequest`
  - [ ] Regenerar código gRPC com `make proto`

- [ ] **Task 0125.2**: Implementar configuração do sistema (`internal/config/config.go`)
  - [ ] Adicionar struct `LLMConfig` com DefaultModel, FallbackModel, Timeout, MaxRetries
  - [ ] Adicionar `TaskTypes map[string]TaskLLMConfig` para configuração por tipo de task
  - [ ] Adicionar `ModelParams map[string]ModelParams` para parâmetros específicos de modelo
  - [ ] Integrar `LLM LLMConfig` na `WorkerConfig`

- [ ] **Task 0125.3**: Criar service de resolução (`internal/llm/resolver.go`)
  - [ ] Implementar `ConfigResolver` struct com cache
  - [ ] Implementar `ResolveConfig(ctx, userID, taskType)` com hierarquia
  - [ ] Implementar cache com TTL e invalidação inteligente
  - [ ] Testes unitários para algoritmo de resolução

- [ ] **Task 0125.4**: Modificar OllamaService (`internal/ai/ollama.go`)
  - [ ] Refatorar `GenerateInsight` para `GenerateInsightWithConfig`
  - [ ] Implementar fallback automático para modelo secundário
  - [ ] Remover modelo hardcoded de todas as funções
  - [ ] Adicionar logs detalhados sobre modelo usado e fallbacks

#### 📋 **Fase 2: Integration (1-2 semanas)**
- [ ] **Task 0125.5**: Database schema para preferências de usuário
  - [ ] Criar migration para tabela `user_llm_preferences`
  - [ ] Campos: user_id, preferred_model, fallback_model, timeout_seconds, max_retries, task_type_configs (JSONB)
  - [ ] Índices de performance e constraints de validação
  - [ ] Trigger para updated_at

- [ ] **Task 0125.6**: Integrar ConfigResolver no Worker (`internal/worker/client.go`)
  - [ ] Modificar `processInsightTask` para usar configuração dinâmica
  - [ ] Implementar `getDefaultLLMConfig` por TaskType
  - [ ] Converter protobuf LLMConfig para configuração interna
  - [ ] Atualizar registro do worker com modelos suportados

- [ ] **Task 0125.7**: Implementar configuração por TaskType
  - [ ] Carregar configurações do arquivo YAML/env vars
  - [ ] Mapear TaskType enum para configurações específicas
  - [ ] Validar modelos disponíveis na startup
  - [ ] Testes de integração com diferentes configurações

#### 📋 **Fase 3: User Configuration (1-2 semanas)**
- [ ] **Task 0125.8**: Service layer para LLM config (`internal/services/llm_service.go`)
  - [ ] Implementar `LLMConfigService` com CRUD operations
  - [ ] `GetUserLLMConfig`, `UpdateUserLLMConfig`, `DeleteUserLLMConfig`
  - [ ] `ResolveConfigForTask`, `ValidateModel`, `ListAvailableModels`
  - [ ] Integração com ConfigResolver

- [ ] **Task 0125.9**: API endpoints (`internal/handlers/llm_config.go`)
  - [ ] `GET /v1/users/llm-config` - Obter configuração atual
  - [ ] `PUT /v1/users/llm-config` - Atualizar configuração
  - [ ] `DELETE /v1/users/llm-config` - Resetar para padrão
  - [ ] `GET /v1/llm/models` - Listar modelos disponíveis
  - [ ] `GET /v1/llm/config/preview` - Preview de configuração

- [ ] **Task 0125.10**: SQLC queries para LLM preferences
  - [ ] Queries CRUD para `user_llm_preferences`
  - [ ] Queries para validação de modelos
  - [ ] Queries para estatísticas de uso
  - [ ] Regenerar com `make sqlc`

#### 📋 **Fase 4: Production Ready (1 semana)**
- [ ] **Task 0125.11**: Observabilidade e monitoring
  - [ ] Métricas de uso por modelo via Prometheus
  - [ ] Alertas para fallbacks frequentes
  - [ ] Dashboard de configurações ativas
  - [ ] Logs estruturados com modelo usado

- [ ] **Task 0125.12**: Testes comprehensivos
  - [ ] Testes unitários para ConfigResolver
  - [ ] Testes de integração com diferentes configurações
  - [ ] Testes end-to-end via Bruno collection
  - [ ] Testes de performance com diferentes modelos
  - [ ] Chaos engineering para cenários de fallback

- [ ] **Task 0125.13**: Configuração de ambiente
  - [ ] Variáveis de ambiente para configuração padrão
  - [ ] Arquivo YAML de configuração (`config/llm.yaml`)
  - [ ] Docker compose com configurações de exemplo
  - [ ] Documentação de deployment

- [ ] **Task 0125.14**: Documentação e deployment
  - [ ] Atualizar documentação da API
  - [ ] Guia de migração do sistema atual
  - [ ] Bruno collection com novos endpoints
  - [ ] Deploy em staging e validação
  - [ ] Deploy em production com rollback plan

**Dependencies**:
- Depende do sistema atual de tasks e insights estar estável
- Requer que Ollama esteja configurado com múltiplos modelos
- Necessita de database migration bem testada

**Acceptance Criteria**:
- [ ] Nenhum modelo hardcoded no código
- [ ] Configuração hierárquica funcional: User > TaskType > System
- [ ] Fallback automático implementado e testado
- [ ] APIs para gerenciamento de configuração
- [ ] Performance igual ou melhor que sistema atual
- [ ] Backward compatibility durante migração
- [ ] Documentação completa e testes abrangentes

**Notes**:
- Implementação incremental com rollback em cada fase
- Manter sistema atual funcionando durante desenvolvimento
- Validar modelos na startup para evitar erros em runtime
- Cache para otimizar performance da resolução hierárquica

### 🎯 Task 0130: Dynamic Professional Role Templates & AI Content Generation
**Priority**: HIGH
**Estimated Effort**: 3-4 weeks
**Due Date**: Q4 2025

**Description**:
Currently, the system is hardcoded for software engineers with activity types like `development`, `code_review`, `debugging`, etc. Implement a flexible role-based template system that allows the application to adapt to different professions (lawyers, managers, accountants, doctors, teachers, etc.) with profession-specific activity types and AI-generated content.

**Requirements**:

1. **Role Template System**:
   - Database schema for profession roles and templates
   - CRUD operations for role management
   - User role assignment and switching
   - Default templates for common professions

2. **Dynamic Activity Types**:
   - Role-specific activity type definitions
   - Template-based activity type generation
   - Migration system for existing data
   - Backward compatibility with current engineer-focused types

3. **AI Content Adaptation**:
   - Role-aware prompt templates for Ollama LLM
   - Profession-specific insight generation
   - Context-aware recommendations
   - Industry-specific terminology and metrics

4. **Template Examples**:
   ```yaml
   # Software Engineer (current)
   activity_types: [development, code_review, debugging, testing, deployment]
   insights_focus: [productivity, code_quality, technical_debt, sprint_goals]

   # Lawyer
   activity_types: [case_research, client_meeting, document_review, court_appearance, legal_writing]
   insights_focus: [billable_hours, case_progress, client_satisfaction, legal_strategy]

   # Manager
   activity_types: [team_meeting, one_on_one, strategic_planning, performance_review, budget_planning]
   insights_focus: [team_productivity, goal_achievement, leadership_effectiveness, resource_allocation]

   # Accountant
   activity_types: [bookkeeping, tax_preparation, financial_analysis, client_consultation, audit_work]
   insights_focus: [accuracy_metrics, deadline_compliance, client_portfolio, seasonal_workload]
   ```

5. **Database Changes**:
   - New table: `profession_roles` (id, name, description, config_json)
   - New table: `role_activity_types` (role_id, name, description, value_factors)
   - New table: `role_ai_templates` (role_id, template_type, prompt_template, config)
   - Update `users` table: add `profession_role_id` foreign key
   - Update `log_entries` table: keep activity type flexible (enum → string + validation)

6. **API Changes**:
   - `/v1/roles` - CRUD operations for profession roles
   - `/v1/roles/{id}/activity-types` - Get role-specific activity types
   - `/v1/users/role` - Update user's profession role
   - Update existing endpoints to be role-aware

7. **AI Template Engine**:
   - Template parser for dynamic prompt generation
   - Role-specific context injection
   - Profession vocabulary and terminology
   - Industry-specific metrics and KPIs

**Technical Implementation**:

```go
// New domain models
type ProfessionRole struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Config      RoleConfig `json:"config"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type RoleConfig struct {
    ActivityTypes    []ActivityTypeConfig `json:"activity_types"`
    AITemplates      map[string]string    `json:"ai_templates"`
    MetricsFocus     []string            `json:"metrics_focus"`
    Terminology      map[string]string    `json:"terminology"`
    ValueFactors     map[string]float64   `json:"value_factors"`
}

type ActivityTypeConfig struct {
    Name        string             `json:"name"`
    Description string             `json:"description"`
    Category    string             `json:"category"`
    ValueFactor float64            `json:"value_factor"`
    Metadata    map[string]interface{} `json:"metadata"`
}
```

**Migration Strategy**:
1. Create new tables with role system
2. Create "Software Engineer" default role with current activity types
3. Migrate existing users to default role
4. Gradually add support for new professions
5. Maintain backward compatibility during transition

**AI Template Examples**:
```yaml
software_engineer:
  productivity_analysis: |
    Analyze the software engineer's work patterns based on these activities:
    Focus on code quality, development velocity, and technical growth.
    Consider factors like debugging time, code review participation, and feature delivery.

lawyer:
  productivity_analysis: |
    Analyze the lawyer's practice patterns based on these activities:
    Focus on billable hour efficiency, case progression, and client service quality.
    Consider factors like research depth, document quality, and case outcomes.

manager:
  productivity_analysis: |
    Analyze the manager's leadership effectiveness based on these activities:
    Focus on team development, strategic execution, and organizational impact.
    Consider factors like meeting efficiency, team engagement, and goal achievement.
```

**Testing Requirements**:
- Unit tests for role template system
- Integration tests for role-aware APIs
- AI template generation testing
- Migration testing with existing data
- User experience testing with different roles

**Documentation Updates**:
- API documentation for new role endpoints
- User guide for profession role selection
- Developer guide for adding new profession templates
- Migration guide for existing users

**Future Enhancements**:
- Custom role creation by enterprise users
- Industry-specific integrations (legal case management, project management tools)
- Role-based analytics dashboards
- Cross-role team collaboration features

---

## Medium Priority Tasks

### 📋 Task 0140: Enhanced Report Generation (PDF/CSV)
**Priority**: MEDIUM
**Estimated Effort**: 2-3 weeks

**Description**: Implement comprehensive report generation with PDF and CSV export capabilities, customizable templates, and automated scheduling.

**Requirements**:
- PDF report generation with professional layouts
- CSV data export for external analysis
- Customizable report templates
- Automated weekly/monthly report scheduling
- Email delivery integration

---

### 📋 Task 0150: Advanced Rate Limiting & Security
**Priority**: MEDIUM
**Estimated Effort**: 1-2 weeks

**Description**: Implement advanced rate limiting with Redis backend and enhance security headers and CORS configuration.

**Requirements**:
- Redis-based distributed rate limiting
- Per-endpoint and per-user rate limits
- Security headers implementation
- Enhanced CORS configuration
- API abuse prevention

---

### 📋 Task 0160: Full-Text Search Implementation
**Priority**: MEDIUM
**Estimated Effort**: 2 weeks

**Description**: Implement PostgreSQL full-text search (FTS) for activity descriptions, titles, and project names.

**Requirements**:
- PostgreSQL FTS configuration
- Search indexing for relevant fields
- Advanced search API endpoints
- Search result ranking and highlighting
- Performance optimization

---

## Low Priority Tasks

### 📋 Task 0170: Team Collaboration Features
**Priority**: LOW
**Estimated Effort**: 4-6 weeks

**Description**: Implement team collaboration features including shared projects, activity sharing, and team analytics.

**Requirements**:
- Team creation and management
- Shared project workspaces
- Activity sharing and collaboration
- Team-level analytics and insights
- Permission and access control

---

### 📋 Task 0180: External Tool Integrations
**Priority**: LOW
**Estimated Effort**: 3-4 weeks per integration

**Description**: Integrate with external tools like GitHub, Jira, Slack for automatic activity import and synchronization.

**Requirements**:
- GitHub integration (commits, PRs, issues)
- Jira integration (task tracking, time logging)
- Slack integration (notifications, status updates)
- Calendar synchronization (Google, Outlook)
- Webhook support for real-time updates

---

### 📋 Task 0190: Mobile Application Development
**Priority**: LOW
**Estimated Effort**: 8-12 weeks

**Description**: Develop mobile applications (iOS/Android) for on-the-go activity tracking and insights viewing.

**Requirements**:
- React Native or Flutter mobile app
- Activity logging and timing
- Offline capability with sync
- Push notifications
- Mobile-optimized UI/UX

---

## Completed Tasks ✅

### ✅ Task 0020: Database Schema and Migrations (COMPLETED)
- Complete PostgreSQL schema with 8 migration files
- Advanced features: materialized views, triggers, functions
- Performance optimization with indexes

### ✅ Task 0030: SQLC Code Generation and Database Layer (COMPLETED)
- 120+ type-safe database queries
- Complete CRUD operations for all entities
- Advanced analytics queries

### ✅ Task 0040: Core Models and Data Structures (COMPLETED)
- Domain models with comprehensive validation
- Enum definitions and constants
- Type safety across the application

### ✅ Task 0060: Core Business Logic Services (COMPLETED)
- Service layer implementation
- Business logic encapsulation
- Error handling and validation

### ✅ Task 0070: HTTP Handlers and API Endpoints (COMPLETED)
- 38 REST API endpoints
- Comprehensive API testing with Bruno collection
- Swagger/OpenAPI documentation

### ✅ Task 0080: API Server Main Application (COMPLETED)
- Complete API server implementation
- Authentication and middleware
- Health monitoring and logging

### ✅ Task 0090: Worker Service Implementation (COMPLETED)
- Background task processing
- Ollama LLM integration
- Task queue and result management

### ✅ Task 0100: gRPC Communication Setup (COMPLETED)
- Protocol buffer definitions
- Server and client implementation
- Streaming and error handling

### ✅ Task 0110: Testing Framework and Quality Assurance (COMPLETED)
- Comprehensive test suite (unit, integration, e2e)
- 71.7% test coverage
- CI/CD pipeline with GitHub Actions

---

## Notes

- **Version**: Current development is on Phase 2 completion, moving to Phase 3
- **Priority Focus**: Dynamic role templates (Task 0130) is the highest priority for expanding the application beyond software engineers
- **Maintenance**: Regular dependency updates and security patches
- **Documentation**: Keep all documentation updated with new features and changes

---

*Last Updated: August 1, 2025*