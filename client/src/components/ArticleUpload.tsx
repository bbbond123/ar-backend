import React, { useState, useRef } from 'react';
import { createArticle } from '../api';
import './ArticleUpload.css';

interface ArticleFormData {
  title: string;
  body_text: string;
  category: string;
  image?: File;
}

interface UploadResponse {
  success: boolean;
  data?: any;
  error_message?: string;
}

const ArticleUpload: React.FC = () => {
  const [formData, setFormData] = useState<ArticleFormData>({
    title: '',
    body_text: '',
    category: '',
  });
  
  const [imagePreview, setImagePreview] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const categories = [
    '旅游',
    '美食',
    '文化',
    '景点',
    '攻略',
    '住宿',
    '交通',
    '购物',
    '其他'
  ];

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      // 检查文件类型
      const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
      if (!allowedTypes.includes(file.type)) {
        setMessage({ type: 'error', text: '请选择有效的图片格式 (JPG, PNG, GIF, WebP)' });
        return;
      }

      // 检查文件大小 (5MB)
      const maxSize = 5 * 1024 * 1024;
      if (file.size > maxSize) {
        setMessage({ type: 'error', text: '图片文件大小不能超过 5MB' });
        return;
      }

      setFormData(prev => ({
        ...prev,
        image: file
      }));

      // 创建预览
      const reader = new FileReader();
      reader.onload = (e) => {
        setImagePreview(e.target?.result as string);
      };
      reader.readAsDataURL(file);
      setMessage(null);
    }
  };

  const removeImage = () => {
    setFormData(prev => {
      const newData = { ...prev };
      delete newData.image;
      return newData;
    });
    setImagePreview('');
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.title.trim() || !formData.body_text.trim()) {
      setMessage({ type: 'error', text: '标题和内容为必填项' });
      return;
    }

    setIsSubmitting(true);
    setMessage(null);

    try {
      // 使用 FormData 上传（支持图片）
      const submitFormData = new FormData();
      submitFormData.append('title', formData.title);
      submitFormData.append('body_text', formData.body_text);
      submitFormData.append('category', formData.category);
      submitFormData.append('like_count', '0');
      submitFormData.append('comment_count', '0');
      
      if (formData.image) {
        submitFormData.append('image', formData.image);
      }

      const result: UploadResponse = await createArticle(submitFormData);

      if (result.success) {
        setMessage({ type: 'success', text: '文章发布成功！' });
        // 重置表单
        setFormData({
          title: '',
          body_text: '',
          category: '',
        });
        setImagePreview('');
        if (fileInputRef.current) {
          fileInputRef.current.value = '';
        }
      } else {
        throw new Error(result.error_message || '发布失败');
      }
    } catch (error) {
      console.error('Upload error:', error);
      setMessage({ 
        type: 'error', 
        text: error instanceof Error ? error.message : '发布失败，请稍后重试' 
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="article-upload">
      <div className="upload-container">
        <h2 className="upload-title">发布文章</h2>
        
        {message && (
          <div className={`message ${message.type}`}>
            {message.text}
          </div>
        )}

        <form onSubmit={handleSubmit} className="upload-form">
          <div className="form-group">
            <label htmlFor="title">文章标题 *</label>
            <input
              type="text"
              id="title"
              name="title"
              value={formData.title}
              onChange={handleInputChange}
              placeholder="请输入文章标题"
              required
              maxLength={255}
            />
          </div>

          <div className="form-group">
            <label htmlFor="category">分类</label>
            <select
              id="category"
              name="category"
              value={formData.category}
              onChange={handleInputChange}
            >
              <option value="">请选择分类</option>
              {categories.map(cat => (
                <option key={cat} value={cat}>{cat}</option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="body_text">文章内容 *</label>
            <textarea
              id="body_text"
              name="body_text"
              value={formData.body_text}
              onChange={handleInputChange}
              placeholder="请输入文章内容"
              required
              rows={10}
            />
          </div>

          <div className="form-group">
            <label htmlFor="image">文章图片</label>
            <div className="image-upload-area">
              {imagePreview ? (
                <div className="image-preview">
                  <img src={imagePreview} alt="预览" />
                  <button 
                    type="button" 
                    className="remove-image"
                    onClick={removeImage}
                  >
                    ✕
                  </button>
                </div>
              ) : (
                <div className="upload-placeholder">
                  <span>📷</span>
                  <p>点击选择图片</p>
                  <p className="file-hint">支持 JPG、PNG、GIF、WebP 格式，最大 5MB</p>
                </div>
              )}
              <input
                ref={fileInputRef}
                type="file"
                id="image"
                name="image"
                accept="image/*"
                onChange={handleImageChange}
                hidden
              />
              <button
                type="button"
                className="select-image-btn"
                onClick={() => fileInputRef.current?.click()}
              >
                {imagePreview ? '更换图片' : '选择图片'}
              </button>
            </div>
          </div>

          <div className="form-actions">
            <button 
              type="submit" 
              className="submit-btn"
              disabled={isSubmitting}
            >
              {isSubmitting ? '发布中...' : '发布文章'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ArticleUpload; 