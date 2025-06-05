export async function fetchMe() {
  const token = localStorage.getItem("access_token");
  const headers: Record<string, string> = {
    "credentials": "include"
  };
  
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }
  
  const res = await fetch("/api/me", {
    credentials: "include", // 关键！带上 cookie
    headers,
  });
  if (!res.ok) throw new Error("Not logged in");
  return res.json();
}

export async function logout() {
  await fetch("/api/logout", {
    method: "POST",
    credentials: "include",
  });
}

// 文章相关API函数
export async function createArticle(formData: FormData) {
  const token = localStorage.getItem("access_token");
  const headers: Record<string, string> = {
    "credentials": "include"
  };
  
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch("/api/articles/with-image", {
    method: "POST",
    credentials: "include",
    headers,
    body: formData
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return response.json();
}

export async function getArticles(page: number = 1, pageSize: number = 10, keyword: string = '') {
  const response = await fetch("/api/articles/list", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "credentials": "include"
    },
    credentials: "include",
    body: JSON.stringify({
      page,
      page_size: pageSize,
      keyword
    })
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return response.json();
}
