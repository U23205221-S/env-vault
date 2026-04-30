package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"filippo.io/age"
	"github.com/spf13/cobra"
)

// Manifest define la estructura del archivo manifest.json
type Manifest struct {
	PublicKeys []string `json:"public_keys"`
}

var addUserCmd = &cobra.Command{
	Use:   "add-user <public_key>",
	Short: "Agrega una clave pública al manifiesto de usuarios autorizados",
	Args:  cobra.ExactArgs(1), // Requiere exactamente 1 argumento
	Run: func(cmd *cobra.Command, args []string) {
		publicKey := args[0]

		// 1. Validar que la clave pública tenga el formato correcto de age (X25519)
		_, err := age.ParseX25519Recipient(publicKey)
		if err != nil {
			log.Fatalf("❌ Error: La clave pública ingresada no es válida. Debe empezar con 'age1...'\nDetalle: %v", err)
		}

		// 2. Verificar que exista el vault
		vaultDir := ".env-vault"
		manifestFile := filepath.Join(vaultDir, "manifest.json")
		if _, err := os.Stat(manifestFile); os.IsNotExist(err) {
			log.Fatalf("❌ Error: No se encontró %s. ¿Ya ejecutaste 'env-vault init' en este repositorio?", manifestFile)
		}

		// 3. Leer el manifest.json actual
		fileData, err := os.ReadFile(manifestFile)
		if err != nil {
			log.Fatalf("❌ Error leyendo %s: %v", manifestFile, err)
		}

		var manifest Manifest
		if err := json.Unmarshal(fileData, &manifest); err != nil {
			log.Fatalf("❌ Error parseando %s. Asegurate de que sea un JSON válido: %v", manifestFile, err)
		}

		// 4. Chequear si la clave ya existe para no duplicarla
		for _, key := range manifest.PublicKeys {
			if key == publicKey {
				fmt.Println("ℹ️  Esta clave pública ya estaba autorizada en el manifiesto.")
				return
			}
		}

		// 5. Agregar la clave y guardar
		manifest.PublicKeys = append(manifest.PublicKeys, publicKey)

		newData, err := json.MarshalIndent(manifest, "", "  ")
		if err != nil {
			log.Fatalf("❌ Error generando nuevo JSON: %v", err)
		}

		// Agregar un salto de línea al final (buena práctica en archivos de texto)
		newData = append(newData, '\n')

		if err := os.WriteFile(manifestFile, newData, 0644); err != nil {
			log.Fatalf("❌ Error guardando %s: %v", manifestFile, err)
		}

		fmt.Printf("✅ Clave pública agregada con éxito al proyecto.\n")
		fmt.Printf("🔑 Total de usuarios autorizados: %d\n", len(manifest.PublicKeys))
		fmt.Println("⚠️  Recordá ejecutar 'env-vault push' para re-cifrar los secretos con los nuevos accesos.")
	},
}

func init() {
	rootCmd.AddCommand(addUserCmd)
}
