#!/bin/bash

# Air Development Helper Script
# Usage: ./scripts/air-dev.sh [api|worker|debug|both]

set -e

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_ROOT"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_colored() {
    color=$1
    message=$2
    echo -e "${color}${message}${NC}"
}

# Function to check if Air is installed
check_air() {
    if ! command -v air &> /dev/null; then
        print_colored $YELLOW "Air not found. Installing..."
        go install github.com/air-verse/air@latest
    fi
}

# Function to check if required env files exist
check_env_files() {
    # Check for API dev environment
    if [[ ! -f ".env.api-dev" ]]; then
        if [[ -f "deployments/environments/development/.env.api-dev" ]]; then
            print_colored $YELLOW "Copying .env.api-dev from deployments/environments/development/"
            cp deployments/environments/development/.env.api-dev .env.api-dev
        elif [[ -f "deployments/environments/development/.env.dev" ]]; then
            print_colored $YELLOW "Warning: .env.api-dev not found. Using .env.dev as fallback."
            cp deployments/environments/development/.env.dev .env.api-dev
        elif [[ -f ".env.dev" ]]; then
            cp .env.dev .env.api-dev
        else
            print_colored $RED "Error: No API development environment found. Run 'make env-api-dev'"
            exit 1
        fi
    fi

    # Check for Worker dev environment
    if [[ ! -f ".env.worker-dev" ]]; then
        if [[ -f "deployments/environments/development/.env.worker-dev" ]]; then
            print_colored $YELLOW "Copying .env.worker-dev from deployments/environments/development/"
            cp deployments/environments/development/.env.worker-dev .env.worker-dev
        elif [[ -f "deployments/environments/development/.env.dev" ]]; then
            print_colored $YELLOW "Warning: .env.worker-dev not found. Using .env.dev as fallback."
            cp deployments/environments/development/.env.dev .env.worker-dev
        elif [[ -f ".env.dev" ]]; then
            cp .env.dev .env.worker-dev
        else
            print_colored $RED "Error: No Worker development environment found. Run 'make env-worker-dev'"
            exit 1
        fi
    fi
}

# Function to create tmp directory if it doesn't exist
ensure_tmp_dir() {
    if [[ ! -d "tmp" ]]; then
        mkdir -p tmp
        print_colored $GREEN "Created tmp directory"
    fi
}

# Function to start API with Air
start_api() {
    print_colored $BLUE "Starting API server with Air..."
    air -c .air.api.toml
}

# Function to start Worker with Air
start_worker() {
    print_colored $BLUE "Starting Worker server with Air..."
    air -c .air.worker.toml
}

# Function to start API in debug mode
start_debug() {
    print_colored $BLUE "Starting API server in debug mode with Air..."
    air -c .air.debug.toml
}

# Function to start both API and Worker
start_both() {
    print_colored $BLUE "Starting both API and Worker servers..."

    # Start API in background
    print_colored $GREEN "Starting API server..."
    air -c .air.api.toml &
    API_PID=$!

    # Wait a moment
    sleep 2

    # Start Worker in background
    print_colored $GREEN "Starting Worker server..."
    air -c .air.worker.toml &
    WORKER_PID=$!

    # Function to handle cleanup
    cleanup() {
        print_colored $YELLOW "Shutting down servers..."
        kill $API_PID $WORKER_PID 2>/dev/null || true
        wait $API_PID $WORKER_PID 2>/dev/null || true
        print_colored $GREEN "Servers stopped"
    }

    # Set trap for cleanup
    trap cleanup SIGINT SIGTERM

    print_colored $GREEN "Both servers are running. Press Ctrl+C to stop both."
    wait
}

# Function to show usage
show_usage() {
    echo "Air Development Helper Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  api     Start API server with live reload"
    echo "  worker  Start Worker server with live reload"
    echo "  debug   Start API server with debug flags"
    echo "  both    Start both API and Worker servers"
    echo "  help    Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 api"
    echo "  $0 worker"
    echo "  $0 debug"
    echo "  $0 both"
}

# Main script logic
main() {
    # Check prerequisites
    check_air
    check_env_files
    ensure_tmp_dir

    case "${1:-help}" in
        "api")
            start_api
            ;;
        "worker")
            start_worker
            ;;
        "debug")
            start_debug
            ;;
        "both")
            start_both
            ;;
        "help"|"-h"|"--help")
            show_usage
            ;;
        *)
            print_colored $RED "Unknown command: $1"
            echo ""
            show_usage
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
