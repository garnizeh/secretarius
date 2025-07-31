# Guia Prático de Deploy da API em VPS

*"Deploying to production is like going live on stage - you've rehearsed a hundred times, but your heart still skips a beat!" 🎭*

## Visão Geral

Este guia fornece instruções práticas e simplificadas para fazer deploy da API englog em um VPS (Virtual Private Server), utilizando Docker e as ferramentas já configuradas no projeto.

## Pré-requisitos

### Requisitos do Servidor
- **OS**: Ubuntu 22.04 LTS ou similar
- **RAM**: 4GB mínimo, 8GB recomendado
- **CPU**: 2 cores mínimo, 4 cores recomendado
- **Storage**: 50GB SSD mínimo
- **Network**: IP público e domínio configurado

### Requisitos de Software
- Docker Engine 24.0+
- Docker Compose 2.0+
- Git
- Acesso SSH ao servidor

## Método 1: Deploy Automatizado (Recomendado)

### 1. Preparação do VPS

#### 1.1 Acesso e Atualização do Sistema
```bash
# Conecte ao VPS via SSH
ssh root@seu-servidor-ip

# Atualize o sistema
sudo apt update && sudo apt upgrade -y
```

#### 1.2 Instalação do Docker
```bash
# Instale Docker usando o script oficial
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Adicione o usuário ao grupo docker
sudo usermod -aG docker $USER

# Instale Docker Compose
sudo apt install -y docker-compose-plugin

# Verifique as instalações
docker --version
docker compose version
```

### 2. Configuração do Projeto

#### 2.1 Clone do Repositório
```bash
# Clone o projeto
git clone https://github.com/garnizeh/englog.git
cd englog

# Torne os scripts executáveis
chmod +x scripts/*.sh
```

#### 2.2 Configuração de Ambiente
```bash
# Copie o template de produção
cp .env.example .env.production

# Edite as configurações
nano .env.production
```

#### 2.3 Variáveis de Ambiente Obrigatórias
Configure no arquivo `.env.production`:

```bash
# === CONFIGURAÇÕES BÁSICAS ===
APP_ENV=production
APP_HOST=0.0.0.0
APP_PORT=8080
LOG_LEVEL=info

# === DOMÍNIO (SUBSTITUA PELO SEU) ===
DOMAIN_NAME=seu-dominio.com

# === BANCO DE DADOS ===
DB_HOST_READ_WRITE=postgres:5432
DB_HOST_READ_ONLY=postgres:5432
DB_SCHEMA=englog
DB_NAME=englog
DB_USER=englog
DB_PASSWORD=SUA_SENHA_POSTGRES_SEGURA

# === REDIS ===
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=SUA_SENHA_REDIS_SEGURA
REDIS_DB=0

# === JWT (GERE UMA CHAVE SEGURA) ===
JWT_SECRET=sua_chave_jwt_muito_segura_com_pelo_menos_32_caracteres
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=168h
JWT_ISSUER=englog-api

# === gRPC ===
GRPC_HOST=0.0.0.0
GRPC_PORT=9090
GRPC_TIMEOUT=30s
```

### 3. Configuração do Proxy Reverso

#### 3.1 Criar Caddyfile
```bash
# Crie o diretório
mkdir -p deployments/caddy

# Crie o arquivo de configuração do Caddy
cat > deployments/caddy/Caddyfile << 'EOF'
{
    email seu-email@exemplo.com
    admin off
}

# Substitua seu-dominio.com pelo seu domínio real
seu-dominio.com {
    reverse_proxy api-server:8080

    # Headers de segurança
    header {
        # Remover headers que revelam informações do servidor
        -Server
        -X-Powered-By

        # Headers de segurança
        X-Content-Type-Options nosniff
        X-Frame-Options DENY
        X-XSS-Protection "1; mode=block"
        Referrer-Policy strict-origin-when-cross-origin
    }

    # Logs estruturados
    log {
        output file /var/log/caddy/access.log {
            roll_size 100mb
            roll_keep 5
            roll_keep_for 168h
        }
        format json
    }

    # Compressão
    encode gzip

    # Rate limiting básico
    rate_limit {
        zone static {
            key {remote_host}
            events 100
            window 1m
        }
    }
}
EOF
```

### 4. Deploy da Aplicação

#### 4.1 Deploy Automatizado
```bash
# Execute o script de deploy automatizado
make deploy-machine1
```

#### 4.2 Verificação do Deploy
```bash
# Verifique se os serviços estão rodando
docker-compose -f docker-compose.api.yml ps

# Verifique os logs
make prod-api-logs

# Teste o endpoint de health
curl -k https://seu-dominio.com/health
```

## Método 2: Deploy Manual Passo a Passo

### 1. Preparação dos Diretórios
```bash
# Crie diretórios para logs
mkdir -p logs/{api,postgres,redis,caddy}
chmod 755 logs/*
```

