package utils

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github-octokit-poc/extractor"
)

// LanguageStats representa estatísticas de uma linguagem
type LanguageStats struct {
	Name       string  `json:"name"`
	Bytes      int     `json:"bytes"`
	Percentage float64 `json:"percentage"`
}

// ActivityMetrics representa métricas de atividade
type ActivityMetrics struct {
	CommitsLastWeek   int     `json:"commits_last_week"`
	CommitsLastMonth  int     `json:"commits_last_month"`
	IssuesLastWeek    int     `json:"issues_last_week"`
	IssuesLastMonth   int     `json:"issues_last_month"`
	PRsLastWeek       int     `json:"prs_last_week"`
	PRsLastMonth      int     `json:"prs_last_month"`
	AvgIssueAge       float64 `json:"avg_issue_age_days"`
	AvgPRAge          float64 `json:"avg_pr_age_days"`
}

// ContributorStats representa estatísticas de colaboradores
type ContributorStats struct {
	TopContributors    []*extractor.Contributor `json:"top_contributors"`
	TotalContributors  int                      `json:"total_contributors"`
	NewContributors    int                      `json:"new_contributors_last_month"`
	CoreTeamSize       int                      `json:"core_team_size"`
}

// RepositoryHealth representa a saúde do repositório
type RepositoryHealth struct {
	HealthScore        float64 `json:"health_score"`
	LastCommitDays     int     `json:"last_commit_days_ago"`
	LastReleaseDays    int     `json:"last_release_days_ago"`
	OpenIssuesRatio    float64 `json:"open_issues_ratio"`
	StaleIssues        int     `json:"stale_issues_count"`
	MaintenanceStatus  string  `json:"maintenance_status"`
}

// AnalyzeLanguages analisa a distribuição de linguagens
func AnalyzeLanguages(data *extractor.RepositoryData) []*LanguageStats {
	if len(data.Languages) == 0 {
		return nil
	}

	total := 0
	for _, bytes := range data.Languages {
		total += bytes
	}

	var stats []*LanguageStats
	for name, bytes := range data.Languages {
		percentage := float64(bytes) / float64(total) * 100
		stats = append(stats, &LanguageStats{
			Name:       name,
			Bytes:      bytes,
			Percentage: percentage,
		})
	}

	// Ordenar por porcentagem (decrescente)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Percentage > stats[j].Percentage
	})

	return stats
}

// AnalyzeActivity analisa a atividade do repositório
func AnalyzeActivity(data *extractor.RepositoryData) *ActivityMetrics {
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)
	monthAgo := now.AddDate(0, -1, 0)

	metrics := &ActivityMetrics{}

	// Analisar commits
	for _, commit := range data.RecentCommits {
		if commit.CreatedAt.After(weekAgo) {
			metrics.CommitsLastWeek++
		}
		if commit.CreatedAt.After(monthAgo) {
			metrics.CommitsLastMonth++
		}
	}

	// Analisar issues
	var issueAges []float64
	for _, issue := range data.RecentIssues {
		if issue.CreatedAt.After(weekAgo) {
			metrics.IssuesLastWeek++
		}
		if issue.CreatedAt.After(monthAgo) {
			metrics.IssuesLastMonth++
		}
		
		// Calcular idade média das issues
		age := now.Sub(issue.CreatedAt).Hours() / 24
		issueAges = append(issueAges, age)
	}

	// Analisar PRs
	var prAges []float64
	for _, pr := range data.RecentPRs {
		if pr.CreatedAt.After(weekAgo) {
			metrics.PRsLastWeek++
		}
		if pr.CreatedAt.After(monthAgo) {
			metrics.PRsLastMonth++
		}
		
		// Calcular idade média dos PRs
		age := now.Sub(pr.CreatedAt).Hours() / 24
		prAges = append(prAges, age)
	}

	// Calcular médias
	if len(issueAges) > 0 {
		sum := 0.0
		for _, age := range issueAges {
			sum += age
		}
		metrics.AvgIssueAge = sum / float64(len(issueAges))
	}

	if len(prAges) > 0 {
		sum := 0.0
		for _, age := range prAges {
			sum += age
		}
		metrics.AvgPRAge = sum / float64(len(prAges))
	}

	return metrics
}

