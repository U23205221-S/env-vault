package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var removeUserCmd = &cobra.Command{
	Use:   "remove-user <public_key>",
	Short: "Elimina una clave pública del manifiesto de usuarios autorizados",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		publicKey := args[0]

		// 1. Verificar que exista el vault
		vaultDir := ".env-vault"
		manifestFile := filepath.Join(vaultDir, "manifest.json")
		if _, err := os.Stat(manifestFile); os.IsNotExist(err) {
			log.Fatalf("❌ Error: No se encontró %s. ¿Ya ejecutaste 'env-vault init'?", manifestFile)
		}

		// 2. Leer el manifest.json actual
		fileData, err := os.ReadFile(manifestFile)
		if err != nil {
			log.Fatalf("❌ Error leyendo %s: %v", manifestFile, err)
		}

		var manifest Manifest
		if err := json.Unmarshal(fileData, &manifest); err != nil {
			log.Fatalf("❌ Error parseando %s: %v", manifestFile, err)
		}

		// 3. Filtrar la clave a eliminar
		found := false
		var newKeys []string
		for _, key := range manifest.PublicKeys {
			if key == publicKey {
				found = true
			} else {
				newKeys = append(newKeys, key)
			}
		}

		if !found {
			fmt.Println("ℹ️  La clave pública no se encontró en el manifiesto. No hay nada que eliminar.")
			return
		}

		// 4. Actualizar el manifiesto
		manifest.PublicKeys = newKeys

		newData, err := json.MarshalIndent(manifest, "", "  ")
		if err != nil {
			log.Fatalf("❌ Error generando nuevo JSON: %v", err)
		}

		// Agregar un salto de línea al final
		newData = append(newData, '\n')

		if err := os.WriteFile(manifestFile, newData, 0644); err != nil {
			log.Fatalf("❌ Error guardando %s: %v", manifestFile, err)
		}

		fmt.Printf("✅ Clave pública eliminada con éxito del proyecto.\n")
		fmt.Printf("🔑 Total de usuarios autorizados restantes: %d\n", len(manifest.PublicKeys))
		fmt.Println("⚠️  CRÍTICO: Ejecutá 'env-vault push' AHORA para re-cifrar los secretos y hacer efectiva la revocación de acceso.")
	},
}

func init() {
	rootCmd.AddCommand(removeUserCmd)
}
