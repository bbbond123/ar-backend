-- 既存のテーブルを削除（依存関係を考慮して逆順で削除）
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

-- 観光施設テーブル生成
CREATE TABLE facilities (
    facility_id SERIAL PRIMARY KEY,                   -- 施設ID: 施設を一意に識別するID
    facility_name VARCHAR(255) NOT NULL,              -- 施設名（全角）
    location VARCHAR(255) NOT NULL,                   -- 所在地: 施設の住所（全角）
    description_text TEXT,                            -- 説明（全角）
    latitude DECIMAL(10,6) NOT NULL,                  -- 緯度: 施設の緯度情報（半角）
    longitude DECIMAL(10,6) NOT NULL,                 -- 経度: 施設の経度情報（半角）
    person_id INTEGER,                                -- 関連人物ID（半角）
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 作成日
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 更新日
);
-- テーブルコメント
COMMENT ON TABLE facilities IS '施設情報管理テーブル';
-- カラムコメント
COMMENT ON COLUMN facilities.facility_id IS '施設ID: 施設を一意に識別するID';
COMMENT ON COLUMN facilities.facility_name IS '施設名（全角）';
COMMENT ON COLUMN facilities.location IS '所在地: 施設の住所（全角）';
COMMENT ON COLUMN facilities.description_text IS '説明（全角）';
COMMENT ON COLUMN facilities.latitude IS '緯度: 施設の緯度情報（半角）';
COMMENT ON COLUMN facilities.longitude IS '経度: 施設の経度情報（半角）';
COMMENT ON COLUMN facilities.person_id IS '関連人物ID（半角）';
COMMENT ON COLUMN facilities.created_at IS '作成日';
COMMENT ON COLUMN facilities.updated_at IS '更新日';

-- ファイルテーブル生成
CREATE TABLE files (
    file_id SERIAL PRIMARY KEY,                       -- ファイルID: ファイルを一意に識別するID
    file_name VARCHAR(255) NOT NULL,                  -- ファイル名（半角）
    file_type VARCHAR(50) NOT NULL,                   -- ファイル種類: MIMEタイプ（例: image/png）（半角）
    file_size INTEGER,                                -- ファイルサイズ（半角）
    file_data BYTEA NOT NULL,                         -- ファイルデータ: 画像データのバイナリ情報
    location VARCHAR(255) NOT NULL,                   -- 所在地（例：東京都北区）（全角）
    related_id INTEGER NOT NULL,                      -- 関連ID: 画像の関連データ（例: ユーザ, 店舗など）（半角）
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 作成日
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 更新日
);
-- テーブルコメント
COMMENT ON TABLE files IS 'ファイル管理テーブル';
-- カラムコメント
COMMENT ON COLUMN files.file_id IS 'ファイルID: ファイルを一意に識別するID';
COMMENT ON COLUMN files.file_name IS 'ファイル名（半角）';
COMMENT ON COLUMN files.file_type IS 'ファイル種類: MIMEタイプ（例: image/png）（半角）';
COMMENT ON COLUMN files.file_size IS 'ファイルサイズ（半角）';
COMMENT ON COLUMN files.file_data IS 'ファイルデータ: 画像データのバイナリ情報';
COMMENT ON COLUMN files.location IS '所在地（例：東京都北区）（全角）';
COMMENT ON COLUMN files.related_id IS '関連ID: 画像の関連データ（例: ユーザ, 店舗など）（半角）';
COMMENT ON COLUMN files.created_at IS '作成日';
COMMENT ON COLUMN files.updated_at IS '更新日';

