package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Comandos auxiliares para Git hooks",
}

var preCommitNoGit bool

var preCommitCmd = &cobra.Command{
	Use:   "pre-commit [paths...]",
	Short: "Bloquea commits que incluyan archivos .env en texto plano",
	Run: func(cmd *cobra.Command, args []string) {
		paths, err := resolvePreCommitPaths(preCommitNoGit, args)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "❌ Error ejecutando git diff --cached --name-only: %v\n", err)
			os.Exit(1)
		}

		if hasEnvFile(paths) {
			errOut := cmd.ErrOrStderr()
			fmt.Fprintln(errOut, "🚨 ERROR: Estás intentando commitear un archivo .env en texto plano.")
			fmt.Fprintln(errOut, "🚨 Usa 'env-vault push' para cifrarlo y agregarlo al vault.")
			fmt.Fprintln(errOut, "🚨 Abortando commit...")
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func resolvePreCommitPaths(noGit bool, args []string) ([]string, error) {
	if noGit {
		return args, nil
	}

	output, err := exec.Command("git", "diff", "--cached", "--name-only").CombinedOutput()
	if err != nil {
		trimmed := strings.TrimSpace(string(output))
		if trimmed == "" {
			return nil, err
		}
		return nil, fmt.Errorf("%v: %s", err, trimmed)
	}

	return splitLines(string(output)), nil
}

func splitLines(output string) []string {
	cleaned := strings.ReplaceAll(output, "\r\n", "\n")
	cleaned = strings.TrimSpace(cleaned)
	if cleaned == "" {
		return []string{}
	}

	lines := strings.Split(cleaned, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return lines
}

func hasEnvFile(paths []string) bool {
	for _, path := range paths {
		cleaned := strings.TrimSpace(path)
		if cleaned == "" {
			continue
		}
		if strings.HasSuffix(strings.ToLower(cleaned), ".env") {
			return true
		}
	}
	return false
}

func init() {
	preCommitCmd.Flags().BoolVar(&preCommitNoGit, "no-git", false, "No ejecutar git diff --cached --name-only; usar paths pasados como argumentos")
	hookCmd.AddCommand(preCommitCmd)
	rootCmd.AddCommand(hookCmd)
}
