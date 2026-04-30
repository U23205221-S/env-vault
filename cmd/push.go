package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"filippo.io/age"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Cifra el archivo .env y lo prepara para Git",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔒 Leyendo variables de entorno locales...")

		// 1. Leer el archivo .env local
		envFile := ".env"
		envContent, err := os.ReadFile(envFile)
		if err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("❌ Error: No se encontró el archivo '%s' en este directorio.", envFile)
			}
			log.Fatalf("❌ Error leyendo %s: %v", envFile, err)
		}

		// 2. Leer el manifest.json
		vaultDir := ".env-vault"
		manifestFile := filepath.Join(vaultDir, "manifest.json")
		manifestData, err := os.ReadFile(manifestFile)
		if err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("❌ Error: No se encontró %s. ¿Ejecutaste 'env-vault init' primero?", manifestFile)
			}
			log.Fatalf("❌ Error leyendo %s: %v", manifestFile, err)
		}

		// Reutilizamos el struct Manifest que definimos en add_user.go
		var manifest Manifest
		if err := json.Unmarshal(manifestData, &manifest); err != nil {
			log.Fatalf("❌ Error parseando %s. El JSON está corrupto: %v", manifestFile, err)
		}

		if len(manifest.PublicKeys) == 0 {
			log.Fatalf("❌ Error: No hay claves públicas en el manifiesto. Usá 'env-vault add-user <clave>' primero.")
		}

		// 3. Preparar los "candados" (recipients) leyendo las claves públicas
		var recipients []age.Recipient
		for _, pubKeyStr := range manifest.PublicKeys {
			recipient, err := age.ParseX25519Recipient(pubKeyStr)
			if err != nil {
				log.Fatalf("❌ Error parseando clave pública '%s': %v", pubKeyStr, err)
			}
			recipients = append(recipients, recipient)
		}

		// 4. Cifrar el contenido del .env
		fmt.Println("⏳ Cifrando archivo con criptografía asimétrica...")
		buf := &bytes.Buffer{}
		
		// Encrypt toma un io.Writer y N "candados" (recipients)
		writer, err := age.Encrypt(buf, recipients...)
		if err != nil {
			log.Fatalf("❌ Error inicializando el motor de cifrado: %v", err)
		}

		// Escribimos el .env adentro del candado y lo cerramos
		if _, err := writer.Write(envContent); err != nil {
			log.Fatalf("❌ Error escribiendo datos cifrados: %v", err)
		}
		if err := writer.Close(); err != nil {
			log.Fatalf("❌ Error finalizando el cifrado: %v", err)
		}

		// 5. Guardar en .env-vault/development.enc
		outputFile := filepath.Join(vaultDir, "development.enc")
		if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
			log.Fatalf("❌ Error guardando archivo cifrado %s: %v", outputFile, err)
		}

		fmt.Printf("✅ Archivo cifrado exitosamente y guardado en:\n   %s\n", outputFile)
		fmt.Printf("🔑 Cifrado para %d usuario(s) autorizado(s).\n", len(recipients))
		fmt.Println("\n🚀 Ya podés hacer 'git add' y 'git commit' de forma segura.")
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
