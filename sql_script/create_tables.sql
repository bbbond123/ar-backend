-- 删除现有表（按照依赖关系顺序）
DROP TABLE IF EXISTS taggings;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS articles;
DROP TABLE IF EXISTS menus;
DROP TABLE IF EXISTS stores;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS languages;
DROP TABLE IF EXISTS visit_history;
DROP TABLE IF EXISTS notices;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS facilities;

-- 创建语言表
CREATE TABLE languages (
    language_id SERIAL PRIMARY KEY,
    language_name VARCHAR(50) NOT NULL,
    display_order INTEGER,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT unique_language_name UNIQUE (language_name)
);

-- 创建用户表
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    name_kana VARCHAR(50),
    birth DATE,
    address VARCHAR(255),
    gender CHAR(1),
    phone_number VARCHAR(15),
    email VARCHAR(255) NOT NULL,
    password VARCHAR(128),
    avatar VARCHAR(255),
    google_id VARCHAR(255),
    apple_id VARCHAR(255),
    provider VARCHAR(20) NOT NULL,
    verify_code VARCHAR(255),
    verify_code_expire TIMESTAMP,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT chk_gender CHECK (gender IN ('1', '2') OR gender IS NULL),
    CONSTRAINT chk_status CHECK (status IN ('pending', 'active', 'disabled')),
    CONSTRAINT chk_provider CHECK (provider IN ('email', 'google', 'apple')),
    CONSTRAINT unique_email UNIQUE (email),
    CONSTRAINT unique_google_id UNIQUE (google_id),
    CONSTRAINT unique_apple_id UNIQUE (apple_id)
);

-- 创建刷新令牌表
CREATE TABLE refresh_tokens (
    token_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT unique_refresh_token UNIQUE (refresh_token)
);

-- 创建设施表
CREATE TABLE facilities (
    facility_id SERIAL PRIMARY KEY,
    facility_name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    description_text TEXT,
    latitude DECIMAL(10,6) NOT NULL,
    longitude DECIMAL(10,6) NOT NULL,
    person_id INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_person_id FOREIGN KEY (person_id) REFERENCES users(user_id) ON DELETE SET NULL
);

