--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4 (Debian 17.4-1.pgdg120+2)
-- Dumped by pg_dump version 17.4 (Debian 17.4-1.pgdg120+2)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: articles; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.articles (
    article_id bigint NOT NULL,
    title character varying(255) NOT NULL,
    body_text text NOT NULL,
    category character varying(100),
    like_count bigint NOT NULL,
    article_image bytea,
    comment_count bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    image_file_id bigint
);


ALTER TABLE public.articles OWNER TO myuser;

--
-- Name: articles_article_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.articles_article_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.articles_article_id_seq OWNER TO myuser;

--
-- Name: articles_article_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.articles_article_id_seq OWNED BY public.articles.article_id;


--
-- Name: comments; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.comments (
    comment_id bigint NOT NULL,
    article_id bigint NOT NULL,
    user_id bigint NOT NULL,
    comment_text text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    is_published boolean NOT NULL,
    reply_to_comment_id bigint
);


ALTER TABLE public.comments OWNER TO myuser;

--
-- Name: comments_comment_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.comments_comment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.comments_comment_id_seq OWNER TO myuser;

--
-- Name: comments_comment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.comments_comment_id_seq OWNED BY public.comments.comment_id;


--
-- Name: facilities; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.facilities (
    facility_id bigint NOT NULL,
    facility_name character varying(255) NOT NULL,
    location character varying(255) NOT NULL,
    description_text text,
    latitude numeric(10,6) NOT NULL,
    longitude numeric(10,6) NOT NULL,
    person_id bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.facilities OWNER TO myuser;

--
-- Name: facilities_facility_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.facilities_facility_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.facilities_facility_id_seq OWNER TO myuser;

--
-- Name: facilities_facility_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.facilities_facility_id_seq OWNED BY public.facilities.facility_id;


--
-- Name: files; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.files (
    file_id bigint NOT NULL,
    file_name character varying(255) NOT NULL,
    file_type character varying(50) NOT NULL,
    file_size bigint,
    file_data bytea,
    location character varying(255) NOT NULL,
    related_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    s3_key character varying(500),
    s3_url character varying(1000)
);


ALTER TABLE public.files OWNER TO myuser;

--
-- Name: files_file_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.files_file_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.files_file_id_seq OWNER TO myuser;

--
-- Name: files_file_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.files_file_id_seq OWNED BY public.files.file_id;


--
-- Name: languages; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.languages (
    language_id bigint NOT NULL,
    language_name character varying(50) NOT NULL,
    display_order bigint,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


ALTER TABLE public.languages OWNER TO myuser;

--
-- Name: languages_language_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.languages_language_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.languages_language_id_seq OWNER TO myuser;

--
-- Name: languages_language_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.languages_language_id_seq OWNED BY public.languages.language_id;


--
-- Name: menus; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.menus (
    menu_id bigint NOT NULL,
    menu_name character varying(100) NOT NULL,
    menu_code character varying(50) NOT NULL,
    display_order bigint,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


ALTER TABLE public.menus OWNER TO myuser;

--
-- Name: menus_menu_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.menus_menu_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.menus_menu_id_seq OWNER TO myuser;

--
-- Name: menus_menu_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.menus_menu_id_seq OWNED BY public.menus.menu_id;


--
-- Name: notices; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.notices (
    notice_id bigint NOT NULL,
    title character varying(255) NOT NULL,
    content text NOT NULL,
    notice_type boolean NOT NULL,
    user_id bigint,
    published_at timestamp with time zone NOT NULL,
    is_active boolean NOT NULL,
    is_read boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


ALTER TABLE public.notices OWNER TO myuser;

--
-- Name: notices_notice_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.notices_notice_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.notices_notice_id_seq OWNER TO myuser;

--
-- Name: notices_notice_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.notices_notice_id_seq OWNED BY public.notices.notice_id;


--
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.refresh_tokens (
    token_id bigint NOT NULL,
    user_id bigint NOT NULL,
    refresh_token character varying(255) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    revoked boolean DEFAULT false NOT NULL
);


ALTER TABLE public.refresh_tokens OWNER TO myuser;

--
-- Name: refresh_tokens_token_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.refresh_tokens_token_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.refresh_tokens_token_id_seq OWNER TO myuser;

--
-- Name: refresh_tokens_token_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.refresh_tokens_token_id_seq OWNED BY public.refresh_tokens.token_id;


--
-- Name: stores; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.stores (
    store_id bigint NOT NULL,
    store_name character varying(255) NOT NULL,
    store_category character varying(100) NOT NULL,
    location character varying(255) NOT NULL,
    description_text text,
    address character varying(255) NOT NULL,
    latitude numeric(10,6) NOT NULL,
    longitude numeric(10,6) NOT NULL,
    business_hours character varying(100) NOT NULL,
    rating_score numeric(3,2) NOT NULL,
    phone_number character varying(20) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.stores OWNER TO myuser;

--
-- Name: stores_store_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.stores_store_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.stores_store_id_seq OWNER TO myuser;

--
-- Name: stores_store_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.stores_store_id_seq OWNED BY public.stores.store_id;


--
-- Name: taggings; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.taggings (
    tagging_id bigint NOT NULL,
    tag_id bigint NOT NULL,
    taggable_type character varying(50) NOT NULL,
    taggable_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


ALTER TABLE public.taggings OWNER TO myuser;

--
-- Name: taggings_tagging_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.taggings_tagging_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.taggings_tagging_id_seq OWNER TO myuser;

--
-- Name: taggings_tagging_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.taggings_tagging_id_seq OWNED BY public.taggings.tagging_id;


--
-- Name: tags; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.tags (
    tag_id bigint NOT NULL,
    tag_name character varying(50) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


ALTER TABLE public.tags OWNER TO myuser;

--
-- Name: tags_tag_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.tags_tag_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tags_tag_id_seq OWNER TO myuser;

--
-- Name: tags_tag_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.tags_tag_id_seq OWNED BY public.tags.tag_id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.users (
    user_id bigint NOT NULL,
    name text,
    name_kana text,
    birth timestamp with time zone,
    address text,
    gender text,
    phone_number text,
    email text NOT NULL,
    password text,
    avatar text,
    google_id text,
    apple_id text,
    provider text NOT NULL,
    status text NOT NULL,
    verify_code text,
    verify_code_expire timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO myuser;

--
-- Name: users_user_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.users_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_user_id_seq OWNER TO myuser;

--
-- Name: users_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.users_user_id_seq OWNED BY public.users.user_id;


--
-- Name: visit_histories; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.visit_histories (
    history_id bigint NOT NULL,
    user_id bigint NOT NULL,
    facility_id bigint NOT NULL,
    scan_at timestamp with time zone NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


ALTER TABLE public.visit_histories OWNER TO myuser;

--
-- Name: visit_histories_history_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.visit_histories_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.visit_histories_history_id_seq OWNER TO myuser;

--
-- Name: visit_histories_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.visit_histories_history_id_seq OWNED BY public.visit_histories.history_id;


--
-- Name: articles article_id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.articles ALTER COLUMN article_id SET DEFAULT nextval('public.articles_article_id_seq'::regclass);


--
-- Name: comments comment_id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.comments ALTER COLUMN comment_id SET DEFAULT nextval('public.comments_comment_id_seq'::regclass);


--
-- Name: facilities facility_id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.facilities ALTER COLUMN facility_id SET DEFAULT nextval('public.facilities_facility_id_seq'::regclass);


--
-- Name: files file_id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.files ALTER COLUMN file_id SET DEFAULT nextval('public.files_file_id_seq'::regclass);


--
-- Name: languages language_id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.languages ALTER COLUMN language_id SET DEFAULT nextval('public.languages_language_id_seq'::regclass);


--
-- Name: menus menu_id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.menus ALTER COLUMN menu_id SET DEFAULT nextval('public.menus_menu_id_seq'::regclass);


--
-- Name: notices notice_id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.notices ALTER COLUMN notice_id SET DEFAULT nextval('public.notices_notice_id_seq'::regclass);


--
-- Name: refresh_tokens token_id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN token_id SET DEFAULT nextval('public.refresh_tokens_token_id_seq'::regclass);


--
-- Name: stores store_id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.stores ALTER COLUMN store_id SET DEFAULT nextval('public.stores_store_id_seq'::regclass);


--
-- Name: taggings tagging_id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.taggings ALTER COLUMN tagging_id SET DEFAULT nextval('public.taggings_tagging_id_seq'::regclass);


--
-- Name: tags tag_id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.tags ALTER COLUMN tag_id SET DEFAULT nextval('public.tags_tag_id_seq'::regclass);


--
-- Name: users user_id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.users ALTER COLUMN user_id SET DEFAULT nextval('public.users_user_id_seq'::regclass);


--
-- Name: visit_histories history_id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.visit_histories ALTER COLUMN history_id SET DEFAULT nextval('public.visit_histories_history_id_seq'::regclass);


--
-- Data for Name: articles; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.articles (article_id, title, body_text, category, like_count, article_image, comment_count, created_at, updated_at, image_file_id) FROM stdin;
\.


--
-- Data for Name: comments; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.comments (comment_id, article_id, user_id, comment_text, created_at, updated_at, is_published, reply_to_comment_id) FROM stdin;
\.


--
-- Data for Name: facilities; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.facilities (facility_id, facility_name, location, description_text, latitude, longitude, person_id, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: files; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.files (file_id, file_name, file_type, file_size, file_data, location, related_id, created_at, updated_at, s3_key, s3_url) FROM stdin;
\.


--
-- Data for Name: languages; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.languages (language_id, language_name, display_order, is_active, created_at, updated_at) FROM stdin;
5	日本語	1	t	2025-06-01 16:10:18.13462+00	2025-06-01 16:10:18.13462+00
6	English	2	t	2025-06-01 16:10:18.13462+00	2025-06-01 16:10:18.13462+00
7	中文	3	t	2025-06-01 16:10:18.13462+00	2025-06-01 16:10:18.13462+00
8	한국어	4	t	2025-06-01 16:10:18.13462+00	2025-06-01 16:10:18.13462+00
\.


--
-- Data for Name: menus; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.menus (menu_id, menu_name, menu_code, display_order, is_active, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: notices; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.notices (notice_id, title, content, notice_type, user_id, published_at, is_active, is_read, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: refresh_tokens; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.refresh_tokens (token_id, user_id, refresh_token, expires_at, created_at, revoked) FROM stdin;
2	20	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyMCwiZXhwIjoxNzQ5NDAwODc5fQ.PbwLXIZCwIKsidOxUfCe4GKmKddyiad51xUQFMjAbaI	2025-06-08 16:41:19.548182+00	2025-06-01 16:41:19.546269+00	f
\.


--
-- Data for Name: stores; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.stores (store_id, store_name, store_category, location, description_text, address, latitude, longitude, business_hours, rating_score, phone_number, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: taggings; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.taggings (tagging_id, tag_id, taggable_type, taggable_id, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: tags; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.tags (tag_id, tag_name, is_active, created_at, updated_at) FROM stdin;
9	観光地	t	2025-06-01 16:10:18.135777+00	2025-06-01 16:10:18.135777+00
10	レストラン	t	2025-06-01 16:10:18.135777+00	2025-06-01 16:10:18.135777+00
11	ホテル	t	2025-06-01 16:10:18.135777+00	2025-06-01 16:10:18.135777+00
12	ショッピング	t	2025-06-01 16:10:18.135777+00	2025-06-01 16:10:18.135777+00
13	文化	t	2025-06-01 16:10:18.135777+00	2025-06-01 16:10:18.135777+00
14	歴史	t	2025-06-01 16:10:18.135777+00	2025-06-01 16:10:18.135777+00
15	自然	t	2025-06-01 16:10:18.135777+00	2025-06-01 16:10:18.135777+00
16	体験	t	2025-06-01 16:10:18.135777+00	2025-06-01 16:10:18.135777+00
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.users (user_id, name, name_kana, birth, address, gender, phone_number, email, password, avatar, google_id, apple_id, provider, status, verify_code, verify_code_expire, created_at, updated_at) FROM stdin;
11	系统管理员	システムカンリシャ	\N	東京都千代田区千代田1-1-1	other	000-0000-0000	admin@ar-backend.com	$2a$10$Td7UYQFuXXj9LKlI2yGc4.3jVY4J9Z5UeGwKbQB6kPwXf8Jj.nLG6	\N	\N	\N	email	active	\N	\N	2025-06-01 16:10:18.136308+00	2025-06-01 16:10:18.136308+00
12	张三	チョウサン	1990-05-15 00:00:00+00	东京都涩谷区1-2-3	男	080-1234-5678	zhangsan@example.com	$2a$10$dVtSQRGR8pdqFU3UK454dekt5X5Q8PLtKJPTBxlhajfTGeDOVgV6q	https://via.placeholder.com/150/0000FF/FFFFFF?text=张三			email	active		\N	2025-06-01 16:29:01.666909+00	2025-06-01 16:29:02.121753+00
13	李四	リシ	1985-08-20 00:00:00+00	大阪府大阪市中央区4-5-6	女	090-2345-6789	lisi@example.com	$2a$10$BfulQ395yVZWgh2oXWx.4uphbFyWF9HQ8UsEcw02wxZnmZv5HmJpS	https://via.placeholder.com/150/FF0000/FFFFFF?text=李四			email	active		\N	2025-06-01 16:29:01.731619+00	2025-06-01 16:29:02.128598+00
14	王五	オウゴ	1992-12-10 00:00:00+00	京都府京都市左京区7-8-9	男	070-3456-7890	wangwu@gmail.com	$2a$10$v2ajuwf6H75IhvTcN//Zcu3/fThIyhy4olgqa3v12e0S4wfsWC7vW	https://via.placeholder.com/150/00FF00/FFFFFF?text=王五	google_123456789		google	active		\N	2025-06-01 16:29:01.796417+00	2025-06-01 16:29:02.129864+00
15	赵六	チョウロク	1988-03-25 00:00:00+00	神奈川県横浜市港北区10-11-12	女	080-4567-8901	zhaoliu@icloud.com	$2a$10$erdaG7Z4/1F84QYWNKAWzOLxUN3.XWIIgw4lzbVQhNCZZjF1w.pGC	https://via.placeholder.com/150/FFFF00/000000?text=赵六		apple_987654321	apple	active		\N	2025-06-01 16:29:01.860593+00	2025-06-01 16:29:02.131363+00
16	孙七	ソンナナ	1995-07-08 00:00:00+00	福岡県福岡市博多区13-14-15	男	090-5678-9012	sunqi@example.com	$2a$10$XkiCGMzJtZmGz1Zauh0exuZVW0edoJiEue2pbhBzGb3Hu3NPYEeKm	https://via.placeholder.com/150/FF00FF/FFFFFF?text=孙七			email	pending	1234	2025-06-01 16:39:01.924996+00	2025-06-01 16:29:01.924997+00	2025-06-01 16:29:02.132205+00
17	周八	シュウハチ	1993-11-30 00:00:00+00	北海道札幌市中央区16-17-18	女	070-6789-0123	zhouba@example.com	$2a$10$aZ1wA/FeS7y73QzeBCm/pOWvLVak84ULUsNI7xW3UBUST31qFvPAC	https://via.placeholder.com/150/00FFFF/000000?text=周八			email	inactive		\N	2025-06-01 16:29:01.989401+00	2025-06-01 16:29:02.133003+00
18	吴九	ゴキュウ	1987-09-12 00:00:00+00	愛知県名古屋市中区19-20-21	男	080-7890-1234	wujiu@gmail.com	$2a$10$tT1GqxGyHVwuIKN.v8SIF.TO9e1bfr34clxEuTcmyuLMGmy5M/xq6	https://via.placeholder.com/150/800080/FFFFFF?text=吴九	google_246810121		google	active		\N	2025-06-01 16:29:02.053935+00	2025-06-01 16:29:02.133806+00
19	郑十	テイジュウ	1991-04-18 00:00:00+00	広島県広島市中区22-23-24	女	090-8901-2345	zhengshi@example.com	$2a$10$Cx3Crv92xS0mcgBws86Nyun1MpneAfrB/j7LVDJveRchN0R1uQkFK	https://via.placeholder.com/150/FFA500/FFFFFF?text=郑十			email	active		\N	2025-06-01 16:29:02.11862+00	2025-06-01 16:29:02.134547+00
20			\N		\N		sadt@gmail.com	$2a$10$hwNbBkJdZYK8Gq1t.cb6POSeowFCy/rDsoZ/NF8iRk8XGrSwzOfj6				email	pending	1612	2025-06-01 16:51:19.530011+00	2025-06-01 16:41:19.53273+00	2025-06-01 16:41:19.53273+00
\.


--
-- Data for Name: visit_histories; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.visit_histories (history_id, user_id, facility_id, scan_at, is_active, created_at, updated_at) FROM stdin;
\.


--
-- Name: articles_article_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.articles_article_id_seq', 1, false);


--
-- Name: comments_comment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.comments_comment_id_seq', 1, false);


--
-- Name: facilities_facility_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.facilities_facility_id_seq', 1, false);


--
-- Name: files_file_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.files_file_id_seq', 1, false);


--
-- Name: languages_language_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.languages_language_id_seq', 8, true);


--
-- Name: menus_menu_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.menus_menu_id_seq', 1, false);


--
-- Name: notices_notice_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.notices_notice_id_seq', 1, false);


--
-- Name: refresh_tokens_token_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.refresh_tokens_token_id_seq', 2, true);


--
-- Name: stores_store_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.stores_store_id_seq', 1, false);


--
-- Name: taggings_tagging_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.taggings_tagging_id_seq', 1, false);


--
-- Name: tags_tag_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.tags_tag_id_seq', 16, true);


--
-- Name: users_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.users_user_id_seq', 20, true);


--
-- Name: visit_histories_history_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.visit_histories_history_id_seq', 1, false);


--
-- Name: articles articles_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_pkey PRIMARY KEY (article_id);


--
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (comment_id);


--
-- Name: facilities facilities_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.facilities
    ADD CONSTRAINT facilities_pkey PRIMARY KEY (facility_id);


--
-- Name: files files_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.files
    ADD CONSTRAINT files_pkey PRIMARY KEY (file_id);


--
-- Name: languages languages_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.languages
    ADD CONSTRAINT languages_pkey PRIMARY KEY (language_id);


--
-- Name: menus menus_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.menus
    ADD CONSTRAINT menus_pkey PRIMARY KEY (menu_id);


--
-- Name: notices notices_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.notices
    ADD CONSTRAINT notices_pkey PRIMARY KEY (notice_id);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (token_id);


--
-- Name: stores stores_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.stores
    ADD CONSTRAINT stores_pkey PRIMARY KEY (store_id);


--
-- Name: taggings taggings_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.taggings
    ADD CONSTRAINT taggings_pkey PRIMARY KEY (tagging_id);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (tag_id);


--
-- Name: users uni_users_email; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_email UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: visit_histories visit_histories_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.visit_histories
    ADD CONSTRAINT visit_histories_pkey PRIMARY KEY (history_id);


--
-- PostgreSQL database dump complete
--

