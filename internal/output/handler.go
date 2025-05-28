package output

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github-octokit-poc/extractor"
)

// Handler gerencia a criação e salvamento de arquivos de saída
type Handler struct {
	baseDir   string
	timestamp string
	owner     string
	repo      string
}

// NewHandler cria um novo handler de output com diretório padrão
func NewHandler(owner, repo string) *Handler {
	return NewHandlerWithDir(owner, repo, "output")
}

// NewHandlerWithDir cria um novo handler de output com diretório customizado
func NewHandlerWithDir(owner, repo, baseDir string) *Handler {
	timestamp := time.Now().Format("20060102_150405")
	
	return &Handler{
		baseDir:   baseDir,
		timestamp: timestamp,
		owner:     owner,
		repo:      repo,
	}
}

// SaveAll salva todos os outputs (JSON e relatório)
func (h *Handler) SaveAll(data *extractor.RepositoryData, report string) error {
	// Criar estrutura de diretórios
	outputDir, err := h.createOutputDirectory()
	if err != nil {
		return err
	}

	// Definir caminhos dos arquivos
	jsonFile := filepath.Join(outputDir, h.getJSONFilename())
	reportFile := filepath.Join(outputDir, h.getReportFilename())

	log.Printf("\n💾 Salvando arquivos:")
	log.Printf("   📊 Dados completos: %s", jsonFile)
	log.Printf("   📋 Relatório: %s", reportFile)

	// Salvar JSON
	if err := h.saveJSON(data, jsonFile); err != nil {
		log.Printf("⚠️ Erro ao salvar JSON: %v", err)
	} else {
		log.Printf("✅ JSON salvo com sucesso!")
	}

	// Salvar relatório
	if err := h.saveReport(report, reportFile); err != nil {
		log.Printf("⚠️ Erro ao salvar relatório: %v", err)
	} else {
		log.Printf("✅ Relatório salvo com sucesso!")
	}

	return nil
}

// createOutputDirectory cria a estrutura de diretórios necessária
func (h *Handler) createOutputDirectory() (string, error) {
	// Criar pasta output base se não existir
	if err := os.MkdirAll(h.baseDir, 0755); err != nil {
		return "", fmt.Errorf("erro ao criar diretório base: %v", err)
	}

	// Criar pasta específica para esta execução
	outputDir := filepath.Join(h.baseDir, h.timestamp)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Printf("⚠️ Erro ao criar diretório específico: %v", err)
		// Fallback para a pasta output base
		return h.baseDir, nil
	}

	return outputDir, nil
}

// saveJSON salva os dados em formato JSON
func (h *Handler) saveJSON(data *extractor.RepositoryData, filename string) error {
	return data.SaveToJSON(filename)
}

// saveReport salva o relatório em formato texto
func (h *Handler) saveReport(report, filename string) error {
	return os.WriteFile(filename, []byte(report), 0644)
}

// getJSONFilename gera o nome do arquivo JSON
func (h *Handler) getJSONFilename() string {
	return fmt.Sprintf("%s_%s_data.json", h.owner, h.repo)
}

// getReportFilename gera o nome do arquivo de relatório
func (h *Handler) getReportFilename() string {
	return fmt.Sprintf("%s_%s_report.txt", h.owner, h.repo)
}