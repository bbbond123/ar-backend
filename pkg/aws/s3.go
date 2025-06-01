package aws

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Service struct {
	session    *session.Session
	s3Client   *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	bucket     string
	region     string
	baseURL    string
}

// NewS3Service 创建新的 S3 服务实例
func NewS3Service() (*S3Service, error) {
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_S3_BUCKET")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	baseURL := os.Getenv("S3_BASE_URL")

	if region == "" || bucket == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("AWS 配置缺失，请检查环境变量")
	}

	// 创建 AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"", // token
		),
	})
	if err != nil {
		return nil, fmt.Errorf("创建 AWS session 失败: %v", err)
	}

	return &S3Service{
		session:    sess,
		s3Client:   s3.New(sess),
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
		bucket:     bucket,
		region:     region,
		baseURL:    baseURL,
	}, nil
}

// UploadFile 上传文件到 S3
func (s *S3Service) UploadFile(fileData []byte, fileName string, contentType string) (string, error) {
	// 生成唯一的文件名，避免冲突
	timestamp := time.Now().Unix()
	extension := filepath.Ext(fileName)
	baseName := fileName[:len(fileName)-len(extension)]
	uniqueFileName := fmt.Sprintf("%s_%d%s", baseName, timestamp, extension)

	// 上传文件
	uploadInput := &s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(uniqueFileName),
		Body:        bytes.NewReader(fileData),
		ContentType: aws.String(contentType),
	}

	result, err := s.uploader.Upload(uploadInput)
	if err != nil {
		return "", fmt.Errorf("上传文件到 S3 失败: %v", err)
	}

	return result.Location, nil
}

// DownloadFile 从 S3 下载文件
func (s *S3Service) DownloadFile(key string) ([]byte, error) {
	buf := aws.NewWriteAtBuffer([]byte{})

	downloadInput := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.downloader.Download(buf, downloadInput)
	if err != nil {
		return nil, fmt.Errorf("从 S3 下载文件失败: %v", err)
	}

	return buf.Bytes(), nil
}

// DeleteFile 从 S3 删除文件
func (s *S3Service) DeleteFile(key string) error {
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.s3Client.DeleteObject(deleteInput)
	if err != nil {
		return fmt.Errorf("从 S3 删除文件失败: %v", err)
	}

	return nil
}

// GetFileURL 获取文件的公共访问 URL
func (s *S3Service) GetFileURL(key string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, key)
}

// GetPresignedURL 获取文件的预签名 URL（用于临时访问私有文件）
func (s *S3Service) GetPresignedURL(key string, expiration time.Duration) (string, error) {
	req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(expiration)
	if err != nil {
		return "", fmt.Errorf("生成预签名 URL 失败: %v", err)
	}

	return urlStr, nil
}

// ListFiles 列出 S3 存储桶中的文件
func (s *S3Service) ListFiles(prefix string, maxKeys int64) ([]*s3.Object, error) {
	listInput := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s.bucket),
		MaxKeys: aws.Int64(maxKeys),
	}

	if prefix != "" {
		listInput.Prefix = aws.String(prefix)
	}

	result, err := s.s3Client.ListObjectsV2(listInput)
	if err != nil {
		return nil, fmt.Errorf("列出 S3 文件失败: %v", err)
	}

	return result.Contents, nil
}

// GetFileInfo 获取文件信息
func (s *S3Service) GetFileInfo(key string) (*s3.HeadObjectOutput, error) {
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	result, err := s.s3Client.HeadObject(headInput)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	return result, nil
}

// TestConnection 测试 S3 连接
func (s *S3Service) TestConnection() error {
	_, err := s.s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return fmt.Errorf("S3 连接测试失败: %v", err)
	}

	return nil
} 