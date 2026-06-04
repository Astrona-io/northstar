export const useMaintenance = () => {
  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  const fetchMaintenanceWindows = (assetId) => {
    return useFetch(`${apiBase}/assets/${assetId}/maintenance`)
  }

  const createMaintenanceWindow = async (assetId, payload) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/assets/${assetId}/maintenance`, {
      method: 'POST',
      body: payload,
      headers: getAuthHeader()
    })
  }

  const updateMaintenanceStatus = async (id, status) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/maintenance/${id}`, {
      method: 'PUT',
      body: { status },
      headers: getAuthHeader()
    })
  }

  const deleteMaintenanceWindow = async (id) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/maintenance/${id}`, {
      method: 'DELETE',
      headers: getAuthHeader()
    })
  }

  return {
    fetchMaintenanceWindows,
    createMaintenanceWindow,
    updateMaintenanceStatus,
    deleteMaintenanceWindow
  }
}
