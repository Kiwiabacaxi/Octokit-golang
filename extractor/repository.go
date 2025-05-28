package extractor

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	ghclient "github-octokit-poc/github"

	"github.com/google/go-github/v57/github"
)

type RepositoryData struct {
	// Informações básicas
	BasicInfo *BasicInfo `json:"basic_info"`
	
	// Estatísticas
	Statistics *Statistics `json:"statistics"`
	
	// Configurações
	Settings *Settings `json:"settings"`
	
	// Linguagens
	Languages map[string]int `json:"languages"`
	
	// Tópicos/Tags
	Topics []string `json:"topics"`
	
	// Colaboradores
	Contributors []*Contributor `json:"contributors"`
	
	// Issues recentes
	RecentIssues []*IssueData `json:"recent_issues"`
	
	// Pull Requests recentes
	RecentPRs []*PullRequestData `json:"recent_prs"`
	
	// Releases
	Releases []*ReleaseData `json:"releases"`
	
	// Commits recentes
	RecentCommits []*CommitData `json:"recent_commits"`
	
	// Eventos recentes
	RecentEvents []*EventData `json:"recent_events"`
	
	// Rate limit info
	RateLimit *RateLimitData `json:"rate_limit"`
	
	// Metadados da extração
	ExtractionMeta *ExtractionMeta `json:"extraction_meta"`
}

type BasicInfo struct {
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Owner           string    `json:"owner"`
	Description     string    `json:"description"`
	URL             string    `json:"url"`
	Homepage        string    `json:"homepage"`
	CloneURL        string    `json:"clone_url"`
	SSHURL          string    `json:"ssh_url"`
	DefaultBranch   string    `json:"default_branch"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
	Size            int       `json:"size_kb"`
	License         string    `json:"license"`
}

type Statistics struct {
	Stars         int `json:"stars"`
	Forks         int `json:"forks"`
	Watchers      int `json:"watchers"`
	Issues        int `json:"open_issues"`
	Subscribers   int `json:"subscribers"`
	NetworkCount  int `json:"network_count"`
}

type Settings struct {
	Private          bool `json:"private"`
	Fork             bool `json:"fork"`
	Archived         bool `json:"archived"`
	Disabled         bool `json:"disabled"`
	HasIssues        bool `json:"has_issues"`
	HasProjects      bool `json:"has_projects"`
	HasWiki          bool `json:"has_wiki"`
	HasPages         bool `json:"has_pages"`
	HasDiscussions   bool `json:"has_discussions"`
	HasDownloads     bool `json:"has_downloads"`
	AllowForking     bool `json:"allow_forking"`
	AllowMergeCommit bool `json:"allow_merge_commit"`
	AllowSquashMerge bool `json:"allow_squash_merge"`
	AllowRebaseMerge bool `json:"allow_rebase_merge"`
}

type Contributor struct {
	Login         string `json:"login"`
	Contributions int    `json:"contributions"`
	AvatarURL     string `json:"avatar_url"`
	Type          string `json:"type"`
}

type IssueData struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Labels    []string  `json:"labels"`
	Comments  int       `json:"comments"`
}

type PullRequestData struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Merged    bool      `json:"merged"`
	Draft     bool      `json:"draft"`
}

type ReleaseData struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
	Author      string    `json:"author"`
}

type CommitData struct {
	SHA       string    `json:"sha"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	URL       string    `json:"url"`
}

type EventData struct {
	Type      string    `json:"type"`
	Actor     string    `json:"actor"`
	CreatedAt time.Time `json:"created_at"`
	Public    bool      `json:"public"`
}

type RateLimitData struct {
	Core      *github.Rate `json:"core"`
	Search    *github.Rate `json:"search"`
	GraphQL   *github.Rate `json:"graphql"`
	Resources *github.Rate `json:"resources"`
}

type ExtractionMeta struct {
	ExtractedAt time.Time `json:"extracted_at"`
	Owner       string    `json:"owner"`
	Repo        string    `json:"repo"`
	Duration    string    `json:"duration"`
	APIVersion  string    `json:"api_version"`
}

