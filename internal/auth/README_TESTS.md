# Authentication Service Test Suite

Esta é uma suíte abrangente de testes para o serviço de autenticação JWT do projeto englog/EngLog.

## Estrutura dos Testes

### `service_test.go` - Testes Unitários Core
Testa as funcionalidades principais do serviço de autenticação:

- ✅ **Criação de Tokens**: Verifica a geração de tokens de acesso e refresh
- ✅ **Validação de Tokens**: Testa a validação de tokens válidos e inválidos
- ✅ **Hash de Senhas**: Verifica hashing e verificação de senhas com bcrypt
- ✅ **Expiração de Tokens**: Testa o comportamento com tokens expirados
- ✅ **Unicidade**: Verifica que tokens únicos são gerados
- ✅ **Geração de JTI**: Testa a geração de identificadores únicos para refresh tokens

### `middleware_test.go` - Testes de Middleware
Testa os middlewares de autenticação do Gin:

- ✅ **RequireAuth**: Middleware obrigatório para rotas protegidas
- ✅ **OptionalAuth**: Middleware opcional que não bloqueia requisições
- ✅ **Formato de Headers**: Testa vários formatos de header Authorization
- ✅ **Case Sensitivity**: Verifica que "Bearer" é case-insensitive
- ✅ **Validação de Tipos**: Garante que apenas access tokens são aceitos em rotas protegidas

### `integration_test.go` - Testes de Integração (Placeholder)
Contém templates para testes de integração com banco de dados real:

- 🔄 **Rotação de Tokens**: Teste completo do fluxo de rotação
- 🚫 **Denylist**: Testes do sistema de denylist de refresh tokens
- ⚡ **Concorrência**: Testes de operações concorrentes
- 🧹 **Cleanup**: Testes de limpeza de tokens expirados

## Executando os Testes

### Testes Unitários (Padrão)
```bash
# Executar todos os testes
go test ./internal/auth

# Executar com verbose
go test -v ./internal/auth

# Executar testes específicos
go test -run TestCreateAccessToken ./internal/auth
```

### Benchmarks
```bash
# Executar todos os benchmarks
go test -bench=. ./internal/auth

# Benchmark específico
go test -bench=BenchmarkCreateAccessToken ./internal/auth
```

### Testes de Integração (Futuro)
```bash
# Quando implementados, usar build tags
go test -tags=integration ./internal/auth
```

## Cobertura de Testes

### Funcionalidades Testadas ✅
- Criação e validação de tokens JWT
- Middleware de autenticação
- Hash e verificação de senhas
- Geração de identificadores únicos
- Tratamento de erros de autenticação
- Performance de operações críticas

### Limitações dos Testes Unitários ⚠️
- **Sem banco de dados**: Testes unitários não testam operações de denylist
- **Sem persistência**: Não testa rotação de tokens com persistência
- **Sem concorrência real**: Não testa condições de corrida em ambiente real

### Para Testes Completos
Para testar completamente o sistema de autenticação, considere:
1. **Testes de integração** com banco PostgreSQL real
2. **Testes de carga** para validar performance sob stress
3. **Testes de segurança** para vulnerabilidades conhecidas
4. **Testes end-to-end** com requisições HTTP reais

## Métricas de Performance

Baseado nos benchmarks atuais (valores aproximados):
- **Token Creation**: ~17-19µs por token
- **Token Validation**: ~24µs por validação
- **Password Hashing**: ~76ms por hash (bcrypt com custo padrão)
- **Password Check**: ~70ms por verificação
- **Middleware Processing**: ~30µs por requisição

## Estrutura de Arquivos

```
internal/auth/
├── service_test.go      # Testes unitários principais
├── middleware_test.go   # Testes de middleware
├── integration_test.go  # Templates para integração
└── README_TESTS.md     # Este arquivo
```

## Comentários nos Testes

Os testes incluem comentários motivacionais e explicativos em formato de emojis:
- 🏗️ "The beginning is the most important part of the work."
- 🔐 "Security is not a product, but a process."
- 🎯 "Trust, but verify."
- ✅ "Validation is the most sincere form of confirmation."
- 🚫 "Error handling is the art of graceful failure."

## Contribuindo

Ao adicionar novos testes:
1. Siga o padrão de nomenclatura `Test{FunctionName}`
2. Use comentários descritivos com emojis
3. Inclua casos de edge cases
4. Adicione benchmarks para operações críticas
5. Mantenha testes unitários independentes de banco de dados

## Futuras Melhorias

- [ ] Adicionar testes de integração com testcontainers
- [ ] Implementar testes de carga com ferramentas adequadas
- [ ] Adicionar testes de segurança automatizados
- [ ] Implementar coverage reports automáticos
- [ ] Adicionar testes de compatibilidade entre versões
