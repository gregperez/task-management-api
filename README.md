# Task Management API

API REST en Golang para gestión de tareas con control de estados y múltiples roles de usuario.

## 🚀 Características

- **Autenticación JWT**: Sistema seguro de autenticación con tokens
- **3 Roles de Usuario**: Administrador, Ejecutor y Auditor
- **CRUD Completo**: Gestión de usuarios y tareas
- **Control de Estados**: Flujo de estados configurado según diagrama
- **Contraseñas Temporales**: Los usuarios nuevos deben cambiar su contraseña en el primer login
- **Validaciones de Negocio**: Restricciones según rol y estado de tareas
- **Arquitectura Limpia**: Separación de capas y responsabilidades
- **Tests Unitarios**: Cobertura de casos críticos

## 🏗️ Arquitectura

> **Documentación completa**: Ver [docs/architecture.md](docs/architecture.md) para diagramas detallados y especificaciones técnicas.

### Principios de Diseño

- ✅ **Clean Architecture**: Separación clara de responsabilidades por capas
- ✅ **Dependency Injection**: Inyección de dependencias en `app.go`
- ✅ **Repository Pattern**: Abstracción de acceso a datos
- ✅ **Service Layer**: Lógica de negocio centralizada
- ✅ **DTO Pattern**: Separación entre modelos de dominio y API
- ✅ **Middleware Chain**: Autenticación, logging y CORS
- ✅ **State Machine**: Control estricto de transiciones de estado
- ✅ **Role-Based Access Control**: Autorización por rol de usuario

## 📋 Requisitos Previos

- Go 1.21+
- PostgreSQL 15+
- Docker y Docker Compose (opcional)

## 🔧 Instalación

### Opción 1: Docker (Recomendado)

```bash
# Clonar el repositorio
git clone https://github.com/yourusername/task-api.git
cd task-api

# Iniciar servicios
Unix -> ./tasks.sh docker-up
Windows -> .\tasks.bat docker-up

# Ver logs
Unix -> ./tasks.sh docker-logs
Windows -> .\tasks.bat docker-logs
```

### Opción 2: Local

```bash
# Clonar el repositorio
git clone https://github.com/yourusername/task-api.git
cd task-api

# Instalar dependencias
go mod download

# Configurar variables de entorno
cp .env.example .env

# Iniciar PostgreSQL localmente

# Ejecutar migraciones
Unix -> ./tasks.sh migrate-up
Windows -> .\tasks.bat migrate-up

# Iniciar servidor
Unix -> ./tasks.sh run
Windows -> .\tasks.bat run
```

## 🔐 Autenticación

### Login
```bash
POST /api/v1/login
Content-Type: application/json

{
  "username": "admin",
  "password": "Admin123!"
}

# Respuesta
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "admin-001",
    "username": "admin",
    "role": "Administrador"
  },
  "must_change_password": false
}
```

### Uso del Token
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

## 📚 API Endpoints

### Autenticación

| Método | Endpoint | Descripción | Autenticado |
|--------|----------|-------------|-------------|
| POST | `/api/v1/login` | Login de usuario | No |
| POST | `/api/v1/logout` | Logout de usuario | Sí |
| POST | `/api/v1/change-password` | Cambiar contraseña | Sí |

### Usuarios (Solo Administrador)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/v1/users` | Crear usuario |
| GET | `/api/v1/users` | Listar usuarios |
| GET | `/api/v1/users/{id}` | Obtener usuario |
| PUT | `/api/v1/users/{id}` | Actualizar usuario |
| DELETE | `/api/v1/users/{id}` | Eliminar usuario |

### Tareas

| Método | Endpoint | Descripción | Roles |
|--------|----------|-------------|-------|
| POST | `/api/v1/tasks` | Crear tarea | Admin |
| PUT | `/api/v1/tasks/{id}` | Actualizar tarea | Admin |
| DELETE | `/api/v1/tasks/{id}` | Eliminar tarea | Admin |
| GET | `/api/v1/tasks/my-tasks` | Listar mis tareas | Ejecutor |
| PATCH | `/api/v1/tasks/{id}/status` | Actualizar estado | Ejecutor |
| POST | `/api/v1/tasks/{id}/comments` | Agregar comentario | Ejecutor |
| GET | `/api/v1/tasks/all` | Listar todas las tareas | Auditor |
| GET | `/api/v1/tasks/{id}` | Obtener tarea | Todos |

## 🎯 Ejemplos de Uso

### Crear Usuario (Admin)
```bash
POST /api/v1/users
Authorization: Bearer {token}
Content-Type: application/json

{
  "username": "executor1",
  "email": "executor1@example.com",
  "role": "Ejecutor"
}

# Respuesta
{
  "user": {
    "id": "uuid-here",
    "username": "executor1",
    "email": "executor1@example.com",
    "role": "Ejecutor",
    "is_temporary_password": true
  },
  "temporary_password": "TempPass123!"
}
```

