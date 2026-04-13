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

# ─── Khởi động PostgreSQL ─────────────────────
echo -e "${YELLOW}🐘 Kiểm tra PostgreSQL...${NC}"

PG_DATA="/opt/homebrew/var/postgresql@16"
PG_LOG="/opt/homebrew/var/log/postgresql@16.log"

if pg_isready -q 2>/dev/null; then
    echo -e "${GREEN}✅ PostgreSQL đã đang chạy${NC}"
else
    echo -e "${YELLOW}   Khởi động PostgreSQL...${NC}"

    # Xóa lock file cũ nếu có
    rm -f "$PG_DATA/postmaster.pid" 2>/dev/null

    # Khởi động bằng pg_ctl
    pg_ctl -D "$PG_DATA" -l "$PG_LOG" start 2>/dev/null

    echo -e "${YELLOW}⏳ Chờ PostgreSQL khởi động...${NC}"
    for i in {1..20}; do
        if pg_isready -q 2>/dev/null; then
            echo -e "${GREEN}✅ PostgreSQL đang chạy${NC}"
            break
        fi
        if [ $i -eq 20 ]; then
            echo -e "${RED}❌ PostgreSQL không khởi động được!${NC}"
            echo -e "${YELLOW}   Xem log: tail -20 $PG_LOG${NC}"
            exit 1
        fi
        sleep 1
    done
fi

# ─── Tạo database nếu chưa có ─────────────────
echo -e "${YELLOW}🗄  Kiểm tra database...${NC}"
DB_NAME=$(grep DB_NAME "$BE_DIR/.env" 2>/dev/null | cut -d '=' -f2 || echo "warehouse_db")

if ! psql -lqt 2>/dev/null | cut -d \| -f 1 | grep -qw "$DB_NAME"; then
    echo -e "${YELLOW}📦 Tạo database $DB_NAME...${NC}"
    createdb "$DB_NAME"
    psql -d "$DB_NAME" -f "$ROOT_DIR/database/schema.sql"
    psql -d "$DB_NAME" -f "$ROOT_DIR/database/seed.sql"
    echo -e "${GREEN}✅ Database đã tạo xong${NC}"
else
    echo -e "${GREEN}✅ Database $DB_NAME đã tồn tại${NC}"
fi

# ─── Tạo .env nếu chưa có ─────────────────────
if [ ! -f "$BE_DIR/.env" ]; then
    echo -e "${YELLOW}📝 Tạo file .env cho backend...${NC}"
    cat > "$BE_DIR/.env" << 'EOF'
DB_HOST=localhost
DB_PORT=5432
DB_USER=anhlt
DB_PASSWORD=
DB_NAME=warehouse_db
DB_SSLMODE=disable
JWT_SECRET=your-secret-key-local
ENV=development
PORT=8080
HOST=0.0.0.0
EOF
    echo -e "${GREEN}✅ Tạo .env xong${NC}"
fi

# ─── Khởi động Backend ────────────────────────
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
    if [ $i -eq 20 ]; then
        echo -e "${RED}❌ Backend không khởi động được! Xem logs/backend.log${NC}"
    fi
    sleep 1
done

# ─── Khởi động Frontend ───────────────────────
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
echo -e "  🐘 PostgreSQL: ${BLUE}$PG_VERSION${NC}"
echo -e "  🔧 Backend:    ${BLUE}http://localhost:8080${NC}"
echo -e "  🌐 Frontend:   ${BLUE}http://localhost:5173${NC}"
echo -e "  📋 Logs:       ${BLUE}$ROOT_DIR/logs/${NC}"
echo ""
echo -e "${YELLOW}Nhấn Ctrl+C để dừng tất cả${NC}"

trap "echo -e '\n${RED}⏹ Đang dừng tất cả...${NC}'; kill $BE_PID $FE_PID 2>/dev/null; exit 0" SIGINT SIGTERM

tail -f "$ROOT_DIR/logs/backend.log" "$ROOT_DIR/logs/frontend.log"