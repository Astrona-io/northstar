export const useUsers = () => {
  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  const fetchUsers = () => {
    const { getAuthHeader } = useAuth()
    return useFetch(`${apiBase}/users/`, {
      headers: getAuthHeader()
    })
  }

  const fetchGroups = () => {
    const { getAuthHeader } = useAuth()
    return useFetch(`${apiBase}/groups/`, {
      headers: getAuthHeader()
    })
  }

  const fetchPermissions = () => {
    const { getAuthHeader } = useAuth()
    return useFetch(`${apiBase}/permissions/`, {
      headers: getAuthHeader()
    })
  }

  const updateUserOverrides = async (id, groupId, permissionIds) => {
    const { getAuthHeader } = useAuth()
    return await $fetch(`${apiBase}/users/${id}`, {
      method: 'PUT',
      body: {
        group_id: groupId,
        permission_ids: permissionIds
      },
      headers: getAuthHeader()
    })
  }

  return {
    fetchUsers,
    fetchGroups,
    fetchPermissions,
    updateUserOverrides
  }
}
