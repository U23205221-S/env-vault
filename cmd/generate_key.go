package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"filippo.io/age"
	"github.com/spf13/cobra"
)

var generateKeyCmd = &cobra.Command{
	Use:   "generate-key",
	Short: "Genera un par de claves (Pública/Privada) para el usuario",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Obtener el directorio "Home" del usuario (~/)
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("❌ Error obteniendo el directorio home: %v", err)
		}

		// 2. Crear la carpeta ~/.env-vault/keys si no existe
		vaultDir := filepath.Join(homeDir, ".env-vault", "keys")
		err = os.MkdirAll(vaultDir, 0700) // 0700: solo el usuario puede leer/escribir/entrar
		if err != nil {
			log.Fatalf("❌ Error creando directorio %s: %v", vaultDir, err)
		}

		keyFile := filepath.Join(vaultDir, "identity.txt")

		// 3. Chequear si ya existe una llave para no pisarla por error
		if _, err := os.Stat(keyFile); err == nil {
			fmt.Println("⚠️  Ya existe una clave privada en tu máquina.")
			fmt.Printf("📂 Ruta: %s\n", keyFile)
			fmt.Println("Si querés generar una nueva, borrá ese archivo primero.")
			return
		}

		// 4. Generar la Identidad (Clave Privada y Pública)
		fmt.Println("⏳ Generando par de claves X25519 (age)...")
		identity, err := age.GenerateX25519Identity()
		if err != nil {
			log.Fatalf("❌ Error generando clave: %v", err)
		}

		// 5. Guardar la Clave Privada en el archivo con permisos estrictos (0600)
		err = os.WriteFile(keyFile, []byte(identity.String()+"\n"), 0600) // 0600: solo lectura/escritura para el dueño
		if err != nil {
			log.Fatalf("❌ Error guardando la clave privada: %v", err)
		}

		// 6. Mostrar la Clave Pública al usuario
		publicKey := identity.Recipient().String()

		fmt.Println("\n✅ ¡Claves generadas con éxito!")
		fmt.Println("--------------------------------------------------")
		fmt.Printf("🔒 Clave Privada guardada en:\n   %s\n", keyFile)
		fmt.Println("   (¡NUNCA compartas este archivo con nadie!)")
		fmt.Println("--------------------------------------------------")
		fmt.Printf("🔓 Tu Clave Pública es:\n   %s\n", publicKey)
		fmt.Println("   (Copiá esta clave pública y mandala por Slack al admin de tu equipo)")
		fmt.Println("--------------------------------------------------")
	},
}

func init() {
	rootCmd.AddCommand(generateKeyCmd)
}
