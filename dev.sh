#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

detect_docker_compose() {
    if command -v "docker-compose" >/dev/null 2>&1; then
        DOCKER_COMPOSE="docker-compose"
    elif docker compose version >/dev/null 2>&1; then
        DOCKER_COMPOSE="docker compose"
    else
        print_error "Neither 'docker-compose' nor 'docker compose' found"
        exit 1
    fi
}

show_help() {
    echo "🚀 Script de Desarrollo para Mattermost Reactions Plugin"
    echo ""
    echo "Uso: $0 [comando]"
    echo ""
    echo "Comandos disponibles:"
    echo "  setup     - Configura el entorno de desarrollo"
    echo "  start     - Inicia Mattermost server"
    echo "  build     - Construye el plugin"
    echo "  deploy    - Despliega el plugin al servidor"
    echo "  logs      - Muestra logs del servidor Mattermost"
    echo "  webhook   - Muestra logs del webhook receiver"
    echo "  stop      - Detiene todos los servicios"
    echo "  clean     - Limpia volúmenes y datos"
    echo "  test      - Ejecuta un test completo"
    echo "  status    - Verifica el estado de los servicios"
    echo "  help      - Muestra esta ayuda"
    echo ""
    echo "Ejemplos:"
    echo "  $0 setup && $0 start     # Primera vez"
    echo "  $0 build && $0 deploy    # Después de cambios en el código"
    echo "  $0 test                  # Prueba completa"
}

check_docker() {
    if ! docker info >/dev/null 2>&1; then
        print_error "Docker no está ejecutándose. Por favor inicia Docker."
        exit 1
    fi
    detect_docker_compose
}

setup_env() {
    print_status "Configurando entorno de desarrollo..."

    mkdir -p dist

    chmod +x "$0"

    print_success "Entorno configurado correctamente"
}

start_services() {
    print_status "Iniciando servicios de desarrollo..."
    check_docker

    if [ ! -f "dist/com.mattermost.reactions-plugin-0.1.0.tar.gz" ]; then
        print_warning "Plugin no encontrado. Construyendo..."
        build_plugin
    fi

    $DOCKER_COMPOSE -f docker-compose.yml up -d

    print_success "Servicios iniciados"
    print_status "🌐 Mattermost: http://localhost:8065"
    print_status "🪝 Webhook receiver: http://localhost:3000"
    print_status "📊 Webhook endpoint: http://localhost:3000/webhook"

    print_status "Esperando que los servicios estén listos..."
    sleep 10

    check_services_status
}

build_plugin() {
    print_status "Construyendo plugin..."
    make server webapp bundle
    print_success "Plugin construido: dist/com.mattermost.reactions-plugin-0.1.0.tar.gz"
}

deploy_plugin() {
    print_status "Desplegando plugin..."

    if [ ! -f "dist/com.mattermost.reactions-plugin-0.1.0.tar.gz" ]; then
        print_error "Plugin no encontrado. Ejecuta 'build' primero."
        exit 1
    fi

    docker cp dist/com.mattermost.reactions-plugin-0.1.0.tar.gz mattermost-reactions_mattermost_1:/tmp/plugin.tar.gz

    docker exec mattermost-reactions_mattermost_1 sh -c "
        cd /tmp &&
        tar -xzf plugin.tar.gz &&
        cp -r com.mattermost.reactions-plugin /mattermost/plugins/ &&
        chown -R mattermost:mattermost /mattermost/plugins/com.mattermost.reactions-plugin
    "

    print_success "Plugin desplegado. Reinicia el servidor para cargar los cambios."
    print_status "💡 O ve a System Console > Plugins para activar el plugin"
}

show_logs() {
    print_status "Mostrando logs de Mattermost..."
    detect_docker_compose
    $DOCKER_COMPOSE -f docker-compose.yml logs -f mattermost
}

show_webhook_logs() {
    print_status "Mostrando logs del webhook receiver..."
    detect_docker_compose
    $DOCKER_COMPOSE -f docker-compose.yml logs -f webhook-receiver
}

stop_services() {
    print_status "Deteniendo servicios..."
    detect_docker_compose
    $DOCKER_COMPOSE -f docker-compose.yml down
    print_success "Servicios detenidos"
}

clean_data() {
    print_warning "⚠️  Esto eliminará TODOS los datos de desarrollo"
    read -p "¿Estás seguro? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Limpiando datos..."
        detect_docker_compose
        $DOCKER_COMPOSE -f docker-compose.yml down -v --remove-orphans
        docker volume prune -f
        print_success "Datos eliminados"
    else
        print_status "Cancelado"
    fi
}

check_services_status() {
    print_status "Verificando estado de servicios..."

    if curl -s http://localhost:8065/api/v4/system/ping >/dev/null 2>&1; then
        print_success "✅ Mattermost server: OK"
    else
        print_warning "❌ Mattermost server: NO DISPONIBLE"
    fi

    if curl -s http://localhost:3000/health >/dev/null 2>&1; then
        print_success "✅ Webhook receiver: OK"
    else
        print_warning "❌ Webhook receiver: NO DISPONIBLE"
    fi
}

run_test() {
    print_status "🧪 Ejecutando test completo del plugin..."

    check_services_status

    print_status "📋 Pasos para probar manualmente:"
    echo "1. Ve a http://localhost:8065"
    echo "2. Crea una cuenta de administrador"
    echo "3. Ve a System Console > Plugins"
    echo "4. Sube el plugin: dist/com.mattermost.reactions-plugin-0.1.0.tar.gz"
    echo "5. Activa el plugin"
    echo "6. Configura webhook URL: http://webhook-receiver:3000/webhook"
    echo "7. Crea un team y channel"
    echo "8. Invita @reactions-bot al channel"
    echo "9. Escribe un mensaje y agrega una reacción"
    echo "10. Verifica los webhooks en: ./dev-test.sh webhook"
    echo ""
    print_status "🔗 Enlaces útiles:"
    echo "   - Mattermost: http://localhost:8065"
    echo "   - Webhook health: http://localhost:3000/health"
    echo "   - Plugin status: http://localhost:8065/plugins/com.mattermost.reactions-plugin/status"
}

main() {
    case "${1:-help}" in
        setup)
            setup_env
            ;;
        start)
            start_services
            ;;
        build)
            build_plugin
            ;;
        deploy)
            deploy_plugin
            ;;
        logs)
            show_logs
            ;;
        webhook)
            show_webhook_logs
            ;;
        stop)
            stop_services
            ;;
        clean)
            clean_data
            ;;
        test)
            run_test
            ;;
        status)
            check_services_status
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "Comando desconocido: $1"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
