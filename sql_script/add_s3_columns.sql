-- 为 files 表添加 S3 存储支持的字段
-- 执行时间: $(date)

-- 添加 S3 相关字段
ALTER TABLE files 
ADD COLUMN IF NOT EXISTS s3_key VARCHAR(500),
ADD COLUMN IF NOT EXISTS s3_url VARCHAR(1000);

-- 修改 file_data 字段为可选（移除 NOT NULL 约束）
ALTER TABLE files 
ALTER COLUMN file_data DROP NOT NULL;

-- 添加字段注释
COMMENT ON COLUMN files.s3_key IS 'S3对象键，用于标识S3中的文件';
COMMENT ON COLUMN files.s3_url IS 'S3文件的访问URL';

-- 创建索引优化查询性能
CREATE INDEX IF NOT EXISTS idx_files_s3_key ON files(s3_key);
CREATE INDEX IF NOT EXISTS idx_files_location ON files(location);
CREATE INDEX IF NOT EXISTS idx_files_related_id ON files(related_id);

-- 显示表结构确认修改
\d files; 