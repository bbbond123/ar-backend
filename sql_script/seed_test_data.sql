-- 用户表
INSERT INTO users (name, name_kana, birth, address, gender, phone_number, email, password, avatar, google_id, apple_id, provider, status, created_at)
VALUES
  ('张三', 'チョウ サン', '1990-01-01', '北京市朝阳区', '1', '13800000001', 'zhangsan@example.com', 'password123', NULL, NULL, NULL, 'email', 'active', NOW()),
  ('李四', 'リ シ', '1985-05-20', '上海市浦东新区', '2', '13900000002', 'lisi@example.com', 'password456', NULL, NULL, NULL, 'google', 'active', NOW());

-- 设施表
INSERT INTO facilities (facility_name, location, description_text, latitude, longitude, person_id, created_at)
VALUES
  ('天安门广场', '北京市东城区', '中国著名地标', 39.9087, 116.3975, 1, NOW()),
  ('东方明珠塔', '上海市浦东新区', '上海地标建筑', 31.2397, 121.4998, 2, NOW());

-- 商店表
INSERT INTO stores (store_name, store_category, location, description_text, address, latitude, longitude, business_hours, rating_score, phone_number, created_at)
VALUES
  ('老北京炸酱面', '餐饮', '北京市', '地道北京风味', '北京市朝阳区建国路', 39.9087, 116.3975, '10:00-22:00', 4.5, '010-88888888', NOW()),
  ('上海小笼包', '餐饮', '上海市', '正宗小笼包', '上海市浦东新区世纪大道', 31.2397, 121.4998, '09:00-21:00', 4.7, '021-66666666', NOW());

-- 文章表
INSERT INTO articles (title, body_text, category, like_count, comment_count, created_at)
VALUES
  ('北京旅游攻略', '这里是北京旅游的详细攻略...', '旅游', 10, 2, NOW()),
  ('上海美食推荐', '上海有哪些必吃美食...', '美食', 8, 1, NOW());

-- 标签表
INSERT INTO tags (tag_name, is_active, created_at)
VALUES
  ('历史', true, NOW()),
  ('美食', true, NOW());

-- 菜单表
INSERT INTO menus (menu_name, menu_code, display_order, is_active, created_at)
VALUES
  ('首页', 'home', 1, true, NOW()),
  ('设置', 'settings', 2, true, NOW());

-- 文件表
INSERT INTO files (file_name, file_type, file_size, file_data, location, related_id, created_at)
VALUES
  ('beijing.jpg', 'image/jpeg', 102400, NULL, '北京市', 1, NOW()),
  ('shanghai.jpg', 'image/jpeg', 204800, NULL, '上海市', 2, NOW());

-- 通知表
INSERT INTO notices (title, content, notice_type, user_id, published_at, is_active, is_read, created_at)
VALUES
  ('欢迎使用', '欢迎来到旅游AR平台！', true, 1, NOW(), true, false, NOW()),
  ('系统维护', '本周末将进行系统维护。', false, 2, NOW(), true, false, NOW());

-- 访问历史表
INSERT INTO visit_history (user_id, facility_id, scan_at, is_active, created_at)
VALUES
  (1, 1, NOW(), true, NOW()),
  (2, 2, NOW(), true, NOW());

-- 评论表
INSERT INTO comments (article_id, user_id, comment_text, created_at, is_published)
VALUES
  (1, 1, '很有用的攻略！', NOW(), true),
  (2, 2, '美食推荐很棒！', NOW(), true);

-- 标签关联表
INSERT INTO taggings (tag_id, taggable_type, taggable_id, created_at)
VALUES
  (1, 'Article', 1, NOW()),
  (2, 'Article', 2, NOW()); 