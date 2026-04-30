package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Descifra el archivo sincronizado de Git al .env local",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔓 Descifrando variables de entorno del repositorio...")
		// TODO: Buscar clave privada en ~/.env-vault/keys/
		// TODO: Leer archivo .env-vault/development.enc cifrado
		// TODO: Descifrar y escribir localmente como .env
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
