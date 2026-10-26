const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api";

/**
 * Shorten a long URL via Go Backend API
 * @param {string} url - The original long URL
 * @returns {Promise<{short_code: string, short_url: string, original_url: string, created_at: string}>}
 */
export async function shortenUrl(url) {
  const response = await fetch(`${API_BASE_URL}/shorten`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ url: url.trim() }),
  });

  const data = await response.json().catch(() => ({}));

  if (!response.ok) {
    throw new Error(data.error || "Failed to shorten URL");
  }

  return data;
}

/**
 * Fetch recently shortened URLs from Go Backend API
 * @returns {Promise<Array<{short_code: string, short_url: string, original_url: string, created_at: string}>>}
 */
export async function getRecentUrls() {
  const response = await fetch(`${API_BASE_URL}/urls`);
  const data = await response.json().catch(() => ([]));

  if (!response.ok) {
    throw new Error(data.error || "Failed to fetch recent URLs");
  }

  return Array.isArray(data) ? data : [];
}
