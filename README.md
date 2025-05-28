# POC Octokit em Go

Esta é uma prova de conceito demonstrando como usar a biblioteca go-github (equivalente ao Octokit) para interagir com a API do GitHub em Go.

## Pré-requisitos

1. **Go 1.21+** instalado
2. **Token de acesso do GitHub** com as seguintes permissões:
   - `repo` (para acessar repositórios)
   - `user` (para informações do usuário)
   - `public_repo` (para repositórios públicos)

## Configuração

1. **Gere um token de acesso pessoal no GitHub:**
   - Vá para Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Clique em "Generate new token"
   - Selecione os escopos necessários (`repo`, `user`, `public_repo`)
   - Copie o token gerado

2. **Configure o arquivo .env:**
   ```bash
   # Copie o arquivo de exemplo
   cp .env.example .env
   
   # Edite o arquivo .env e adicione seu token
   # GITHUB_TOKEN=ghp_seu_token_aqui
   ```

   **Ou use variáveis de ambiente tradicionais:**
   ```bash
   export GITHUB_TOKEN="seu_token_aqui"
   ```

## Instalação e Execução

1. **Clone ou crie o projeto:**
   ```bash
   mkdir github-octokit-poc
   cd github-octokit-poc
   ```

2. **Copie os arquivos:** `main.go`, `go.mod`, `.env.example`, `.gitignore`

3. **Configure o ambiente:**
   ```bash
   # Copie e configure o .env
   cp .env.example .env
   # Edite o .env com seu token do GitHub
   ```

4. **Baixe as dependências:**
   ```bash
   go mod tidy
   ```

5. **Execute a POC:**
   ```bash
   go run main.go
   ```

## Estrutura do Projeto

```
github-octokit-poc/
├── main.go              # Ponto de entrada da aplicação
├── go.mod               # Dependências do Go
├── .env.example         # Template de configuração
├── .env                 # Configurações locais (não commitado)
├── .gitignore           # Arquivos ignorados pelo Git
├── github/
│   └── client.go        # Cliente GitHub configurável
├── extractor/
│   └── repository.go    # Extrator completo de dados
├── utils/
│   └── analyzer.go      # Analisador e gerador de relatórios
└── README.md            # Este arquivo
```

## Variáveis de Ambiente Disponíveis

| Variável | Obrigatória | Descrição |
|----------|-------------|-----------|
| `GITHUB_TOKEN` | ✅ Sim | Token de acesso pessoal do GitHub |
| `GITHUB_DEFAULT_USER` | ❌ Não | Usuário padrão para buscar repositórios |
| `GITHUB_DEFAULT_REPO` | ❌ Não | Repositório padrão para buscar eventos |
| `GITHUB_API_BASE_URL` | ❌ Não | URL base para GitHub Enterprise |
| `GITHUB_REQUEST_TIMEOUT` | ❌ Não | Timeout das requisições (segundos) |
| `GITHUB_MAX_PER_PAGE` | ❌ Não | Máximo de itens por página |
| `DEBUG` | ❌ Não | Modo debug (true/false) |

### ✅ Funcionalidades Implementadas

1. **Arquitetura modular** com separação de responsabilidades
2. **Extração completa de dados** do repositório kubernetes/kubernetes:
   - Informações básicas e metadados
   - Estatísticas (stars, forks, watchers, issues)
   - Configurações e permissões
   - Distribuição de linguagens de programação
   - Lista completa de colaboradores
   - Issues e Pull Requests recentes
   - Histórico de releases
   - Commits recentes
   - Eventos de atividade
   - Rate limits da API
3. **Análise inteligente** dos dados extraídos:
   - Score de saúde do repositório
   - Métricas de atividade (commits, issues, PRs)
   - Análise de colaboradores (top contributors, core team)
   - Detecção de issues obsoletas
   - Tendências de manutenção
4. **Relatórios detalhados** em múltiplos formatos:
   - Resumo na tela com insights
   - Arquivo JSON com dados completos
   - Relatório texto formatado
5. **Configuração via .env** para facilitar o desenvolvimento
6. **Suporte a GitHub Enterprise** (configurável via .env)

### 📋 Principais operações da API

**Módulo GitHub Client (`github/client.go`):**
- Configuração automática com OAuth2
- Suporte a GitHub Enterprise
- Gerenciamento de rate limits

