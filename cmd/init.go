package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicializa el vault en el repositorio Git actual",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Inicializando .env-vault en el repositorio...")
		// TODO: Crear directorio .env-vault/
		// TODO: Crear manifest.json con las claves públicas
		// TODO: Configurar hook de Git pre-commit
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
