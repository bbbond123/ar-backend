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