**Módulo Extrator (`extractor/repository.go`):**
- `ExtractRepositoryData()` - Extração completa de dados
- `client.Repositories.Get()` - Informações básicas
- `client.Repositories.ListLanguages()` - Linguagens de programação
- `client.Repositories.ListContributors()` - Colaboradores
- `client.Issues.ListByRepo()` - Issues do repositório
- `client.PullRequests.List()` - Pull requests
- `client.Repositories.ListReleases()` - Releases
- `client.Repositories.ListCommits()` - Commits
- `client.Activity.ListRepositoryEvents()` - Eventos
- `client.RateLimits()` - Verificação de limites

**Módulo Analisador (`utils/analyzer.go`):**
- `AnalyzeLanguages()` - Distribuição de linguagens
- `AnalyzeActivity()` - Métricas de atividade
- `AnalyzeContributors()` - Análise de colaboradores
- `AnalyzeHealth()` - Score de saúde do repositório
- `GenerateReport()` - Relatório completo formatado

## Exemplo de saída

```
🔍 Iniciando extração completa do repositório kubernetes/kubernetes
📋 Extraindo informações básicas...
💻 Extraindo linguagens...
👥 Extraindo colaboradores...
🎯 Extraindo issues recentes...
🔄 Extraindo pull requests recentes...
🚀 Extraindo releases...
📝 Extraindo commits recentes...
⚡ Extraindo eventos recentes...
📊 Verificando rate limits...
✅ Extração concluída em 8.2s

================================================================================
📊 RESUMO DA EXTRAÇÃO - kubernetes/kubernetes
================================================================================
🏷️  Nome: kubernetes
👤 Proprietário: kubernetes
📝 Descrição: Production-Grade Container Scheduling and Management
🌐 URL: https://github.com/kubernetes/kubernetes
📅 Criado em: 07/06/2014
🔄 Última atualização: 28/05/2025 14:30

📈 ESTATÍSTICAS:
⭐ Stars: 110.2K
🍴 Forks: 39.5K
👀 Watchers: 3.8K
🎯 Issues abertas: 2.1K

💻 LINGUAGENS:
   Go: 96.2%
   Shell: 1.8%
   Python: 1.2%
   Makefile: 0.5%
   Dockerfile: 0.3%

👥 COLABORADORES: 3847 encontrados
🎯 ISSUES RECENTES: 10 encontradas
🔄 PULL REQUESTS: 10 encontrados
🚀 RELEASES: 10 encontrados
📝 COMMITS RECENTES: 10 encontrados
⚡ EVENTOS RECENTES: 10 encontrados

📊 RATE LIMITS:
   Core API: 4850/5000 (reset em 16:45:30)

⏱️  Extração concluída em: 8.234567s
================================================================================

📊 RELATÓRIO COMPLETO DE ANÁLISE
================================================================================

🏥 SAÚDE DO REPOSITÓRIO
----------------------------------------
Score de saúde: 95.0/100
Status: Excelente
Último commit: 0 dias atrás
Último release: 15 dias atrás
Issues obsoletas: 2
Ratio de issues abertas: 45.0%

💾 Salvando arquivos:
   📊 Dados completos: kubernetes_kubernetes_data_20250528_143045.json
   📋 Relatório: kubernetes_kubernetes_report_20250528_143045.txt
✅ JSON salvo com sucesso!
✅ Relatório salvo com sucesso!
```

## Próximos passos

Para expandir esta POC, você pode implementar:

- ✨ **Webhooks** para receber eventos do GitHub
- 📊 **Análise de dados** dos repositórios
- 🔄 **Sincronização** com banco de dados local
- 📈 **Relatórios** de atividade
- 🤖 **Automação** de tarefas no GitHub

## Documentação útil

- [go-github Documentation](https://pkg.go.dev/github.com/google/go-github/v57/github)
- [GitHub API Documentation](https://docs.github.com/en/rest)
- [OAuth2 em Go](https://pkg.go.dev/golang.org/x/oauth2)

## Troubleshooting

**Erro de autenticação:**
- Verifique se o token está correto no arquivo `.env`
- Confirme que a variável `GITHUB_TOKEN` está definida
- Teste o token no navegador: https://api.github.com/user (com header Authorization: token SEU_TOKEN)

**Arquivo .env não encontrado:**
- Copie o `.env.example` para `.env`
- O programa funcionará com variáveis de ambiente do sistema se o .env não existir

**Rate limiting:**
- A API do GitHub tem limites de taxa
- Use `client.RateLimits()` para verificar o status
- Considere implementar retry com backoff

**Dependências:**
- Execute `go mod tidy` se houver problemas com módulos
- Verifique se está usando Go 1.21 ou superior