// ExtractRepositoryData extrai todos os dados possíveis de um repositório
func ExtractRepositoryData(client *ghclient.Client, owner, repo string) (*RepositoryData, error) {
	startTime := time.Now()
	
	log.Printf("🔍 Iniciando extração completa do repositório %s/%s", owner, repo)
	
	data := &RepositoryData{
		ExtractionMeta: &ExtractionMeta{
			ExtractedAt: startTime,
			Owner:       owner,
			Repo:        repo,
			APIVersion:  "v3",
		},
	}

	// 1. Informações básicas do repositório
	log.Println("📋 Extraindo informações básicas...")
	if err := extractBasicInfo(client, owner, repo, data); err != nil {
		return nil, fmt.Errorf("erro ao extrair informações básicas: %v", err)
	}

	// 2. Linguagens do repositório
	log.Println("💻 Extraindo linguagens...")
	if err := extractLanguages(client, owner, repo, data); err != nil {
		log.Printf("⚠️ Erro ao extrair linguagens: %v", err)
	}

	// 3. Colaboradores
	log.Println("👥 Extraindo colaboradores...")
	if err := extractContributors(client, owner, repo, data); err != nil {
		log.Printf("⚠️ Erro ao extrair colaboradores: %v", err)
	}

	// 4. Issues recentes
	log.Println("🎯 Extraindo issues recentes...")
	if err := extractRecentIssues(client, owner, repo, data); err != nil {
		log.Printf("⚠️ Erro ao extrair issues: %v", err)
	}

	// 5. Pull Requests recentes
	log.Println("🔄 Extraindo pull requests recentes...")
	if err := extractRecentPRs(client, owner, repo, data); err != nil {
		log.Printf("⚠️ Erro ao extrair PRs: %v", err)
	}

	// 6. Releases
	log.Println("🚀 Extraindo releases...")
	if err := extractReleases(client, owner, repo, data); err != nil {
		log.Printf("⚠️ Erro ao extrair releases: %v", err)
	}

	// 7. Commits recentes
	log.Println("📝 Extraindo commits recentes...")
	if err := extractRecentCommits(client, owner, repo, data); err != nil {
		log.Printf("⚠️ Erro ao extrair commits: %v", err)
	}

	// 8. Eventos recentes
	log.Println("⚡ Extraindo eventos recentes...")
	if err := extractRecentEvents(client, owner, repo, data); err != nil {
		log.Printf("⚠️ Erro ao extrair eventos: %v", err)
	}

	// 9. Rate limit
	log.Println("📊 Verificando rate limits...")
	if err := extractRateLimit(client, data); err != nil {
		log.Printf("⚠️ Erro ao verificar rate limits: %v", err)
	}

	// Finalização
	data.ExtractionMeta.Duration = time.Since(startTime).String()
	
	log.Printf("✅ Extração concluída em %s", data.ExtractionMeta.Duration)
	
	return data, nil
}

func extractBasicInfo(client *ghclient.Client, owner, repo string, data *RepositoryData) error {
	repository, _, err := client.GitHub.Repositories.Get(client.Ctx, owner, repo)
	if err != nil {
		return err
	}

	license := ""
	if repository.License != nil {
		license = repository.License.GetName()
	}

	data.BasicInfo = &BasicInfo{
		Name:          repository.GetName(),
		FullName:      repository.GetFullName(),
		Owner:         repository.GetOwner().GetLogin(),
		Description:   repository.GetDescription(),
		URL:           repository.GetHTMLURL(),
		Homepage:      repository.GetHomepage(),
		CloneURL:      repository.GetCloneURL(),
		SSHURL:        repository.GetSSHURL(),
		DefaultBranch: repository.GetDefaultBranch(),
		CreatedAt:     repository.GetCreatedAt().Time,
		UpdatedAt:     repository.GetUpdatedAt().Time,
		PushedAt:      repository.GetPushedAt().Time,
		Size:          repository.GetSize(),
		License:       license,
	}

	data.Statistics = &Statistics{
		Stars:        repository.GetStargazersCount(),
		Forks:        repository.GetForksCount(),
		Watchers:     repository.GetWatchersCount(),
		Issues:       repository.GetOpenIssuesCount(),
		Subscribers:  repository.GetSubscribersCount(),
		NetworkCount: repository.GetNetworkCount(),
	}

	data.Settings = &Settings{
		Private:          repository.GetPrivate(),
		Fork:             repository.GetFork(),
		Archived:         repository.GetArchived(),
		Disabled:         repository.GetDisabled(),
		HasIssues:        repository.GetHasIssues(),
		HasProjects:      repository.GetHasProjects(),
		HasWiki:          repository.GetHasWiki(),
		HasPages:         repository.GetHasPages(),
		HasDiscussions:   repository.GetHasDiscussions(),
		HasDownloads:     repository.GetHasDownloads(),
		AllowForking:     repository.GetAllowForking(),
		AllowMergeCommit: repository.GetAllowMergeCommit(),
		AllowSquashMerge: repository.GetAllowSquashMerge(),
		AllowRebaseMerge: repository.GetAllowRebaseMerge(),
	}

	data.Topics = repository.Topics

	return nil
}

