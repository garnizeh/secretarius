# EngLog: Request Insight Generation API Documentation

> "The goal is to turn data into information, and information into insight." - Carly Fiorina 🔍

## Visão Geral

Este documento fornece uma documentação abrangente sobre o endpoint `POST /v1/tasks/insights` do EngLog, que é responsável por solicitar a geração de insights baseados em IA para atividades de trabalho dos usuários. O sistema utiliza uma arquitetura distribuída com comunicação gRPC entre o API Server (Machine 1) e o Worker Service (Machine 2) para processamento assíncrono via Ollama LLM.

## Índice

1. [Arquitetura do Sistema](#arquitetura-do-sistema)
2. [Fluxo de Requisição](#fluxo-de-requisição)
3. [Especificação da API](#especificação-da-api)
4. [Estrutura de Dados](#estrutura-de-dados)
5. [Diagramas de Sequência](#diagramas-de-sequência)
6. [Processamento Interno](#processamento-interno)
7. [Estados e Transições](#estados-e-transições)
8. [Tratamento de Erros](#tratamento-de-erros)
9. [Exemplos de Uso](#exemplos-de-uso)
10. [Monitoramento e Observabilidade](#monitoramento-e-observabilidade)

## Arquitetura do Sistema

### Visão Geral da Arquitetura Distribuída

```mermaid
graph TB
    subgraph "Machine 1 - API Server (Public)"
        Client[Cliente/Frontend]
        API[API Server :8080]
        DB[(PostgreSQL)]
        Redis[(Redis Cache)]
        GRPCServer[gRPC Server :50051]
    end

    subgraph "Machine 2 - Worker Server (Private)"
        Worker[Worker Service]
        Ollama[Ollama LLM :11434]
        GRPCClient[gRPC Client]
    end

    Client -->|POST /v1/tasks/insights| API
    API -->|Store Task| DB
    API -->|Session Cache| Redis
    API -->|Queue Task| GRPCServer
    GRPCServer -.->|gRPC/TLS| GRPCClient
    GRPCClient -->|Process Task| Worker
    Worker -->|Generate Insight| Ollama
    Ollama -->|AI Response| Worker
    Worker -->|Result| GRPCClient
    GRPCClient -.->|Report Result| GRPCServer
    GRPCServer -->|Update Task| DB

    style API fill:#e1f5fe
    style Worker fill:#f3e5f5
    style Ollama fill:#fff3e0
    style DB fill:#e8f5e8
    style Redis fill:#ffebee
```

### Componentes Principais

- **API Server (Machine 1)**: Expõe endpoints REST, gerencia autenticação e orquestra tarefas
- **Worker Service (Machine 2)**: Processa tarefas de IA usando Ollama LLM
- **gRPC Communication**: Comunicação segura entre máquinas via TLS
- **Banco de Dados**: PostgreSQL para persistência de dados e resultados
- **Cache**: Redis para sessões e cache de resultados

## Fluxo de Requisição

### Diagrama de Fluxo Geral

```mermaid
flowchart TD
    Start([Início: Cliente faz requisição]) --> Auth{Autenticação JWT válida?}
    Auth -->|Não| AuthError[Retorna 401 Unauthorized]
    Auth -->|Sim| Validate{Validação do payload?}
    Validate -->|Falha| ValidationError[Retorna 400 Bad Request]
    Validate -->|Sucesso| QueueTask[Cria Task no gRPC Manager]
    QueueTask --> GenerateID[Gera Task ID único]
    GenerateID --> StoreDB[(Armazena task no DB)]
    StoreDB --> QueueGRPC[Enfileira task via gRPC]
    QueueGRPC --> ReturnResponse[Retorna 202 Accepted + task_id]
    ReturnResponse --> WorkerReceive[Worker recebe task via stream]
    WorkerReceive --> ProcessTask[Processa task com Ollama]
    ProcessTask --> UpdateProgress[Atualiza progresso da task]
    UpdateProgress --> AIGeneration[Gera insight com IA]
    AIGeneration --> ReportResult[Reporta resultado via gRPC]
    ReportResult --> UpdateDB[(Atualiza status no DB)]
    UpdateDB --> End([Fim: Task completada])

    AuthError --> End
    ValidationError --> End

    style Start fill:#c8e6c9
    style End fill:#ffcdd2
    style ProcessTask fill:#fff3e0
    style AIGeneration fill:#f3e5f5
```

## Especificação da API

### Endpoint

```
POST /v1/tasks/insights
```

### Headers Obrigatórios

```http
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

### Request Body

```json
{
  "user_id": "string (UUID, obrigatório)",
  "entry_ids": ["string (UUID array, obrigatório)"],
  "insight_type": "string (obrigatório)",
  "context": "string | object (opcional, flexível)"
}
```

### Contexto Flexível (Context Field)

O campo `context` aceita tanto strings simples quanto objetos estruturados para máxima flexibilidade:

#### Contexto Simples (String)
```json
{
  "context": "Weekly productivity analysis for performance review"
}
```

#### Contexto Estruturado (Object)
```json
{
  "context": {
    "analysis_type": "weekly_productivity",
    "purpose": "performance_review",
    "date_range": {
      "start": "2025-07-01",
      "end": "2025-07-31"
    },
    "filters": {
      "include_meetings": true,
      "min_value_rating": "medium",
      "activity_types": ["development", "code_review"]
    },
    "preferences": {
      "detail_level": "high",
      "include_recommendations": true,
      "focus_areas": ["time_optimization", "skill_development"]
    }
  }
}
```

### Response

#### Sucesso (202 Accepted)
```json
{
  "task_id": "insight_dd72f29f-a51b-4a2d-add8-496537c8e078_1722470400",
  "message": "Insight generation task queued successfully"
}
```

#### Erro de Validação (400 Bad Request)
```json
{
  "error": "Key: 'user_id' Error:Field validation for 'user_id' failed on the 'required' tag"
}
```

#### Erro de Autenticação (401 Unauthorized)
```json
{
  "error": "Invalid or expired token"
}
```

#### Erro Interno (500 Internal Server Error)
```json
{
  "error": "Failed to queue insight generation task: connection timeout"
}
```

## Estrutura de Dados

### Estrutura da Requisição (Go)

```go
type InsightGenerationRequest struct {
    UserID      string   `json:"user_id" binding:"required"`
    EntryIDs    []string `json:"entry_ids" binding:"required"`
    InsightType string   `json:"insight_type" binding:"required"`
    Context     string   `json:"context"`
}
```

### Payload da Task (JSON)

```go
type TaskPayload struct {
    UserID      string   `json:"user_id"`
    EntryIDs    []string `json:"entry_ids"`
    InsightType string   `json:"insight_type"`
    Context     string   `json:"context"`
}
```

### Estrutura do Insight (AI Response)

```go
type Insight struct {
    Content    string   `json:"content"`
    Tags       []string `json:"tags"`
    Confidence float32  `json:"confidence"`
}
```

## Diagramas de Sequência

### Sequência Completa de Processamento

```mermaid
sequenceDiagram
    participant C as Cliente
    participant API as API Server
    participant DB as PostgreSQL
    participant GM as gRPC Manager
    participant W as Worker Service
    participant O as Ollama LLM

    Note over C,O: Fase 1: Recepção e Validação
    C->>+API: POST /v1/tasks/insights
    Note right of C: Headers: Authorization, Content-Type<br/>Body: user_id, entry_ids, insight_type, context

    API->>API: Validar JWT Token
    API->>API: Validar Request Body

    Note over C,O: Fase 2: Criação e Enfileiramento da Task
    API->>+GM: QueueInsightGenerationTask()
    GM->>GM: Gerar Task ID único
    Note right of GM: Format: insight_{userID}_{timestamp}

    GM->>GM: Criar TaskRequest (gRPC)
    Note right of GM: TaskType: TASK_TYPE_INSIGHT_GENERATION<br/>Priority: 5, Deadline: 5 min

    GM->>DB: Armazenar task metadata
    GM->>+W: StreamTasks (gRPC)
    Note right of GM: Via gRPC bidirectional stream

    GM-->>-API: Retorna Task ID
    API-->>-C: 202 Accepted + task_id

    Note over C,O: Fase 3: Processamento Assíncrono
    W->>+W: Receber task via stream
    W->>W: Verificar semáforo de concorrência
    W->>W: Adicionar à lista de tasks ativas

    W->>GM: UpdateTaskProgress(25%, "Starting AI insight generation")

    W->>W: Deserializar payload JSON
    Note right of W: Extrair InsightRequest do payload

    W->>+O: GenerateInsight(prompt)
    Note right of O: Modelo: llama3.2:3b<br/>Timeout: 120s, Retry: 3x
    O-->>-W: Insight com confidence score

    W->>GM: UpdateTaskProgress(100%, "Insight generation completed")

    W->>W: Serializar resultado para JSON
    W->>+GM: ReportTaskResult()
    Note right of W: Status: COMPLETED<br/>Result: JSON insight<br/>Timestamps

    GM->>DB: Atualizar task status
    GM-->>-W: Confirmação

    W->>W: Remover da lista de tasks ativas
    W->>W: Atualizar estatísticas

    Note over C,O: Fase 4: Recuperação do Resultado
    C->>+API: GET /v1/tasks/{task_id}/result
    API->>GM: GetTaskResult(task_id)
    GM->>DB: Buscar resultado da task
    GM-->>API: TaskResult
    API-->>-C: 200 OK + resultado completo
```

### Sequência de Error Handling

```mermaid
sequenceDiagram
    participant C as Cliente
    participant API as API Server
    participant GM as gRPC Manager
    participant W as Worker Service
    participant CB as Circuit Breaker
    participant O as Ollama LLM

    Note over C,O: Cenário: Falha no Processamento de IA
    C->>+API: POST /v1/tasks/insights
    API->>+GM: QueueInsightGenerationTask()
    GM->>+W: StreamTasks (gRPC)
    API-->>-C: 202 Accepted + task_id

    W->>+CB: Execute com Circuit Breaker
    CB->>W: Permitir execução

    W->>GM: UpdateTaskProgress(25%, "Starting AI insight generation")

    W->>+O: GenerateInsight(prompt)
    Note right of O: Cenário: Ollama indisponível
    O-->>-W: Timeout Error

    W->>W: Retry Logic (3x com backoff)
    W->>+O: GenerateInsight(prompt) - Retry 1
    O-->>-W: Timeout Error

    W->>+O: GenerateInsight(prompt) - Retry 2
    O-->>-W: Timeout Error

    W->>+O: GenerateInsight(prompt) - Retry 3
    O-->>-W: Timeout Error

    W-->>-CB: Error: AI generation failed
    CB->>CB: Incrementar failure count

    W->>+GM: ReportTaskResult()
    Note right of W: Status: FAILED<br/>Error: "AI insight generation failed: timeout"
    GM->>DB: Atualizar task com erro
    GM-->>-W: Confirmação

    Note over C,O: Cliente consulta resultado
    C->>+API: GET /v1/tasks/{task_id}/result
    API->>GM: GetTaskResult(task_id)
    GM-->>API: TaskResult (status: FAILED)
    API-->>-C: 200 OK + error details
```

### Sequência de Recovery e Reconnection

```mermaid
sequenceDiagram
    participant W as Worker Service
    participant CM as Connection Manager
    participant CB as Circuit Breaker
    participant API as API Server

    Note over W,API: Cenário: Perda de Conexão gRPC
    W->>+CM: Health Check Routine
    CM->>CM: Verificar estado da conexão
    CM->>+API: gRPC Health Check
    Note right of API: Conexão perdida
    API-->>-CM: Connection Error

    CM-->>-W: IsConnected() = false
    W->>W: setConnected(false)

    Note over W,API: Tentativa de Reconnection
    W->>+CM: Reconnect()
    CM->>CM: Exponential Backoff
    CM->>+API: Estabelecer nova conexão gRPC
    API-->>-CM: Connection Success
    CM-->>-W: Reconnection Success

    Note over W,API: Re-registration do Worker
    W->>+CB: Execute re-registration
    CB->>W: Permitir execução
    W->>+API: RegisterWorker()
    API-->>-W: Session Token
    W->>W: setConnected(true)
    W-->>-CB: Registration Success

    Note over W,API: Retomar Task Streaming
    W->>+API: StreamTasks() - Nova stream
    API-->>-W: Task Stream Ready
    W->>W: Continuar processamento normal
```

## Processamento Interno

### Fluxo do gRPC Manager

```mermaid
graph TD
    subgraph "gRPC Manager - QueueInsightGenerationTask"
        A[Receber Parâmetros] --> B[Gerar Task ID]
        B --> C[Criar Task Payload JSON]
        C --> D[Criar TaskRequest protobuf]
        D --> E[Definir Prioridade e Deadline]
        E --> F[Chamar server.QueueTask]
        F --> G{Task Enfileirada?}
        G -->|Sim| H[Retornar Task ID]
        G -->|Não| I[Retornar Erro]
    end

    subgraph "Task Payload Structure"
        J[user_id: string]
        K[entry_ids: string array]
        L[insight_type: string]
        M[context: string]
    end

    C --> J
    C --> K
    C --> L
    C --> M

    style A fill:#e3f2fd
    style H fill:#c8e6c9
    style I fill:#ffcdd2
```

### Fluxo do Worker Service

```mermaid
graph TD
    subgraph "Worker Service - processInsightTask"
        A[Receber TaskRequest via Stream] --> B[Deserializar Payload JSON]
        B --> C[Criar InsightRequest]
        C --> D[UpdateTaskProgress 25%]
        D --> E[Chamar aiService.GenerateInsight]
        E --> F{AI Generation OK?}
        F -->|Sim| G[UpdateTaskProgress 100%]
        F -->|Não| H[Retry Logic com Backoff]
        H --> I{Retry Successful?}
        I -->|Sim| G
        I -->|Não| J[ReportTaskResult FAILED]
        G --> K[Serializar Insight para JSON]
        K --> L[ReportTaskResult COMPLETED]
        L --> M[Limpar Task Ativa]
        J --> M
        M --> N[Atualizar Estatísticas]
    end

    subgraph "AI Service Integration"
        O[Ollama HTTP Client]
        P[GenerateRequest]
        Q[Model: llama3.2:3b]
        R[Timeout: 120s]
        S[Retry: 3x]
    end

    E --> O
    O --> P
    P --> Q
    P --> R
    P --> S

    style A fill:#e3f2fd
    style L fill:#c8e6c9
    style J fill:#ffcdd2
    style O fill:#fff3e0
```

## Estados e Transições

### Diagrama de Estados da Task

```mermaid
stateDiagram-v2
    [*] --> Created: POST /v1/tasks/insights
    Created --> Queued: QueueInsightGenerationTask()
    Queued --> Streaming: Worker recebe via stream
    Streaming --> Processing: processInsightTask()
    Processing --> AI_Generation: aiService.GenerateInsight()

    AI_Generation --> Completed: Insight gerado com sucesso
    AI_Generation --> Retrying: Falha temporária (timeout/connection)
    AI_Generation --> Failed: Falha permanente (3 retries)

    Retrying --> AI_Generation: Retry attempt
    Retrying --> Failed: Max retries exceeded

    Processing --> Failed: Erro de deserialização
    Streaming --> Failed: Erro de stream

    Completed --> [*]: ReportTaskResult(COMPLETED)
    Failed --> [*]: ReportTaskResult(FAILED)

    note right of Created
        Task ID gerado
        Payload criado
        Deadline definida (5 min)
    end note

    note right of AI_Generation
        Ollama LLM Processing
        Model: llama3.2:3b
        Timeout: 120s
    end note

    note right of Retrying
        Exponential Backoff
        Max 3 retries
        Circuit Breaker monitoring
    end note
```

### Estados de Conexão do Worker

```mermaid
stateDiagram-v2
    [*] --> Disconnected
    Disconnected --> Connecting: connectionManager.Connect()
    Connecting --> Connected: Connection established
    Connecting --> Disconnected: Connection failed

    Connected --> Registering: registerWorkerWithRetry()
    Registering --> Registered: Registration successful
    Registering --> Disconnected: Registration failed

    Registered --> Streaming: StreamTasks()
    Streaming --> Processing: Task received
    Processing --> Streaming: Task completed

    Streaming --> Reconnecting: Connection lost
    Processing --> Reconnecting: Connection lost
    Registered --> Reconnecting: Heartbeat failed

    Reconnecting --> Connecting: Retry connection
    Reconnecting --> Disconnected: Max retries exceeded

    note right of Registered
        Session token received
        Heartbeat active (30s)
        Ready for tasks
    end note

    note right of Processing
        Max concurrent: 5 tasks
        Progress updates
        Circuit breaker active
    end note
```

## Tratamento de Erros

### Hierarquia de Erros

```mermaid
graph TD
    subgraph "API Layer Errors"
        E1[400 Bad Request]
        E2[401 Unauthorized]
        E3[500 Internal Server Error]
    end

    subgraph "gRPC Communication Errors"
        E4[Connection Timeout]
        E5[Service Unavailable]
        E6[Authentication Failed]
    end

    subgraph "Worker Processing Errors"
        E7[Deserialization Error]
        E8[AI Generation Timeout]
        E9[Circuit Breaker Open]
    end

    subgraph "AI Service Errors"
        E10[Ollama Unavailable]
        E11[Model Loading Error]
        E12[Prompt Processing Error]
    end

    E1 --> R1[Retorna erro imediato]
    E2 --> R1
    E3 --> R1

    E4 --> R2[Retry com backoff]
    E5 --> R2
    E6 --> R3[Re-authentication]

    E7 --> R4[Report FAILED status]
    E8 --> R5[Retry with circuit breaker]
    E9 --> R4

    E10 --> R5
    E11 --> R5
    E12 --> R5

    R2 --> R6{Max retries?}
    R5 --> R6
    R6 -->|No| R7[Continue retry]
    R6 -->|Yes| R4

    style R1 fill:#ffcdd2
    style R4 fill:#ffcdd2
    style R7 fill:#fff3e0
```

### Circuit Breaker States

```mermaid
graph LR
    subgraph "Circuit Breaker - Task Processing"
        CB1[Closed] -->|Failure threshold reached| CB2[Open]
        CB2 -->|Timeout period| CB3[Half-Open]
        CB3 -->|Success| CB1
        CB3 -->|Failure| CB2
    end

    subgraph "Configuration"
        CF1[Failure Threshold: 5]
        CF2[Success Threshold: 3]
        CF3[Timeout: 30s]
    end

    CB1 -.-> CF1
    CB2 -.-> CF3
    CB3 -.-> CF2

    style CB1 fill:#c8e6c9
    style CB2 fill:#ffcdd2
    style CB3 fill:#fff3e0
```

## Exemplos de Uso

### Exemplo 1: Requisição Básica (cURL)

```bash
curl --request POST \
  --url http://localhost:8080/v1/tasks/insights \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...' \
  --header 'content-type: application/json' \
  --data '{
    "user_id": "dd72f29f-a51b-4a2d-add8-496537c8e078",
    "entry_ids": [
      "550e8400-e29b-41d4-a716-446655440001",
      "550e8400-e29b-41d4-a716-446655440002"
    ],
    "insight_type": "productivity",
    "context": "Weekly productivity analysis for performance review"
  }'
```

**Response:**
```json
{
  "task_id": "insight_dd72f29f-a51b-4a2d-add8-496537c8e078_1722470400",
  "message": "Insight generation task queued successfully"
}
```

### Exemplo 2: Verificação do Resultado

```bash
curl --request GET \
  --url http://localhost:8080/v1/tasks/insight_dd72f29f-a51b-4a2d-add8-496537c8e078_1722470400/result \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
```

**Response (Task Completed):**
```json
{
  "task_id": "insight_dd72f29f-a51b-4a2d-add8-496537c8e078_1722470400",
  "worker_id": "worker-1",
  "status": "TASK_STATUS_COMPLETED",
  "result": "{\"content\":\"Based on the productivity analysis...\",\"tags\":[\"productivity\",\"efficiency\"],\"confidence\":0.85}",
  "error": "",
  "started_at": "2025-08-01T10:30:00Z",
  "completed_at": "2025-08-01T10:32:30Z"
}
```

### Exemplo 3: Contexto Estruturado Avançado

#### Análise de Produtividade com Contexto Estruturado
```bash
curl --request POST \
  --url http://localhost:8080/v1/tasks/insights \
  --header 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...' \
  --header 'content-type: application/json' \
  --data '{
    "user_id": "dd72f29f-a51b-4a2d-add8-496537c8e078",
    "entry_ids": ["uuid1", "uuid2", "uuid3"],
    "insight_type": "productivity",
    "context": {
      "analysis_type": "weekly_productivity",
      "purpose": "performance_review",
      "date_range": {
        "start": "2025-07-01",
        "end": "2025-07-31"
      },
      "filters": {
        "include_meetings": true,
        "min_value_rating": "medium",
        "activity_types": ["development", "code_review", "debugging"]
      },
      "preferences": {
        "detail_level": "high",
        "include_recommendations": true,
        "focus_areas": ["time_optimization", "efficiency_patterns"]
      }
    }
  }'
```

#### Análise de Desenvolvimento de Habilidades
```json
{
  "user_id": "dd72f29f-a51b-4a2d-add8-496537c8e078",
  "entry_ids": ["uuid1", "uuid2"],
  "insight_type": "skill_development",
  "context": {
    "career_goal": "senior_architect",
    "focus_areas": ["golang", "system_design", "leadership"],
    "current_level": "senior_developer",
    "time_frame": "6_months",
    "learning_preferences": {
      "hands_on": true,
      "mentorship": true,
      "certification": false
    },
    "skill_assessment": {
      "golang": "advanced",
      "architecture": "intermediate",
      "leadership": "beginner"
    }
  }
}
```

#### Análise de Gestão de Tempo
```json
{
  "user_id": "dd72f29f-a51b-4a2d-add8-496537c8e078",
  "entry_ids": ["uuid1", "uuid2", "uuid3"],
  "insight_type": "time_management",
  "context": {
    "optimization_goal": "reduce_context_switching",
    "time_blocks": ["morning", "afternoon"],
    "problem_areas": ["too_many_meetings", "fragmented_coding_time"],
    "constraints": {
      "meeting_free_hours": ["09:00-11:00", "14:00-16:00"],
      "deep_work_preference": "morning"
    },
    "metrics_focus": ["time_in_flow_state", "interruption_frequency"],
    "action_preferences": {
      "calendar_restructuring": true,
      "notification_management": true,
      "task_batching": true
    }
  }
}
```

### Exemplo 4: Compatibilidade com Contexto String (Backward Compatible)

```json
{
  "user_id": "dd72f29f-a51b-4a2d-add8-496537c8e078",
  "entry_ids": ["uuid1", "uuid2"],
  "insight_type": "productivity",
  "context": "Focus on identifying time optimization opportunities and efficiency patterns"
}
```## Monitoramento e Observabilidade

### Métricas de Performance

```mermaid
graph TB
    subgraph "API Server Metrics"
        M1[Request Rate /min]
        M2[Response Time P95]
        M3[Error Rate %]
        M4[Task Queue Length]
    end

    subgraph "Worker Service Metrics"
        M5[Active Tasks]
        M6[Completed Tasks]
        M7[Failed Tasks]
        M8[AI Generation Time]
    end

    subgraph "gRPC Communication Metrics"
        M9[Connection Status]
        M10[Message Latency]
        M11[Reconnection Count]
        M12[Circuit Breaker State]
    end

    subgraph "AI Service Metrics"
        M13[Ollama Response Time]
        M14[Model Loading Time]
        M15[Confidence Score Avg]
        M16[AI Availability %]
    end

    style M1 fill:#e3f2fd
    style M5 fill:#f3e5f5
    style M9 fill:#fff3e0
    style M13 fill:#e8f5e8
```

### Logs Estruturados

```json
{
  "timestamp": "2025-08-01T10:30:00Z",
  "level": "INFO",
  "component": "insight_generation",
  "operation": "request_received",
  "user_id": "dd72f29f-a51b-4a2d-add8-496537c8e078",
  "task_id": "insight_dd72f29f-a51b-4a2d-add8-496537c8e078_1722470400",
  "entry_count": 2,
  "insight_type": "productivity",
  "request_id": "req-123456789"
}
```

```json
{
  "timestamp": "2025-08-01T10:32:30Z",
  "level": "INFO",
  "component": "worker_service",
  "operation": "ai_generation_completed",
  "task_id": "insight_dd72f29f-a51b-4a2d-add8-496537c8e078_1722470400",
  "worker_id": "worker-1",
  "duration_ms": 150000,
  "confidence_score": 0.85,
  "ai_model": "llama3.2:3b"
}
```

### Health Checks

```mermaid
graph TD
    subgraph "Health Check Endpoints"
        H1[GET /health - Basic health]
        H2[GET /ready - Readiness check]
        H3[GET /v1/workers/health - Worker system health]
    end

    subgraph "Health Status Components"
        S1[API Server Status]
        S2[Database Connectivity]
        S3[Redis Connectivity]
        S4[gRPC Server Status]
        S5[Active Workers Count]
        S6[AI Service Status]
    end

    H1 --> S1
    H1 --> S2
    H2 --> S3
    H2 --> S4
    H3 --> S5
    H3 --> S6

    style H1 fill:#c8e6c9
    style H3 fill:#e3f2fd
```

## Considerações de Performance

### Otimizações Implementadas

1. **Concorrência Limitada**: Máximo 5 tasks simultâneas por worker
2. **Connection Pooling**: Pool de conexões gRPC reutilizáveis
3. **Circuit Breaker**: Proteção contra falhas em cascata
4. **Exponential Backoff**: Retry inteligente com jitter
5. **Progress Updates**: Atualizações incrementais de progresso
6. **Resource Management**: Semáforos para controle de recursos

### Benchmarks Esperados

- **Request Response Time**: < 100ms (para enfileiramento)
- **AI Generation Time**: 30-180s (dependendo da complexidade)
- **gRPC Latency**: < 50ms (comunicação inter-service)
- **Throughput**: 50-100 requests/min por worker
- **Success Rate**: > 99% (com retry logic)

## Conclusão

O sistema de Request Insight Generation do EngLog implementa uma arquitetura robusta e escalável para processamento assíncrono de insights baseados em IA. Com comunicação gRPC segura, retry logic avançado, circuit breakers e monitoramento abrangente, o sistema garante alta disponibilidade e performance para análise de produtividade dos usuários.

A implementação distributed permite isolamento de responsabilidades, onde o API Server foca na interface e orquestração, enquanto o Worker Service se dedica ao processamento intensivo de IA, garantindo que a experiência do usuário permaneça responsiva mesmo durante operações computacionalmente intensivas.
