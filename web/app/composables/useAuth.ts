export const useAuth = () => {
  const config = useRuntimeConfig()
  const apiUrl = config.public.apiUrl

  // Reactive state
  const isAuthenticated = useState<boolean>('auth.isAuthenticated', () => false)
  const user = useState<any>('auth.user', () => null)
  const isLoading = useState<boolean>('auth.isLoading', () => false)
  
  // Cache timestamp to avoid repeated calls
  const lastCheck = useState<number>('auth.lastCheck', () => 0)
  const CACHE_DURATION = 5 * 60 * 1000 // 5 minutes

  /**
   * Check if user is authenticated by calling the /me endpoint
   * Uses cache if available and not expired
   */
  const checkAuth = async (force = false): Promise<boolean> => {
    // On server side, don't use cache - always check for security
    // On client side, use cache to avoid repeated calls
    const now = Date.now()
    const cacheValid = !process.server && lastCheck.value > 0 && (now - lastCheck.value) < CACHE_DURATION

    // Return cached result if available and not forcing refresh (client side only)
    if (!force && cacheValid) {
      return isAuthenticated.value
    }

    isLoading.value = true

    try {
      // On server side, we need to forward cookies from the request
      const headers: Record<string, string> = {}
      if (process.server) {
        const requestHeaders = useRequestHeaders(['cookie'])
        if (requestHeaders.cookie) {
          headers.cookie = requestHeaders.cookie
        }
      }

      const userData = await $fetch(`${apiUrl}/me`, {
        method: 'GET',
        credentials: 'include',
        headers,
        retry: false
      })

      // User is authenticated
      isAuthenticated.value = true
      user.value = userData
      lastCheck.value = now
      return true
    } catch (error) {
      // User is not authenticated
      isAuthenticated.value = false
      user.value = null
      lastCheck.value = now
      return false
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Clear authentication state (for logout)
   */
  const clearAuth = () => {
    isAuthenticated.value = false
    user.value = null
    lastCheck.value = 0
  }

  /**
   * Logout user by calling the API endpoint
   */
  const logout = async (): Promise<void> => {
    isLoading.value = true

    try {
      // On server side, we need to forward cookies from the request
      const headers: Record<string, string> = {}
      if (process.server) {
        const requestHeaders = useRequestHeaders(['cookie'])
        if (requestHeaders.cookie) {
          headers.cookie = requestHeaders.cookie
        }
      }

      await $fetch(`${apiUrl}/logout`, {
        method: 'POST',
        credentials: 'include',
        headers,
        retry: false
      })
    } catch (error) {
    } finally {
      clearAuth()
      isLoading.value = false
    }
  }

  /**
   * Refresh authentication state (force check)
   */
  const refreshAuth = async (): Promise<boolean> => {
    return await checkAuth(true)
  }

  return {
    isAuthenticated: readonly(isAuthenticated),
    user: readonly(user),
    isLoading: readonly(isLoading),
    checkAuth,
    clearAuth,
    refreshAuth,
    logout
  }
}

