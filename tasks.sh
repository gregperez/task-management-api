#!/bin/bash
# Script de tareas para Linux/Mac - Equivalente al Makefile

case "$1" in
  help|"")
    echo ""
    echo "Comandos disponibles:"
    echo "  ./tasks.sh help           - Muestra esta ayuda"
    echo "  ./tasks.sh build          - Compila la aplicación"
    echo "  ./tasks.sh run            - Ejecuta la aplicación"
    echo "  ./tasks.sh test           - Ejecuta los tests"
    echo "  ./tasks.sh test-coverage  - Muestra cobertura de tests"
    echo "  ./tasks.sh clean          - Limpia archivos generados"
    echo "  ./tasks.sh migrate-up     - Ejecuta migraciones hacia arriba"
    echo "  ./tasks.sh migrate-down   - Ejecuta migraciones hacia abajo"
    echo "  ./tasks.sh docker-up      - Inicia contenedores Docker"
    echo "  ./tasks.sh docker-down    - Detiene contenedores Docker"
    echo "  ./tasks.sh docker-logs    - Muestra logs de Docker"
    echo "  ./tasks.sh lint           - Ejecuta linter"
    echo ""
    ;;
  
  build)
    echo "Compilando la aplicación..."
    go build -o bin/api ./cmd/api
    ;;
  
  run)
    echo "Ejecutando la aplicación..."
    go run ./cmd/api/main.go
    ;;
  
  test)
    echo "Ejecutando tests..."
    go test -v -coverprofile=coverage.out ./...
    ;;
  
  test-coverage)
    echo "Ejecutando tests y generando cobertura..."
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out
    ;;
  
  clean)
    echo "Limpiando archivos generados..."
    rm -rf bin/
    rm -f coverage.out
    echo "Limpieza completada."
    ;;
  
  migrate-up)
    echo "Ejecutando migraciones hacia arriba..."
    docker run --rm -v "$(pwd)/migrations:/migrations" --network host migrate/migrate -path=/migrations/ -database "postgres://postgres:postgres@localhost:5432/taskmanager?sslmode=disable" up
    ;;
  
  migrate-down)
    echo "Ejecutando migraciones hacia abajo..."
    docker run --rm -v "$(pwd)/migrations:/migrations" --network host migrate/migrate -path=/migrations/ -database "postgres://postgres:postgres@localhost:5432/taskmanager?sslmode=disable" down
    ;;
  
  docker-up)
    echo "Iniciando contenedores Docker..."
    docker-compose up -d
    ;;
  
  docker-down)
    echo "Deteniendo contenedores Docker..."
    docker-compose down
    ;;
  
  docker-logs)
    echo "Mostrando logs de Docker..."
    docker-compose logs -f
    ;;
  
  lint)
    echo "Ejecutando linter..."
    golangci-lint run
    ;;
  
  *)
    echo "Comando desconocido: $1"
    echo "Usa './tasks.sh help' para ver los comandos disponibles"
    exit 1
    ;;
esac
