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

  return {
    fetchDatacenters,
    fetchDatacenterTypes,
    fetchRacks,
    fetchRackDetails
  }
}