-- お知らせテーブル生成
CREATE TABLE notices (
    notice_id SERIAL PRIMARY KEY,                      -- お知らせID: お知らせを一意に識別するID
    title VARCHAR(255) NOT NULL,                      -- タイトル（全角）
    content TEXT NOT NULL,                            -- 内容（全角）
    notice_type BOOLEAN NOT NULL,                     -- お知らせ種別: 1=公開、0=個人（半角）
    user_id INTEGER,                                  -- 受信ユーザID（半角）（ユーザテーブルのFK）
    published_at TIMESTAMP NOT NULL,                  -- 公開日時
    is_active BOOLEAN NOT NULL,                       -- 有効フラグ: 1=公開中、0=非公開（半角）
    is_read BOOLEAN DEFAULT FALSE,                    -- 既読フラグ: 1=既読、0=未読（半角）
    created_at TIMESTAMP NOT NULL,                    -- 作成日
    updated  updated_at TIMESTAMP                             -- 更新日
);
-- テーブルコメント
COMMENT ON TABLE notices IS 'お知らせテーブル';
-- カラムコメント
COMMENT ON COLUMN notices.notice_id IS 'お知らせID: お知らせを一意に識別するID';
COMMENT ON COLUMN notices.title IS 'タイトル（全角）';
COMMENT ON COLUMN notices.content IS '内容（全角）';
COMMENT ON COLUMN notices.notice_type IS 'お知らせ種別: 1=公開、0=個人（半角）';
COMMENT ON COLUMN notices.user_id IS '受信ユーザID（半角）（ユーザテーブルのFK）';
COMMENT ON COLUMN notices.published_at IS '公開日時';
COMMENT ON COLUMN notices.is_active IS '有効フラグ: 1=公開中、0=非公開（半角）';
COMMENT ON COLUMN notices.is_read IS '既読フラグ: 1=既読、0=未読（半角）';
COMMENT ON COLUMN notices.created_at IS '作成日';
COMMENT ON COLUMN notices.updated_at IS '更新日';

-- 履歴テーブル生成
CREATE TABLE visit_history (
    history_id SERIAL PRIMARY KEY,                    -- 履歴ID: 訪問履歴を一意に識別するID
    user_id INTEGER NOT NULL,                        -- ユーザーID（半角）
    facility_id INTEGER NOT NULL,                    -- 観光施設ID（半角）
    scan_at TIMESTAMP NOT NULL,                      -- スキャン日時
    is_active BOOLEAN NOT NULL DEFAULT TRUE,         -- 有効フラグ: 1=有効、0=無効（半角）
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 作成日
    updated_at TIMESTAMP                            -- 更新日
);
-- テーブルコメント
COMMENT ON TABLE visit_history IS '施設訪問履歴テーブル';
-- カラムコメント
COMMENT ON COLUMN visit_history.history_id IS '履歴ID: 訪問履歴を一意に識別するID';
COMMENT ON COLUMN visit_history.user_id IS 'ユーザーID（半角）';
COMMENT ON COLUMN visit_history.facility_id IS '観光施設ID（半角）';
COMMENT ON COLUMN visit_history.scan_at IS 'スキャン日時';
COMMENT ON COLUMN visit_history.is_active IS '有効フラグ: 1=有効、0=無効（半角）';
COMMENT ON COLUMN visit_history.created_at IS '作成日';
COMMENT ON COLUMN visit_history.updated_at IS '更新日';

-- 言語テーブル生成
CREATE TABLE languages (
    language_id SERIAL PRIMARY KEY,                   -- 言語ID: 言語を一意に識別するID
    language_name VARCHAR(50) NOT NULL,               -- 言語名（全角）
    display_order INTEGER,                            -- 表示順: UIでの言語選択時の表示順（半角）
    is_active BOOLEAN NOT NULL DEFAULT TRUE,          -- 有効フラグ: 使用可能な言語かどうかを管理（1=有効、0=無効）（半角）
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 作成日
    updated_at TIMESTAMP                             -- 更新日
);
-- テーブルコメント
COMMENT ON TABLE languages IS '言語マスタテーブル';
-- カラムコメント
COMMENT ON COLUMN languages.language_id IS '言語ID: 言語を一意に識別するID';
COMMENT ON COLUMN languages.language_name IS '言語名（全角）';
COMMENT ON COLUMN languages.display_order IS '表示順: UIでの言語選択時の表示順（半角）';
COMMENT ON COLUMN languages.is_active IS '有効フラグ: 使用可能な言語かどうかを管理（1=有効、0=無効）（半角）';
COMMENT ON COLUMN languages.created_at IS '作成日';
COMMENT ON COLUMN languages.updated_at IS '更新日';

