#!/bin/bash

set -x
LOG_DIR="/home/tuantech/workspace/oeg/log.nginx.allRequest/logs"

# Kiểm tra thư mục tồn tại
if [ ! -d "$LOG_DIR" ]; then
    echo "Lỗi: Thư mục $LOG_DIR không tồn tại."
    exit 1
fi

# Kiểm tra quyền ghi
if [ ! -w "$LOG_DIR" ]; then
    echo "Lỗi: Không có quyền ghi trong $LOG_DIR."
    exit 1
fi

echo "Bắt đầu xóa log lúc $(date)"

# Kiểm tra file thỏa mãn điều kiện
echo "Danh sách file log cũ hơn 7 ngày:"
find "$LOG_DIR" -type f -name "log_*.log" -mtime +7 -ls

# Xóa file và ghi log
find "$LOG_DIR" -type f -name "log_*.log" -mtime +7 -exec rm -v {} \;

echo "Kết thúc xóa log."