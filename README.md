# gh-flow

CLI tool para automatizar flujos de trabajo Git en sistemas con múltiples repositorios.

## Instalación

```bash
# Compilar localmente
go build -o gh-flow .

# O instalar globalmente:
go install .
```

## Uso

### Preparar repositorios para desarrollo

```bash
gh-flow start
```

Este comando:
1. Detecta todos los repositorios git en subdirectorios
2. Verifica cambios sin commitear
3. Permite seleccionar repositorios (TUI interactivo)
4. Pregunta por nivel de urgencia (main vs dev)
5. Configura el tipo y nombre de rama
6. Ejecuta: checkout → pull → create branch

### Finalizar desarrollo

```bash
gh-flow finish
```

Este comando:
1. Detecta repositorios con cambios
2. Permite seleccionar repositorios
3. Pide mensaje de commit
4. Ejecuta: add → commit → push
5. Crea PRs automáticamente con `gh`
6. Retorna a la rama original

## Flujo de trabajo

### Caso 1: Feature normal (rama dev)

```bash
# Preparar
gh-flow start
# Seleccionar: normal, feat, nombre-del-feature
# → Crea feat/nombre-del-feature desde dev

# Desarrollar...

# Finalizar
gh-flow finish
# → Commit, push, PR a dev
```

### Caso 2: Hotfix urgente (rama main)

```bash
# Preparar
gh-flow start
# Seleccionar: urgent, hotfix, nombre-del-fix
# → Crea hotfix/nombre-del-fix desde main

# Desarrollar...

# Finalizar
gh-flow finish
# → Commit, push, PR a main y dev
```

## Makefile

El proyecto incluye un Makefile con varios targets útiles:

```bash
# Ver todos los targets disponibles
make help

# Compilar el binario
make build

# Ejecutar tests
make test

# Limpiar archivos compilados
make clean

# Instalar globalmente
make install

# Formatear código
make fmt

# Ejecutar linter
make lint

# Actualizar dependencias
make tidy
```

## Testing Manual

Para probar la herramienta sin afectar repositorios reales:

### 1. Crear repositorios de prueba

```bash
make setup-test
```

Esto crea:
```
test-repos/
├── frontend/    (repo git con ramas main y dev)
└── backend/     (repo git con ramas main y dev)
```

### 2. Probar comando `start`

```bash
make test-start
```

O manualmente:
```bash
cd test-repos && ../gh-flow start
```

**Flujo de prueba:**
1. Seleccionar repositorios (Space para marcar, Enter para continuar)
2. Elegir "normal" o "urgent"
3. Seleccionar tipo de rama (feat/hotfix)
4. Ingresar nombre de la rama
5. Verificar que se crearon las ramas: `git branch` en cada repo

### 3. Probar con cambios sin commitear

```bash
cd test-repos/frontend
echo "// test" >> README.md
cd ../..
make test-start
# Debería detectar cambios y ofrecer opciones de stash
```

### 4. Probar comando `finish`

```bash
make test-finish
```

**Flujo de prueba:**
1. Crea cambios de prueba en los repos
2. Ejecuta `gh-flow finish`
3. Selecciona repositorios
4. Ingresa mensaje de commit
5. Verifica que se hizo commit y push

### 5. Limpiar repositorios de prueba

```bash
make clean-test
```

### Prueba completa paso a paso

```bash
# 1. Setup
make setup-test

# 2. Ir a test-repos
cd test-repos

# 3. Ejecutar start
../gh-flow start
#   - Seleccionar: frontend y backend (Space + Enter)
#   - Elegir: "normal"
#   - Seleccionar: "feat"
#   - Nombre: "test-feature"

# 4. Verificar ramas creadas
cd frontend && git branch  # Debe mostrar feat/test-feature
cd ../backend && git branch  # Debe mostrar feat/test-feature

# 5. Crear cambios
cd ../frontend
echo "Nuevo feature" >> README.md
cd ../backend
echo "Nuevo feature" >> README.md

# 6. Ejecutar finish
cd ..
../gh-flow finish
#   - Seleccionar: frontend y backend
#   - Mensaje: "feat: add test feature"

# 7. Verificar
cd frontend && git log --oneline -1  # Debe mostrar el commit
cd ../backend && git log --oneline -1

# 8. Limpiar
cd ../..
make clean-test
```

## Requisitos

- Go 1.21+
- Git
- GitHub CLI (`gh`) instalado y autenticado

## Características

- TUI elegante con Bubbletea
- Detección automática de repositorios
- Manejo de cambios sin commitear (stash)
- Selección interactiva de repositorios
- Creación automática de PRs
- Soporte para múltiples repositorios simultáneamente

## Navegación del TUI

- **↑/↓** o **j/k**: Moverse entre opciones
- **Space** o **Enter**: Seleccionar/deseleccionar
- **Ctrl+D**: Confirmar selección
- **Ctrl+C** o **q**: Salir
- **Esc**: Volver atrás (en formularios)
