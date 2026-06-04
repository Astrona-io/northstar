export const useDiscovery = () => {
  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  const scanSubnet = (subnet) => {
    return useFetch(`${apiBase}/discovery/scan`, {
      params: { subnet },
      lazy: true,
      server: false
    })
  }

  const importDevice = async (device) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/discovery/import`, {
      method: 'POST',
      body: device,
      headers: getAuthHeader()
    })
  }

  return {
    scanSubnet,
    importDevice
  }
}
