# Middleware Integration Tests with Redis

Este arquivo demonstra como executar testes de integração para o middleware de rate limiting usando testcontainers com Redis real.

## Visão Geral

O arquivo `middleware_integration_test.go` fornece testes de integração abrangentes para:

1. **Rate Limiting com Redis**: Testa o comportamento real do rate limiting usando Redis
2. **Stack Completo de Middleware**: Testa a integração entre todos os middleware components
3. **Fallback Gracioso**: Testa o comportamento quando o Redis está indisponível

## Estrutura dos Testes

### 🔧 Configuração de Container
```go
func setupRedisContainer(ctx context.Context) (*RedisContainer, error)
```
- Inicia um container Redis usando testcontainers
- Aguarda o Redis estar pronto para conexões
- Retorna URI de conexão para os testes

### 🚀 Testes de Rate Limiting

#### 1. Requisições Dentro do Limite
- Testa que requisições dentro do limite são permitidas
- Verifica headers de rate limiting corretos
- Confirma logging estruturado

#### 2. Excesso de Requisições
- Testa que requisições além do limite são bloqueadas
- Verifica resposta HTTP 429 (Too Many Requests)
- Confirma headers de rate limiting apropriados

#### 3. Clientes Independentes
- Verifica que diferentes IPs têm limites independentes
- Testa isolamento entre clientes

#### 4. Janela Deslizante
- Testa comportamento de sliding window
- Verifica persistência de dados no Redis

### 🔄 Stack Completo de Middleware
```go
func TestMiddlewareStack_WithRedis_Integration
```
- Testa integração completa: Request Logger + CORS + Rate Limiting + Validation
- Verifica que todos os middleware funcionam juntos
- Confirma order de execução correto

### 🛡️ Teste de Tolerância a Falhas
```go
func TestRedisConnection_FailureHandling
```
- Testa comportamento quando Redis está indisponível
- Verifica fallback gracioso (permite requisições)
- Confirma logging de erros

## Como Executar

### Todos os Testes de Integração
```bash
go test -v ./internal/middleware/ -run Integration
```

### Testes Específicos de Rate Limiting
```bash
go test -v ./internal/middleware/ -run TestRateLimitMiddleware_WithRedis_Integration
```

### Stack Completo
```bash
go test -v ./internal/middleware/ -run TestMiddlewareStack_WithRedis_Integration
```

### Pular Testes de Integração (modo rápido)
```bash
go test -short ./internal/middleware/
```

## Requisitos

- **Docker**: Para executar containers Redis
- **testcontainers-go**: Gerenciamento de containers em testes
- **Redis**: Container Redis 7-alpine

## Configuração Automática

Os testes automaticamente:
1. ⬇️  Baixam a imagem Redis se necessário
2. 🚀 Iniciam container Redis temporário
3. ⏳ Aguardam Redis estar pronto
4. 🧪 Executam testes
5. 🧹 Limpam containers após testes

## Headers de Rate Limiting

Os testes verificam os seguintes headers:
- `X-Rate-Limit-Limit`: Limite máximo de requisições
- `X-Rate-Limit-Remaining`: Requisições restantes no período
- `X-Rate-Limit-Reset`: Timestamp quando o limite reseta

## Logging Estruturado

Todos os testes verificam:
- ✅ Log de requisições HTTP
- ✅ Log de verificações de rate limiting
- ✅ Log de violações de rate limit
- ✅ Log de erros de conexão Redis

## Benefícios dos Testes de Integração

1. **Testa Dependências Reais**: Usa Redis real em vez de mocks
2. **Verifica Timing**: Testa comportamento temporal real
3. **Detecta Problemas de Rede**: Identifica issues de conectividade
4. **Valida Configuração**: Confirma configuração Redis correta
5. **Testa Fallbacks**: Verifica tolerância a falhas

## Exemplo de Uso em CI/CD

```yaml
test:
  script:
    - docker info  # Verifica Docker disponível
    - go test -v ./internal/middleware/  # Inclui testes de integração
```

Para ambientes sem Docker:
```yaml
test-unit:
  script:
    - go test -short ./internal/middleware/  # Pula integração
```

## Debugging

Para debug detalhado, adicione:
```bash
go test -v -count=1 ./internal/middleware/ -run TestRateLimitMiddleware_WithRedis_Integration
```

Os testes incluem logging detalhado mostrando:
- 🐳 Criação e inicialização de containers
- 🔗 Estabelecimento de conexões Redis
- 📊 Contadores de rate limiting
- 🚦 Decisões de allow/deny

## Performance

Os testes de integração são mais lentos que testes unitários devido a:
- Inicialização de containers Docker
- Rede entre teste e Redis
- Cleanup de containers

Tempo típico: ~1-3 segundos por suite de testes
