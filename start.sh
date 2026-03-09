#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Запуск full-stack приложения (Vue + Go) ===${NC}"

# 1. Проверка наличия Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go не найден. Установите Go (https://golang.org/dl/) и добавьте в PATH.${NC}"
    exit 1
else
    echo "✅ Go найден: $(go version)"
fi

# 2. Проверка наличия Node.js и npm
if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ Node.js не найден. Установите Node.js (https://nodejs.org/).${NC}"
    exit 1
else
    echo "✅ Node.js найден: $(node -v)"
fi

if ! command -v npm &> /dev/null; then
    echo -e "${RED}❌ npm не найден. Убедитесь, что npm установлен вместе с Node.js.${NC}"
    exit 1
fi

# 3. Проверка структуры папок
if [ ! -d "backend" ]; then
    echo -e "${RED}❌ Папка 'backend' не найдена в текущем каталоге.${NC}"
    exit 1
fi

if [ ! -d "frontend" ]; then
    echo -e "${RED}❌ Папка 'frontend' не найдена в текущем каталоге.${NC}"
    exit 1
fi

# 4. Проверка наличия go.mod и package.json
if [ ! -f "backend/go.mod" ]; then
    echo -e "${RED}❌ В папке backend нет go.mod. Возможно, не инициализирован модуль.${NC}"
    exit 1
fi

if [ ! -f "frontend/package.json" ]; then
    echo -e "${RED}❌ В папке frontend нет package.json. Возможно, проект Vue не создан.${NC}"
    exit 1
fi

# 5. Установка зависимостей
echo -e "${GREEN}--- Установка зависимостей backend (go mod tidy) ---${NC}"
cd backend
go mod tidy
cd ..

echo -e "${GREEN}--- Установка зависимостей frontend (npm install) ---${NC}"
cd frontend
npm install
cd ..

# 6. Функция для убийства процессов на портах
kill_port() {
    local port=$1
    local pids=$(lsof -ti:$port)
    if [ ! -z "$pids" ]; then
        echo -e "${RED}   Останавливаю процессы на порту $port...${NC}"
        echo "$pids" | xargs kill -9 2>/dev/null
        sleep 1
    fi
}

# 7. Убиваем процессы на нужных портах перед запуском
echo -e "${GREEN}--- Очистка портов ---${NC}"
kill_port 8080
kill_port 5173

# 8. Запуск бэкенда
echo -e "${GREEN}--- Запуск бэкенда (Go) на порту 8080 ---${NC}"
cd backend
go run cmd/api/main.go &
BACKEND_PID=$!
cd ..

# 9. Запуск фронтенда
echo -e "${GREEN}--- Запуск фронтенда (Vite) ---${NC}"
cd frontend
npm run dev &
FRONTEND_PID=$!
cd ..

# 10. Функция для завершения при выходе
cleanup() {
    echo -e "\n${GREEN}--- Остановка приложения ---${NC}"
    
    # Убиваем процессы по PID
    if [ ! -z "$BACKEND_PID" ]; then
        echo -e "${RED}   Останавливаю бэкенд (PID: $BACKEND_PID)${NC}"
        kill -9 $BACKEND_PID 2>/dev/null
    fi
    
    if [ ! -z "$FRONTEND_PID" ]; then
        echo -e "${RED}   Останавливаю фронтенд (PID: $FRONTEND_PID)${NC}"
        kill -9 $FRONTEND_PID 2>/dev/null
    fi
    
    # Дополнительно убиваем все процессы на портах
    kill_port 8080
    kill_port 5173
    
    echo -e "${GREEN}✅ Приложение остановлено.${NC}"
    exit 0
}

# Устанавливаем обработчик сигналов
trap cleanup SIGINT SIGTERM EXIT

# 11. Вывод информации
echo -e "${GREEN}✅ Приложение запущено!${NC}"
echo -e "   Бэкенд: http://localhost:8080"
echo -e "   Фронтенд: http://localhost:5173"
echo -e "   Нажмите Ctrl+C для остановки\n"

# Ждем
wait