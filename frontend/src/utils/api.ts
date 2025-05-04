const API_BASE_URL = "http://localhost:8080";

export async function fetchLocations() {
  const res = await fetch(`${API_BASE_URL}/locations`);
  if (!res.ok) throw new Error("Failed to fetch locations");
  return res.json();
}

export async function createLocation(data: {
  name: string;
  latitude: number;
  longitude: number;
}) {
  const res = await fetch(`${API_BASE_URL}/locations`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });

  if (!res.ok) {
    const err = await res.json();
    throw new Error(err.error || "Failed to create location");
  }

  return res.json();
}
