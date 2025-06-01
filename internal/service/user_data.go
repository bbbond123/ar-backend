package server

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// InitializeAdminUser 初始化管理员用户
func InitializeAdminUser() {
	db := database.GetDB()

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	// 如果没有设置管理员环境变量，跳过初始化
	if adminEmail == "" || adminPassword == "" {
		fmt.Println("未设置管理员账户环境变量，跳过管理员初始化")
		return
	}

	// 检查管理员是否已存在
	var existingAdmin model.User
	err := db.Where("email = ?", adminEmail).First(&existingAdmin).Error
	if err == nil {
		fmt.Printf("管理员账户已存在: %s\n", adminEmail)
		return
	}

	// 创建管理员账户
	hashedPassword := hashPassword(adminPassword)
	admin := model.User{
		Name:        "系统管理员",
		NameKana:    "システムカンリシャ",
		Email:       adminEmail,
		Password:    hashedPassword,
		Provider:    "email",
		Status:      "active",
		Address:     "系统管理",
		PhoneNumber: "000-0000-0000",
		CreatedAt:   time.Now(),
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Printf("创建管理员账户失败: %v\n", err)
	} else {
		fmt.Printf("成功创建管理员账户: %s\n", adminEmail)
	}
}

// InitializeSampleUsers 初始化示例用户数据
func InitializeSampleUsers() {
	// 首先初始化管理员用户
	InitializeAdminUser()

	db := database.GetDB()

	// 检查是否已有用户数据（排除管理员）
	var userCount int64
	adminEmail := os.Getenv("ADMIN_EMAIL")
	query := db.Model(&model.User{})
	if adminEmail != "" {
		query = query.Where("email != ?", adminEmail)
	}
	query.Count(&userCount)

	if userCount > 0 {
		fmt.Printf("数据库中已有 %d 个用户，跳过示例数据初始化\n", userCount)
		return
	}

	fmt.Println("开始初始化示例用户数据...")

	sampleUsers := []model.User{
		{
			Name:        "张三",
			NameKana:    "チョウサン",
			Birth:       parseTime("1990-05-15"),
			Address:     "东京都涩谷区1-2-3",
			Gender:      stringPtr("男"),
			PhoneNumber: "080-1234-5678",
			Email:       "zhangsan@example.com",
			Password:    hashPassword("123456"),
			Avatar:      "https://via.placeholder.com/150/0000FF/FFFFFF?text=张三",
			Provider:    "email",
			Status:      "active",
			CreatedAt:   time.Now(),
		},
		{
			Name:        "李四",
			NameKana:    "リシ",
			Birth:       parseTime("1985-08-20"),
			Address:     "大阪府大阪市中央区4-5-6",
			Gender:      stringPtr("女"),
			PhoneNumber: "090-2345-6789",
			Email:       "lisi@example.com",
			Password:    hashPassword("123456"),
			Avatar:      "https://via.placeholder.com/150/FF0000/FFFFFF?text=李四",
			Provider:    "email",
			Status:      "active",
			CreatedAt:   time.Now(),
		},
		{
			Name:        "王五",
			NameKana:    "オウゴ",
			Birth:       parseTime("1992-12-10"),
			Address:     "京都府京都市左京区7-8-9",
			Gender:      stringPtr("男"),
			PhoneNumber: "070-3456-7890",
			Email:       "wangwu@gmail.com",
			Password:    hashPassword("123456"),
			Avatar:      "https://via.placeholder.com/150/00FF00/FFFFFF?text=王五",
			GoogleID:    "google_123456789",
			Provider:    "google",
			Status:      "active",
			CreatedAt:   time.Now(),
		},
		{
			Name:        "赵六",
			NameKana:    "チョウロク",
			Birth:       parseTime("1988-03-25"),
			Address:     "神奈川県横浜市港北区10-11-12",
			Gender:      stringPtr("女"),
			PhoneNumber: "080-4567-8901",
			Email:       "zhaoliu@icloud.com",
			Password:    hashPassword("123456"),
			Avatar:      "https://via.placeholder.com/150/FFFF00/000000?text=赵六",
			AppleID:     "apple_987654321",
			Provider:    "apple",
			Status:      "active",
			CreatedAt:   time.Now(),
		},
		{
			Name:             "孙七",
			NameKana:         "ソンナナ",
			Birth:            parseTime("1995-07-08"),
			Address:          "福岡県福岡市博多区13-14-15",
			Gender:           stringPtr("男"),
			PhoneNumber:      "090-5678-9012",
			Email:            "sunqi@example.com",
			Password:         hashPassword("123456"),
			Avatar:           "https://via.placeholder.com/150/FF00FF/FFFFFF?text=孙七",
			Provider:         "email",
			Status:           "pending",
			VerifyCode:       "1234",
			VerifyCodeExpire: timePtr(time.Now().Add(10 * time.Minute)),
			CreatedAt:        time.Now(),
		},
		{
			Name:        "周八",
			NameKana:    "シュウハチ",
			Birth:       parseTime("1993-11-30"),
			Address:     "北海道札幌市中央区16-17-18",
			Gender:      stringPtr("女"),
			PhoneNumber: "070-6789-0123",
			Email:       "zhouba@example.com",
			Password:    hashPassword("123456"),
			Avatar:      "https://via.placeholder.com/150/00FFFF/000000?text=周八",
			Provider:    "email",
			Status:      "inactive",
			CreatedAt:   time.Now(),
		},
		{
			Name:        "吴九",
			NameKana:    "ゴキュウ",
			Birth:       parseTime("1987-09-12"),
			Address:     "愛知県名古屋市中区19-20-21",
			Gender:      stringPtr("男"),
			PhoneNumber: "080-7890-1234",
			Email:       "wujiu@gmail.com",
			Password:    hashPassword("123456"),
			Avatar:      "https://via.placeholder.com/150/800080/FFFFFF?text=吴九",
			GoogleID:    "google_246810121",
			Provider:    "google",
			Status:      "active",
			CreatedAt:   time.Now(),
		},
		{
			Name:        "郑十",
			NameKana:    "テイジュウ",
			Birth:       parseTime("1991-04-18"),
			Address:     "広島県広島市中区22-23-24",
			Gender:      stringPtr("女"),
			PhoneNumber: "090-8901-2345",
			Email:       "zhengshi@example.com",
			Password:    hashPassword("123456"),
			Avatar:      "https://via.placeholder.com/150/FFA500/FFFFFF?text=郑十",
			Provider:    "email",
			Status:      "active",
			CreatedAt:   time.Now(),
		},
	}

	// 批量插入用户
	for i, user := range sampleUsers {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("创建用户 %d 失败: %v\n", i+1, err)
		} else {
			fmt.Printf("创建用户: %s (%s) - %s\n", user.Name, user.Email, user.Status)
		}
	}

	fmt.Printf("成功初始化 %d 个示例用户\n", len(sampleUsers))
}

