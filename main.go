package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github-octokit-poc/extractor"
	"github-octokit-poc/github"
	"github-octokit-poc/utils"
)

func main() {
	log.Println("🚀 Iniciando POC do GitHub Octokit em Go")
	
	// Criar cliente GitHub
	client, err := github.NewClient()
	if err != nil {
		log.Fatalf("❌ Erro ao criar cliente GitHub: %v", err)
	}

	log.Println("✅ Cliente GitHub configurado com sucesso")

	// Definir repositório alvo
	owner := "kubernetes"
	repo := "kubernetes"

	log.Printf("🎯 Alvo: %s/%s", owner, repo)

	// Extrair todos os dados do repositório
	data, err := extractor.ExtractRepositoryData(client, owner, repo)
	if err != nil {
		log.Fatalf("❌ Erro na extração: %v", err)
	}

	// Mostrar resumo na tela
	data.PrintSummary()

	// Gerar relatório detalhado
	report := utils.GenerateReport(data)
	fmt.Println("\n" + report)

	// Criar pasta output se não existir
	outputBaseDir := "output"
	if err := os.MkdirAll(outputBaseDir, 0755); err != nil {
		log.Printf("⚠️ Erro ao criar diretório output: %v", err)
	}

	// Criar pasta específica para esta execução baseada na data
	timestamp := time.Now().Format("20060102_150405")
	outputDir := filepath.Join(outputBaseDir, timestamp)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Printf("⚠️ Erro ao criar diretório específico: %v", err)
		// Fallback para a pasta output se não conseguir criar a subpasta
		outputDir = outputBaseDir
	}

	// Definir caminhos dos arquivos
	jsonFilename := filepath.Join(outputDir, fmt.Sprintf("%s_%s_data.json", owner, repo))
	reportFilename := filepath.Join(outputDir, fmt.Sprintf("%s_%s_report.txt", owner, repo))
	
	log.Printf("\n💾 Salvando arquivos:")
	log.Printf("   📊 Dados completos: %s", jsonFilename)
	log.Printf("   📋 Relatório: %s", reportFilename)
	
	if err := data.SaveToJSON(jsonFilename); err != nil {
		log.Printf("⚠️ Erro ao salvar JSON: %v", err)
	} else {
		log.Printf("✅ JSON salvo com sucesso!")
	}

	// Salvar relatório em arquivo texto
	if err := os.WriteFile(reportFilename, []byte(report), 0644); err != nil {
		log.Printf("⚠️ Erro ao salvar relatório: %v", err)
	} else {
		log.Printf("✅ Relatório salvo com sucesso!")
	}

	// Mostrar alguns insights específicos
	showInsights(data)
}

func showInsights(data *extractor.RepositoryData) {
    fmt.Println(strings.Repeat("=", 80))
    fmt.Println("🔍 INSIGHTS ESPECÍFICOS")
    fmt.Println(strings.Repeat("=", 80))

	// Top colaboradores
	fmt.Println("\n👑 TOP 5 COLABORADORES:")
	for i, contrib := range data.Contributors {
		if i >= 5 {
			break
		}
		fmt.Printf("   %d. %s (%d contribuições)\n", i+1, contrib.Login, contrib.Contributions)
	}

	// Linguagem dominante
	if len(data.Languages) > 0 {
		maxLang := ""
		maxBytes := 0
		total := 0
		for lang, bytes := range data.Languages {
			total += bytes
			if bytes > maxBytes {
				maxBytes = bytes
				maxLang = lang
			}
		}
		percentage := float64(maxBytes) / float64(total) * 100
		fmt.Printf("\n💻 LINGUAGEM DOMINANTE: %s (%.1f%%)\n", maxLang, percentage)
	}

	// Release mais recente
	if len(data.Releases) > 0 {
		latest := data.Releases[0]
		fmt.Printf("\n🚀 RELEASE MAIS RECENTE: %s\n", latest.TagName)
		fmt.Printf("   📅 Publicado em: %s\n", latest.PublishedAt.Format("02/01/2006"))
		fmt.Printf("   👤 Por: %s\n", latest.Author)
	}

	// Atividade recente
	if len(data.RecentCommits) > 0 {
		lastCommit := data.RecentCommits[0]
		fmt.Printf("\n📝 ÚLTIMO COMMIT: %s\n", lastCommit.SHA[:8])
		fmt.Printf("   📅 Em: %s\n", lastCommit.CreatedAt.Format("02/01/2006 15:04"))
		fmt.Printf("   👤 Por: %s\n", lastCommit.Author)
		fmt.Printf("   💬 Mensagem: %.100s...\n", lastCommit.Message)
	}

	// Issues vs PRs
	openIssues := 0
	openPRs := 0
	for _, issue := range data.RecentIssues {
		if issue.State == "open" {
			openIssues++
		}
	}
	for _, pr := range data.RecentPRs {
		if pr.State == "open" {
			openPRs++
		}
	}

	fmt.Printf("\n📊 ATIVIDADE ATUAL:\n")
	fmt.Printf("   🎯 Issues abertas (amostra): %d\n", openIssues)
	fmt.Printf("   🔄 PRs abertos (amostra): %d\n", openPRs)

	// Tópicos/Tags
	if len(data.Topics) > 0 {
		fmt.Printf("\n🏷️  TÓPICOS: %v\n", data.Topics)
	}

	// Configurações interessantes
	fmt.Printf("\n⚙️  CONFIGURAÇÕES:\n")
	fmt.Printf("   📝 Wiki habilitada: %v\n", data.Settings.HasWiki)
	fmt.Printf("   🎯 Issues habilitadas: %v\n", data.Settings.HasIssues)
	fmt.Printf("   🚀 Pages habilitadas: %v\n", data.Settings.HasPages)
	fmt.Printf("   💬 Discussions habilitadas: %v\n", data.Settings.HasDiscussions)

}