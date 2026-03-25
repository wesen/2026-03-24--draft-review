export function getBackendOrigin(): string {
  const configuredOrigin = import.meta.env.VITE_BACKEND_ORIGIN as string | undefined;
  if (configuredOrigin && configuredOrigin.trim() !== "") {
    return configuredOrigin;
  }
  if (typeof window !== "undefined" && window.location.origin) {
    return window.location.origin;
  }
  return "http://127.0.0.1:8080";
}
