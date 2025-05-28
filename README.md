# 🚀 GitHub Repository Analyzer

Ferramenta completa para análise de repositórios GitHub desenvolvida em Go. Extrai informações detalhadas, estatísticas e gera relatórios abrangentes sobre qualquer repositório público.

## ✨ Funcionalidades

- 🔍 **Análise completa** de repositórios GitHub
- 📊 **Estatísticas detalhadas** (stars, forks, issues, PRs)
- 👥 **Análise de colaboradores** e contribuições
- 💻 **Distribuição de linguagens** de programação
- 🏥 **Score de saúde** do repositório
- 📈 **Métricas de atividade** (commits, issues, PRs)
- 📋 **Relatórios** em JSON e texto formatado
- 🎯 **Interface CLI** intuitiva
- ⚙️ **Configuração flexível** via .env

## 🛠️ Instalação

### Pré-requisitos

- **Go 1.21+** instalado
- **Token de acesso do GitHub** ([como obter](https://github.com/settings/tokens))

### Setup rápido

```bash
# 1. Clone o projeto
git clone <seu-repositorio>
cd github-octokit-poc

# 2. Configure o token
cp .env.example .env
# Edite o .env e adicione seu GITHUB_TOKEN

# 3. Baixe dependências
go mod tidy

# 4. Execute
go run main.go https://github.com/kubernetes/kubernetes
```

## 🎮 Como usar

### Formato básico
```bash
go run main.go [opções] [url-do-repositório]
```

### 📝 Exemplos práticos

**Analisar repositório por URL completa:**
```bash
go run main.go https://github.com/kubernetes/kubernetes
```

**Formato owner/repo:**
```bash
go run main.go kubernetes/kubernetes
```

**URL SSH também funciona:**
```bash
go run main.go git@github.com:facebook/react.git
```

**Usando flags específicas:**
```bash
go run main.go -o microsoft -r vscode
```

**Diretório de saída customizado:**
```bash
go run main.go --output /tmp/analise https://github.com/golang/go
```

**Ver ajuda:**
```bash
go run main.go --help
```

### 🎛️ Opções disponíveis

| Flag | Descrição | Exemplo |
|------|-----------|---------|
| `-u, --url` | URL do repositório | `--url https://github.com/owner/repo` |
| `-o, --owner` | Proprietário do repositório | `--owner kubernetes` |
| `-r, --repo` | Nome do repositório | `--repo kubernetes` |
| `--output` | Diretório de saída | `--output /tmp/results` |
| `-h, --help` | Mostrar ajuda | `--help` |
| `-v, --version` | Mostrar versão | `--version` |

### 🌐 Formatos de URL suportados

✅ `https://github.com/owner/repo`  
✅ `https://github.com/owner/repo.git`  
✅ `git@github.com:owner/repo.git`  
✅ `owner/repo`  

## 📁 Estrutura do projeto

```
github-octokit-poc/
├── main.go                    # 🎯 Ponto de entrada
├── cmd/
│   └── runner.go             # 🎬 Orquestrador principal
├── internal/
│   ├── cli/
│   │   └── parser.go         # 🎛️ Parser de argumentos CLI
│   ├── config/
│   │   └── config.go         # ⚙️ Gerenciamento de configurações
│   ├── output/
│   │   └── handler.go        # 💾 Gerenciamento de arquivos
│   └── insights/
│       └── display.go        # 🔍 Exibição de insights
├── extractor/
│   └── repository.go         # 📥 Extração de dados GitHub
├── github/
│   └── clients.go            # 🐙 Cliente GitHub
├── utils/
│   └── analyzer.go           # 🧮 Análises e relatórios
└── README.md                 # 📖 Documentação
```

## ⚙️ Configuração

### Arquivo .env

```bash
# Obrigatório
GITHUB_TOKEN=ghp_seu_token_aqui

# Opcionais
GITHUB_DEFAULT_USER=kubernetes
GITHUB_DEFAULT_REPO=kubernetes
OUTPUT_DIR=output
DEBUG=false
```

### Variáveis de ambiente

| Variável | Obrigatória | Descrição |
|----------|-------------|-----------|
| `GITHUB_TOKEN` | ✅ | Token de acesso do GitHub |
| `GITHUB_DEFAULT_USER` | ❌ | Usuário padrão |
| `GITHUB_DEFAULT_REPO` | ❌ | Repositório padrão |
| `GITHUB_API_BASE_URL` | ❌ | URL para GitHub Enterprise |
| `OUTPUT_DIR` | ❌ | Diretório de saída padrão |
| `DEBUG` | ❌ | Modo debug (true/false) |

## 📊 Exemplo de saída

```
🚀 Iniciando GitHub Repository Analyzer
✅ Cliente GitHub configurado com sucesso
🎯 Alvo (via CLI): kubernetes/kubernetes
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

👥 COLABORADORES: 3847 encontrados
🎯 ISSUES RECENTES: 8 encontradas
🔄 PULL REQUESTS: 10 encontrados
🚀 RELEASES: 10 encontrados

💾 Salvando arquivos:
   📊 Dados completos: output/20250528_143045/kubernetes_kubernetes_data.json
   📋 Relatório: output/20250528_143045/kubernetes_kubernetes_report.txt
✅ JSON salvo com sucesso!
✅ Relatório salvo com sucesso!

================================================================================
🔍 INSIGHTS ESPECÍFICOS
================================================================================

👑 TOP 5 COLABORADORES:
   1. k8s-ci-robot (12,567 contribuições)
   2. liggitt (4,234 contribuições)
   3. brendandburns (3,891 contribuições)
   4. smarterclayton (3,456 contribuições)
   5. wojtek-t (2,987 contribuições)

💻 LINGUAGEM DOMINANTE: Go (96.2%)

🚀 RELEASE MAIS RECENTE: v1.30.1
   📅 Publicado em: 15/05/2025
   👤 Por: k8s-release-robot
```

## 🔧 Build para produção

```bash
# Build do binário
go build -o github-analyzer main.go

# Executar o binário
./github-analyzer https://github.com/kubernetes/kubernetes

# Cross-compilation (exemplo para Windows)
GOOS=windows GOARCH=amd64 go build -o github-analyzer.exe main.go
```

## 🚀 Próximos passos

- [ ] 📊 Dashboard web interativo
- [ ] 🔄 Comparação entre repositórios
- [ ] 📈 Análise histórica de crescimento
- [ ] 🔔 Sistema de notificações
- [ ] 🐳 Container Docker
- [ ] 📦 Packaging para diferentes OS

## 🛟 Troubleshooting

**Erro de autenticação:**
```bash
# Verificar se o token está correto
curl -H "Authorization: token SEU_TOKEN" https://api.github.com/user
```

**Rate limiting:**
```bash
# Verificar limits restantes
curl -H "Authorization: token SEU_TOKEN" https://api.github.com/rate_limit
```

**Dependências:**
```bash
# Limpar e reinstalar módulos
go clean -modcache
go mod tidy
```

## 📚 Documentação adicional

- [GitHub API Documentation](https://docs.github.com/en/rest)
- [go-github Library](https://pkg.go.dev/github.com/google/go-github/v57/github)
- [OAuth2 em Go](https://pkg.go.dev/golang.org/x/oauth2)

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

---

**Desenvolvido com ❤️ em Go**