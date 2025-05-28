package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config representa as configurações da aplicação
type Config struct {
	DefaultOwner string
	DefaultRepo  string
	OutputDir    string
	Debug        bool
}

// Load carrega as configurações do .env e variáveis de ambiente
func Load() (*Config, error) {
	// Carrega variáveis do arquivo .env
	if err := godotenv.Load(); err != nil {
		// Não é erro crítico, pode usar variáveis de ambiente do sistema
	}

	return &Config{
		DefaultOwner: getEnvOrDefault("GITHUB_DEFAULT_USER", "kubernetes"),
		DefaultRepo:  getEnvOrDefault("GITHUB_DEFAULT_REPO", "kubernetes"),
		OutputDir:    getEnvOrDefault("OUTPUT_DIR", "output"),
		Debug:        os.Getenv("DEBUG") == "true",
	}, nil
}

// GetTarget retorna o owner e repo alvo para análise
func (c *Config) GetTarget() (owner, repo string) {
	// Aqui você pode adicionar lógica para pegar de argumentos da linha de comando,
	// ou de outras fontes. Por enquanto, usa os padrões.
	return c.DefaultOwner, c.DefaultRepo
}

// getEnvOrDefault retorna o valor da variável de ambiente ou um valor padrão
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}