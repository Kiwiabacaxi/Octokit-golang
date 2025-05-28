package github

import (
	"context"
	"log"
	"net/url"
	"os"

	"github.com/google/go-github/v57/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type Client struct {
	GitHub *github.Client
	Ctx    context.Context
}

type Config struct {
	Token      string
	BaseURL    string
	Debug      bool
}

// NewClient cria um novo cliente GitHub configurado
func NewClient() (*Client, error) {
	// Carrega variáveis do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	config := &Config{
		Token:   os.Getenv("GITHUB_TOKEN"),
		BaseURL: os.Getenv("GITHUB_API_BASE_URL"),
		Debug:   os.Getenv("DEBUG") == "true",
	}

	if config.Token == "" {
		log.Fatal("GITHUB_TOKEN é obrigatório")
	}

	ctx := context.Background()
	
	// Configuração OAuth2
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Cliente GitHub
	client := github.NewClient(tc)

	// URL base personalizada (GitHub Enterprise)
	if config.BaseURL != "" {
		client.BaseURL, err = url.Parse(config.BaseURL)
		if err != nil {
			return nil, err
		}
	}

	if config.Debug {
		log.Println("Cliente GitHub configurado com sucesso")
	}

	return &Client{
		GitHub: client,
		Ctx:    ctx,
	}, nil
}