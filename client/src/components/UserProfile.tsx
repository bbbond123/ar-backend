import React, { useEffect, useState } from 'react';
import { getUserInfo } from '../api';
import './UserProfile.css';

interface UserData {
  user_id: number;
  name: string;
  name_kana: string;
  birth: string | null;
  address: string;
  gender: string | null;
  phone_number: string;
  email: string;
  avatar: string;
  provider: string;
  status: string;
  created_at: string;
  updated_at: string;
}

const UserProfile: React.FC = () => {
  const [userData, setUserData] = useState<UserData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await getUserInfo();
        if (response.success && response.data) {
          setUserData(response.data);
        } else {
          setError('获取用户信息失败');
        }
      } catch (err) {
        setError('获取用户信息时发生错误');
      } finally {
        setLoading(false);
      }
    };

    fetchUserData();
  }, []);

  if (loading) {
    return <div className="loading">加载中...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  if (!userData) {
    return <div className="error">未找到用户信息</div>;
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('ja-JP', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className="user-profile">
      <div className="profile-header">
        <div className="avatar-container">
          <img src={userData.avatar || '/default-avatar.png'} alt="用户头像" className="avatar" />
        </div>
        <div className="user-status">
          <span className={`status-badge ${userData.status}`}>
            {userData.status === 'active' ? '已激活' : '待验证'}
          </span>
        </div>
      </div>

      <div className="profile-content">
        <div className="profile-section">
          <h2>基本信息</h2>
          <div className="info-grid">
            <div className="info-item">
              <label>用户名</label>
              <span>{userData.name || '未设置'}</span>
            </div>
            <div className="info-item">
              <label>假名</label>
              <span>{userData.name_kana || '未设置'}</span>
            </div>
            <div className="info-item">
              <label>邮箱</label>
              <span>{userData.email}</span>
            </div>
            <div className="info-item">
              <label>电话</label>
              <span>{userData.phone_number || '未设置'}</span>
            </div>
            <div className="info-item">
              <label>地址</label>
              <span>{userData.address || '未设置'}</span>
            </div>
            <div className="info-item">
              <label>性别</label>
              <span>{userData.gender || '未设置'}</span>
            </div>
            <div className="info-item">
              <label>生日</label>
              <span>{userData.birth ? formatDate(userData.birth) : '未设置'}</span>
            </div>
          </div>
        </div>

        <div className="profile-section">
          <h2>账户信息</h2>
          <div className="info-grid">
            <div className="info-item">
              <label>登录方式</label>
              <span>{userData.provider === 'email' ? '邮箱登录' : 'Google登录'}</span>
            </div>
            <div className="info-item">
              <label>注册时间</label>
              <span>{formatDate(userData.created_at)}</span>
            </div>
            <div className="info-item">
              <label>最后更新</label>
              <span>{formatDate(userData.updated_at)}</span>
            </div>
          </div>
        </div>
      </div>

      <div className="profile-actions">
        <button className="edit-button">编辑资料</button>
        {userData.status === 'pending' && (
          <button className="verify-button">验证邮箱</button>
        )}
      </div>
    </div>
  );
};

export default UserProfile; 