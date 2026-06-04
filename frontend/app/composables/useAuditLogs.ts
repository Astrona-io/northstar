export const useAuditLogs = () => {
  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  const fetchAuditLogs = (assetId) => {
    return useFetch(`${apiBase}/assets/${assetId}/logs`)
  }

  return {
    fetchAuditLogs
  }
}
