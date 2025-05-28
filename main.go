package main

import (
	"log"

	"github-octokit-poc/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Erro na execução: %v", err)
	}
}