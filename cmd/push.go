package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Cifra el archivo .env y lo prepara para Git",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔒 Cifrando variables de entorno locales...")
		// TODO: Leer .env local
		// TODO: Leer manifest.json para obtener claves públicas autorizadas
		// TODO: Cifrar y guardar en .env-vault/development.enc
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
