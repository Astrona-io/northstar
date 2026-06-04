export const useInterfaces = () => {
  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  const fetchInterfaces = (assetId) => {
    return useFetch(`${apiBase}/assets/${assetId}/interfaces`)
  }

  const createInterface = async (assetId, payload) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/assets/${assetId}/interfaces`, {
      method: 'POST',
      body: payload,
      headers: getAuthHeader()
    })
  }

  const updateInterface = async (id, payload) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/interfaces/${id}`, {
      method: 'PUT',
      body: payload,
      headers: getAuthHeader()
    })
  }

  const deleteInterface = async (id) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/interfaces/${id}`, {
      method: 'DELETE',
      headers: getAuthHeader()
    })
  }

  return {
    fetchInterfaces,
    createInterface,
    updateInterface,
    deleteInterface
  }
}
