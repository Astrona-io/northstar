export const useRelationships = () => {
  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  const fetchRelationships = (assetId) => {
    return useFetch(`${apiBase}/assets/${assetId}/relationships`)
  }

  const syncRelationships = async (assetId, uplinkIds, downlinkIds) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/assets/${assetId}/relationships`, {
      method: 'POST',
      body: {
        uplink_ids: uplinkIds,
        downlink_ids: downlinkIds
      },
      headers: getAuthHeader()
    })
  }

  return {
    fetchRelationships,
    syncRelationships
  }
}
