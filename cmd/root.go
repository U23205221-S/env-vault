package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "env-vault",
	Short: "Sincroniza archivos .env de forma segura usando Git",
	Long: `env-vault es una herramienta CLI serverless.
Permite cifrar y compartir variables de entorno en equipos
usando criptografía asimétrica (age) y el repositorio de Git como transporte.
Nadie comparte secretos por Slack, y no se requiere servidor VPS.`,
}

// Execute corre el comando principal de la CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
