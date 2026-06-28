package main

import (
	"fb-management-system/data-worker/database"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Đang khởi động Data Worker...")

	dsn := "host=localhost user=postgres password=Tienanh0108! dbname=fanpage_db port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"

	db := database.ConnectDB(dsn)

	if db != nil {
		fmt.Println("Sẵn sàng! Worker đã kết nối Database và tạo bảng thành công.")
	} else {
		log.Fatal("Lỗi: Không thể khởi tạo Database. Worker đang dừng lại.")
	}
}
