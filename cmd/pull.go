package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Descifra el archivo sincronizado de Git al .env local",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔓 Buscando tu clave privada local...")

		// 1. Buscar la clave privada del usuario en ~/.env-vault/keys/identity.txt
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("❌ Error obteniendo directorio home: %v", err)
		}

		keyFile := filepath.Join(homeDir, ".env-vault", "keys", "identity.txt")
		keyData, err := os.ReadFile(keyFile)
		if err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("❌ Error: No se encontró tu clave privada en %s.\n¿Ya ejecutaste 'env-vault generate-key' en esta máquina?", keyFile)
			}
			log.Fatalf("❌ Error leyendo clave privada: %v", err)
		}

		// La clave privada es un string (generalmente empieza con AGE-SECRET-KEY-...)
		secretKey := strings.TrimSpace(string(keyData))
		
		// Parsear la identidad usando age
		identity, err := age.ParseX25519Identity(secretKey)
		if err != nil {
			log.Fatalf("❌ Error parseando tu clave privada. Asegurate de que el archivo no esté corrupto: %v", err)
		}

		fmt.Println("⏳ Descifrando archivo desde Git...")

		// 2. Leer el archivo cifrado del repo
		vaultDir := ".env-vault"
		encryptedFile := filepath.Join(vaultDir, "development.enc")
		encryptedData, err := os.ReadFile(encryptedFile)
		if err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("❌ Error: No se encontró %s. ¿Alguien en el equipo ya ejecutó 'env-vault push'?", encryptedFile)
			}
			log.Fatalf("❌ Error leyendo archivo cifrado: %v", err)
		}

		// 3. Descifrar el contenido usando la clave privada del usuario
		// Decrypt toma un io.Reader y N "Identidades" (llaves privadas)
		reader, err := age.Decrypt(bytes.NewReader(encryptedData), identity)
		if err != nil {
			log.Fatalf("❌ Error descifrando archivo. Probablemente tu clave pública no está autorizada en este proyecto o alguien revocó tu acceso: %v", err)
		}

		// Leemos todo el stream desencriptado
		decryptedContent, err := io.ReadAll(reader)
		if err != nil {
			log.Fatalf("❌ Error extrayendo contenido descifrado: %v", err)
		}

		// 4. Guardarlo como .env local
		envFile := ".env"
		
		// Advertir si ya existe un .env
		if _, err := os.Stat(envFile); err == nil {
			fmt.Println("⚠️  Atención: Ya tenés un archivo .env local. Se va a sobrescribir.")
		}

		if err := os.WriteFile(envFile, decryptedContent, 0600); err != nil {
			log.Fatalf("❌ Error escribiendo el archivo .env final: %v", err)
		}

		fmt.Printf("✅ ¡Éxito! Variables de entorno descifradas y guardadas en '%s'.\n", envFile)
		fmt.Println("🎉 Ya podés correr tu proyecto con las variables del equipo.")
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
