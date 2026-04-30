# Product Requirements Document (PRD): env-vault-sync

## 1. Propósito y Visión
**El problema:** Los equipos comparten variables de entorno (`.env`) copiando y pegando por Slack (inseguro y caótico) o usando gestores de contraseñas de terceros (fricción alta y fuera del flujo de trabajo).

**La solución:** Una CLI *serverless* que usa Git como capa de transporte. Permite a los equipos sincronizar archivos `.env` directamente en el repositorio de manera 100% cifrada, usando criptografía asimétrica para que nadie tenga que compartir secretos por chats, y permitiendo inyección en memoria.

## 2. Arquitectura y Seguridad
El proyecto NO usará un servidor centralizado (VPS). Toda la lógica de sincronización ocurrirá a través del repositorio de Git.

**Principios Criptográficos:**
- **Criptografía Asimétrica:** Cada desarrollador tiene un par de claves (Pública/Privada).
- **Algoritmo Recomendado:** X25519 (a través de herramientas estándar como `age` o `libsodium`). NUNCA reinventaremos algoritmos criptográficos.
- **Revocación Granular:** El archivo cifrado se genera para un conjunto específico de claves públicas. Quitar una clave pública del manifiesto y volver a hacer push revoca el acceso a ese desarrollador para los futuros cambios.

## 3. Flujos de Usuario (User Journeys)

### Flujo 1: Inicialización del Proyecto (Admin)
1. El líder técnico clona el repo y ejecuta `env-vault init`.
2. La CLI genera su par de claves (si no las tenía), crea el archivo `.env-vault/manifest.json` (que guarda las claves públicas autorizadas) y lo agrega al tracking de Git.
3. Configura automáticamente un `pre-commit hook` para evitar que alguien commitee un `.env` en texto plano por error.
4. Ejecuta `env-vault push`, cifrando el `.env` local usando su propia clave pública.

### Flujo 2: Onboarding de un Nuevo Developer
1. El dev nuevo instala la CLI y ejecuta `env-vault generate-key`.
2. La CLI crea la clave privada en `~/.env-vault/keys` y le muestra la clave pública por consola.
3. El dev manda su clave pública por Slack (es 100% seguro hacerlo).
4. El Admin ejecuta `env-vault add-user <clave-publica>` y hace `env-vault push`.
5. El archivo `.env` ahora está cifrado para el Admin Y para el Dev nuevo.
6. El dev nuevo hace `git pull` y luego `env-vault pull` (o `env-vault run -- npm run dev`) y ya tiene acceso.

### Flujo 3: Offboarding (Revocación)
1. Un dev se va del equipo.
2. El Admin ejecuta `env-vault remove-user <clave-publica-del-ex-dev>`.
3. El Admin regenera los secretos o actualiza variables (ej. cambia contraseñas de BD por seguridad) y hace `env-vault push`.
4. El ex-dev ya no podrá descifrar los nuevos `.env-vault/*.enc` generados.

## 4. Interfaz de Línea de Comandos (CLI)
- `env-vault generate-key`: Genera y guarda localmente el par de claves del usuario.
- `env-vault init`: Inicializa el vault en el repo actual (crea el directorio `.env-vault` y los hooks).
- `env-vault add-user <pub-key>`: Agrega una clave pública al manifiesto de usuarios permitidos.
- `env-vault remove-user <pub-key>`: Elimina una clave pública del manifiesto.
- `env-vault push [-e environment]`: Lee el `.env` local (o `.env.production`), lo cifra con todas las claves públicas del manifiesto y lo guarda en `.env-vault/development.enc`.
- `env-vault pull [-e environment]`: Toma el `.env-vault/*.enc`, lo descifra usando la clave privada local y genera el archivo `.env`.
- `env-vault run [-e environment] -- <comando>`: Descifra el archivo `.env` e inyecta las variables **directamente en memoria** para el comando especificado, sin tocar el disco.