// 辅助函数
func hashPassword(password string) string {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("密码加密失败:", err)
	}
	return string(hashedPwd)
}

func parseTime(timeStr string) *time.Time {
	t, err := time.Parse("2006-01-02", timeStr)
	if err != nil {
		log.Printf("时间解析失败: %v", err)
		return nil
	}
	return &t
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// GetUserStatistics 获取用户统计信息
func GetUserStatistics() map[string]interface{} {
	db := database.GetDB()

	var totalUsers int64
	var activeUsers int64
	var pendingUsers int64
	var inactiveUsers int64
	var emailUsers int64
	var googleUsers int64
	var appleUsers int64

	db.Model(&model.User{}).Count(&totalUsers)
	db.Model(&model.User{}).Where("status = ?", "active").Count(&activeUsers)
	db.Model(&model.User{}).Where("status = ?", "pending").Count(&pendingUsers)
	db.Model(&model.User{}).Where("status = ?", "inactive").Count(&inactiveUsers)
	db.Model(&model.User{}).Where("provider = ?", "email").Count(&emailUsers)
	db.Model(&model.User{}).Where("provider = ?", "google").Count(&googleUsers)
	db.Model(&model.User{}).Where("provider = ?", "apple").Count(&appleUsers)

	return map[string]interface{}{
		"total_users":    totalUsers,
		"active_users":   activeUsers,
		"pending_users":  pendingUsers,
		"inactive_users": inactiveUsers,
		"email_users":    emailUsers,
		"google_users":   googleUsers,
		"apple_users":    appleUsers,
	}
}
