# 🛒 Dự án Lesson 12: Shopping Cart API

Đây là dự án **RESTful API** được xây dựng bằng **Golang (framework Gin)** để quản lý hệ thống giỏ hàng (Shopping Cart).  
Dự án được thiết kế theo **kiến trúc 3 tầng chuẩn (Handler → Service → Repository)**, sử dụng nhiều công nghệ thực tế như:

- **PostgreSQL** (kết hợp với `sqlc` và `pgx`)
- **Redis** (dùng làm cache)
- **RabbitMQ** (xử lý nền - background worker)
- **Docker Compose** (đóng gói và triển khai)
- **JWT + bcrypt** (xác thực và mã hóa mật khẩu)

---


---

## ⚙️ Tính năng chính

| Chức năng | Mô tả |
|------------|--------|
| 🔐 **Xác thực người dùng** | Đăng ký, đăng nhập bằng JWT + bcrypt |
| 🛍️ **Quản lý sản phẩm** | CRUD sản phẩm |
| 🛒 **Quản lý giỏ hàng** | Thêm, xoá, cập nhật sản phẩm trong giỏ |
| 📦 **Xử lý đơn hàng** | Gửi yêu cầu qua RabbitMQ để xử lý nền |
| 🚀 **Hiệu năng cao** | Tích hợp Redis để cache dữ liệu truy cập nhiều |
| 🧱 **Kiến trúc sạch** | Tách biệt rõ logic, repository, service |
| 🐳 **Triển khai dễ dàng** | Docker Compose tự động hóa toàn bộ môi trường |

---

## 🧰 Công nghệ sử dụng

| Thành phần | Công nghệ |
|-------------|------------|
| Framework | [Gin](https://github.com/gin-gonic/gin) |
| CSDL chính | PostgreSQL (thông qua `sqlc` + `pgx`) |
| Cache | Redis |
| Message Queue | RabbitMQ |
| Authentication | JWT + bcrypt |
| Triển khai | Docker Compose |

---

