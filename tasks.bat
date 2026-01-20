@echo off
REM Script de tareas para Windows - Equivalente al Makefile

if "%1"=="" goto help
if "%1"=="help" goto help
if "%1"=="build" goto build
if "%1"=="run" goto run
if "%1"=="test" goto test
if "%1"=="test-coverage" goto test-coverage
if "%1"=="clean" goto clean
if "%1"=="migrate-up" goto migrate-up
if "%1"=="migrate-down" goto migrate-down
if "%1"=="docker-up" goto docker-up
if "%1"=="docker-down" goto docker-down
if "%1"=="docker-logs" goto docker-logs
if "%1"=="lint" goto lint
goto help

:help
echo.
echo Comandos disponibles:
echo   tasks.bat help           - Muestra esta ayuda
echo   tasks.bat build          - Compila la aplicacion
echo   tasks.bat run            - Ejecuta la aplicacion
echo   tasks.bat test           - Ejecuta los tests
echo   tasks.bat test-coverage  - Muestra cobertura de tests
echo   tasks.bat clean          - Limpia archivos generados
echo   tasks.bat migrate-up     - Ejecuta migraciones hacia arriba
echo   tasks.bat migrate-down   - Ejecuta migraciones hacia abajo
echo   tasks.bat docker-up      - Inicia contenedores Docker
echo   tasks.bat docker-down    - Detiene contenedores Docker
echo   tasks.bat docker-logs    - Muestra logs de Docker
echo   tasks.bat lint           - Ejecuta linter
echo.
goto end

:build
echo Compilando la aplicacion...
go build -o bin\api.exe .\cmd\api
goto end

:run
echo Ejecutando la aplicacion...
go run .\cmd\api\main.go
goto end

:test
echo Ejecutando tests...
go test -v -race -coverprofile=coverage.out ./...
goto end

:test-coverage
echo Ejecutando tests y generando cobertura...
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
goto end

:clean
echo Limpiando archivos generados...
if exist bin rmdir /s /q bin
if exist coverage.out del /f coverage.out
echo Limpieza completada.
goto end

:migrate-up
echo Ejecutando migraciones hacia arriba...
docker run --rm -v "%cd%\migrations:/migrations" --network host migrate/migrate -path=/migrations/ -database "postgres://postgres:postgres@localhost:5432/taskmanager?sslmode=disable" up
goto end

:migrate-down
echo Ejecutando migraciones hacia abajo...
docker run --rm -v "%cd%\migrations:/migrations" --network host migrate/migrate -path=/migrations/ -database "postgres://postgres:postgres@localhost:5432/taskmanager?sslmode=disable" down
goto end

:docker-up
echo Iniciando contenedores Docker...
docker-compose up -d
goto end

:docker-down
echo Deteniendo contenedores Docker...
docker-compose down
goto end

:docker-logs
echo Mostrando logs de Docker...
docker-compose logs -f
goto end

:lint
echo Ejecutando linter...
golangci-lint run
goto end

:end
