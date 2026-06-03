package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// Banner ASCII de la herramienta. Se muestra en --help y al ejecutar
// el binario sin argumentos.
const banner = "                                          ____\n" +
	"  ___  ____ _   __     _   ______ ___  __/ / /_\n" +
	" / _ \\/ __ \\ | / /____| | / / __ `/ / / / / __/\n" +
	"/  __/ / / / |/ /_____/ |/ / /_/ / /_/ / / /_  \n" +
	"\\___/_/ /_/|___/      |___/\\__,_/\\__,_/_/\\__/  "

// ANSI color codes. Se aplican solo cuando stdout es una terminal
// interactiva, para no contaminar logs o redirecciones (>, |, etc.).
const (
	colorReset   = "\033[0m"
	colorBold    = "\033[1m"
	colorDim     = "\033[2m"
	colorCyan    = "\033[36m"
	colorMagenta = "\033[35m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorRed     = "\033[31m"
	colorOrange  = "\033[38;5;208m"
)

func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// colorize envuelve un string con códigos ANSI solo si estamos en una
// terminal interactiva.
func colorize(c, s string) string {
	if !isTerminal() {
		return s
	}
	return c + s + colorReset
}

var rootCmd = &cobra.Command{
	Use:           "env-vault",
	Short:        "Sincroniza archivos .env cifrados con tu equipo, vía Git",
	Long:         colorize(colorMagenta, banner) + "\n\nSincroniza archivos .env cifrados con tu equipo, vía Git.\nServerless. Sin servidores. Sin pasar secretos por Slack.",
	Version:      "1.3.0",
	SilenceUsage: true,
}

// Execute corre el comando principal de la CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpTemplate(helpTemplate())
	rootCmd.SetUsageTemplate(usageTemplate())
}

// helpTemplate estiliza la salida de --help. El banner ya viene en
// rootCmd.Long, así que acá no lo repetimos.
func helpTemplate() string {
	orange := colorOrange
	bold := colorBold
	dim := colorDim
	cyan := colorCyan

	out := colorize(bold+cyan, "Sincroniza archivos .env cifrados con tu equipo, vía Git") + "\n" +
		colorize(dim, "Versión {{.Version}}") + "\n\n" +
		"{{.Long}}" + "\n\n" +
		colorize(orange+bold, "Comandos:") + "\n" +
		"{{range .Commands}}{{if not .Hidden}}" +
		"  " + colorize(orange, "{{rpad .Name .NamePadding}}") + "{{.Short}}\n" +
		"{{end}}{{end}}\n" +
		colorize(orange+bold, "Flags:") + "\n" +
		"{{.Flags.FlagUsages | trimTrailingWhitespaces}}\n\n" +
		colorize(dim, "Usá \"{{.CommandPath}} [comando] --help\" para más información sobre un comando.") + "\n"

	if runtime.GOOS == "darwin" {
		out += colorize(colorYellow, "⚠ Aviso macOS: este binario no está firmado. Si Gatekeeper lo bloquea, corré: xattr -d com.apple.quarantine $(which env-vault)") + "\n"
	}

	return out
}

// usageTemplate estiliza la salida cuando se escriben argumentos
// incorrectos.
func usageTemplate() string {
	orange := colorOrange
	bold := colorBold
	dim := colorDim
	red := colorRed

	return colorize(red+bold, "Uso:") + " {{.UseLine}}\n" +
		"{{if .HasAvailableSubCommands}}" +
		colorize(orange+bold, "Comandos disponibles:") + " {{.CommandPath}} [comando]\n" +
		"{{end}}" +
		"{{if .HasAvailableLocalFlags}}" +
		colorize(orange+bold, "Flags:") + "\n" +
		"{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}\n" +
		"{{end}}" +
		"{{if .HasAvailableInheritedFlags}}" +
		colorize(orange+bold, "Flags globales:") + "\n" +
		"{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}\n" +
		"{{end}}" +
		colorize(dim, "Para ayuda: "+os.Args[0]+" --help") + "\n"
}
