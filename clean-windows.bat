@echo off
title Example.com - Clean System Data
echo ==========================================
echo [WARNING] Stopping system and destroying all database volumes...
echo ==========================================
docker-compose down -v
echo ==========================================
echo [SUCCESS] All containers stopped and database volumes wiped clean!
echo ==========================================
pause
