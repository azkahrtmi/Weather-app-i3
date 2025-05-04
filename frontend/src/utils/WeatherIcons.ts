export type WeatherCode = number;

export const weatherDescriptions: Record<WeatherCode, string> = {
  1: "Not available",
  2: "Sunny",
  3: "Mostly sunny",
  4: "Partly sunny",
  5: "Mostly cloudy",
  6: "Cloudy",
  7: "Overcast",
  8: "Overcast with low clouds",
  9: "Fog",
  10: "Light rain",
  11: "Rain",
  12: "Possible rain",
  13: "Rain shower",
  14: "Thunderstorm",
  15: "Local thunderstorms",
  16: "Light snow",
  17: "Snow",
  18: "Possible snow",
  19: "Snow shower",
  20: "Rain and snow",
  21: "Possible rain and snow",
  22: "Rain and snow",
  23: "Freezing rain",
  24: "Possible freezing rain",
  25: "Hail",
  26: "Clear",
  27: "Mostly clear",
  28: "Partly clear",
  29: "Mostly cloudy",
  30: "Cloudy",
  31: "Overcast with low clouds",
  32: "Rain shower",
  33: "Local thunderstorms",
  34: "Snow shower",
  35: "Rain and snow",
  36: "Possible freezing rain",
};

export function getWeatherCodeFromSummary(summary: string): WeatherCode {
  const entry = Object.entries(weatherDescriptions).find(
    ([, desc]) => desc.toLowerCase() === summary.toLowerCase()
  );
  return entry ? parseInt(entry[0], 10) : 1; // Default ke 1: "Not available"
}

export function getWeatherIconUrl(code: WeatherCode): string {
  return `/assets/weather-icons/${code}.png`;
}

export function getWeatherDescription(code: WeatherCode): string {
  return weatherDescriptions[code] || "Unknown";
}
