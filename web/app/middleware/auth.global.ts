export default defineNuxtRouteMiddleware(async (to) => {
  const { checkAuth } = useAuth()
  
  // Check authentication (uses cache if available)
  // This works on both server and client
  const authenticated = await checkAuth()

  // If trying to access signin/signup pages while authenticated, redirect to dashboard
  if ((to.path === '/signin' || to.path === '/signup') && authenticated) {
    return navigateTo('/dashboard')
  }

  // If trying to access dashboard without authentication, redirect to signin
  if (to.path === '/dashboard' && !authenticated) {
    return navigateTo('/signin')
  }
})

