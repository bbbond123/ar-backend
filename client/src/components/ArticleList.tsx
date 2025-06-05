import React, { useState, useEffect } from 'react';
import { getArticles } from '../api';
import './ArticleList.css';

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

interface ArticleListResponse {
  success: boolean;
  total: number;
  list: Article[];
  error_message?: string;
}

interface ArticleListProps {
  onArticleClick?: (articleId: number) => void;
}

const ArticleList: React.FC<ArticleListProps> = ({ onArticleClick }) => {
  const [articles, setArticles] = useState<Article[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [currentPage, setCurrentPage] = useState(1);
  const [totalArticles, setTotalArticles] = useState(0);
  const [keyword, setKeyword] = useState('');
  const [searchInput, setSearchInput] = useState('');
  
  const pageSize = 6; // 每页显示6篇文章
  const totalPages = Math.ceil(totalArticles / pageSize);

  const fetchArticleList = async (page: number, searchKeyword: string = '') => {
    try {
      setLoading(true);
      setError('');
      
      const response: ArticleListResponse = await getArticles(page, pageSize, searchKeyword);
      
      if (response.success) {
        setArticles(response.list);
        setTotalArticles(response.total);
      } else {
        throw new Error(response.error_message || '获取文章列表失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取文章列表失败');
      setArticles([]);
      setTotalArticles(0);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchArticleList(currentPage, keyword);
  }, [currentPage, keyword]);

  const handleSearch = () => {
    setKeyword(searchInput);
    setCurrentPage(1);
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    window.scrollTo({ top: 0, behavior: 'smooth' });
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

  const truncateText = (text: string, maxLength: number) => {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
  };

  const renderPagination = () => {
    if (totalPages <= 1) return null;

    const pages = [];
    const showPages = 5; // 显示的页码数量
    let startPage = Math.max(1, currentPage - Math.floor(showPages / 2));
    let endPage = Math.min(totalPages, startPage + showPages - 1);

    if (endPage - startPage + 1 < showPages) {
      startPage = Math.max(1, endPage - showPages + 1);
    }

    // 上一页
    if (currentPage > 1) {
      pages.push(
        <button
          key="prev"
          onClick={() => handlePageChange(currentPage - 1)}
          className="pagination-btn"
        >
          ‹
        </button>
      );
    }

    // 第一页
    if (startPage > 1) {
      pages.push(
        <button
          key={1}
          onClick={() => handlePageChange(1)}
          className="pagination-btn"
        >
          1
        </button>
      );
      if (startPage > 2) {
        pages.push(<span key="dots1" className="pagination-dots">...</span>);
      }
    }

    // 页码
    for (let i = startPage; i <= endPage; i++) {
      pages.push(
        <button
          key={i}
          onClick={() => handlePageChange(i)}
          className={`pagination-btn ${i === currentPage ? 'active' : ''}`}
        >
          {i}
        </button>
      );
    }

    // 最后一页
    if (endPage < totalPages) {
      if (endPage < totalPages - 1) {
        pages.push(<span key="dots2" className="pagination-dots">...</span>);
      }
      pages.push(
        <button
          key={totalPages}
          onClick={() => handlePageChange(totalPages)}
          className="pagination-btn"
        >
          {totalPages}
        </button>
      );
    }

    // 下一页
    if (currentPage < totalPages) {
      pages.push(
        <button
          key="next"
          onClick={() => handlePageChange(currentPage + 1)}
          className="pagination-btn"
        >
          ›
        </button>
      );
    }

    return <div className="pagination">{pages}</div>;
  };

  return (
    <div className="article-list">
      <div className="article-list-container">
        <div className="list-header">
          <h2 className="list-title">文章列表</h2>
          <div className="search-box">
            <input
              type="text"
              placeholder="搜索文章标题或内容..."
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
              className="search-input"
            />
            <button onClick={handleSearch} className="search-btn">
              🔍
            </button>
          </div>
        </div>

        {loading && (
          <div className="loading-container">
            <div className="loading-spinner"></div>
            <p>加载中...</p>
          </div>
        )}

        {error && (
          <div className="error-container">
            <p className="error-message">❌ {error}</p>
            <button onClick={() => fetchArticleList(currentPage, keyword)} className="retry-btn">
              重试
            </button>
          </div>
        )}

        {!loading && !error && articles.length === 0 && (
          <div className="empty-container">
            <p className="empty-message">📝 暂无文章</p>
            {keyword && (
              <button 
                onClick={() => {
                  setKeyword('');
                  setSearchInput('');
                  setCurrentPage(1);
                }} 
                className="clear-search-btn"
              >
                清除搜索
              </button>
            )}
          </div>
        )}

        {!loading && !error && articles.length > 0 && (
          <>
            <div className="articles-grid">
              {articles.map((article) => (
                <div 
                  key={article.article_id} 
                  className="article-card"
                  onClick={() => onArticleClick?.(article.article_id)}
                >
                  {article.image_url && (
                    <div className="article-image">
                      <img 
                        src={article.image_url} 
                        alt={article.title}
                        onError={(e) => {
                          const target = e.target as HTMLImageElement;
                          target.style.display = 'none';
                        }}
                      />
                    </div>
                  )}
                  
                  <div className="article-content">
                    <div className="article-meta">
                      {article.category && (
                        <span className="category-tag">{article.category}</span>
                      )}
                      <span className="publish-date">{formatDate(article.created_at)}</span>
                    </div>
                    
                    <h3 className="article-title">{article.title}</h3>
                    
                    <p className="article-preview">
                      {truncateText(article.body_text, 120)}
                    </p>
                    
                    <div className="article-stats">
                      <span className="stat-item">
                        ❤️ {article.like_count}
                      </span>
                      <span className="stat-item">
                        💬 {article.comment_count}
                      </span>
                    </div>
                  </div>
                </div>
              ))}
            </div>

            <div className="list-footer">
              <div className="result-info">
                共 {totalArticles} 篇文章，第 {currentPage} 页 / 共 {totalPages} 页
              </div>
              {renderPagination()}
            </div>
          </>
        )}
      </div>
    </div>
  );
};

export default ArticleList; 