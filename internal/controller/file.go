package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/aws"
	"ar-backend/pkg/database"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UploadFile godoc
// @Summary 上传文件到 S3
// @Description 上传文件到 S3 存储并保存记录到数据库
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param location formData string true "位置"
// @Param related_id formData int true "关联ID"
// @Success 200 {object} model.Response[model.File]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/files/upload [post]
func UploadFile(c *gin.Context) {
	// 解析表单数据
	var req model.FileUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	// 获取上传的文件
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: "未找到文件"})
		return
	}

	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: "打开文件失败"})
		return
	}
	defer file.Close()

	// 读取文件数据
	fileData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: "读取文件失败"})
		return
	}

	// 获取文件信息
	fileName := fileHeader.Filename
	fileSize := int(fileHeader.Size)
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		// 根据文件扩展名猜测 MIME 类型
		ext := filepath.Ext(fileName)
		contentType = getContentType(ext)
	}

	// 初始化 S3 服务
	s3Service, err := aws.NewS3Service()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: "S3 服务初始化失败: " + err.Error()})
		return
	}

	// 上传到 S3
	s3URL, err := s3Service.UploadFile(fileData, fileName, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: "文件上传失败: " + err.Error()})
		return
	}

	// 从 S3 URL 提取 S3 Key
	s3Key := extractS3Key(s3URL)

	// 保存文件记录到数据库
	fileRecord := model.File{
		FileName:  fileName,
		FileType:  contentType,
		FileSize:  fileSize,
		S3Key:     s3Key,
		S3URL:     s3URL,
		Location:  req.Location,
		RelatedID: req.RelatedID,
	}

	db := database.GetDB()
	if err := db.Create(&fileRecord).Error; err != nil {
		// 如果数据库保存失败，删除已上传的 S3 文件
		s3Service.DeleteFile(s3Key)
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response[model.File]{Success: true, Data: fileRecord})
}

// CreateFile godoc
// @Summary 新建文件
// @Description 新建一个文件（支持 S3 存储）
// @Tags Files
// @Accept json
// @Produce json
// @Param file body model.FileReqCreate true "文件信息"
// @Success 200 {object} model.Response[model.File]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/files [post]
func CreateFile(c *gin.Context) {
	var req model.FileReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	// 如果有文件数据且配置了 S3，则上传到 S3
	var s3Key, s3URL string
	if len(req.FileData) > 0 {
		s3Service, err := aws.NewS3Service()
		if err == nil {
			// S3 可用，上传到 S3
			s3URL, err = s3Service.UploadFile(req.FileData, req.FileName, req.FileType)
			if err == nil {
				s3Key = extractS3Key(s3URL)
				// 清空 FileData，只保存 S3 信息
				req.FileData = nil
			}
		}
	}

	file := model.File{
		FileName:  req.FileName,
		FileType:  req.FileType,
		FileSize:  req.FileSize,
		FileData:  req.FileData, // 如果 S3 上传成功则为 nil
		S3Key:     s3Key,
		S3URL:     s3URL,
		Location:  req.Location,
		RelatedID: req.RelatedID,
	}

	db := database.GetDB()
	if err := db.Create(&file).Error; err != nil {
		// 如果数据库保存失败且已上传到 S3，删除 S3 文件
		if s3Key != "" {
			if s3Service, err := aws.NewS3Service(); err == nil {
				s3Service.DeleteFile(s3Key)
			}
		}
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.File]{Success: true, Data: file})
}

// UpdateFile godoc
// @Summary 更新文件
// @Description 更新文件信息
// @Tags Files
// @Accept json
// @Produce json
// @Param file body model.FileReqEdit true "文件信息"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/files [put]
func UpdateFile(c *gin.Context) {
	var req model.FileReqEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var file model.File
	if err := db.First(&file, req.FileID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "文件不存在"})
		return
	}

	// 如果更新文件数据且配置了 S3
	if len(req.FileData) > 0 {
		s3Service, err := aws.NewS3Service()
		if err == nil {
			// 删除旧的 S3 文件
			if file.S3Key != "" {
				s3Service.DeleteFile(file.S3Key)
			}

			// 上传新文件到 S3
			fileName := req.FileName
			if fileName == "" {
				fileName = file.FileName
			}
			fileType := req.FileType
			if fileType == "" {
				fileType = file.FileType
			}

			s3URL, err := s3Service.UploadFile(req.FileData, fileName, fileType)
			if err == nil {
				file.S3Key = extractS3Key(s3URL)
				file.S3URL = s3URL
				// 清空数据库中的 FileData
				file.FileData = nil
			}
		}
	}

	db.Model(&file).Updates(req)
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// DeleteFile godoc
// @Summary 删除文件
// @Description 删除一个文件（同时删除 S3 文件）
// @Tags Files
// @Accept json
// @Produce json
// @Param file_id path int true "文件ID"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/files/{file_id} [delete]
func DeleteFile(c *gin.Context) {
	id := c.Param("file_id")
	fileID, _ := strconv.Atoi(id)
	
	db := database.GetDB()
	var file model.File
	if err := db.First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "文件不存在"})
		return
	}

	// 如果文件存储在 S3，先删除 S3 文件
	if file.S3Key != "" {
		s3Service, err := aws.NewS3Service()
		if err == nil {
			if err := s3Service.DeleteFile(file.S3Key); err != nil {
				// S3 删除失败，记录日志但继续删除数据库记录
				// 可以添加日志记录
			}
		}
	}

	if err := db.Delete(&model.File{}, fileID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// GetFile godoc
// @Summary 获取文件
// @Description 获取单个文件信息
// @Tags Files
// @Accept json
// @Produce json
// @Param file_id path int true "文件ID"
// @Success 200 {object} model.Response[model.File]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/files/{file_id} [get]
func GetFile(c *gin.Context) {
	id := c.Param("file_id")
	fileID, _ := strconv.Atoi(id)
	db := database.GetDB()
	var file model.File
	if err := db.First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "文件不存在"})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.File]{Success: true, Data: file})
}

