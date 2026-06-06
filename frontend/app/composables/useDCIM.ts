export const useDCIM = () => {
  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  const fetchDatacenters = () => {
    return useFetch(`${apiBase}/datacenters/`)
  }

  const fetchDatacenterTypes = () => {
    return useFetch(`${apiBase}/datacenter-types/`)
  }

  const fetchRacks = () => {
    return useFetch(`${apiBase}/racks/`)
  }

  const fetchRackDetails = (rackId) => {
    return useFetch(`${apiBase}/racks/${rackId}`)
  }

  const createDatacenter = async (payload) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/datacenters/`, {
      method: 'POST',
      body: payload,
      headers: getAuthHeader()
    })
  }

  const updateDatacenter = async (id, payload) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/datacenters/${id}/`, {
      method: 'PUT',
      body: payload,
      headers: getAuthHeader()
    })
  }

  const deleteDatacenter = async (id) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/datacenters/${id}/`, {
      method: 'DELETE',
      headers: getAuthHeader()
    })
  }

  return {
    fetchDatacenters,
    fetchDatacenterTypes,
    fetchRacks,
    fetchRackDetails,
    createDatacenter,
    updateDatacenter,
    deleteDatacenter
  }
}