func extractLanguages(client *ghclient.Client, owner, repo string, data *RepositoryData) error {
	languages, _, err := client.GitHub.Repositories.ListLanguages(client.Ctx, owner, repo)
	if err != nil {
		return err
	}

	data.Languages = languages
	return nil
}

func extractContributors(client *ghclient.Client, owner, repo string, data *RepositoryData) error {
	opts := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}

	contributors, _, err := client.GitHub.Repositories.ListContributors(client.Ctx, owner, repo, opts)
	if err != nil {
		return err
	}

	data.Contributors = make([]*Contributor, len(contributors))
	for i, contrib := range contributors {
		data.Contributors[i] = &Contributor{
			Login:         contrib.GetLogin(),
			Contributions: contrib.GetContributions(),
			AvatarURL:     contrib.GetAvatarURL(),
			Type:          contrib.GetType(),
		}
	}

	return nil
}

func extractRecentIssues(client *ghclient.Client, owner, repo string, data *RepositoryData) error {
	opts := &github.IssueListByRepoOptions{
		State:       "all",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 10},
	}

	issues, _, err := client.GitHub.Issues.ListByRepo(client.Ctx, owner, repo, opts)
	if err != nil {
		return err
	}

	data.RecentIssues = make([]*IssueData, 0)
	for _, issue := range issues {
		if issue.PullRequestLinks == nil { // Apenas issues, não PRs
			labels := make([]string, len(issue.Labels))
			for i, label := range issue.Labels {
				labels[i] = label.GetName()
			}

			data.RecentIssues = append(data.RecentIssues, &IssueData{
				Number:    issue.GetNumber(),
				Title:     issue.GetTitle(),
				State:     issue.GetState(),
				Author:    issue.GetUser().GetLogin(),
				CreatedAt: issue.GetCreatedAt().Time,
				UpdatedAt: issue.GetUpdatedAt().Time,
				Labels:    labels,
				Comments:  issue.GetComments(),
			})
		}
	}

	return nil
}

func extractRecentPRs(client *ghclient.Client, owner, repo string, data *RepositoryData) error {
	opts := &github.PullRequestListOptions{
		State:       "all",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 10},
	}

	prs, _, err := client.GitHub.PullRequests.List(client.Ctx, owner, repo, opts)
	if err != nil {
		return err
	}

	data.RecentPRs = make([]*PullRequestData, len(prs))
	for i, pr := range prs {
		data.RecentPRs[i] = &PullRequestData{
			Number:    pr.GetNumber(),
			Title:     pr.GetTitle(),
			State:     pr.GetState(),
			Author:    pr.GetUser().GetLogin(),
			CreatedAt: pr.GetCreatedAt().Time,
			UpdatedAt: pr.GetUpdatedAt().Time,
			Merged:    pr.GetMerged(),
			Draft:     pr.GetDraft(),
		}
	}

	return nil
}

func extractReleases(client *ghclient.Client, owner, repo string, data *RepositoryData) error {
	opts := &github.ListOptions{PerPage: 10}

	releases, _, err := client.GitHub.Repositories.ListReleases(client.Ctx, owner, repo, opts)
	if err != nil {
		return err
	}

	data.Releases = make([]*ReleaseData, len(releases))
	for i, release := range releases {
		data.Releases[i] = &ReleaseData{
			TagName:     release.GetTagName(),
			Name:        release.GetName(),
			CreatedAt:   release.GetCreatedAt().Time,
			PublishedAt: release.GetPublishedAt().Time,
			Prerelease:  release.GetPrerelease(),
			Draft:       release.GetDraft(),
			Author:      release.GetAuthor().GetLogin(),
		}
	}

	return nil
}

func extractRecentCommits(client *ghclient.Client, owner, repo string, data *RepositoryData) error {
	opts := &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	commits, _, err := client.GitHub.Repositories.ListCommits(client.Ctx, owner, repo, opts)
	if err != nil {
		return err
	}

	data.RecentCommits = make([]*CommitData, len(commits))
	for i, commit := range commits {
		data.RecentCommits[i] = &CommitData{
			SHA:       commit.GetSHA(),
			Message:   commit.GetCommit().GetMessage(),
			Author:    commit.GetCommit().GetAuthor().GetName(),
			CreatedAt: commit.GetCommit().GetAuthor().GetDate().Time,
			URL:       commit.GetHTMLURL(),
		}
	}

	return nil
}

