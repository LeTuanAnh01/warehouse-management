#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
BE_DIR="$ROOT_DIR/backend-nodejs"
FE_DIR="$ROOT_DIR/frontend-web"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   🏪 Khởi động Quản lý Kho App        ${NC}"
echo -e "${BLUE}========================================${NC}"

mkdir -p "$ROOT_DIR/logs"

# Kill process cũ
echo -e "${YELLOW}⏹  Dừng process cũ nếu có...${NC}"
pkill -f "node dist/main" 2>/dev/null
pkill -f "vite" 2>/dev/null
sleep 1

# ─── Kiểm tra Docker ──────────────────────────
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker chưa chạy! Hãy mở Docker Desktop rồi thử lại.${NC}"
    exit 1
fi

# ─── Khởi động PostgreSQL bằng Docker ─────────
echo -e "${YELLOW}🐘 Kiểm tra PostgreSQL (Docker)...${NC}"

PG_CONTAINER="warehouse_db_local"
PG_USER="admin"
PG_PASSWORD="123456"
PG_DB="warehouse"
PG_PORT="5432"

if docker ps --format '{{.Names}}' | grep -q "^${PG_CONTAINER}$"; then
    echo -e "${GREEN}✅ PostgreSQL container đã đang chạy${NC}"
elif docker ps -a --format '{{.Names}}' | grep -q "^${PG_CONTAINER}$"; then
    echo -e "${YELLOW}   Khởi động lại container PostgreSQL...${NC}"
    docker start "$PG_CONTAINER" > /dev/null
else
    echo -e "${YELLOW}   Tạo container PostgreSQL mới...${NC}"
    docker run -d \
        --name "$PG_CONTAINER" \
        -e POSTGRES_USER="$PG_USER" \
        -e POSTGRES_PASSWORD="$PG_PASSWORD" \
        -e POSTGRES_DB="$PG_DB" \
        -p ${PG_PORT}:5432 \
        -v warehouse_pg_data:/var/lib/postgresql/data \
        postgres:16-alpine > /dev/null
fi

echo -e "${YELLOW}⏳ Chờ PostgreSQL khởi động...${NC}"
for i in {1..30}; do
    if docker exec "$PG_CONTAINER" pg_isready -U "$PG_USER" -d "$PG_DB" -q 2>/dev/null; then
        echo -e "${GREEN}✅ PostgreSQL đang chạy${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ PostgreSQL không khởi động được!${NC}"
        docker logs --tail 20 "$PG_CONTAINER"
        exit 1
    fi
    sleep 1
done

# ─── Nạp schema + seed nếu database mới ───────
TABLE_COUNT=$(docker exec "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -tAc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public';" 2>/dev/null | tr -d '[:space:]')
if [ "$TABLE_COUNT" = "0" ] || [ -z "$TABLE_COUNT" ]; then
    echo -e "${YELLOW}📦 Nạp schema và seed data...${NC}"
    docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" < "$ROOT_DIR/database/schema.sql"
    docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" < "$ROOT_DIR/database/seed.sql"
    echo -e "${GREEN}✅ Database đã tạo xong${NC}"
else
    echo -e "${GREEN}✅ Database đã có sẵn ($TABLE_COUNT tables)${NC}"
fi

# ─── Tạo .env nếu chưa có ─────────────────────
if [ ! -f "$BE_DIR/.env" ]; then
    echo -e "${YELLOW}📝 Tạo file .env cho backend...${NC}"
    cat > "$BE_DIR/.env" << EOF
NODE_ENV=development
PORT=3000
DB_HOST=localhost
DB_PORT=${PG_PORT}
DB_USERNAME=${PG_USER}
DB_PASSWORD=${PG_PASSWORD}
DB_NAME=${PG_DB}
JWT_SECRET=your-secret-key-local
JWT_EXPIRES_IN=7d
EOF
    echo -e "${GREEN}✅ Tạo .env xong${NC}"
fi

# ─── Khởi động Backend ────────────────────────
echo -e "${YELLOW}🚀 Khởi động Backend Node.js (port 3000)...${NC}"
cd "$BE_DIR"

# Build nếu chưa có dist
if [ ! -d "dist" ]; then
    echo -e "${YELLOW}🔨 Build backend lần đầu...${NC}"
    npm run build 2>&1
fi

node dist/main > "$ROOT_DIR/logs/backend.log" 2>&1 &
BE_PID=$!
echo "$BE_PID" > "$ROOT_DIR/logs/backend.pid"
echo -e "${GREEN}✅ Backend PID: $BE_PID${NC}"

# Chờ backend ready
echo -e "${YELLOW}⏳ Chờ Backend khởi động...${NC}"
for i in {1..20}; do
    if curl -s http://localhost:3000 > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Backend đã sẵn sàng!${NC}"
        break
    fi
    if [ $i -eq 20 ]; then
        echo -e "${RED}❌ Backend không khởi động được! Xem logs/backend.log${NC}"
    fi
    sleep 1
done

# ─── Khởi động Frontend ───────────────────────
echo -e "${YELLOW}🌐 Khởi động Frontend...${NC}"
cd "$FE_DIR"
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}📦 Cài dependencies frontend...${NC}"
    npm install 2>&1
fi
npm run dev > "$ROOT_DIR/logs/frontend.log" 2>&1 &
FE_PID=$!
echo "$FE_PID" > "$ROOT_DIR/logs/frontend.pid"
echo -e "${GREEN}✅ Frontend PID: $FE_PID${NC}"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   ✅ Ứng dụng đã khởi động xong!      ${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "  🐘 PostgreSQL: ${BLUE}Docker container ($PG_CONTAINER)${NC}"
echo -e "  🔧 Backend:    ${BLUE}http://localhost:3000${NC}"
echo -e "  🌐 Frontend:   ${BLUE}http://localhost:5173${NC}"
echo -e "  📋 Logs:       ${BLUE}$ROOT_DIR/logs/${NC}"
echo ""
echo -e "${YELLOW}Nhấn Ctrl+C để dừng tất cả${NC}"

trap "echo -e '\n${RED}⏹ Đang dừng tất cả...${NC}'; kill $BE_PID $FE_PID 2>/dev/null; docker stop $PG_CONTAINER > /dev/null 2>&1; exit 0" SIGINT SIGTERM

tail -f "$ROOT_DIR/logs/backend.log" "$ROOT_DIR/logs/frontend.log"