// DownloadFile godoc
// @Summary 下载文件
// @Description 从 S3 或数据库下载文件
// @Tags Files
// @Param file_id path int true "文件ID"
// @Success 200 {file} binary "文件内容"
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/files/{file_id}/download [get]
func DownloadFile(c *gin.Context) {
	id := c.Param("file_id")
	fileID, _ := strconv.Atoi(id)
	
	db := database.GetDB()
	var file model.File
	if err := db.First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "文件不存在"})
		return
	}

	var fileData []byte
	var err error

	// 优先从 S3 下载
	if file.S3Key != "" {
		s3Service, s3Err := aws.NewS3Service()
		if s3Err == nil {
			fileData, err = s3Service.DownloadFile(file.S3Key)
		} else {
			err = s3Err
		}
	} else if len(file.FileData) > 0 {
		// 从数据库获取
		fileData = file.FileData
	} else {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "文件数据不存在"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: "下载文件失败: " + err.Error()})
		return
	}

	// 设置响应头
	c.Header("Content-Type", file.FileType)
	c.Header("Content-Disposition", "attachment; filename=\""+file.FileName+"\"")
	c.Header("Content-Length", strconv.Itoa(len(fileData)))

	// 返回文件数据
	c.Data(http.StatusOK, file.FileType, fileData)
}

// ListFiles godoc
// @Summary 获取文件列表
// @Description 获取文件分页列表
// @Tags Files
// @Accept json
// @Produce json
// @Param req body model.FileReqList true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.File]
// @Failure 400 {object} model.BaseResponse
// @Router /api/files/list [post]
func ListFiles(c *gin.Context) {
	var req model.FileReqList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var files []model.File
	var total int64

	query := db.Model(&model.File{})
	if req.Keyword != "" {
		query = query.Where("file_name LIKE ?", "%"+req.Keyword+"%")
	}
	
	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&files)

	c.JSON(http.StatusOK, model.ListResponse[model.File]{
		Success: true,
		Total:   total,
		List:    files,
	})
}

// TestS3Connection godoc
// @Summary 测试 S3 连接
// @Description 测试 AWS S3 连接状态
// @Tags Files
// @Accept json
// @Produce json
// @Success 200 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/files/test-s3 [get]
func TestS3Connection(c *gin.Context) {
	s3Service, err := aws.NewS3Service()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: "S3 服务初始化失败: " + err.Error()})
		return
	}

	if err := s3Service.TestConnection(); err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: "S3 连接测试失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.BaseResponse{Success: true, ErrMessage: "S3 连接正常"})
}

// 辅助函数：根据文件扩展名获取 MIME 类型
func getContentType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	case ".zip":
		return "application/zip"
	case ".mp4":
		return "video/mp4"
	case ".mp3":
		return "audio/mp3"
	default:
		return "application/octet-stream"
	}
}

// 辅助函数：从 S3 URL 提取 S3 Key
func extractS3Key(s3URL string) string {
	// S3 URL 格式通常是: https://bucket.s3.region.amazonaws.com/key
	// 或者: https://s3.region.amazonaws.com/bucket/key
	if len(s3URL) == 0 {
		return ""
	}
	
	// 简单的实现：从最后一个 '/' 后面提取
	lastSlash := len(s3URL) - 1
	for i := len(s3URL) - 1; i >= 0; i-- {
		if s3URL[i] == '/' {
			lastSlash = i
			break
		}
	}
	
	if lastSlash < len(s3URL)-1 {
		return s3URL[lastSlash+1:]
	}
	
	return ""
}