-- ユーザテーブル生成（更新済み：サポートGoogle、Apple、メールログイン）
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,                        -- ユーザID: ユーザを一意に識別するID
    name VARCHAR(50),                                 -- 氏名（全角、可空）
    name_kana VARCHAR(50),                           -- 氏名（カナ）（全角、可空）
    birth DATE,                                       -- 生年月日（可空）
    address VARCHAR(255),                             -- 住所（全角、可空）
    gender CHAR(1),                                   -- 性別（半角、可空）
    phone_number VARCHAR(15),                         -- 電話番号（半角、可空）
    email VARCHAR(255) NOT NULL,                      -- メールアドレス（半角、必須）
    password VARCHAR(128),                            -- パスワード（半角、可空、メールログイン用）
    google_id VARCHAR(255),                           -- GoogleログインのユニークID（半角、可空）
    apple_id VARCHAR(255),                            -- AppleログインのユニークID（半角、可空）
    provider VARCHAR(20) NOT NULL,                    -- ログイン方式: email, google, apple（半角）
    status VARCHAR(1) NOT NULL,                       -- アカウント状態: 1=有効、0=無効、2=凍結（半角）
    created_at TIMESTAMP NOT NULL,                     -- 登録日
    updated_at TIMESTAMP,                              -- 更新日
    CONSTRAINT chk_gender CHECK (gender IN ('1', '2') OR gender IS NULL), -- 性別チェック制約
    CONSTRAINT chk_status CHECK (status IN ('0', '1', '2')), -- ステータスチェック制約
    CONSTRAINT chk_provider CHECK (provider IN ('email', 'google', 'apple')), -- ログイン方式チェック制約
    CONSTRAINT unique_email UNIQUE (email)             -- メールアドレスのユニーク制約
);
-- テーブルコメント
COMMENT ON TABLE users IS 'ユーザ情報テーブル';
-- カラムコメント
COMMENT ON COLUMN users.user_id IS 'ユーザID: ユーザを一意に識別するID';
COMMENT ON COLUMN users.name IS '氏名（全角、可空）';
COMMENT ON COLUMN users.name_kana IS '氏名（カナ）（全角、可空）';
COMMENT ON COLUMN users.birth IS '生年月日（可空）';
COMMENT ON COLUMN users.address IS '住所（全角、可空）';
COMMENT ON COLUMN users.gender IS '性別（半角、可空）';
COMMENT ON COLUMN users.phone_number IS '電話番号（半角、可空）';
COMMENT ON COLUMN users.email IS 'メールアドレス（半角、必須）';
COMMENT ON COLUMN users.password IS 'パスワード（半角、可空、メールログイン用）';
COMMENT ON COLUMN users.google_id IS 'GoogleログインのユニークID（半角、可空）';
COMMENT ON COLUMN users.apple_id IS 'AppleログインのユニークID（半角、可空）';
COMMENT ON COLUMN users.provider IS 'ログイン方式: email, google, apple（半角）';
COMMENT ON COLUMN users.status IS 'アカウント状態: 1=有効、0=無効、2=凍結（半角）';
COMMENT ON COLUMN users.created_at IS '登録日';
COMMENT ON COLUMN users.updated_at IS '更新日';

-- リフレッシュトークンテーブル生成
CREATE TABLE refresh_tokens (
    token_id SERIAL PRIMARY KEY,                      -- トークンID: リフレッシュトークンを一意に識別するID
    user_id INTEGER NOT NULL,                        -- ユーザID（ユーザテーブルのFK）
    refresh_token VARCHAR(255) NOT NULL,              -- リフレッシュトークン（半角）
    expires_at TIMESTAMP NOT NULL,                    -- 有効期限
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 作成日
    revoked BOOLEAN NOT NULL DEFAULT FALSE,           -- 無効化フラグ: true=無効、false=有効
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) -- ユーザIDの外部キー制約
);
-- テーブルコメント
COMMENT ON TABLE refresh_tokens IS 'リフレッシュトークン管理テーブル';
-- カラムコメント
COMMENT ON COLUMN refresh_tokens.token_id IS 'トークンID: リフレッシュトークンを一意に識別するID';
COMMENT ON COLUMN refresh_tokens.user_id IS 'ユーザID（ユーザテーブルのFK）';
COMMENT ON COLUMN refresh_tokens.refresh_token IS 'リフレッシュトークン（半角）';
COMMENT ON COLUMN refresh_tokens.expires_at IS '有効期限';
COMMENT ON COLUMN refresh_tokens.created_at IS '作成日';
COMMENT ON COLUMN refresh_tokens.revoked IS '無効化フラグ: true=無効、false=有効';

