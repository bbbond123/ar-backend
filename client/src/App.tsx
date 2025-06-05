import { useEffect, useState } from "react";
import { fetchMe, logout } from "./api";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import GoogleLoginButton from "./components/GoogleLoginButton";
import ArticleUpload from "./components/ArticleUpload";
import ArticleList from "./components/ArticleList";
import ArticleDetail from "./components/ArticleDetail";
import "./App.css";

type User = {
  user_id: number;
  email: string;
  name: string;
  avatar?: string;
  provider: string;
};

function App() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [currentView, setCurrentView] = useState<'profile' | 'upload' | 'articles' | 'article-detail'>('profile');
  const [selectedArticleId, setSelectedArticleId] = useState<number | null>(null);

  useEffect(() => {
    // 检查URL参数中是否有token
    const urlParams = new URLSearchParams(window.location.search);
    const token = urlParams.get("token");

    if (token) {
      // 存储token到localStorage
      localStorage.setItem("access_token", token);

      // 清除URL中的token参数
      window.history.replaceState({}, document.title, window.location.pathname);
    }

    // 获取用户信息
    fetchMe()
      .then(setUser)
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, []);

  // const handleLogin = () => {
  //   // window.location.href = "http://localhost:3000/api/auth/google";
  //   window.location.href = process.env.VITE_API_URL + "api/auth/google";
  // };

  if (loading) return <div>Loading...</div>;

  if (!user) {
    // 未登录时显示登录按钮
    return (
      <>
        <div>
          <a href="https://vite.dev" target="_blank">
            <img src={viteLogo} className="logo" alt="Vite logo" />
          </a>
          <a href="https://react.dev" target="_blank">
            <img src={reactLogo} className="logo react" alt="React logo" />
          </a>
        </div>
        {/* <div>
          <button onClick={handleLogin}>使用 Google 登录</button>
        </div> */}
        <GoogleLoginButton />
      </>
    );
  }

  // 已登录时显示导航和对应页面
  return (
    <div>
      {/* 导航栏 */}
      <nav style={{ 
        padding: '1rem', 
        background: '#f8f9fa', 
        borderBottom: '1px solid #dee2e6',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
      }}>
        <div style={{ display: 'flex', gap: '1rem' }}>
          <button
            onClick={() => setCurrentView('profile')}
            style={{
              padding: '0.5rem 1rem',
              border: 'none',
              background: currentView === 'profile' ? '#007bff' : 'transparent',
              color: currentView === 'profile' ? 'white' : '#007bff',
              borderRadius: '4px',
              cursor: 'pointer'
            }}
          >
            个人资料
          </button>
          <button
            onClick={() => setCurrentView('articles')}
            style={{
              padding: '0.5rem 1rem',
              border: 'none',
              background: currentView === 'articles' ? '#007bff' : 'transparent',
              color: currentView === 'articles' ? 'white' : '#007bff',
              borderRadius: '4px',
              cursor: 'pointer'
            }}
          >
            文章列表
          </button>
          <button
            onClick={() => setCurrentView('upload')}
            style={{
              padding: '0.5rem 1rem',
              border: 'none',
              background: currentView === 'upload' ? '#007bff' : 'transparent',
              color: currentView === 'upload' ? 'white' : '#007bff',
              borderRadius: '4px',
              cursor: 'pointer'
            }}
          >
            发布文章
          </button>
        </div>
        
        <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
          <span>欢迎, {user.name || user.email}</span>
          {user.avatar && (
            <img
              src={user.avatar || "https://www.gravatar.com/avatar/?d=mp"}
              alt="avatar"
              width={32}
              height={32}
              style={{ borderRadius: '50%' }}
              referrerPolicy="no-referrer"
            />
          )}
          <button
            onClick={async () => {
              await logout();
              localStorage.removeItem("access_token");
              setUser(null);
            }}
            style={{
              padding: '0.5rem 1rem',
              border: '1px solid #dc3545',
              background: 'white',
              color: '#dc3545',
              borderRadius: '4px',
              cursor: 'pointer'
            }}
          >
            退出登录
          </button>
        </div>
      </nav>

      {/* 页面内容 */}
      {currentView === 'profile' && (
        <div style={{ padding: '2rem' }}>
          <h2>个人资料</h2>
          <div style={{ background: 'white', padding: '1.5rem', borderRadius: '8px', boxShadow: '0 2px 4px rgba(0,0,0,0.1)' }}>
            {user.avatar && (
              <img
                src={user.avatar || "https://www.gravatar.com/avatar/?d=mp"}
                alt="avatar"
                width={64}
                height={64}
                style={{ borderRadius: '50%', marginBottom: '1rem' }}
                referrerPolicy="no-referrer"
              />
            )}
            <p><strong>姓名:</strong> {user.name || '未设置'}</p>
            <p><strong>邮箱:</strong> {user.email}</p>
            <p><strong>登录方式:</strong> {user.provider}</p>
            <p><strong>用户ID:</strong> {user.user_id}</p>
          </div>
        </div>
      )}

      {currentView === 'articles' && (
        <ArticleList 
          onArticleClick={(articleId) => {
            setSelectedArticleId(articleId);
            setCurrentView('article-detail');
          }}
        />
      )}

      {currentView === 'article-detail' && selectedArticleId && (
        <ArticleDetail 
          articleId={selectedArticleId}
          onBack={() => {
            setCurrentView('articles');
            setSelectedArticleId(null);
          }}
        />
      )}

      {currentView === 'upload' && <ArticleUpload />}
    </div>
  );
}

export default App;
