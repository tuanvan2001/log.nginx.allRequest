# Server Log sử dụng Gin Framework

Đây là một server log đơn giản được xây dựng bằng Gin Framework của Go. Server này có khả năng ghi lại thông tin của mọi request đến bất kỳ endpoint nào và trả về nội dung log của một ngày cụ thể.

## Tính năng

1. Ghi log cho mọi request đến server (trừ endpoint `/api/v1/log`)
2. Trả về nội dung log của một ngày cụ thể thông qua endpoint `/api/v1/log?date=YYYY-MM-DD`
3. Tự động tạo file log mới mỗi ngày với định dạng `log_YYYY-MM-DD.log`

## Cấu trúc log

Mỗi entry trong file log có định dạng như sau:

```
#---
time: 2025-06-02 15:04:05
method: GET
path: /api/v1/user?name=abc
query: {}
header: {}
body: {}
```

## Cài đặt và chạy

### Yêu cầu

- Go 1.24.2 hoặc cao hơn

### Cài đặt

```bash
# Clone repository
git clone <repository-url>
cd log.nginx.allRequest

# Cài đặt dependencies
go mod download

# Chạy server
go run main.go
```

Server sẽ chạy trên cổng 8080 theo mặc định.

## Sử dụng

### Ghi log

Gửi request đến bất kỳ endpoint nào (trừ `/api/v1/log`) để ghi log:

```bash
curl -X GET "http://localhost:8080/api/v1/user?name=abc"
```

Server sẽ trả về response `200 OK` và ghi thông tin request vào file log.

### Xem log

Để xem log của một ngày cụ thể:

```bash
curl -X GET "http://localhost:8080/api/v1/log?date=2025-06-02"
```

Server sẽ trả về nội dung của file log cho ngày đó.

## Cấu trúc thư mục

```
.
├── main.go       # Mã nguồn chính của server
├── go.mod        # File quản lý dependencies
├── go.sum        # File checksum của dependencies
└── logs/         # Thư mục chứa các file log (tự động tạo)
    └── log_YYYY-MM-DD.log  # File log cho mỗi ngày
```