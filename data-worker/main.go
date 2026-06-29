package main

import (
	"fb-management-system/data-worker/database"
	"fb-management-system/data-worker/tasks"

	"fmt"
	"log"
)

func main() {
	fmt.Println("Đang khởi động Data Worker...")

	// 1. Kết nối Database (Thay đổi mật khẩu thật của bạn tại đây)
	dsn := "host=localhost user=postgres password=Tienanh0108! dbname=fanpage_db port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"
	db := database.ConnectDB(dsn)

	if db == nil {
		log.Fatal("Lỗi khởi tạo Database.")
	}

	fmt.Println("Kích hoạt tiến trình quét dữ liệu kiểm tra...")

	// 2. Chạy thử nghiệm hàm quét dữ liệu luôn
	tasks.SyncPageData(db)

	fmt.Println("Tiến trình kiểm tra hoàn tất!")
}
