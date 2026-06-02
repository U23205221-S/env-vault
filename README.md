# 🔐 env-vault

**Comparte archivos `.env` con tu equipo, cifrados, sin servidor central y sin pasar secretos por canales inseguros.**

---

## Instalación

Un solo comando por plataforma. Sin dependencias, sin compiladores, sin Docker.

| Plataforma | Comando |
|---|---|
| **Fedora, Arch, Ubuntu, Debian, macOS** | `curl -sSf https://raw.githubusercontent.com/U23205221-S/env-vault/main/install.sh \| sh` |
| **Windows (PowerShell)** | `irm https://raw.githubusercontent.com/U23205221-S/env-vault/main/install.ps1 \| iex` |

Para una versión específica: agregá `v1.2.0` al final del comando en Linux/macOS, o usá `$env:ENV_VAULT_VERSION = "v1.2.0"` antes del comando en Windows.

Después de instalar, abrí una terminal nueva y verificá con `env-vault --help`.

---

## ¿Qué hace?

`env-vault` cifra tu archivo `.env` con criptografía asimétrica y lo guarda cifrado dentro del propio repositorio Git del proyecto. Cada desarrollador tiene su propio par de claves. El servidor de Git nunca ve el texto plano. Nadie comparte contraseñas.

| Sin env-vault | Con env-vault |
|---|---|
| ❌ Compartir el `.env` por Slack, email o Drive — queda en texto plano y fuera de control | ✅ Cifrado en el repo, accesible solo para los miembros autorizados |
| ❌ Vault centralizado o AWS Secrets Manager — infraestructura que mantener y pagar | ✅ Serverless, sin servicios externos, sin costos recurrentes |
| ❌ 1Password / LastPass — copiar y pegar manualmente cada vez que hay un cambio | ✅ Sincronizado vía Git, igual que cualquier archivo del proyecto |
| ❌ `.env` commiteado por accidente — un dev sube sus claves a un repo público | ✅ Hook `pre-commit` bloquea automáticamente el commit de `.env` sin cifrar |

---

## 🛡️ Seguridad

`env-vault` está construido sobre primitivas criptográficas estándar y auditables:

- **Cifrado asimétrico X25519**, a través de [`filippo.io/age`](https://github.com/FiloSottile/age) (el estándar moderno de cifrado de archivos).
- **Las claves privadas nunca salen de tu máquina.** Se generan localmente con `env-vault generate-key` y se guardan en `~/.env-vault/keys/identity.txt` con permisos `0600` (solo vos podés leerlo).
- **Solo se comparten claves públicas** (`age1...`). Mandarlas por Slack o email es seguro: sirven para cifrar, no para descifrar.
- **El servidor de Git ve blobs cifrados**, no texto plano. Si alguien clona el repo sin estar autorizado, los archivos `.enc` son inútiles sin la clave privada.
- **Revocación granular:** quitar una clave pública del manifiesto y volver a hacer `push` invalida el acceso a ese desarrollador para todos los cambios futuros. No hay forma de "re-abrir" un `.env` ya cifrado con la clave revocada.
- **Hook `pre-commit` instalado automáticamente** durante `env-vault init`. Bloquea `git commit` de cualquier `.env` en texto plano, incluso si el dev tiene prisa y se olvida.

El modelo de amenaza que cubre: robo accidental del repo (GitHub hackeado, dev que se va con copia, repo público por error). **No** cubre: una máquina comprometida con la clave privada adentro — para eso, usá `env-vault run` y nunca escribas el `.env` en disco.

---

## Flujo de uso

### 1. El admin inicializa el vault en el proyecto

En la raíz del repositorio Git del proyecto del equipo:

```bash
env-vault init
```

Esto crea la carpeta `.env-vault/`, el archivo `manifest.json` (lista de claves públicas autorizadas) y un hook `pre-commit` que bloquea subir un `.env` en texto plano.

### 2. Cada desarrollador genera su par de claves

```bash
env-vault generate-key
```

Guarda una clave privada local y muestra por consola tu clave pública (`age1...`). Compartís la pública con el admin por el canal que prefieras — no es un secreto.

### 3. El admin te autoriza

```bash
env-vault add-user age1...clave-que-te-mandaron...
```

Tu clave pública queda en el manifiesto. A partir de ahora podés descifrar los archivos del equipo.

### 4. Alguien modifica el `.env` y lo comparte con el equipo

```bash
env-vault push
git add .env-vault/
git commit -m "chore: actualiza variables de entorno"
git push
```

El archivo se cifra para todos los del manifiesto y se sube al repo. El servidor de Git ve bytes cifrados.

### 5. Después de hacer `git pull`, descifrá el `.env` actualizado

```bash
env-vault pull
```

Tu clave privada local descifra el archivo y crea el `.env` físico en el proyecto.

### 6. Para ejecutar el proyecto sin que el `.env` toque el disco

```bash
env-vault run -- npm run dev
# o
env-vault run -- go run main.go
```

Las variables se inyectan directamente en la memoria del proceso. El `.env` no se escribe a disco.

### 7. Para revocar el acceso a alguien que dejó el equipo

```bash
env-vault remove-user age1...clave-del-ex-dev...
env-vault push
```

La próxima vez que se cifre el `.env`, esa persona ya no podrá descifrarlo. Los `.env` que vio antes siguen siendo los mismos (no se re-cifra retroactivamente) — para mayor seguridad, rotá las credenciales.

---

## 📋 Comandos disponibles

| Comando | Para qué sirve |
|---|---|
| `env-vault init` | Inicializa el vault en el proyecto (una vez por proyecto) |
| `env-vault generate-key` | Genera tu par de claves local (una vez por máquina) |
| `env-vault add-user <pubkey>` | Agrega una clave pública al manifiesto (lo corre el admin) |
| `env-vault remove-user <pubkey>` | Quita una clave pública del manifiesto |
| `env-vault push` | Cifra el `.env` local para todos los del manifiesto |
| `env-vault pull` | Descifra el archivo del equipo y crea el `.env` local |
| `env-vault run -- <cmd>` | Descifra en memoria y ejecuta el comando con las variables inyectadas |

---

## 🖥️ Notas por plataforma

### Todas las plataformas

- **Después de instalar, abrí una terminal nueva** antes de ejecutar `env-vault`. El script de instalación agrega el binario al `PATH`, pero los terminales ya abiertos no recargan el `PATH` automáticamente.

### macOS — Gatekeeper la primera vez

Los binarios no están firmados por Apple (no compensa pagar la cuenta de developer solo para una CLI open source). La primera vez que ejecutes `env-vault`, macOS lo va a bloquear. Solución, una sola vez por máquina:

```bash
xattr -d com.apple.quarantine "$(which env-vault)"
```

O graficamente: clic derecho sobre el binario → Abrir → confirmar en el diálogo.

### Windows — Hook `pre-commit` necesita Git for Windows

El hook `pre-commit` es un script de shell (convención de Git). Git for Windows incluye `sh.exe` y lo usa para ejecutarlo. Si tu Git viene de otra fuente, instalá [Git for Windows](https://git-scm.com/download/win) (es el estándar en Windows).

---

## ✅ Plataformas soportadas

- **Linux:** Fedora, Arch, Ubuntu, Debian y cualquier distro moderna con `bash`, `curl`, `tar` y `sha256sum` (casi todas).
- **macOS:** Intel y Apple Silicon (M1, M2, M3, M4).
- **Windows:** 10 y 11 con PowerShell 5.1 o superior (ya viene preinstalado en ambos).
