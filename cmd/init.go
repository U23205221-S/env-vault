package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicializa el vault en el repositorio Git actual",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Inicializando env-vault en el repositorio...")

		// 1. Verificar si estamos en un repositorio Git
		if _, err := os.Stat(".git"); os.IsNotExist(err) {
			log.Fatalf("❌ Error: No se encontró el directorio .git. Debes ejecutar este comando en la raíz de un repositorio Git.")
		}

		// 2. Crear directorio .env-vault/
		vaultDir := ".env-vault"
		if err := os.MkdirAll(vaultDir, 0755); err != nil {
			log.Fatalf("❌ Error creando directorio %s: %v", vaultDir, err)
		}
		fmt.Printf("✅ Directorio %s creado.\n", vaultDir)

		// 3. Crear manifest.json inicial si no existe
		manifestFile := filepath.Join(vaultDir, "manifest.json")
		if _, err := os.Stat(manifestFile); os.IsNotExist(err) {
			initialManifest := []byte(`{
  "public_keys": []
}
`)
			if err := os.WriteFile(manifestFile, initialManifest, 0644); err != nil {
				log.Fatalf("❌ Error creando %s: %v", manifestFile, err)
			}
			fmt.Printf("✅ Archivo %s creado.\n", manifestFile)
		} else {
			fmt.Printf("ℹ️  Archivo %s ya existe, se mantiene intacto.\n", manifestFile)
		}

		// 4. Configurar el Git Hook (pre-commit)
		hooksDir := filepath.Join(".git", "hooks")
		if err := os.MkdirAll(hooksDir, 0755); err != nil {
			log.Fatalf("❌ Error comprobando directorio de hooks %s: %v", hooksDir, err)
		}

		preCommitFile := filepath.Join(hooksDir, "pre-commit")
		hookScript := `#!/bin/sh
# Hook generado por env-vault
# Evita hacer commit de un archivo llamado .env (en texto plano)

if git diff --cached --name-only | grep -E '(^|/)\.env$'; then
  echo "🚨 ERROR: Estás intentando commitear un archivo .env en texto plano."
  echo "🚨 Usa 'env-vault push' para cifrarlo y agregarlo al vault."
  echo "🚨 Abortando commit..."
  exit 1
fi
`
		
		// Revisar si ya existe el hook y si ya tiene nuestra protección
		hookContent, err := os.ReadFile(preCommitFile)
		if err == nil {
			// Si ya existe, nos fijamos si ya tiene nuestra palabra clave
			if !containsEnvVaultHook(string(hookContent)) {
				// Hacemos un append
				f, err := os.OpenFile(preCommitFile, os.O_APPEND|os.O_WRONLY, 0755)
				if err != nil {
					log.Fatalf("❌ Error abriendo %s: %v", preCommitFile, err)
				}
				defer f.Close()
				if _, err := f.WriteString("\n" + hookScript); err != nil {
					log.Fatalf("❌ Error escribiendo hook: %v", err)
				}
				fmt.Printf("✅ Hook agregado exitosamente a %s.\n", preCommitFile)
			} else {
				fmt.Printf("ℹ️  El hook %s ya contenía protección de env-vault.\n", preCommitFile)
			}
		} else if os.IsNotExist(err) {
			// Si no existe, lo creamos de cero y le damos permisos de ejecución (0755)
			if err := os.WriteFile(preCommitFile, []byte(hookScript), 0755); err != nil {
				log.Fatalf("❌ Error creando hook %s: %v", preCommitFile, err)
			}
			fmt.Printf("✅ Hook %s creado exitosamente.\n", preCommitFile)
		} else {
			log.Fatalf("❌ Error verificando hook %s: %v", preCommitFile, err)
		}

		fmt.Println("\n🎉 ¡Vault inicializado con éxito! El repositorio ahora está protegido.")
	},
}

func containsEnvVaultHook(content string) bool {
	// Simple chequeo para ver si ya inyectamos el código de env-vault antes
	return len(content) > 0 && ( /* un check rústico */ len(content) > 0 && 
		(stringContains(content, "env-vault push") || stringContains(content, "Estás intentando commitear un archivo .env")))
}

func stringContains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(initCmd)
}
