# Projeto de Monitoramento Completo

Serviço HTTP containerizado em Golang com monitoramento e observabilidade via Prometheus e Grafana, orquestrado por Docker Compose e provisionado automaticamente com Ansible, atingindo mais de 100000 requisições por segundo.

## Performance

```
/FullMonitoringProject$ wrk -t16 -c550 -d60s --latency http://127.0.0.1:80/projeto
Running 1m test @ http://127.0.0.1:80/projeto
  16 threads and 550 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     5.31ms    4.26ms  56.69ms   78.06%
    Req/Sec     7.18k   386.90    18.11k    82.27%
  Latency Distribution
     50%    4.26ms
     75%    7.08ms
     90%   10.85ms
     99%   20.14ms
  6867333 requests in 1.00m, 1.35GB read
Requests/sec: 114267.86
Transfer/sec:     22.99MB
```


## Visão Geral

O Projeto implementa um microsserviço que demonstra boas práticas de:
- **Desenvolvimento**: aplicação HTTP em Go com métricas expostas
- **Containerização**: Docker com imagens otimizadas e Docker Compose para orquestração
- **Redes**: configuração de rede bridge para comunicação inter-container
- **Proxy Reverso**: NGINX como reverse proxy com keepalive e otimizações
- **Observabilidade**: Prometheus para coleta de métricas e Grafana para visualização
- **Automação**: Ansible para provisionamento completo do ambiente

## Arquitetura
```
                          ┌─────────────────────────────────────────────────────────┐
                          │                   Host Linux                            │
                          │                                                         │
  ┌──────────┐            │  ┌─────────────────────────────────────────────────┐    │
  │  Cliente │            │  │              Docker Network (bridge)            │    │
  │  (curl / │            │  │                                                 │    │
  │ browser) │            │  │  ┌─────────────┐       ┌──────────────────────┐ │    │
  └────┬─────┘            │  │  │    NGINX    │       │ http-server-projeto  │ │    │
       │                  │  │  │             │──────▶│                      │ │    │
       │  :80             │  │  │  Proxy      │       │                      │ │    │
       └──────────────────┼──┼─▶│  Reverso    │       │  Go HTTP Server      │ │    │
                          │  │  │             │       │  porta 8080          │ │    │
                          │  │  └─────────────┘       │                      │ │    │
                          │  │                        │  GET /projeto        │ │    │
                          │  │                        │  GET /metrics        │ │    │
                          │  │                        └──────────┬───────────┘ │    │
                          │  │                                   │             │    │
                          │  │                        ┌──────────▼───────────┐ │    │
                          │  │                        │      Prometheus      │ │    │
                          │  │                        │                      │ │    │
                          │  │   :3000                │  Coleta métricas     │ │    │
       ┌──────────────────┼──┼─▶ Grafana ◀─────────── │  a cada 15s          │ │    │
       │  Dashboard       │  │  │        │            │  porta 9090          │ │    │
       │                  │  │  └────────┘            └──────────────────────┘ │    │
       │                  │  └─────────────────────────────────────────────────┘    │
       │                  │                                                         │
       │                  │  ┌────────────────────────────────────────────────────┐ │
       │                  │  │                  Ansible Playbook                  │ │
       │                  │  │  ✔ Instala Docker        ✔ Configura NGINX         │ │
       │                  │  │  ✔ Cria rede bridge      ✔ Sobe containers         │ │
       │                  │  │  ✔ Build da imagem        ✔ Valida o serviço       │ │
       │                  │  └────────────────────────────────────────────────────┘ │
       └──────────────────────────────────────────────────────────────────────────▶ │
                          └─────────────────────────────────────────────────────────┘
```


### Fluxo de Requisição

1. Cliente faz requisição para `http://localhost:80/projeto`
2. NGINX (proxy reverso) recebe a requisição
3. NGINX encaminha para `http-server-projeto:8080/projeto`
4. Serviço Go processa e retorna JSON com timestamp UTC
5. Serviço registra métrica da requisição
6. Prometheus coleta métricas a cada 10 segundos
7. Grafana consulta Prometheus e exibe no dashboard

## Como Usar

### Pré-requisitos

- **Docker**
- **Docker Compose**
- **Ansible**
- **Git** (para clonar o repositório)
- Sistema operacional: **Linux** (Debian)

### Opção 1: Docker Compose (Manual)

1. Clone o repositório:
```bash
git clone https://github.com/Gansblaidx/FullMonitoringProject.git
cd FullMonitoringProject
```

2. Inicie os containers:
```bash
docker compose up -d --build
```

3. Aguarde os serviços iniciarem e teste:
```bash
curl http://localhost:80/projeto
```

