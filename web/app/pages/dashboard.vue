<template>
  <div class="dashboard-page">
    <div class="container">
      <h1>Dashboard</h1>
      <div v-if="isLoading" class="loading">
        Loading...
      </div>
      <div v-else-if="isAuthenticated && user" class="content">
        <p class="welcome-text">Currently, <strong>{{ user.email }}</strong> is the email address of the user.</p>
      </div>
      <div v-else class="error">
        <p>You are not authenticated. Please sign in.</p>
        <NuxtLink to="/signin" class="link-button">Sign In</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { isAuthenticated, user, isLoading, checkAuth } = useAuth()

// Check authentication on page load
onMounted(async () => {
  await checkAuth()
})
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.dashboard-page {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  padding: 40px;
  width: 100%;
  max-width: 600px;
}

h1 {
  color: #333;
  margin-bottom: 30px;
  font-size: 32px;
  text-align: center;
}

.content {
  margin-top: 20px;
}

.welcome-text {
  color: #333;
  font-size: 18px;
  line-height: 1.6;
  text-align: center;
}

.welcome-text strong {
  color: #667eea;
  font-weight: 600;
}

.loading {
  text-align: center;
  color: #666;
  font-size: 16px;
  padding: 20px;
}

.error {
  text-align: center;
  color: #721c24;
  padding: 20px;
}

.error p {
  margin-bottom: 20px;
  font-size: 16px;
}

.link-button {
  display: inline-block;
  padding: 12px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  text-decoration: none;
  border-radius: 8px;
  font-weight: 600;
  transition: transform 0.2s, box-shadow 0.2s;
}

.link-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(102, 126, 234, 0.4);
}
</style>