// AnalyzeContributors analisa os colaboradores
func AnalyzeContributors(data *extractor.RepositoryData) *ContributorStats {
	stats := &ContributorStats{
		TotalContributors: len(data.Contributors),
	}

	if len(data.Contributors) == 0 {
		return stats
	}

	// Top 10 colaboradores
	topCount := 10
	if len(data.Contributors) < topCount {
		topCount = len(data.Contributors)
	}
	stats.TopContributors = data.Contributors[:topCount]

	// Estimar time principal (colaboradores com mais de 100 contribuições)
	for _, contrib := range data.Contributors {
		if contrib.Contributions >= 100 {
			stats.CoreTeamSize++
		}
	}

	return stats
}

// AnalyzeHealth analisa a saúde do repositório
func AnalyzeHealth(data *extractor.RepositoryData) *RepositoryHealth {
	now := time.Now()
	health := &RepositoryHealth{}

	// Dias desde o último commit
	if len(data.RecentCommits) > 0 {
		lastCommit := data.RecentCommits[0].CreatedAt
		health.LastCommitDays = int(now.Sub(lastCommit).Hours() / 24)
	}

	// Dias desde o último release
	if len(data.Releases) > 0 {
		lastRelease := data.Releases[0].PublishedAt
		health.LastReleaseDays = int(now.Sub(lastRelease).Hours() / 24)
	}

	// Ratio de issues abertas
	if data.Statistics.Issues > 0 {
		// Estimativa baseada nas issues da amostra
		openCount := 0
		for _, issue := range data.RecentIssues {
			if issue.State == "open" {
				openCount++
			}
		}
		if len(data.RecentIssues) > 0 {
			health.OpenIssuesRatio = float64(openCount) / float64(len(data.RecentIssues))
		}
	}

	// Issues obsoletas (mais de 90 dias sem atividade)
	staleThreshold := now.AddDate(0, 0, -90)
	for _, issue := range data.RecentIssues {
		if issue.State == "open" && issue.UpdatedAt.Before(staleThreshold) {
			health.StaleIssues++
		}
	}

	// Calcular score de saúde (0-100)
	score := 100.0

	// Penalizar por inatividade
	if health.LastCommitDays > 30 {
		score -= 20
	} else if health.LastCommitDays > 7 {
		score -= 10
	}

	// Penalizar por releases antigas
	if health.LastReleaseDays > 365 {
		score -= 15
	} else if health.LastReleaseDays > 180 {
		score -= 10
	}

	// Penalizar por muitas issues abertas
	if health.OpenIssuesRatio > 0.8 {
		score -= 15
	} else if health.OpenIssuesRatio > 0.6 {
		score -= 10
	}

	// Penalizar por issues obsoletas
	if health.StaleIssues > 10 {
		score -= 10
	} else if health.StaleIssues > 5 {
		score -= 5
	}

	health.HealthScore = score

	// Determinar status de manutenção
	switch {
	case score >= 90:
		health.MaintenanceStatus = "Excelente"
	case score >= 80:
		health.MaintenanceStatus = "Muito Bom"
	case score >= 70:
		health.MaintenanceStatus = "Bom"
	case score >= 60:
		health.MaintenanceStatus = "Regular"
	case score >= 50:
		health.MaintenanceStatus = "Precisa Atenção"
	default:
		health.MaintenanceStatus = "Crítico"
	}

	return health
}

