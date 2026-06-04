export const useAssets = () => {
  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  const fetchAssets = (params = {}) => {
    return useFetch(`${apiBase}/assets/`, { params })
  }

  const fetchAsset = (id) => {
    return useFetch(`${apiBase}/assets/${id}`)
  }

  const createAsset = async (payload) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/assets/`, {
      method: 'POST',
      body: payload,
      headers: getAuthHeader()
    })
  }

  const updateAsset = async (id, payload) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/assets/${id}`, {
      method: 'PUT',
      body: payload,
      headers: getAuthHeader()
    })
  }

  const deleteAsset = async (id) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/assets/${id}`, {
      method: 'DELETE',
      headers: getAuthHeader()
    })
  }

  return {
    fetchAssets,
    fetchAsset,
    createAsset,
    updateAsset,
    deleteAsset
  }
}
