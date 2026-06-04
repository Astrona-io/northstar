import { ref, computed } from 'vue'

export const useAuth = () => {
  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  // Keep session in a cookie so it persists across page reloads
  const sessionCookie = useCookie('cmdb-session', {
    maxAge: 60 * 60 * 24 // 1 day
  })

  const loginError = ref('')
  const isLoggingIn = ref(false)

  const user = computed(() => {
    return sessionCookie.value || null
  })

  const isAuthenticated = computed(() => !!user.value?.token)
  
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isOperator = computed(() => user.value?.role === 'operator')
  const isViewer = computed(() => user.value?.role === 'viewer')

  const canMutate = computed(() => {
    return user.value?.role === 'admin' || user.value?.role === 'operator'
  })

  const canDelete = computed(() => {
    return user.value?.role === 'admin'
  })

  const login = async (username, password) => {
    loginError.value = ''
    isLoggingIn.value = true
    try {
      const data = await $fetch(`${apiBase}/auth/login`, {
        method: 'POST',
        body: { username, password }
      })
      
      if (data && data.token) {
        sessionCookie.value = {
          token: data.token,
          username: data.username,
          role: data.role
        }
        return true
      }
    } catch (err) {
      console.error('Login failed:', err)
      loginError.value = 'Invalid username or password'
    } finally {
      isLoggingIn.value = false
    }
    return false
  }

  const logout = () => {
    sessionCookie.value = null
    navigateTo('/login')
  }

  // Helper to append Authorization header to API calls (if authenticated)
  const getAuthHeader = () => {
    if (isAuthenticated.value) {
      return { Authorization: `Bearer ${user.value.token}` }
    }
    return {}
  }

  return {
    user,
    isAuthenticated,
    isAdmin,
    isOperator,
    isViewer,
    canMutate,
    canDelete,
    loginError,
    isLoggingIn,
    login,
    logout,
    getAuthHeader
  }
}
