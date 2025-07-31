# TLS Certificates

Esta pasta contém os certificados TLS para comunicação gRPC segura no EngLog.

## 🔐 Certificados de Desenvolvimento

Os certificados nesta pasta são **auto-assinados** e destinados apenas para **desenvolvimento local**.

### Arquivos:
- `server.crt` - Certificado público do servidor
- `server.key` - Chave privada do servidor (**não deve ser versionada**)

## 🚀 Gerando Certificados

### Método 1: Script Automático
```bash
./scripts/generate-certs.sh
```

### Método 2: Manual
```bash
openssl req -x509 -newkey rsa:4096 \
    -keyout certs/server.key \
    -out certs/server.crt \
    -days 365 \
    -nodes \
    -subj "/C=BR/ST=SP/L=São Paulo/O=EngLog/OU=Development/CN=localhost/emailAddress=dev@englog.local"
```

## ⚙️ Configuração

### Variáveis de Ambiente
```bash
# Habilitar TLS
TLS_ENABLED=true

# Caminhos dos certificados (padrão)
TLS_CERT_FILE=./certs/server.crt
TLS_KEY_FILE=./certs/server.key

# Nome do servidor para validação
GRPC_SERVER_NAME=localhost
```

### Desabilitar TLS (desenvolvimento)
```bash
TLS_ENABLED=false
```

## 🏭 Produção

⚠️ **IMPORTANTE**: Para produção, use certificados de uma Autoridade Certificadora (CA) confiável:

1. **Let's Encrypt** (gratuito)
2. **Certificados corporativos**
3. **Certificados de CA comercial**

### Exemplo com Let's Encrypt:
```bash
# Obter certificado
certbot certonly --standalone -d seu-dominio.com

# Configurar paths
TLS_CERT_FILE=/etc/letsencrypt/live/seu-dominio.com/fullchain.pem
TLS_KEY_FILE=/etc/letsencrypt/live/seu-dominio.com/privkey.pem
```

## 🛡️ Segurança

- ✅ Chaves privadas estão no `.gitignore`
- ✅ Permissões adequadas (600 para .key, 644 para .crt)
- ✅ Certificados auto-assinados apenas para desenvolvimento
- ⚠️ Renovar certificados antes do vencimento (365 dias)

## 🔍 Verificar Certificado

```bash
# Detalhes do certificado
openssl x509 -in certs/server.crt -text -noout

# Datas de validade
openssl x509 -in certs/server.crt -dates -noout

# Verificar chave privada
openssl rsa -in certs/server.key -check
```

## 🚨 Problemas Comuns

### "certificate signed by unknown authority"
- Esperado para certificados auto-assinados
- Cliente deve aceitar certificados inseguros em desenvolvimento
- Use `TLS_ENABLED=false` se necessário

### "x509: cannot validate certificate"
- Verifique se `GRPC_SERVER_NAME` está correto
- Use `localhost` para desenvolvimento local

### Permissões negadas
```bash
chmod 600 certs/server.key
chmod 644 certs/server.crt
```
