export const useCustomFields = () => {
  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  const fetchCustomFields = (assetType = '') => {
    return useFetch(`${apiBase}/custom-fields`, {
      params: { asset_type: assetType }
    })
  }

  const createCustomField = async (payload) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/custom-fields/`, {
      method: 'POST',
      body: payload,
      headers: getAuthHeader()
    })
  }

  const deleteCustomField = async (id) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/custom-fields/${id}`, {
      method: 'DELETE',
      headers: getAuthHeader()
    })
  }

  return {
    fetchCustomFields,
    createCustomField,
    deleteCustomField
  }
}
