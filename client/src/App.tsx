import React, { useEffect, useState } from "react";
import { fetchMe, logout } from "./api";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import GoogleLoginButton from "./components/GoogleLoginButton";
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

  useEffect(() => {
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

  // 已登录时显示用户信息
  return (
    <div>
      <h2>欢迎, {user.name || user.email}</h2>
      {user.avatar && (
        <img
          src={user.avatar || "https://www.gravatar.com/avatar/?d=mp"}
          alt="avatar"
          width={48}
          referrerPolicy="no-referrer"
        />
      )}
      <div>邮箱: {user.email}</div>
      <div>登录方式: {user.provider}</div>
      <button
        onClick={async () => {
          await logout();
          setUser(null); // 清空本地用户状态
        }}
      >
        退出登录
      </button>
    </div>
  );
}

export default App;