4. Acesse os serviços:
   - **API do Serviço**: http://localhost:80/projeto
   - **Health Check**: http://localhost/health
   - **Métricas Prometheus**: http://localhost:9090
   - **Grafana Dashboard**: http://localhost:3000 (admin/admin)

### Opção 2: Ansible (Automação Completa)

1. Clone o repositório:
```bash
git clone https://github.com/Gansblaidx/FullMonitoringProject.git
cd FullMonitoringProject
```

2. Execute o playbook Ansible:
```bash
ansible-playbook ansible/playbook.yml
```

3. O playbook irá:
   - Instalar Docker (se necessário)
   - Criar rede bridge `project-network`
   - Build e deploy de todos os containers
   - Validar o funcionamento do serviço
   - Exibir informações de acesso

## Métricas e Observabilidade

### Métricas Expostas

O serviço expõe as seguintes métricas via Prometheus:

| Métrica | Descrição |
|---------|-----------|
| `http_requests_total` | Total de requisições no endpoint `/projeto` / Usado para medir req/s |
| `service_up` | Indicador de disponibilidade do serviço (UP, DOWN) |
| `request` | Métrica auto-gerada pelo Prometheus (status do target) |

### Dashboard Grafana

O dashboard padrão inclui:

**Seção: Go Service**
- Status de disponibilidade do serviço
- Total de requisições # Como gráfico e contador
- Taxa de requisições por segundo

### Acessar Métricas

- **Raw Prometheus**: http://localhost:9090
  - Consultar métrica: `/api/v1/query?query=http_requests_total`

- **Grafana**: http://localhost:3000
  - Credenciais padrão: `admin` / `admin`
  - Dashboard provisionado automaticamente

## Endpoints da Aplicação

### GET `/projeto`

Retorna informações do projeto com timestamp UTC.

**Exemplo de requisição:**
```bash
curl http://localhost:80/projeto
```

**Resposta:**
```json
{
  "nome": "Projeto de Monitoramento",
  "horario": "2026-06-04T11:20:15Z"
}
```

### GET `/health`

Verificação de saúde da aplicação.

**Exemplo de requisição:**
```bash
curl http://localhost/health
```

**Resposta:**
```json
{
  "status": "healthy"
}
```

### GET `/metrics`

Métricas no formato Prometheus (apenas acesso interno).

**Exemplo de requisição:**
```bash
curl http://localhost/metrics
```

## Docker Compose Services

### http-server-projeto-korp
- **Imagem**: Build local da aplicação Go
- **Porta**: 8080 (interna, não exposta)
- **Rede**: project-network
- **Reinicialização**: unless-stopped

### nginx
- **Imagem**: nginx:latest
- **Portas**: 80:80 (host:container)
- **Volumes**: Configurações de proxy reverso
- **Rede**: project-network
- **Dependências**: http-server-projeto
- **Reinicialização**: unless-stopped

### prometheus
- **Imagem**: prom/prometheus:latest
- **Portas**: 9090:9090 (host:container)
- **Volumes**: Configuração de scrape
- **Rede**: project-network
- **Intervalo de coleta**: 10 segundos
- **Reinicialização**: unless-stopped

### grafana
- **Imagem**: grafana/grafana:latest
- **Portas**: 3000:3000 (host:container)
- **Volumes**: Dashboards e datasources provisionados
- **Rede**: project-network
- **Senha admin**: admin (configurável via variável de ambiente)
- **Reinicialização**: unless-stopped

## Configuração Detalhada

### Dockerfile (Go)

Usa build multi-stage para otimizar o tamanho da imagem:
- **Stage 1**: Compila a aplicação com `golang:1.22-alpine`
- **Stage 2**: Runtime mínimo com `alpine:3.19`

Resultado: imagem pequena e segura (~20MB)

### NGINX

**Upstream:**
- Configuração de keepalive para melhor performance
- Timeout de 60 segundos

**Location blocks:**
- `/` → proxy para Go service na porta 8080
- `/stub_status` → status do NGINX (acesso restrito à rede local)

### Prometheus

**Scrape Config:**
- Coleta métricas do job `http-server-projeto`
- Intervalo de 10 segundos
- Endpoint: `/metrics`

### Grafana

**Provisionamento Automático:**
- Datasource Prometheus configurado via `datasources.yml`
- Dashboard carregado via `dashboards.yml`
- Arquivos JSON do dashboard em `/var/lib/grafana/dashboards`

## Dependências

### Go
- `github.com/prometheus/client_golang` v1.19.1

### Docker
- golang:1.22-alpine (build)
- alpine:3.19 (runtime)
- nginx:latest
- prom/prometheus:latest
- grafana/grafana:latest
