import React, { useState, useEffect } from 'react';
import { getArticle } from '../api';
import './ArticleDetail.css';

interface Article {
  article_id: number;
  title: string;
  body_text: string;
  category: string;
  like_count: number;
  comment_count: number;
  image_url?: string;
  created_at: string;
  updated_at?: string;
}

interface ArticleDetailResponse {
  success: boolean;
  data: Article;
  error_message?: string;
}

interface ArticleDetailProps {
  articleId: number;
  onBack: () => void;
}

const ArticleDetail: React.FC<ArticleDetailProps> = ({ articleId, onBack }) => {
  const [article, setArticle] = useState<Article | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [liked, setLiked] = useState(false);
  const [localLikeCount, setLocalLikeCount] = useState(0);

  useEffect(() => {
    fetchArticleDetail();
  }, [articleId]);

  const fetchArticleDetail = async () => {
    try {
      setLoading(true);
      setError('');
      
      const response: ArticleDetailResponse = await getArticle(articleId);
      
      if (response.success) {
        setArticle(response.data);
        setLocalLikeCount(response.data.like_count);
      } else {
        throw new Error(response.error_message || '获取文章详情失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取文章详情失败');
    } finally {
      setLoading(false);
    }
  };

  const handleLike = () => {
    if (!liked) {
      setLiked(true);
      setLocalLikeCount(prev => prev + 1);
      // TODO: 调用后端点赞API
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const formatContent = (content: string) => {
    // 简单的内容格式化，将换行符转换为<br>
    return content.split('\n').map((line, index) => (
      <React.Fragment key={index}>
        {line}
        {index < content.split('\n').length - 1 && <br />}
      </React.Fragment>
    ));
  };

  if (loading) {
    return (
      <div className="article-detail">
        <div className="article-detail-container">
          <div className="loading-container">
            <div className="loading-spinner"></div>
            <p>加载中...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="article-detail">
        <div className="article-detail-container">
          <div className="error-container">
            <p className="error-message">❌ {error}</p>
            <div className="error-actions">
              <button onClick={fetchArticleDetail} className="retry-btn">
                重试
              </button>
              <button onClick={onBack} className="back-btn">
                返回列表
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!article) {
    return (
      <div className="article-detail">
        <div className="article-detail-container">
          <div className="error-container">
            <p className="error-message">📝 文章不存在</p>
            <button onClick={onBack} className="back-btn">
              返回列表
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="article-detail">
      <div className="article-detail-container">
        {/* 顶部导航 */}
        <div className="article-header">
          <button onClick={onBack} className="back-button">
            ← 返回列表
          </button>
          <div className="article-actions">
            <button 
              onClick={handleLike}
              className={`like-button ${liked ? 'liked' : ''}`}
              disabled={liked}
            >
              ❤️ {localLikeCount}
            </button>
            <button className="share-button">
              🔗 分享
            </button>
          </div>
        </div>

        {/* 文章内容 */}
        <article className="article-content">
          {/* 文章元信息 */}
          <div className="article-meta">
            {article.category && (
              <span className="category-tag">{article.category}</span>
            )}
            <span className="publish-info">
              发布于 {formatDate(article.created_at)}
              {article.updated_at && (
                <span className="update-info">
                  · 更新于 {formatDate(article.updated_at)}
                </span>
              )}
            </span>
          </div>

          {/* 文章标题 */}
          <h1 className="article-title">{article.title}</h1>

          {/* 文章配图 */}
          {article.image_url && (
            <div className="article-image-container">
              <img 
                src={article.image_url} 
                alt={article.title}
                className="article-image"
                onError={(e) => {
                  const target = e.target as HTMLImageElement;
                  target.style.display = 'none';
                }}
              />
            </div>
          )}

          {/* 文章正文 */}
          <div className="article-body">
            <p className="article-text">
              {formatContent(article.body_text)}
            </p>
          </div>

          {/* 文章统计 */}
          <div className="article-stats">
            <div className="stats-item">
              <span className="stats-icon">❤️</span>
              <span className="stats-text">{localLikeCount} 点赞</span>
            </div>
            <div className="stats-item">
              <span className="stats-icon">💬</span>
              <span className="stats-text">{article.comment_count} 评论</span>
            </div>
            <div className="stats-item">
              <span className="stats-icon">👁️</span>
              <span className="stats-text">阅读</span>
            </div>
          </div>
        </article>

        {/* 评论区域（预留） */}
        <div className="comments-section">
          <h3 className="comments-title">评论 ({article.comment_count})</h3>
          <div className="comments-placeholder">
            <p>💬 评论功能开发中...</p>
          </div>
        </div>

        {/* 相关文章（预留） */}
        <div className="related-articles">
          <h3 className="related-title">相关文章</h3>
          <div className="related-placeholder">
            <p>📚 相关文章推荐功能开发中...</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ArticleDetail; 