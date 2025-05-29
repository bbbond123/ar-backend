import React, { useEffect } from "react";

const GoogleLoginButton: React.FC = () => {
  useEffect(() => {
    // 初始化 Google 登录
    // console.log("🚀 ~ useEffect ~ window.google:", window.google)
    if (!window.google) {
      alert("Google is not loaded");
      return;
    }

    if (!import.meta.env.VITE_GOOGLE_CLIENT_ID) {
      alert("Google Client ID is not set");
      return;
    }

    window.google?.accounts.id.initialize({
      client_id: import.meta.env.VITE_GOOGLE_CLIENT_ID,
      callback: handleCredentialResponse,
    });
    window.google?.accounts.id.renderButton(
      document.getElementById("google-login-btn"),
      { theme: "outline", size: "large" }
    );
  }, []);

  const handleCredentialResponse = async (response: any) => {
    const id_token = response.credential;
    // 调用后端接口
    const res = await fetch(import.meta.env.VITE_API_URL + "api/auth/google", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id_token }),
    });
    const data = await res.json();
    if (data.success) {
      // 保存token，跳转主页面
      localStorage.setItem("access_token", data.data.access_token);
      localStorage.setItem("refresh_token", data.data.refresh_token);
      // ...跳转或刷新页面
    } else {
      alert(data.errMessage || "登录失败");
    }
  };

  return <div id="google-login-btn"></div>;
};

export default GoogleLoginButton;