-- 店舗テーブル生成
CREATE TABLE stores (
    store_id SERIAL PRIMARY KEY,                      -- 店舗ID: 店舗を一意に識別するID
    store_name VARCHAR(255) NOT NULL,                 -- 店舗名（全角）
    store_category VARCHAR(100) NOT NULL,             -- 店舗カテゴリ（例：飲食店, お土産）（半角）
    location VARCHAR(255) NOT NULL,                   -- 所在地（例：東京都北区）（全角）
    description_text TEXT,                            -- 説明（全角）
    address VARCHAR(255) NOT NULL,                    -- 住所（全角）
    latitude DECIMAL(10,6) NOT NULL,                  -- 緯度: 店舗の緯度情報（半角）
    longitude DECIMAL(10,6) NOT NULL,                 -- 経度: 店舗の経度情報（半角）
    business_hours VARCHAR(100) NOT NULL,             -- 営業時間（半角）
    rating_score DECIMAL(3,2) NOT NULL,               -- 評価点数（半角）
    phone_number VARCHAR(20) NOT NULL,                -- 電話番号（半角）
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 作成日
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 更新日
);
-- テーブルコメント
COMMENT ON TABLE stores IS '店舗情報管理テーブル';
-- カラムコメント
COMMENT ON COLUMN stores.store_id IS '店舗ID: 店舗を一意に識別するID';
COMMENT ON COLUMN stores.store_name IS '店舗名（全角）';
COMMENT ON COLUMN stores.store_category IS '店舗カテゴリ（例：飲食店, お土産）（半角）';
COMMENT ON COLUMN stores.location IS '所在地（例：東京都北区）（全角）';
COMMENT ON COLUMN stores.description_text IS '説明（全角）';
COMMENT ON COLUMN stores.address IS '住所（全角）';
COMMENT ON COLUMN stores.latitude IS '緯度: 店舗の緯度情報（半角）';
COMMENT ON COLUMN stores.longitude IS '経度: 店舗の経度情報（半角）';
COMMENT ON COLUMN stores.business_hours IS '営業時間（半角）';
COMMENT ON COLUMN stores.rating_score IS '評価点数（半角）';
COMMENT ON COLUMN stores.phone_number IS '電話番号（半角）';
COMMENT ON COLUMN stores.created_at IS '作成日';
COMMENT ON COLUMN stores.updated_at IS '更新日';

-- メニューテーブル生成
CREATE TABLE menus (
    menu_id SERIAL PRIMARY KEY,                      -- メニューID: メニューを一意に識別するID
    menu_name VARCHAR(100) NOT NULL,                 -- メニュー名（全角）
    menu_code VARCHAR(50) NOT NULL,                  -- メニューコード: メニューを識別するためのコード（例: home, settings）（半角）
    display_order INTEGER,                           -- 表示順（半角）
    is_active BOOLEAN NOT NULL DEFAULT TRUE,         -- 有効フラグ: 1=有効、0=無効（半角）
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 作成日
    updated_at TIMESTAMP                            -- 更新日
);
-- テーブルコメント
COMMENT ON TABLE menus IS 'メニュー管理テーブル';
-- カラムコメント
COMMENT ON COLUMN menus.menu_id IS 'メニューID: メニューを一意に識別するID';
COMMENT ON COLUMN menus.menu_name IS 'メニュー名（全角）';
COMMENT ON COLUMN menus.menu_code IS 'メニューコード: メニューを識別するためのコード（例: home, settings）（半角）';
COMMENT ON COLUMN menus.display_order IS '表示順（半角）';
COMMENT ON COLUMN menus.is_active IS '有効フラグ: 1=有効、0=無効（半角）';
COMMENT ON COLUMN menus.created_at IS '作成日';
COMMENT ON COLUMN menus.updated_at IS '更新日';

-- 文章テーブル生成
CREATE TABLE articles (
    article_id SERIAL PRIMARY KEY,                     -- 文章ID: 文章を一意に識別するID
    title VARCHAR(255) NOT NULL,                       -- タイトル
    body_text TEXT NOT NULL,                           -- 本文
    category VARCHAR(100),                             -- カテゴリー
    like_count INTEGER NOT NULL,                       -- いいね数
    article_image BYTEA,                               -- 文章写真
    comment_count INTEGER NOT NULL,                    -- コメント数
    created_at TIMESTAMP NOT NULL,                     -- 作成日
    updated_at TIMESTAMP                               -- 更新日
);
-- テーブルコメント
COMMENT ON TABLE articles IS '文章テーブル';
-- カラムコメント
COMMENT ON COLUMN articles.article_id IS '文章ID: 文章を一意に識別するID';
COMMENT ON COLUMN articles.title IS 'タイトル（全角）';
COMMENT ON COLUMN articles.body_text IS '本文（全角）';
COMMENT ON COLUMN articles.category IS 'カテゴリー（全角）';
COMMENT ON COLUMN articles.like_count IS 'いいね数（半角）';
COMMENT ON COLUMN articles.article_image IS '文章写真（半角）';
COMMENT ON COLUMN articles.comment_count IS 'コメント数（半角）';
COMMENT ON COLUMN articles.created_at IS '作成日';
COMMENT ON COLUMN articles.updated_at IS '更新日';

