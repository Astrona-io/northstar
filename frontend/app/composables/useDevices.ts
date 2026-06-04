export const useDevices = () => {
  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  const fetchDevices = () => {
    return useFetch(`${apiBase}/devices/`)
  }

  const fetchDevice = (id) => {
    return useFetch(`${apiBase}/devices/${id}`)
  }

  const createDevice = async (payload) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/devices/`, {
      method: 'POST',
      body: payload,
      headers: getAuthHeader()
    })
  }

  const updateDevice = async (id, payload) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/devices/${id}`, {
      method: 'PUT',
      body: payload,
      headers: getAuthHeader()
    })
  }

  const deleteDevice = async (id) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/devices/${id}`, {
      method: 'DELETE',
      headers: getAuthHeader()
    })
  }

  return {
    fetchDevices,
    fetchDevice,
    createDevice,
    updateDevice,
    deleteDevice
  }
}
