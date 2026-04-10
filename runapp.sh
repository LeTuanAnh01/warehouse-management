#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
BE_DIR="$ROOT_DIR/backend-go"
FE_DIR="$ROOT_DIR/frontend-web"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   🏪 Khởi động Quản lý Kho App        ${NC}"
echo -e "${BLUE}========================================${NC}"

mkdir -p "$ROOT_DIR/logs"

# Kill process cũ
echo -e "${YELLOW}⏹  Dừng process cũ nếu có...${NC}"
pkill -f "go run ./cmd/api/main.go" 2>/dev/null
pkill -f "vite" 2>/dev/null
sleep 1

# Kiểm tra PostgreSQL
echo -e "${YELLOW}🔍 Kiểm tra PostgreSQL...${NC}"
if ! pg_isready -q 2>/dev/null; then
    echo -e "${YELLOW}⚡ Khởi động PostgreSQL...${NC}"
    brew services start postgresql 2>/dev/null
    sleep 3
fi
echo -e "${GREEN}✅ PostgreSQL đang chạy${NC}"

# Tạo database nếu chưa có
echo -e "${YELLOW}🗄  Kiểm tra database...${NC}"
if ! psql -lqt 2>/dev/null | cut -d \| -f 1 | grep -qw warehouse; then
    echo -e "${YELLOW}📦 Tạo database warehouse...${NC}"
    createdb warehouse
    psql -d warehouse -f "$ROOT_DIR/database/schema.sql"
    psql -d warehouse -f "$ROOT_DIR/database/seed.sql"
    echo -e "${GREEN}✅ Database đã tạo xong${NC}"
else
    echo -e "${GREEN}✅ Database warehouse đã tồn tại${NC}"
fi

# Tạo .env nếu chưa có
if [ ! -f "$BE_DIR/.env" ]; then
    echo -e "${YELLOW}📝 Tạo file .env cho backend...${NC}"
    cat > "$BE_DIR/.env" << 'EOF'
DB_HOST=localhost
DB_PORT=5432
DB_USER=anhlt
DB_PASSWORD=
DB_NAME=warehouse
DB_SSLMODE=disable
JWT_SECRET=your-secret-key-local
ENV=development
PORT=8080
HOST=0.0.0.0
EOF
    echo -e "${GREEN}✅ Tạo .env xong${NC}"
fi

# Khởi động Backend
echo -e "${YELLOW}🚀 Khởi động Backend (port 8080)...${NC}"
cd "$BE_DIR"
go run ./cmd/api/main.go > "$ROOT_DIR/logs/backend.log" 2>&1 &
BE_PID=$!
echo "$BE_PID" > "$ROOT_DIR/logs/backend.pid"
echo -e "${GREEN}✅ Backend PID: $BE_PID${NC}"

# Chờ backend ready
echo -e "${YELLOW}⏳ Chờ Backend khởi động...${NC}"
for i in {1..20}; do
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Backend đã sẵn sàng!${NC}"
        break
    fi
    sleep 1
done

# Khởi động Frontend
echo -e "${YELLOW}🌐 Khởi động Frontend...${NC}"
cd "$FE_DIR"
npm run dev > "$ROOT_DIR/logs/frontend.log" 2>&1 &
FE_PID=$!
echo "$FE_PID" > "$ROOT_DIR/logs/frontend.pid"
echo -e "${GREEN}✅ Frontend PID: $FE_PID${NC}"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   ✅ Ứng dụng đã khởi động xong!      ${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "  🔧 Backend:  ${BLUE}http://localhost:8080${NC}"
echo -e "  🌐 Frontend: ${BLUE}http://localhost:5173${NC}"
echo -e "  📋 Logs:     ${BLUE}$ROOT_DIR/logs/${NC}"
echo ""
echo -e "${YELLOW}Nhấn Ctrl+C để dừng tất cả${NC}"

trap "echo -e '\n${RED}⏹ Đang dừng tất cả...${NC}'; kill $BE_PID $FE_PID 2>/dev/null; exit 0" SIGINT SIGTERM

tail -f "$ROOT_DIR/logs/backend.log" "$ROOT_DIR/logs/frontend.log"
