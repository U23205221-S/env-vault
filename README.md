# 🔐 env-vault

**Sincroniza archivos `.env` entre miembros de tu equipo de forma 100% cifrada usando Git. Sin servidores, sin compartir secretos por Slack, y con criptografía de grado militar.**

---

## 🤔 El Problema
¿Cómo comparte tu equipo las variables de entorno locales (`.env`)?
- ❌ **Por Slack/WhatsApp:** Quedan registradas para siempre en texto plano.
- ❌ **Gestores de contraseñas (1Password/LastPass):** Lento, requiere copiar y pegar manualmente cada vez que hay un cambio.
- ❌ **Servidores centralizados (Vault/AWS Secrets):** Requiere infraestructura, mantenimiento y pagar licencias.

## 💡 La Solución: GitOps + Criptografía Asimétrica
`env-vault` es una CLI *serverless* escrita en Go. 
Usa **criptografía asimétrica (X25519)** a través de la librería estándar `age`. Cada desarrollador tiene su propia clave privada local. El archivo `.env` se cifra con las claves públicas de todo el equipo y se sube directamente al repositorio Git (`.env-vault/development.enc`).

Nadie comparte contraseñas. El servidor de Git nunca ve el texto plano.

---

## 🚀 Instalación

### Opción 1: Compilar desde el código fuente (Recomendado para Linux/macOS)
Requiere tener [Go 1.23+](https://golang.org/doc/install) instalado.

```bash
git clone https://github.com/U23205221-S/env-vault.git
cd env-vault
go build -o env-vault main.go
sudo mv env-vault /usr/local/bin/
```

### Opción 2: Usar Make
```bash
make build
sudo mv bin/env-vault /usr/local/bin/
```

### Windows
Requiere tener [Go 1.23+](https://golang.org/doc/install) instalado.

**Instalar con Go:**
```powershell
go install github.com/U23205221-S/env-vault@latest
```

**Compilar desde el repo (Git Bash recomendado):**
```powershell
git clone https://github.com/U23205221-S/env-vault.git
cd env-vault
make build
```

**Notas Windows:**
- Para `env-vault run --`, usá `env-vault run -- cmd /c dir` (u otro comando de `cmd`).
- El hook `pre-commit` requiere Git for Windows (incluye `sh.exe`, que Git usa por defecto para hooks).

---

## 🛠️ Flujo de Trabajo del Equipo (Workflow)

### 1. El Admin inicializa el proyecto
En la raíz de tu repositorio Git:
```bash
env-vault init
```
*Esto crea la carpeta `.env-vault/`, el `manifest.json` y un hook de Git (`pre-commit`) que bloquea accidentalmente subir un `.env` en texto plano.*

### 2. El nuevo desarrollador genera sus llaves
En la máquina del desarrollador nuevo:
```bash
env-vault generate-key
```
*Esto guarda una clave privada en `~/.env-vault/keys/` y le muestra una Clave Pública (`age1...`). El dev manda esta clave pública por Slack al Admin (es 100% seguro).*

### 3. El Admin autoriza al desarrollador
```bash
env-vault add-user age1...la-clave-del-dev...
```

### 4. Empaquetar y Cifrar (Push)
Cuando alguien modifica el `.env` local y quiere compartirlo con el equipo:
```bash
env-vault push
git add .env-vault/
git commit -m "chore: actualiza variables de entorno"
git push
```
*El archivo se cifra para todos los usuarios del manifiesto.*

### 5. Recuperar y Descifrar (Pull)
Cuando un dev hace `git pull` y ve que hay cambios en las variables:
```bash
env-vault pull
```
*La herramienta usa la clave privada local del dev para abrir el archivo y crear el `.env` físico.*

### 6. Modo "Zero Trust" (Run)
Si no querés que el `.env` toque el disco duro por seguridad, podés inyectar las variables directamente en la memoria RAM de tu proceso:
```bash
env-vault run -- npm run dev
# o
env-vault run -- go run main.go
```

### 7. Revocar accesos (Offboarding)
Si un dev se va del equipo:
```bash
env-vault remove-user age1...clave...
env-vault push
```
*Su acceso queda revocado inmediatamente para todos los futuros cambios.*

---

## 🛡️ Seguridad
- **Algoritmo:** X25519 (Curva Elíptica).
- **Motor:** [filippo.io/age](https://github.com/FiloSottile/age) (Estándar moderno de encriptación asimétrica).
- **Protección de fugas:** El hook `pre-commit` impide que los humanos cometan el error de hacer `git commit` de un archivo `.env` sin cifrar.

---

## 🧪 Platform Support
Probado en Windows, macOS, Fedora y Arch Linux. Debería funcionar en cualquier distro moderna con Go 1.23+.
