#!/bin/bash
# =====================================================
# Microservices Live Reload Dev Runner using Air
# =====================================================
# Runs product, cart, and order services concurrently
# with automatic rebuild on file changes.
# =====================================================

# Colors for clarity
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to start each service in background
run_service() {
  local name=$1
  local path=$2
  echo -e "${YELLOW}→ Starting ${name} service...${NC}"
  cd "$path" || exit
  air &>/dev/null &
  pid=$!
  echo -e "${GREEN}${name} started with PID ${pid}${NC}"
  cd - >/dev/null || exit
}

# Start all microservices
run_service "Product" "product-service"
run_service "Cart" "cart-service"
run_service "Order" "order-service"

# Keep script running until Ctrl+C
echo -e "${YELLOW}All services are running with live reload via Air.${NC}"
echo -e "Press Ctrl+C to stop all.\n"

trap "echo 'Stopping all services...'; pkill air; exit 0" SIGINT
while true; do sleep 1; done