### 2. Build das Imagens
```bash
# Build da imagem Docker
docker-compose -f docker-compose.api.yml build
```

### 3. Inicialização dos Serviços
```bash
# Inicie os serviços em background
docker-compose -f docker-compose.api.yml up -d

# Aguarde os serviços ficarem saudáveis
sleep 30

# Verifique o status
docker-compose -f docker-compose.api.yml ps
```

## Configurações de Segurança

### 1. Firewall
```bash
# Configure o firewall UFW
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22      # SSH
sudo ufw allow 80      # HTTP
sudo ufw allow 443     # HTTPS
sudo ufw enable
```

### 2. SSL/TLS
O Caddy configurará automaticamente certificados Let's Encrypt para seu domínio.

### 3. Senhas Seguras
```bash
# Gere senhas seguras para o banco
openssl rand -base64 32

# Gere chave JWT segura
openssl rand -hex 64
```

## Comandos de Gerenciamento

### Comandos do Makefile
```bash
# Iniciar serviços
make prod-api-up

# Parar serviços
make prod-api-down

# Ver logs
make prod-api-logs

# Reiniciar serviços
make prod-api-down && make prod-api-up
```

### Comandos Docker Compose Diretos
```bash
# Status dos containers
docker-compose -f docker-compose.api.yml ps

# Logs de um serviço específico
docker-compose -f docker-compose.api.yml logs -f api-server

# Restart de um serviço específico
docker-compose -f docker-compose.api.yml restart api-server

# Executar comandos dentro do container
docker-compose -f docker-compose.api.yml exec api-server sh
```

## Monitoramento e Logs

### 1. Health Checks
```bash
# API Health
curl https://seu-dominio.com/health

# Readiness Check
curl https://seu-dominio.com/ready

# Verificar via Docker
docker-compose -f docker-compose.api.yml exec api-server wget -q -O- http://localhost:8080/health
```

### 2. Logs Estruturados
```bash
# Logs da API
tail -f logs/api/app.log

# Logs do Caddy
tail -f logs/caddy/access.log

# Logs do PostgreSQL
docker-compose -f docker-compose.api.yml logs postgres

# Logs do Redis
docker-compose -f docker-compose.api.yml logs redis
```

## Backup e Manutenção

### 1. Backup do Banco de Dados
```bash
# Criar script de backup
cat > scripts/backup-db.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/backup/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

docker-compose -f docker-compose.api.yml exec postgres pg_dump -U englog englog > $BACKUP_DIR/backup_$DATE.sql

# Manter apenas os últimos 7 backups
find $BACKUP_DIR -name "backup_*.sql" -mtime +7 -delete
EOF

chmod +x scripts/backup-db.sh

# Configurar cron para backup automático (diário às 2h)
(crontab -l 2>/dev/null; echo "0 2 * * * /caminho/para/englog/scripts/backup-db.sh") | crontab -
```

### 2. Atualizações
```bash
# Fazer backup antes de atualizar
./scripts/backup-db.sh

# Pull das últimas alterações
git pull origin main

# Rebuild e restart
docker-compose -f docker-compose.api.yml build --no-cache
make prod-api-down
make prod-api-up
```

## Solução de Problemas

### 1. Problemas Comuns
```bash
# Container não inicia
docker-compose -f docker-compose.api.yml logs api-server

# Problemas de conexão com banco
docker-compose -f docker-compose.api.yml exec postgres psql -U englog -d englog -c "\l"

# Problemas com SSL
docker-compose -f docker-compose.api.yml logs caddy
```

### 2. Reset Completo (Use com Cuidado)
```bash
# Parar tudo e limpar volumes
make prod-api-down
docker-compose -f docker-compose.api.yml down -v
docker system prune -a

# Restart completo
make prod-api-up
```

## Métricas e Performance

### 1. Endpoints de Métricas
- `GET /health` - Status geral da aplicação
- `GET /ready` - Readiness para receber tráfego
- `GET /metrics` - Métricas Prometheus (se habilitado)

### 2. Recursos Recomendados
- **CPU**: 2-4 cores para aplicação média
- **RAM**: 4-8GB para aplicação média
- **Storage**: SSD com pelo menos 50GB
- **Network**: Bandwidth adequado para seu tráfego esperado

## Conclusão

Este guia fornece uma base sólida para deploy da API englog em produção. O projeto já vem com:

✅ **Dockerfile otimizado** com multi-stage build
✅ **Health checks** configurados
✅ **SSL automático** via Caddy
✅ **Logs estruturados**
✅ **Scripts de automação**
✅ **Configuração de segurança básica**

Para ambientes de alta disponibilidade, considere implementar:
- Load balancer
- Cluster de banco de dados
- Monitoring avançado (Prometheus + Grafana)
- CI/CD automatizado

---

**Lembre-se**: Sempre teste em ambiente de desenvolvimento antes de aplicar em produção!
