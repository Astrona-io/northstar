export default defineNuxtRouteMiddleware((to) => {
  const { isAuthenticated, isAdmin } = useAuth()

  // Avoid infinite redirect loop by skipping middleware on login page
  if (to.path === '/login') {
    return
  }

  // If user is anonymous, redirect them immediately to the Sign-in page (Auth Requirement)
  if (!isAuthenticated.value) {
    return navigateTo('/login')
  }

  // Admin page guards (Phase 1 RBAC UI protection)
  if (to.path.startsWith('/settings') && !isAdmin.value) {
    return navigateTo('/')
  }
})