// GenerateReport gera um relatório completo de análise
func GenerateReport(data *extractor.RepositoryData) string {
	var report strings.Builder

	report.WriteString("📊 RELATÓRIO COMPLETO DE ANÁLISE\n")
	report.WriteString(strings.Repeat("=", 80) + "\n\n")

	// Informações básicas
	report.WriteString("📋 INFORMAÇÕES BÁSICAS\n")
	report.WriteString(strings.Repeat("-", 40) + "\n")
	report.WriteString(fmt.Sprintf("Nome: %s\n", data.BasicInfo.FullName))
	report.WriteString(fmt.Sprintf("Descrição: %s\n", data.BasicInfo.Description))
	report.WriteString(fmt.Sprintf("Criado em: %s\n", data.BasicInfo.CreatedAt.Format("02/01/2006")))
	report.WriteString(fmt.Sprintf("Licença: %s\n", data.BasicInfo.License))
	report.WriteString(fmt.Sprintf("Tamanho: %d KB\n\n", data.BasicInfo.Size))

	// Estatísticas
	report.WriteString("📈 ESTATÍSTICAS\n")
	report.WriteString(strings.Repeat("-", 40) + "\n")
	report.WriteString(fmt.Sprintf("⭐ Stars: %s\n", formatNumber(data.Statistics.Stars)))
	report.WriteString(fmt.Sprintf("🍴 Forks: %s\n", formatNumber(data.Statistics.Forks)))
	report.WriteString(fmt.Sprintf("👀 Watchers: %s\n", formatNumber(data.Statistics.Watchers)))
	report.WriteString(fmt.Sprintf("🎯 Issues: %s\n\n", formatNumber(data.Statistics.Issues)))

	// Análise de linguagens
	languages := AnalyzeLanguages(data)
	if len(languages) > 0 {
		report.WriteString("💻 DISTRIBUIÇÃO DE LINGUAGENS\n")
		report.WriteString(strings.Repeat("-", 40) + "\n")
		for i, lang := range languages {
			if i >= 5 { // Top 5
				break
			}
			report.WriteString(fmt.Sprintf("%s: %.1f%%\n", lang.Name, lang.Percentage))
		}
		report.WriteString("\n")
	}

	// Análise de atividade
	activity := AnalyzeActivity(data)
	report.WriteString("⚡ ATIVIDADE RECENTE\n")
	report.WriteString(strings.Repeat("-", 40) + "\n")
	report.WriteString(fmt.Sprintf("Commits (última semana): %d\n", activity.CommitsLastWeek))
	report.WriteString(fmt.Sprintf("Commits (último mês): %d\n", activity.CommitsLastMonth))
	report.WriteString(fmt.Sprintf("Issues (última semana): %d\n", activity.IssuesLastWeek))
	report.WriteString(fmt.Sprintf("PRs (última semana): %d\n", activity.PRsLastWeek))
	report.WriteString(fmt.Sprintf("Idade média das issues: %.1f dias\n", activity.AvgIssueAge))
	report.WriteString(fmt.Sprintf("Idade média dos PRs: %.1f dias\n\n", activity.AvgPRAge))

	// Análise de colaboradores
	contributors := AnalyzeContributors(data)
	report.WriteString("👥 COLABORADORES\n")
	report.WriteString(strings.Repeat("-", 40) + "\n")
	report.WriteString(fmt.Sprintf("Total de colaboradores: %d\n", contributors.TotalContributors))
	report.WriteString(fmt.Sprintf("Time principal (100+ commits): %d\n", contributors.CoreTeamSize))
	report.WriteString("Top 5 colaboradores:\n")
	for i, contrib := range contributors.TopContributors {
		if i >= 5 {
			break
		}
		report.WriteString(fmt.Sprintf("  %d. %s (%d contribuições)\n", i+1, contrib.Login, contrib.Contributions))
	}
	report.WriteString("\n")

	// Análise de saúde
	health := AnalyzeHealth(data)
	report.WriteString("🏥 SAÚDE DO REPOSITÓRIO\n")
	report.WriteString(strings.Repeat("-", 40) + "\n")
	report.WriteString(fmt.Sprintf("Score de saúde: %.1f/100\n", health.HealthScore))
	report.WriteString(fmt.Sprintf("Status: %s\n", health.MaintenanceStatus))
	report.WriteString(fmt.Sprintf("Último commit: %d dias atrás\n", health.LastCommitDays))
	report.WriteString(fmt.Sprintf("Último release: %d dias atrás\n", health.LastReleaseDays))
	report.WriteString(fmt.Sprintf("Issues obsoletas: %d\n", health.StaleIssues))
	report.WriteString(fmt.Sprintf("Ratio de issues abertas: %.1f%%\n\n", health.OpenIssuesRatio*100))

	// Releases recentes
	if len(data.Releases) > 0 {
		report.WriteString("🚀 RELEASES RECENTES\n")
		report.WriteString(strings.Repeat("-", 40) + "\n")
		for i, release := range data.Releases {
			if i >= 3 { // Top 3
				break
			}
			report.WriteString(fmt.Sprintf("%s (%s) - %s\n", 
				release.TagName, 
				release.PublishedAt.Format("02/01/2006"),
				release.Author))
		}
		report.WriteString("\n")
	}

	report.WriteString(strings.Repeat("=", 80) + "\n")
	report.WriteString(fmt.Sprintf("Relatório gerado em: %s\n", time.Now().Format("02/01/2006 15:04:05")))

	return report.String()
}

// formatNumber formata números grandes com separadores
func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	} else if n < 1000000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	} else {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
}