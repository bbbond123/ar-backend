-- 为 articles 表添加图片文件关联字段
-- 执行时间: $(date)

-- 添加 image_file_id 字段
ALTER TABLE articles 
ADD COLUMN IF NOT EXISTS image_file_id INTEGER;

-- 添加外键约束
ALTER TABLE articles 
ADD CONSTRAINT IF NOT EXISTS fk_articles_image_file 
FOREIGN KEY (image_file_id) REFERENCES files(file_id) ON DELETE SET NULL;

-- 添加字段注释
COMMENT ON COLUMN articles.image_file_id IS '关联的图片文件ID，引用files表';

-- 创建索引优化查询性能
CREATE INDEX IF NOT EXISTS idx_articles_image_file_id ON articles(image_file_id);
CREATE INDEX IF NOT EXISTS idx_articles_category ON articles(category);
CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at);

-- 显示表结构确认修改
\d articles; 