func extractRecentEvents(client *ghclient.Client, owner, repo string, data *RepositoryData) error {
	opts := &github.ListOptions{PerPage: 10}

	events, _, err := client.GitHub.Activity.ListRepositoryEvents(client.Ctx, owner, repo, opts)
	if err != nil {
		return err
	}

	data.RecentEvents = make([]*EventData, len(events))
	for i, event := range events {
		data.RecentEvents[i] = &EventData{
			Type:      event.GetType(),
			Actor:     event.GetActor().GetLogin(),
			CreatedAt: event.GetCreatedAt().Time,
			Public:    event.GetPublic(),
		}
	}

	return nil
}

func extractRateLimit(client *ghclient.Client, data *RepositoryData) error {
	rates, _, err := client.GitHub.RateLimits(client.Ctx)
	if err != nil {
		return err
	}

	data.RateLimit = &RateLimitData{
		Core:      rates.Core,
		Search:    rates.Search,
		GraphQL:   rates.GraphQL,
		// Resources: rates.Resources,
	}

	return nil
}

// PrintSummary imprime um resumo dos dados extraídos
func (rd *RepositoryData) PrintSummary() {
	strings.Repeat("=", 80)
	fmt.Printf("📊 RESUMO DA EXTRAÇÃO - %s\n", rd.BasicInfo.FullName)
	strings.Repeat("=", 80)
	
	fmt.Printf("🏷️  Nome: %s\n", rd.BasicInfo.Name)
	fmt.Printf("👤 Proprietário: %s\n", rd.BasicInfo.Owner)
	fmt.Printf("📝 Descrição: %s\n", rd.BasicInfo.Description)
	fmt.Printf("🌐 URL: %s\n", rd.BasicInfo.URL)
	fmt.Printf("📅 Criado em: %s\n", rd.BasicInfo.CreatedAt.Format("02/01/2006"))
	fmt.Printf("🔄 Última atualização: %s\n", rd.BasicInfo.UpdatedAt.Format("02/01/2006 15:04"))
	
	fmt.Println("\n📈 ESTATÍSTICAS:")
	fmt.Printf("⭐ Stars: %d\n", rd.Statistics.Stars)
	fmt.Printf("🍴 Forks: %d\n", rd.Statistics.Forks)
	fmt.Printf("👀 Watchers: %d\n", rd.Statistics.Watchers)
	fmt.Printf("🎯 Issues abertas: %d\n", rd.Statistics.Issues)
	
	fmt.Println("\n💻 LINGUAGENS:")
	total := 0
	for _, bytes := range rd.Languages {
		total += bytes
	}
	for lang, bytes := range rd.Languages {
		percentage := float64(bytes) / float64(total) * 100
		fmt.Printf("   %s: %.1f%%\n", lang, percentage)
	}
	
	fmt.Printf("\n👥 COLABORADORES: %d encontrados\n", len(rd.Contributors))
	fmt.Printf("🎯 ISSUES RECENTES: %d encontradas\n", len(rd.RecentIssues))
	fmt.Printf("🔄 PULL REQUESTS: %d encontrados\n", len(rd.RecentPRs))
	fmt.Printf("🚀 RELEASES: %d encontrados\n", len(rd.Releases))
	fmt.Printf("📝 COMMITS RECENTES: %d encontrados\n", len(rd.RecentCommits))
	fmt.Printf("⚡ EVENTOS RECENTES: %d encontrados\n", len(rd.RecentEvents))
	
	fmt.Println("\n📊 RATE LIMITS:")
	if rd.RateLimit.Core != nil {
		fmt.Printf("   Core API: %d/%d (reset em %s)\n", 
			rd.RateLimit.Core.Remaining, 
			rd.RateLimit.Core.Limit,
			rd.RateLimit.Core.Reset.Format("15:04:05"))
	}
	
	fmt.Printf("\n⏱️  Extração concluída em: %s\n", rd.ExtractionMeta.Duration)
	strings.Repeat("=", 80)
}

// SaveToJSON salva os dados em um arquivo JSON
func (rd *RepositoryData) SaveToJSON(filename string) error {
	data, err := json.MarshalIndent(rd, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filename, data, 0644)
}