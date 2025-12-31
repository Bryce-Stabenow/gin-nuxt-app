export default defineNuxtRouteMiddleware(async (to) => {
  const { checkAuth } = useAuth();

  // Allow share routes without authentication - they handle their own auth flow
  if (to.path.startsWith("/lists/share/")) {
    return;
  }

  // Check authentication (uses cache if available)
  // This works on both server and client
  const authenticated = await checkAuth();

  // If trying to access signin/signup pages while authenticated, redirect to dashboard
  if ((to.path === "/signin" || to.path === "/signup") && authenticated) {
    return navigateTo("/dashboard");
  }

  // If trying to access dashboard without authentication, redirect to signin
  if (to.path === "/dashboard" && !authenticated) {
    return navigateTo("/signin");
  }

  // If trying to access list pages without authentication, redirect to signin
  if (to.path.startsWith("/lists/") && !authenticated) {
    return navigateTo("/signin");
  }
});
