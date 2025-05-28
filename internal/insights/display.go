package insights

import (
	"fmt"
	"strings"

	"github-octokit-poc/extractor"
)

// ShowDetailedInsights exibe insights específicos e detalhados sobre o repositório
func ShowDetailedInsights(data *extractor.RepositoryData) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 INSIGHTS ESPECÍFICOS")
	fmt.Println(strings.Repeat("=", 80))

	showTopContributors(data)
	showDominantLanguage(data)
	showLatestRelease(data)
	showRecentActivity(data)
	showCurrentActivity(data)
	showTopics(data)
	showRepositorySettings(data)

	fmt.Println("\n" + strings.Repeat("=", 80))
}

// showTopContributors mostra os principais colaboradores
func showTopContributors(data *extractor.RepositoryData) {
	if len(data.Contributors) == 0 {
		return
	}

	fmt.Println("\n👑 TOP 5 COLABORADORES:")
	limit := 5
	if len(data.Contributors) < limit {
		limit = len(data.Contributors)
	}

	for i := 0; i < limit; i++ {
		contrib := data.Contributors[i]
		fmt.Printf("   %d. %s (%d contribuições)\n", i+1, contrib.Login, contrib.Contributions)
	}
}

// showDominantLanguage mostra a linguagem dominante do repositório
func showDominantLanguage(data *extractor.RepositoryData) {
	if len(data.Languages) == 0 {
		return
	}

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

	if total > 0 {
		percentage := float64(maxBytes) / float64(total) * 100
		fmt.Printf("\n💻 LINGUAGEM DOMINANTE: %s (%.1f%%)\n", maxLang, percentage)
	}
}

// showLatestRelease mostra informações sobre o release mais recente
func showLatestRelease(data *extractor.RepositoryData) {
	if len(data.Releases) == 0 {
		return
	}

	latest := data.Releases[0]
	fmt.Printf("\n🚀 RELEASE MAIS RECENTE: %s\n", latest.TagName)
	fmt.Printf("   📅 Publicado em: %s\n", latest.PublishedAt.Format("02/01/2006"))
	fmt.Printf("   👤 Por: %s\n", latest.Author)
}

// showRecentActivity mostra atividade recente de commits
func showRecentActivity(data *extractor.RepositoryData) {
	if len(data.RecentCommits) == 0 {
		return
	}

	lastCommit := data.RecentCommits[0]
	fmt.Printf("\n📝 ÚLTIMO COMMIT: %s\n", lastCommit.SHA[:8])
	fmt.Printf("   📅 Em: %s\n", lastCommit.CreatedAt.Format("02/01/2006 15:04"))
	fmt.Printf("   👤 Por: %s\n", lastCommit.Author)
	
	// Truncar mensagem se for muito longa
	message := lastCommit.Message
	if len(message) > 100 {
		message = message[:100] + "..."
	}
	fmt.Printf("   💬 Mensagem: %s\n", message)
}

// showCurrentActivity mostra estatísticas de atividade atual
func showCurrentActivity(data *extractor.RepositoryData) {
	openIssues := countOpenItems(data.RecentIssues)
	openPRs := countOpenPRs(data.RecentPRs)

	fmt.Printf("\n📊 ATIVIDADE ATUAL:\n")
	fmt.Printf("   🎯 Issues abertas (amostra): %d\n", openIssues)
	fmt.Printf("   🔄 PRs abertos (amostra): %d\n", openPRs)
}

// showTopics mostra os tópicos/tags do repositório
func showTopics(data *extractor.RepositoryData) {
	if len(data.Topics) > 0 {
		fmt.Printf("\n🏷️  TÓPICOS: %v\n", data.Topics)
	}
}

// showRepositorySettings mostra configurações importantes do repositório
func showRepositorySettings(data *extractor.RepositoryData) {
	if data.Settings == nil {
		return
	}

	fmt.Printf("\n⚙️  CONFIGURAÇÕES:\n")
	fmt.Printf("   📝 Wiki habilitada: %v\n", data.Settings.HasWiki)
	fmt.Printf("   🎯 Issues habilitadas: %v\n", data.Settings.HasIssues)
	fmt.Printf("   🚀 Pages habilitadas: %v\n", data.Settings.HasPages)
	fmt.Printf("   💬 Discussions habilitadas: %v\n", data.Settings.HasDiscussions)
}

// countOpenItems conta quantas issues estão abertas
func countOpenItems(issues []*extractor.IssueData) int {
	count := 0
	for _, issue := range issues {
		if issue.State == "open" {
			count++
		}
	}
	return count
}

// countOpenPRs conta quantos PRs estão abertos
func countOpenPRs(prs []*extractor.PullRequestData) int {
	count := 0
	for _, pr := range prs {
		if pr.State == "open" {
			count++
		}
	}
	return count
}