### Crear Tarea (Admin)
```bash
POST /api/v1/tasks
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Desarrollar módulo de reportes",
  "description": "Crear sistema de generación de reportes en PDF",
  "due_date": "2026-02-01T23:59:59Z",
  "assigned_to": "executor-uuid"
}
```

### Actualizar Estado de Tarea (Ejecutor)
```bash
PATCH /api/v1/tasks/{task-id}/status
Authorization: Bearer {token}
Content-Type: application/json

{
  "status": "Iniciado"
}
```

### Agregar Comentario a Tarea Vencida (Ejecutor)
```bash
POST /api/v1/tasks/{task-id}/comments
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "La tarea se retrasó debido a problemas técnicos con el servidor"
}
```

## 🧪 Testing

```bash
# Ejecutar todos los tests
Unix -> ./tasks.sh test
Windows -> .\tasks.bat test

# Ver cobertura
Unix -> ./tasks.sh test-coverage
Windows -> .\tasks.bat test-coverage

# Ejecutar linter
Unix -> ./tasks.sh lint
Windows -> .\tasks.bat lint
```

## 🔒 Seguridad

### Implementadas

- **Autenticación JWT**: Tokens firmados con HS256
- **Hash de Contraseñas**: Bcrypt con cost 10
- **Validación de Entrada**: Sanitización de datos
- **RBAC**: Control de acceso basado en roles
- **SQL Injection Protection**: Uso de prepared statements
- **Rate Limiting**: Middleware de recuperación de panics

### Recomendaciones para Producción

- Usar HTTPS/TLS
- Implementar rate limiting por IP
- Agregar registro de auditoría
- Rotar secrets regularmente
- Implementar 2FA
- Usar variables de entorno seguras

## 📊 Modelos de Datos

### Usuario
```go
type User struct {
    ID                 string
    Username           string
    Email              string
    Password           string (hashed)
    Role               UserRole (Administrador|Ejecutor|Auditor)
    IsTemporaryPassword bool
    CreatedAt          time.Time
    UpdatedAt          time.Time
}
```

### Tarea
```go
type Task struct {
    ID          string
    Title       string
    Description string
    Status      TaskStatus
    DueDate     time.Time
    AssignedTo  string (user_id)
    CreatedBy   string (user_id)
    Comments    []Comment
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

## 🚦 Estados de Tarea

- **Asignado**: Estado inicial al crear la tarea
- **Iniciado**: El ejecutor comenzó a trabajar
- **En espera**: Trabajo pausado temporalmente
- **Finalizado Éxito**: Tarea completada exitosamente
- **Finalizado Error**: Tarea completada con errores

## 📝 Reglas de Negocio

### Usuarios
- Solo administradores pueden crear usuarios
- No se pueden crear más administradores
- Los usuarios nuevos reciben contraseña temporal
- La contraseña debe cambiarse en el primer login
- Solo se pueden crear usuarios tipo Ejecutor o Auditor

### Tareas
- Solo administradores pueden crear/actualizar/eliminar tareas
- Solo se pueden modificar tareas en estado "Asignado"
- Los ejecutores solo pueden actualizar estado de sus tareas
- No se puede actualizar estado de tareas vencidas
- Solo se pueden comentar tareas vencidas
- Las transiciones de estado deben seguir el flujo definido

### Permisos por Rol

**Administrador**:
- Gestión completa de usuarios (CRUD)
- Gestión completa de tareas (CRUD)
- Asignar tareas a ejecutores

**Ejecutor**:
- Ver sus tareas asignadas
- Actualizar estado de sus tareas
- Agregar comentarios a tareas vencidas

**Auditor**:
- Ver todas las tareas del sistema
- Ver estado de cualquier tarea

## 🐛 Manejo de Errores

La API retorna códigos HTTP estándar:

- `200 OK`: Operación exitosa
- `201 Created`: Recurso creado
- `400 Bad Request`: Datos inválidos
- `401 Unauthorized`: No autenticado
- `403 Forbidden`: Sin permisos
- `404 Not Found`: Recurso no encontrado
- `500 Internal Server Error`: Error del servidor

Formato de error:
```json
{
  "error": "descripción del error"
}
```

## 📦 Dependencias Principales

- `chi/v5`: Router HTTP
- `pgx/v5`: Driver PostgreSQL
- `jwt/v5`: JSON Web Tokens
- `bcrypt`: Hash de contraseñas
- `testify`: Testing framework
- `uuid`: Generación de UUIDs

## 🤝 Contribución

1. Fork el proyecto
2. Crear rama feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crear Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT.

## 👥 Autor

Greg Perez - [@gregperez](https://github.com/gregperez)

## 🙏 Agradecimientos

- Equipo de Ionix por la oportunidad
- Comunidad de Go por las excelentes librerías