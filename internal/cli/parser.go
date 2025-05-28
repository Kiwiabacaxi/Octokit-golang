package cli

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Args representa os argumentos da linha de comando
type Args struct {
	RepoURL     string
	Owner       string
	Repo        string
	OutputDir   string
	ShowHelp    bool
	ShowVersion bool
}

// Parse analisa os argumentos da linha de comando
func Parse() (*Args, error) {
	args := &Args{}

	// Definir flags
	flag.StringVar(&args.RepoURL, "url", "", "URL do repositório GitHub (ex: https://github.com/owner/repo)")
	flag.StringVar(&args.RepoURL, "u", "", "URL do repositório GitHub (formato curto)")
	flag.StringVar(&args.Owner, "owner", "", "Proprietário do repositório")
	flag.StringVar(&args.Owner, "o", "", "Proprietário do repositório (formato curto)")
	flag.StringVar(&args.Repo, "repo", "", "Nome do repositório")
	flag.StringVar(&args.Repo, "r", "", "Nome do repositório (formato curto)")
	flag.StringVar(&args.OutputDir, "output", "output", "Diretório de saída")
	flag.BoolVar(&args.ShowHelp, "help", false, "Mostrar ajuda")
	flag.BoolVar(&args.ShowHelp, "h", false, "Mostrar ajuda (formato curto)")
	flag.BoolVar(&args.ShowVersion, "version", false, "Mostrar versão")
	flag.BoolVar(&args.ShowVersion, "v", false, "Mostrar versão (formato curto)")

	// Personalizar usage
	flag.Usage = func() {
		showUsage()
	}

	// Parse dos argumentos
	flag.Parse()

	// Verificar se precisa mostrar help ou version
	if args.ShowHelp {
		showUsage()
		os.Exit(0)
	}

	if args.ShowVersion {
		showVersion()
		os.Exit(0)
	}

	// Se não há argumentos, verificar se tem argumentos posicionais
	if args.RepoURL == "" && args.Owner == "" && args.Repo == "" {
		positionalArgs := flag.Args()
		if len(positionalArgs) > 0 {
			args.RepoURL = positionalArgs[0]
		}
	}

	// Parse da URL se fornecida
	if args.RepoURL != "" {
		owner, repo, err := parseGitHubURL(args.RepoURL)
		if err != nil {
			return nil, err
		}
		args.Owner = owner
		args.Repo = repo
	}

	return args, nil
}

// parseGitHubURL extrai owner e repo de uma URL do GitHub
func parseGitHubURL(url string) (owner, repo string, err error) {
	// Limpar a URL
	url = strings.TrimSpace(url)
	
	// Padrões de URL do GitHub que suportamos
	patterns := []string{
		// HTTPS URLs
		`^https?://github\.com/([^/]+)/([^/]+?)(?:\.git)?/?$`,
		`^https?://github\.com/([^/]+)/([^/]+)/.*$`,
		// SSH URLs
		`^git@github\.com:([^/]+)/([^/]+?)(?:\.git)?$`,
		// Formato simplificado: owner/repo
		`^([^/]+)/([^/]+)$`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		
		if len(matches) >= 3 {
			owner = matches[1]
			repo = matches[2]
			
			// Remover extensão .git se presente
			repo = strings.TrimSuffix(repo, ".git")
			
			return owner, repo, nil
		}
	}

	return "", "", fmt.Errorf("formato de URL inválido: %s\n\nFormatos suportados:\n  - https://github.com/owner/repo\n  - git@github.com:owner/repo.git\n  - owner/repo", url)
}

// IsEmpty verifica se os argumentos estão vazios
func (a *Args) IsEmpty() bool {
	return a.Owner == "" || a.Repo == ""
}

// GetTarget retorna owner e repo
func (a *Args) GetTarget() (string, string) {
	return a.Owner, a.Repo
}

// showUsage exibe a ajuda do comando
func showUsage() {
	fmt.Printf(`🚀 GitHub Repository Analyzer

DESCRIÇÃO:
    Ferramenta para análise completa de repositórios GitHub.
    Extrai informações detalhadas, estatísticas e gera relatórios.

USO:
    %s [opções] [url-do-repositório]

ARGUMENTOS:
    url-do-repositório    URL do repositório GitHub a ser analisado

OPÇÕES:
    -u, --url string     URL do repositório GitHub
                         (ex: https://github.com/kubernetes/kubernetes)
    
    -o, --owner string   Proprietário do repositório
    -r, --repo string    Nome do repositório
    
    --output string      Diretório de saída (padrão: "output")
    
    -h, --help          Mostrar esta ajuda
    -v, --version       Mostrar versão

EXEMPLOS:
    # Analisar repositório pelo URL completo
    %s https://github.com/kubernetes/kubernetes
    
    # Analisar repositório pelo formato owner/repo
    %s kubernetes/kubernetes
    
    # Analisar usando flags separadas
    %s -o kubernetes -r kubernetes
    
    # Analisar com diretório de saída customizado
    %s -u https://github.com/kubernetes/kubernetes --output /tmp/analise
    
    # URL SSH também funciona
    %s git@github.com:kubernetes/kubernetes.git

FORMATOS DE URL SUPORTADOS:
    ✅ https://github.com/owner/repo
    ✅ https://github.com/owner/repo.git
    ✅ git@github.com:owner/repo.git
    ✅ owner/repo

CONFIGURAÇÃO:
    As configurações podem ser definidas via arquivo .env:
    
    GITHUB_TOKEN=ghp_seu_token_aqui
    GITHUB_DEFAULT_USER=owner_padrao
    GITHUB_DEFAULT_REPO=repo_padrao

Para mais informações, visite: https://github.com/seu-usuario/github-octokit-poc
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

// showVersion exibe a versão
func showVersion() {
	fmt.Println("GitHub Repository Analyzer v1.0.0")
	fmt.Println("Desenvolvido com Go e github.com/google/go-github")
}