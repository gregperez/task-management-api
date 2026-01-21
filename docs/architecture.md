# Arquitectura - Task Management API

## Tabla de Contenidos

1. [Diagrama de Contexto](#1-diagrama-de-contexto)
2. [Arquitectura de Capas](#2-arquitectura-de-capas)
3. [Diagrama de Componentes](#3-diagrama-de-componentes)
4. [Flujo de Autenticación](#4-flujo-de-autenticación)
5. [Máquina de Estados](#5-máquina-de-estados)
6. [Modelo de Datos](#6-modelo-de-datos)
7. [Deployment Docker](#7-deployment-docker)
8. [Decisiones de Arquitectura](#8-decisiones-de-arquitectura)

---

## 1. Diagrama de Contexto

Vista de alto nivel del sistema mostrando actores externos y sus interacciones.

<image src="diagrams/diagrama-contexto.png" alt="Diagrama de Contexto" />

### Descripción

- **Usuarios**: Sistema multi-rol (Administrador, Ejecutor, Auditor)
- **API REST**: Aplicación Go con autenticación JWT
- **PostgreSQL**: Persistencia de usuarios y tareas

---

## 2. Arquitectura de Capas

Implementación de **Clean Architecture** con separación clara de responsabilidades.

<image src="diagrams/arquitectura-capas.png" alt="Arquitectura de Capas" />

### Responsabilidades por Capa

#### Capa de Presentación (Handlers)
- Manejo de HTTP requests/responses
- Validación de entrada (formato)
- Serialización JSON
- Gestión de códigos HTTP

#### Capa de Negocio (Services)
- Lógica de negocio
- Validaciones de reglas de negocio
- Autorización por roles
- Orquestación de repositorios

#### Capa de Datos (Repositories)
- Abstracción de acceso a datos
- Queries SQL
- Mapeo de resultados a entidades
- Transacciones

#### Dominio
- Entidades de negocio
- Reglas de dominio
- Tipos y constantes
- Errores de negocio

---

## 3. Diagrama de Componentes

Vista detallada de todos los módulos y sus dependencias.

<image src="diagrams/diagrama-componentes.png" alt="Diagrama de Componentes" />

### Middlewares Implementados

1. **CORS Handler**: Permite requests cross-origin con configuración permisiva (configurable)
2. **Request Logger**: Log automático de todas las requests (método, ruta, tiempo)
3. **Panic Recovery**: Captura panics y retorna 500 sin crashear el servidor
4. **JWT Auth**: Valida tokens en rutas protegidas
5. **Request ID**: Añade ID único a cada request para trazabilidad
6. **Real IP**: Extrae la IP real del cliente detrás de proxies

---

## 4. Flujo de Autenticación

Secuencia de login y autorización de requests.

<image src="diagrams/flujo-autenticacion.png" alt="Flujo de Autenticación" />

### Detalles de Implementación

- **Hash**: bcrypt con cost factor 10
- **JWT**: HS256 con secret configurable
- **Token Expiry**: Configurable vía variable de entorno
- **Claims**: Incluye userID, role y expiration

---

## 5. Máquina de Estados

Flujo de estados de tareas con transiciones permitidas.

<image src="diagrams/maquina-estados.png" alt="Máquina de Estados" />

### Implementación

La lógica de transiciones está implementada en `internal/domain/task.go`:

```go
func (t *Task) CanTransitionTo(newStatus TaskStatus) error {
    validTransitions := map[TaskStatus][]TaskStatus{
        StatusAssigned:   {StatusInProgress},
        StatusInProgress: {StatusSuccess, StatusError, StatusPending},
        StatusPending:    {StatusInProgress, StatusSuccess, StatusError},
        StatusSuccess:    {}, // Estado final
        StatusError:      {}, // Estado final
    }
    // Validación...
}
```

---

## 6. Modelo de Datos

Esquema de base de datos relacional.

<image src="diagrams/modelo-datos.png" alt="Modelo de Datos" />

### Migraciones

Las migraciones están en `/migrations`:

- `001_create_users.up.sql`: Crea tabla users
- `002_create_tasks.up.sql`: Crea tabla tasks
- Archivos `.down.sql` para rollback

---

## 7. Deployment Docker

Arquitectura de contenedores y networking.

<image src="diagrams/deployment-docker.png" alt="Deployment Docker" />

### Configuración Docker Compose

```yaml
services:
  postgres:
    image: postgres:15-alpine
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d

  api:
    build: .
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped
```

---

## 8. Decisiones de Arquitectura

### Clean Architecture

**Decisión**: Separar la aplicación en capas con dependencias unidireccionales.

**Razones**:
- Facilita testing con mocks
- Permite cambiar implementaciones sin afectar negocio
- Código más mantenible y escalable

**Trade-offs**:
- Mayor cantidad de archivos
- Curva de aprendizaje inicial

### Repository Pattern

**Decisión**: Abstraer acceso a datos mediante interfaces.

**Razones**:
- Permite testear servicios sin base de datos
- Facilita cambiar de PostgreSQL a otro motor
- Centraliza queries SQL

**Trade-offs**:
- Capa adicional de abstracción
- No aprovecha características específicas del motor

### Chi Router

**Decisión**: Usar `chi` en lugar de `net/http` puro o frameworks más pesados.

**Razones**:
- Ligero y compatible con `net/http`
- Excelente soporte de middleware
- Routing basado en contexto

**Alternativas consideradas**: Gin, Echo

### State Machine en Dominio

**Decisión**: Implementar validación de transiciones en la entidad `Task`.

**Razones**:
- La lógica de estados es regla de negocio core
- Garantiza consistencia en todo el sistema
- Facilita testing de reglas de negocio

**Implementación**:
```go
// domain/task.go
func (t *Task) CanTransitionTo(newStatus TaskStatus) error {
    // Mapa de transiciones válidas
    validTransitions := map[TaskStatus][]TaskStatus{...}
    // Validación
}
```

### JWT para Autenticación

**Decisión**: Usar JWT stateless en lugar de sesiones.

**Razones**:
- No requiere storage de sesiones
- Facilita escalabilidad horizontal
- Cliente puede guardar token localmente

**Trade-offs**:
- No se pueden revocar tokens antes de expiración
- Payload visible (no sensible)

**Mitigación**: Tokens de corta duración + refresh tokens (futuro)

### pgx/v5 Driver

**Decisión**: Usar `pgx` en lugar de `database/sql` + `pq`.

**Razones**:
- Mejor performance (pool nativo)
- Soporte de tipos PostgreSQL específicos
- API moderna con contexto

### Validación en Múltiples Capas

**Decisión**: Validar en Handler (formato) y Service (negocio).

**Razones**:
- Handlers validan estructura JSON
- Services validan reglas de negocio
- Separación de responsabilidades

### Passwords Temporales

**Decisión**: Forzar cambio de password en primer login.

**Razones**:
- Mejor seguridad (admin no conoce password final)
- Cumple buenas prácticas de seguridad
- Flag `must_change_password` en base de datos

---

## Tecnologías Utilizadas

### Backend
- **Go 1.21+**: Lenguaje principal
- **Chi v5**: Router HTTP
- **pgx/v5**: Driver PostgreSQL
- **jwt/v5**: JSON Web Tokens
- **bcrypt**: Hash de passwords

### Base de Datos
- **PostgreSQL 15**: Motor de base de datos
- **UUID**: Identificadores únicos

### DevOps
- **Docker**: Containerización
- **Docker Compose**: Orquestación local

### Testing
- **testify**: Assertions y mocks
- **Go testing**: Tests unitarios nativos

---

## Métricas de Calidad

- **Cobertura de Tests**: > 80% en servicios críticos
- **Separación de Capas**: 4 capas claramente definidas
- **Tamaño de Handlers**: < 100 líneas promedio
- **Tamaño de Servicios**: Separados por responsabilidad (validation, authorization, builder)

---

## Futuras Mejoras

1. **Refresh Tokens**: Implementar renovación de tokens
2. **Rate Limiting**: Limitar requests por IP
3. **Audit Log**: Registrar todas las operaciones
4. **Observability**: Métricas con Prometheus
5. **API Versioning**: Versionado semántico
6. **GraphQL**: Alternativa a REST
7. **Webhooks**: Notificaciones de eventos
8. **Búsqueda Full-Text**: PostgreSQL FTS o ElasticSearch

---

**Última actualización**: 2026-01-20