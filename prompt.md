Tôi muốn xây dựng một server log sử dụng gin với yêu cầu như sau:
có 1 api /*

khi gọi bất kỳ endpoint nào đến server log nó sẽ lưu lại log vào file log với format như sau:

# ---

time: 2025-06-02 15:04:05
method: GET
path: /api/v1/user?name=abc
query: {}
header: {}
body: {}

khi đó server log sẽ trả về response là 200 OK

file log sẽ được lưu trong thư mục logs với tên file là log_2025-06-02.log
và khi gọi endpoint /api/v1/log?date=2025-06-02 thì server log sẽ trả về response là nội dung của file log_2025-06-02.log

mỗi ngày sẽ sinh ra 1 file log mới
