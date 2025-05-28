package cmd

import (
	"fmt"
	"log"

	"github-octokit-poc/extractor"
	"github-octokit-poc/github"
	"github-octokit-poc/internal/config"
	"github-octokit-poc/internal/insights"
	"github-octokit-poc/internal/output"
	"github-octokit-poc/utils"
)

// Run é o ponto de entrada principal da aplicação
func Run() error {
	log.Println("🚀 Iniciando POC do GitHub Octokit em Go")

	// 1. Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// 2. Criar cliente GitHub
	client, err := github.NewClient()
	if err != nil {
		return err
	}
	log.Println("✅ Cliente GitHub configurado com sucesso")

	// 3. Definir repositório alvo
	owner, repo := cfg.GetTarget()
	log.Printf("🎯 Alvo: %s/%s", owner, repo)

	// 4. Extrair dados do repositório
	data, err := extractor.ExtractRepositoryData(client, owner, repo)
	if err != nil {
		return err
	}

	// 5. Exibir resumo
	data.PrintSummary()

	// 6. Gerar relatório detalhado
	report := utils.GenerateReport(data)
	fmt.Println("\n" + report)

	// 7. Salvar outputs
	outputHandler := output.NewHandler(owner, repo)
	if err := outputHandler.SaveAll(data, report); err != nil {
		log.Printf("⚠️ Erro ao salvar outputs: %v", err)
	}

	// 8. Mostrar insights específicos
	insights.ShowDetailedInsights(data)

	return nil
}