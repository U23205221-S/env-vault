# 🔐 env-vault

**Sincroniza archivos `.env` entre los miembros de tu equipo de forma 100 % cifrada usando Git. Sin servidores, sin compartir secretos por Slack y con criptografía de grado militar.**

---

## 🤔 El problema

¿Cómo comparte tu equipo las variables de entorno locales (`.env`)?

- ❌ **Por Slack o WhatsApp:** quedan registradas para siempre en texto plano.
- ❌ **Gestores de contraseñas (1Password, LastPass):** son lentos y obligan a copiar y pegar manualmente cada vez que hay un cambio.
- ❌ **Servidores centralizados (Vault, AWS Secrets):** requieren infraestructura, mantenimiento y pagar licencias.

## 💡 La solución: GitOps + criptografía asimétrica

`env-vault` es una CLI *serverless* escrita en Go. Usa **criptografía asimétrica (X25519)** a través de la biblioteca estándar `age`. Cada desarrollador tiene su propia clave privada local. El archivo `.env` se cifra con las claves públicas de todo el equipo y se sube directamente al repositorio Git (`.env-vault/development.enc`).

Nadie comparte contraseñas. El servidor de Git nunca ve el texto plano.

---

## 🚀 Instalación

### Instalación rápida (Linux y macOS)

```bash
curl -sSf https://raw.githubusercontent.com/U23205221-S/env-vault/main/install.sh | sh
```

El script detecta el sistema operativo y la arquitectura, descarga el binario correspondiente desde la página de [Releases](https://github.com/U23205221-S/env-vault/releases), verifica el SHA256 y lo instala en `/usr/local/bin` (o en `~/.local/bin` si no tienes permisos de escritura en `/usr/local/bin`).

Para instalar una versión específica:

```bash
curl -sSf https://raw.githubusercontent.com/U23205221-S/env-vault/main/install.sh | sh -s -- v1.1.1
```

### Windows

Los binarios se publican en [Releases](https://github.com/U23205221-S/env-vault/releases) como `.zip`. Sigue estos pasos:

1. Descarga el `.zip` correspondiente a tu arquitectura (`windows-amd64` para la mayoría de los equipos, `windows-arm64` para Surface Pro X y máquinas con chip ARM).
2. Extrae el `.exe` en una carpeta permanente (por ejemplo, `C:\tools\env-vault\`).
3. Agrega esa carpeta a la variable de entorno `PATH` (Panel de control → Sistema → Configuración avanzada → Variables de entorno).
4. Abre una terminal nueva y ejecuta `env-vault --help` para confirmar.

**Notas sobre Windows:**

- El hook `pre-commit` requiere Git for Windows, que incluye `sh.exe`. Git usa este intérprete por defecto al ejecutar los hooks.
- El comando `env-vault run -- <comando>` invoca automáticamente el shell del sistema, por lo que puedes usar builtins como `dir` o `type` directamente, sin prefijar con `cmd /c`.

### Desde el código fuente (avanzado)

Requiere [Go 1.24+](https://golang.org/doc/install).

```bash
git clone https://github.com/U23205221-S/env-vault.git
cd env-vault
go build -o env-vault main.go
sudo mv env-vault /usr/local/bin/
```

O usando Make:

```bash
make build
sudo mv bin/env-vault /usr/local/bin/
```

### Verificación de descargas

Cada release publica un archivo `checksums.txt` con los SHA256 de todos los artefactos. El script de instalación lo verifica de forma automática. Si descargas el binario a mano, puedes verificarlo con:

```bash
sha256sum -c checksums.txt
```

---

## 🛠️ Flujo de trabajo del equipo

### 1. El administrador inicializa el proyecto

En la raíz del repositorio Git del proyecto:

```bash
env-vault init
```

Esto crea la carpeta `.env-vault/`, el archivo `manifest.json` y un hook de Git (`pre-commit`) que bloquea la subida accidental de un `.env` en texto plano.

### 2. El nuevo desarrollador genera sus claves

En la máquina del desarrollador:

```bash
env-vault generate-key
```

Esto guarda una clave privada en `~/.env-vault/keys/identity.txt` y muestra por consola una clave pública (`age1...`). El desarrollador envía esta clave pública por Slack al administrador (es seguro hacerlo porque la clave es pública).

### 3. El administrador autoriza al desarrollador

```bash
env-vault add-user age1...clave-del-desarrollador...
```

### 4. Empaquetar y cifrar (push)

Cuando alguien modifica el `.env` local y quiere compartirlo con el equipo:

```bash
env-vault push
git add .env-vault/
git commit -m "chore: actualiza variables de entorno"
git push
```

El archivo se cifra para todos los usuarios del manifiesto.

### 5. Recuperar y descifrar (pull)

Cuando un desarrollador hace `git pull` y ve cambios en las variables:

```bash
env-vault pull
```

La herramienta usa la clave privada local del desarrollador para descifrar el archivo y crear el `.env` físico.

### 6. Modo "Zero Trust" (run)

Si no quieres que el `.env` toque el disco duro, puedes inyectar las variables directamente en la memoria del proceso:

```bash
env-vault run -- npm run dev
# o
env-vault run -- go run main.go
```

En Windows, este comando invoca el shell del sistema de forma automática, por lo que puedes usar builtins como `dir` sin prefijar con `cmd /c`.

### 7. Revocar accesos (offboarding)

Si un desarrollador deja el equipo:

```bash
env-vault remove-user age1...clave...
env-vault push
```

Su acceso queda revocado de inmediato para todos los cambios futuros.

---

## 🛡️ Seguridad

- **Algoritmo:** X25519 (curva elíptica).
- **Motor:** [filippo.io/age](https://github.com/FiloSottile/age), estándar moderno de cifrado asimétrico.
- **Protección de fugas:** el hook `pre-commit` impide que alguien haga `git commit` de un archivo `.env` sin cifrar por accidente.

---

## 🧪 Plataformas soportadas

Probado en Windows, macOS, Fedora y Arch Linux. La herramienta funciona en cualquier distribución moderna con Go 1.24 o superior.