-- 创建文件表
CREATE TABLE files (
    file_id SERIAL PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size INTEGER,
    file_data BYTEA,
    s3_key VARCHAR(500),
    s3_url VARCHAR(1000),
    location VARCHAR(255) NOT NULL,
    related_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建通知表
CREATE TABLE notices (
    notice_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    notice_type BOOLEAN NOT NULL,
    user_id INTEGER,
    published_at TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_notice_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE SET NULL
);

-- 创建访问历史表
CREATE TABLE visit_history (
    history_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    facility_id INTEGER NOT NULL,
    scan_at TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_visit_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_visit_facility_id FOREIGN KEY (facility_id) REFERENCES facilities(facility_id) ON DELETE CASCADE
);

-- 创建商店表
CREATE TABLE stores (
    store_id SERIAL PRIMARY KEY,
    store_name VARCHAR(255) NOT NULL,
    store_category VARCHAR(100) NOT NULL,
    location VARCHAR(255) NOT NULL,
    description_text TEXT,
    address VARCHAR(255) NOT NULL,
    latitude DECIMAL(10,6) NOT NULL,
    longitude DECIMAL(10,6) NOT NULL,
    business_hours VARCHAR(100) NOT NULL,
    rating_score DECIMAL(3,2) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建菜单表
CREATE TABLE menus (
    menu_id SERIAL PRIMARY KEY,
    menu_name VARCHAR(100) NOT NULL,
    menu_code VARCHAR(50) NOT NULL,
    display_order INTEGER,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT unique_menu_code UNIQUE (menu_code)
);

-- 创建文章表
CREATE TABLE articles (
    article_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    body_text TEXT NOT NULL,
    category VARCHAR(100),
    like_count INTEGER NOT NULL DEFAULT 0,
    article_image BYTEA,
    image_file_id INTEGER,
    comment_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_article_image_file_id FOREIGN KEY (image_file_id) REFERENCES files(file_id) ON DELETE SET NULL
);

-- 创建评论表
CREATE TABLE comments (
    comment_id SERIAL PRIMARY KEY,
    article_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    comment_text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    is_published BOOLEAN NOT NULL DEFAULT TRUE,
    reply_to_comment_id INTEGER,
    CONSTRAINT fk_comment_article_id FOREIGN KEY (article_id) REFERENCES articles(article_id) ON DELETE CASCADE,
    CONSTRAINT fk_comment_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_reply_comment_id FOREIGN KEY (reply_to_comment_id) REFERENCES comments(comment_id) ON DELETE SET NULL
);

-- 创建标签表
CREATE TABLE tags (
    tag_id SERIAL PRIMARY KEY,
    tag_name VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT unique_tag_name UNIQUE (tag_name)
);

-- 创建标签关联表
CREATE TABLE taggings (
    tagging_id SERIAL PRIMARY KEY,
    tag_id INTEGER NOT NULL,
    taggable_type VARCHAR(50) NOT NULL,
    taggable_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_tag_id FOREIGN KEY (tag_id) REFERENCES tags(tag_id) ON DELETE CASCADE,
    CONSTRAINT unique_taggable UNIQUE (tag_id, taggable_type, taggable_id)
);

-- 创建触发器函数：检查taggable_id的有效性
CREATE OR REPLACE FUNCTION check_taggable_id() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.taggable_type = 'Article' AND NOT EXISTS (
        SELECT 1 FROM articles WHERE article_id = NEW.taggable_id
    ) THEN
        RAISE EXCEPTION 'Invalid taggable_id % for taggable_type Article', NEW.taggable_id;
    ELSIF NEW.taggable_type = 'History' AND NOT EXISTS (
        SELECT 1 FROM visit_history WHERE history_id = NEW.taggable_id
    ) THEN
        RAISE EXCEPTION 'Invalid taggable_id % for taggable_type History', NEW.taggable_id;
    ELSIF NEW.taggable_type = 'Store' AND NOT EXISTS (
        SELECT 1 FROM stores WHERE store_id = NEW.taggable_id
    ) THEN
        RAISE EXCEPTION 'Invalid taggable_id % for taggable_type Store', NEW.taggable_id;
    ELSIF NEW.taggable_type = 'Comment' AND NOT EXISTS (
        SELECT 1 FROM comments WHERE comment_id = NEW.taggable_id
    ) THEN
        RAISE EXCEPTION 'Invalid taggable_id % for taggable_type Comment', NEW.taggable_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建触发器
CREATE TRIGGER taggings_check_taggable_id
    BEFORE INSERT OR UPDATE ON taggings
    FOR EACH ROW
    EXECUTE FUNCTION check_taggable_id();

-- 添加表注释
COMMENT ON TABLE languages IS '语言管理表';
COMMENT ON TABLE users IS '用户信息表';
COMMENT ON TABLE refresh_tokens IS '刷新令牌管理表';
COMMENT ON TABLE facilities IS '设施信息管理表';
COMMENT ON TABLE files IS '文件管理表';
COMMENT ON TABLE notices IS '通知管理表';
COMMENT ON TABLE visit_history IS '访问历史记录表';
COMMENT ON TABLE stores IS '商店信息管理表';
COMMENT ON TABLE menus IS '菜单管理表';
COMMENT ON TABLE articles IS '文章管理表';
COMMENT ON TABLE comments IS '评论管理表';
COMMENT ON TABLE tags IS '标签管理表';
COMMENT ON TABLE taggings IS '标签关联表';

-- 创建索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_google_id ON users(google_id);
CREATE INDEX idx_users_apple_id ON users(apple_id);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(refresh_token);
CREATE INDEX idx_facilities_location ON facilities(location);
CREATE INDEX idx_files_related_id ON files(related_id);
CREATE INDEX idx_notices_user_id ON notices(user_id);
CREATE INDEX idx_visit_history_user_id ON visit_history(user_id);
CREATE INDEX idx_visit_history_facility_id ON visit_history(facility_id);
CREATE INDEX idx_stores_location ON stores(location);
CREATE INDEX idx_articles_category ON articles(category);
CREATE INDEX idx_comments_article_id ON comments(article_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_taggings_tag_id ON taggings(tag_id);
CREATE INDEX idx_taggings_taggable ON taggings(taggable_type, taggable_id); 