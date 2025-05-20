export async function fetchMe() {
  const res = await fetch("/api/me", {
    credentials: "include", // 关键！带上 cookie
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
