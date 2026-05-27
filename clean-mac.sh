#!/bin/bash
echo "=========================================="
echo "[WARNING] Stopping system and destroying all database volumes..."
echo "=========================================="
docker compose down -v
echo "=========================================="
echo "[SUCCESS] All containers stopped and database volumes wiped clean!"
echo "=========================================="
read -p "Press Enter to close this window..."
