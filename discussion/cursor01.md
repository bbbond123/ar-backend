我来帮您分析一下这个项目。首先让我们看看项目的目录结构。

Ran tool
让我们先看一下主要的配置文件和入口文件，以了解项目的基本设置。

Read file: go.mod

Read file: main.go
让我们再看看项目的配置文件：

Read file: config.yaml
让我们看看项目的目录结构中比较重要的几个目录：

Ran tool

Ran tool
基于以上分析，我可以为您总结这个项目的主要特点：

1. **技术栈**：
   - 后端框架：使用 Go 语言开发，基于 Gin Web 框架
   - 数据库：PostgreSQL（通过 GORM ORM 框架操作）
   - API 文档：集成了 Swagger 用于 API 文档生成
   - 认证：使用 JWT（JSON Web Token）进行身份验证

2. **项目结构**：
   - 采用标准的 Go 项目布局
   - `internal/`: 存放内部代码，包含：
     - `model/`: 数据模型定义
     - `controller/`: 控制器层
     - `service/`: 业务逻辑层
     - `repository/`: 数据访问层
     - `middleware/`: 中间件
     - `auth/`: 认证相关
     - `router/`: 路由配置
   - `pkg/`: 可共享的工具包
   - `docs/`: 文档目录
   - `sql_script/`: SQL 脚本文件

3. **主要功能模块**：
   根据数据模型可以看出，系统包含以下主要功能：
   - 用户管理（User）
   - 设施管理（Facility）
   - 文件管理（File）
   - 通知系统（Notice）
   - 访问历史记录（VisitHistory）
   - 多语言支持（Language）
   - 商店管理（Store）
   - 菜单管理（Menu）
   - 文章系统（Article）
   - 评论系统（Comment）
   - 标签系统（Tag, Tagging）

4. **项目特点**：
   - 支持 OAuth 认证（通过 goth 包）
   - 实现了 CORS 跨域支持
   - 使用 Viper 进行配置管理
   - 集成了 Swagger 用于 API 文档
   - 实现了刷新令牌机制（RefreshToken）
   - 使用 Docker 支持容器化部署（有 docker-compose.yml）

5. **开发工具支持**：
   - 包含 `.vscode/` 配置，支持 VS Code 开发
   - 使用 Air 实现热重载（有 .air.toml 配置）
   - 包含完整的依赖管理（go.mod 和 go.sum）

这是一个比较完整的后端项目，看起来是一个旅游相关的 AR（增强现实）应用的后端服务。项目结构清晰，采用了标准的分层架构，并集成了多个现代化的功能和工具。系统支持多语言、文件上传、用户认证等功能，并且考虑到了开发效率（热重载）和部署便利性（Docker支持）。
