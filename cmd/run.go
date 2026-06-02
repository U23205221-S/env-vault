package cmd

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"filippo.io/age"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run -- <comando>",
	Short: "Ejecuta un comando inyectando las variables de entorno en memoria",
	Long: `run descifra el archivo .env-vault/development.enc en memoria
e inyecta las variables directamente en el proceso del comando especificado,
sin escribir el archivo .env en el disco.

Ejemplo:
  env-vault run -- npm run dev
  env-vault run -- go run main.go`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Obtener la clave privada del usuario
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("❌ Error obteniendo directorio home: %v", err)
		}

		keyFile := filepath.Join(homeDir, ".env-vault", "keys", "identity.txt")
		keyData, err := os.ReadFile(keyFile)
		if err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("❌ Error: No se encontró tu clave privada. Ejecutá 'env-vault generate-key'.")
			}
			log.Fatalf("❌ Error leyendo clave privada: %v", err)
		}

		secretKey := strings.TrimSpace(string(keyData))
		identity, err := age.ParseX25519Identity(secretKey)
		if err != nil {
			log.Fatalf("❌ Error parseando tu clave privada: %v", err)
		}

		// 2. Leer el archivo cifrado
		vaultDir := ".env-vault"
		encryptedFile := filepath.Join(vaultDir, "development.enc")
		encryptedData, err := os.ReadFile(encryptedFile)
		if err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("❌ Error: No se encontró %s.", encryptedFile)
			}
			log.Fatalf("❌ Error leyendo archivo cifrado: %v", err)
		}

		// 3. Descifrar en memoria (RAM)
		reader, err := age.Decrypt(bytes.NewReader(encryptedData), identity)
		if err != nil {
			log.Fatalf("❌ Error descifrando archivo. Probablemente no estás autorizado: %v", err)
		}

		decryptedContent, err := io.ReadAll(reader)
		if err != nil {
			log.Fatalf("❌ Error extrayendo contenido descifrado: %v", err)
		}

		// 4. Parsear las variables de entorno (KEY=VALUE)
		lines := strings.Split(string(decryptedContent), "\n")
		envVars := []string{}
		
		for _, line := range lines {
			line = strings.TrimSpace(line)
			// Ignorar líneas vacías o comentarios (#)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			// Verificar que tenga formato KEY=VALUE
			if strings.Contains(line, "=") {
				envVars = append(envVars, line)
			}
		}

		// 5. Preparar el comando a ejecutar.
		//
		// En Linux/macOS invocamos el binario directamente con sus argumentos
		// (path lookup + execve). En Windows hacemos lo mismo pero, si la
		// plataforma requiere un interprete de comandos para los builtins
		// de la shell (por ejemplo `dir`, `type`, `copy`), envolvemos la
		// llamada con `cmd /c` para que el usuario pueda usar builtins sin
		// tener que escribirlos a mano.
		commandName := args[0]
		commandArgs := args[1:]

		var execCmd *exec.Cmd
		if runtime.GOOS == "windows" {
			// `cmd /c <comando>` recibe todo el comando como un único
			// string; pasamos commandName + argumentos como un solo
			// argumento, que cmd.exe se encarga de parsear.
			execCmd = exec.Command("cmd", append([]string{"/c", commandName}, commandArgs...)...)
		} else {
			execCmd = exec.Command(commandName, commandArgs...)
		}

		// Heredar las variables del sistema operativo + inyectar las nuestras
		execCmd.Env = append(os.Environ(), envVars...)

		// Conectar la entrada/salida estándar para que parezca que el comando corre nativamente
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Stdin = os.Stdin

		// 6. Ejecutar el comando
		if err := execCmd.Run(); err != nil {
			// Si el comando hijo falla (ej. npm run dev falla), salimos con error
			os.Exit(execCmd.ProcessState.ExitCode())
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