-- コメントテーブル生成
CREATE TABLE comments (
    comment_id SERIAL PRIMARY KEY,                     -- コメントID: コメントを一意に識別するID
    article_id INTEGER NOT NULL,                       -- 文章ID: コメントが紐づく文章ID（文章テーブルのFK）
    user_id INTEGER NOT NULL,                          -- ユーザID: コメント投稿者のユーザID（ユーザテーブルのFK）
    comment_text TEXT NOT NULL,                        -- コメント本文（全角）
    created_at TIMESTAMP NOT NULL,                     -- 作成日
    updated_at TIMESTAMP,                              -- 更新日
    is_published BOOLEAN NOT NULL,                     -- 公開フラグ（半角）
    reply_to_comment_id INTEGER                       -- 返信先コメントID（半角）
);
-- テーブルコメント
COMMENT ON TABLE comments IS 'コメントテーブル';
-- カラムコメント
COMMENT ON COLUMN comments.comment_id IS 'コメントID: コメントを一意に識別するID';
COMMENT ON COLUMN comments.article_id IS '文章ID: コメントが紐づく文章ID（文章テーブルのFK）';
COMMENT ON COLUMN comments.user_id IS 'ユーザID: コメント投稿者のユーザID（ユーザテーブルのFK）';
COMMENT ON COLUMN comments.comment_text IS 'コメント本文（全角）';
COMMENT ON COLUMN comments.created_at IS '作成日';
COMMENT ON COLUMN comments.updated_at IS '更新日';
COMMENT ON COLUMN comments.is_published IS '公開フラグ（半角）';
COMMENT ON COLUMN comments.reply_to_comment_id IS '返信先コメントID（半角）';

-- タグテーブル生成
CREATE TABLE tags (
    tag_id SERIAL PRIMARY KEY,                        -- タグID: タグを一意に識別するID
    tag_name VARCHAR(50) NOT NULL,                   -- タグ名（全角）
    is_active BOOLEAN NOT NULL DEFAULT TRUE,         -- 有効フラグ: 1=有効、0=無効（半角）
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 作成日
    updated_at TIMESTAMP                              -- 更新日
);
-- テーブルコメント
COMMENT ON TABLE tags IS 'タグ管理テーブル';
-- カラムコメント
COMMENT ON COLUMN tags.tag_id IS 'タグID: タグを一意に識別するID';
COMMENT ON COLUMN tags.tag_name IS 'タグ名（全角）';
COMMENT ON COLUMN tags.is_active IS '有効フラグ: 1=有効、0=無効（半角）';
COMMENT ON COLUMN tags.created_at IS '作成日';
COMMENT ON COLUMN tags.updated_at IS '更新日';

-- タグとエンティティの関連テーブル生成（多対多）
CREATE TABLE taggings (
    tagging_id SERIAL PRIMARY KEY,                    -- 関連ID: タグとエンティティの関連を一意に識別するID
    tag_id INTEGER NOT NULL,                         -- タグID（タグテーブルのFK）
    taggable_type VARCHAR(50) NOT NULL,              -- エ部分、エンティティ種別: 関連するテーブルの種類（例: Article, History, Store, Comment）
    taggable_id INTEGER NOT NULL,                    -- エンティティID: 関連するテーブルのレコードID
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 作成日
    updated_at TIMESTAMP,                            -- 更新日
    CONSTRAINT fk_tag_id FOREIGN KEY (tag_id) REFERENCES tags(tag_id) -- タグIDの外部キー制約
);
-- テーブルコメント
COMMENT ON TABLE taggings IS 'タグとエンティティ（文章、履歴、店舗、コメント）の関連管理テーブル';
-- カラムコメント
COMMENT ON COLUMN taggings.tagging_id IS '関連ID: タグとエンティティの関連を一意に識別するID';
COMMENT ON COLUMN taggings.tag_id IS 'タグID（タグテーブルのFK）';
COMMENT ON COLUMN taggings.taggable_type IS 'エンティティ種別: 関連するテーブルの種類（例: Article, History, Store, Comment）';
COMMENT ON COLUMN taggings.taggable_id IS 'エンティティID: 関連するテーブルのレコードID';
COMMENT ON COLUMN taggings.created_at IS '作成日';
COMMENT ON COLUMN taggings.updated_at IS '更新日';

-- トリガー関数: taggable_idの整合性をチェック
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

-- トリガーの作成
CREATE TRIGGER taggings_check_taggable_id
    BEFORE INSERT OR UPDATE ON taggings
    FOR EACH ROW
    EXECUTE FUNCTION check_taggable_id();