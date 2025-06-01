package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser godoc
// @Summary 新建用户
// @Description 新建一个用户
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.UserReqCreate true "用户信息"
// @Success 200 {object} model.Response[model.User]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/users [post]
func CreateUser(c *gin.Context) {
	var req model.UserReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	user := model.User{
		Name:        req.Name,
		NameKana:    req.NameKana,
		Address:     req.Address,
		Gender:      &req.Gender,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    req.Password,
		GoogleID:    req.GoogleID,
		AppleID:     req.AppleID,
		Provider:    req.Provider,
		Status:      req.Status,
	}

	db := database.GetDB()
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response[model.User]{Success: true, Data: user})
}

// UpdateUser godoc
// @Summary 更新用户
// @Description 更新一个用户信息
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.UserReqEdit true "用户信息"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/users [put]
func UpdateUser(c *gin.Context) {
	var req model.UserReqEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	db := database.GetDB()
	var user model.User
	if err := db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "用户不存在"})
		return
	}

	db.Model(&user).Updates(req)
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// DeleteUser godoc
// @Summary 删除用户
// @Description 删除一个用户
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/users/{user_id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("user_id")
	db := database.GetDB()
	if err := db.Delete(&model.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// GetUser godoc
// @Summary 获取用户信息
// @Description 获取单个用户信息
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Success 200 {object} model.Response[model.User]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/users/{user_id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("user_id")
	db := database.GetDB()
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, model.Response[model.User]{Success: true, Data: user})
}

// ListUsers godoc
// @Summary 获取用户列表
// @Description 获取用户分页列表
// @Tags Users
// @Accept json
// @Produce json
// @Param req body model.UserReqList true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.User]
// @Failure 400 {object} model.BaseResponse
// @Router /api/users/list [post]
func ListUsers(c *gin.Context) {
	var req model.UserReqList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	db := database.GetDB()
	var users []model.User
	var total int64

	query := db.Model(&model.User{})
	if req.Keyword != "" {
		query = query.Where("name ILIKE ?", "%"+req.Keyword+"%")
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&users)

	c.JSON(http.StatusOK, model.ListResponse[model.User]{
		Success: true,
		Total:   total,
		List:    users,
	})
}

// UserProfile godoc
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} model.Response[model.User]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Security ApiKeyAuth
// @Router /api/auth/user/profile [get]
func UserProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	db := database.GetDB()
	var user model.User
	db.First(&user, userID)
	c.JSON(http.StatusOK, model.Response[model.User]{Success: true, Data: user})
}

// GetUserStatistics godoc
// @Summary 获取用户统计信息
// @Description 获取用户总数、活跃用户数、各种状态和提供商的统计信息
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} model.Response[map[string]interface{}]
// @Failure 500 {object} model.BaseResponse
// @Router /api/users/statistics [get]
func GetUserStatistics(c *gin.Context) {
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

	statistics := map[string]interface{}{
		"total_users":    totalUsers,
		"active_users":   activeUsers,
		"pending_users":  pendingUsers,
		"inactive_users": inactiveUsers,
		"email_users":    emailUsers,
		"google_users":   googleUsers,
		"apple_users":    appleUsers,
		"timestamp":      time.Now(),
	}

	c.JSON(http.StatusOK, model.Response[map[string]interface{}]{
		Success: true,
		Data:    statistics,
	})
}

// InitializeSampleUsers godoc
// @Summary 初始化示例用户数据
// @Description 创建示例用户数据用于测试和演示
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/users/init-sample [post]
func InitializeSampleUsers(c *gin.Context) {
	db := database.GetDB()

	// 检查是否已有用户数据
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)

	if userCount > 0 {
		c.JSON(http.StatusOK, model.BaseResponse{
			Success:    true,
			ErrMessage: fmt.Sprintf("数据库中已有 %d 个用户，跳过初始化", userCount),
		})
		return
	}

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
	var successCount int
	for _, user := range sampleUsers {
		if err := db.Create(&user).Error; err != nil {
			log.Printf("创建用户失败: %v\n", err)
		} else {
			successCount++
		}
	}

	c.JSON(http.StatusOK, model.BaseResponse{
		Success:    true,
		ErrMessage: fmt.Sprintf("成功初始化 %d 个示例用户", successCount),
	})
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
