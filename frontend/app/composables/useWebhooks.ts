export const useWebhooks = () => {
  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  const fetchWebhooks = () => {
    const { getAuthHeader } = useAuth()
    return useFetch(`${apiBase}/webhooks/`, {
      headers: getAuthHeader()
    })
  }

  const createWebhook = async (payload) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/webhooks/`, {
      method: 'POST',
      body: payload,
      headers: getAuthHeader()
    })
  }

  const deleteWebhook = async (id) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/webhooks/${id}`, {
      method: 'DELETE',
      headers: getAuthHeader()
    })
  }

  return {
    fetchWebhooks,
    createWebhook,
    deleteWebhook
  }
}
