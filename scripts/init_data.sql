-- 清理现有数据
DELETE FROM refresh_tokens;
DELETE FROM users;
DELETE FROM languages;
DELETE FROM tags;

-- 插入语言数据 (根据实际表结构)
INSERT INTO languages (language_name, display_order, is_active, created_at, updated_at) VALUES 
('日本語', 1, true, NOW(), NOW()),
('English', 2, true, NOW(), NOW()),
('中文', 3, true, NOW(), NOW()),
('한국어', 4, true, NOW(), NOW());

-- 插入标签数据 (根据实际表结构)
INSERT INTO tags (tag_name, is_active, created_at, updated_at) VALUES 
('観光地', true, NOW(), NOW()),
('レストラン', true, NOW(), NOW()),
('ホテル', true, NOW(), NOW()),
('ショッピング', true, NOW(), NOW()),
('文化', true, NOW(), NOW()),
('歴史', true, NOW(), NOW()),
('自然', true, NOW(), NOW()),
('体験', true, NOW(), NOW());

-- 插入管理员用户 (密码是: admin123 的bcrypt哈希)
INSERT INTO users (
    name, name_kana, email, password, provider, status, 
    phone_number, address, gender,
    created_at, updated_at
) VALUES (
    '系统管理员', 
    'システムカンリシャ', 
    'admin@ar-backend.com', 
    '$2a$10$Td7UYQFuXXj9LKlI2yGc4.3jVY4J9Z5UeGwKbQB6kPwXf8Jj.nLG6', 
    'email', 
    'active',
    '000-0000-0000',
    '東京都千代田区千代田1-1-1',
    'other',
    NOW(),
    NOW()
);

-- 显示结果
SELECT 'Data initialization completed!' as message; 