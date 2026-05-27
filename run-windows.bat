@echo off
title Example.com - Start Queue System
echo ==========================================
echo Starting Example.com Queue Ticket System...
echo ==========================================
docker-compose up -d --build
echo ==========================================
echo [SUCCESS] System is now up and running!
echo ------------------------------------------
echo Front-end UI : http://localhost:4200
echo Backend API  : http://localhost:3000
echo DB Adminer   : http://localhost:8080
echo ==========================================
pause
