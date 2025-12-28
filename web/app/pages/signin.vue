<template>
  <div class="signin-page">
    <div class="container">
      <h1>Sign In</h1>
      <p class="subtitle">Welcome back! Please sign in to your account</p>
      <form id="signinForm" @submit.prevent="handleSubmit">
        <div class="form-group">
          <label for="email">Email</label>
          <input
            type="email"
            id="email"
            v-model="email"
            required
          />
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <input
            type="password"
            id="password"
            v-model="password"
            required
          />
        </div>
        <button type="submit">Sign In</button>
      </form>
      <div v-if="message" :class="['message', messageType, 'show']">
        <div v-html="message"></div>
      </div>
      <div class="link">
        Don't have an account? <NuxtLink to="/signup">Sign Up</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const config = useRuntimeConfig()
const apiUrl = config.public.apiUrl
const { refreshAuth } = useAuth()

const email = ref('')
const password = ref('')
const message = ref('')
const messageType = ref<'success' | 'error'>('success')

const handleSubmit = async () => {
  message.value = ''
  
  try {
    const response = await $fetch<{ token?: string }>(`${apiUrl}/signin`, {
      method: 'POST',
      body: {
        email: email.value,
        password: password.value
      },
      credentials: 'include'
    })
    
    messageType.value = 'success'
    message.value = 'Signed in successfully! Cookie set. <br><button onclick="testMe()" style="margin-top: 10px; padding: 8px 16px; background: #28a745; color: white; border: none; border-radius: 4px; cursor: pointer;">Test /me endpoint</button>'
    
    if (response.token) {
      if (typeof window !== 'undefined') {
        localStorage.setItem('token', response.token)
      }
    }

    // Refresh auth state to update the flag
    await refreshAuth()
    
    // Redirect to dashboard after successful signin
    await navigateTo('/dashboard')
  } catch (error: any) {
    messageType.value = 'error'
    message.value = 'Error: ' + (error.data?.error || error.message || 'Invalid email or password')
  }
}

const testMe = async () => {
  try {
    const response = await $fetch(`${apiUrl}/me`, {
      method: 'GET',
      credentials: 'include'
    })
    
    messageType.value = 'success'
    message.value = 'Success! User info: <pre style="margin-top: 10px; background: #f5f5f5; padding: 10px; border-radius: 4px; overflow-x: auto;">' + JSON.stringify(response, null, 2) + '</pre>'
  } catch (error: any) {
    messageType.value = 'error'
    message.value = 'Error: ' + (error.data?.error || error.message || 'Failed to fetch user info')
  }
}

// Expose testMe to window for onclick handler
if (typeof window !== 'undefined') {
  ;(window as any).testMe = testMe
}
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.signin-page {
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
  max-width: 400px;
}

h1 {
  color: #333;
  margin-bottom: 10px;
  font-size: 28px;
}

.subtitle {
  color: #666;
  margin-bottom: 30px;
  font-size: 14px;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  color: #333;
  margin-bottom: 8px;
  font-weight: 500;
  font-size: 14px;
}

input {
  width: 100%;
  padding: 12px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.3s;
}

input:focus {
  outline: none;
  border-color: #667eea;
}

button {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

button:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(102, 126, 234, 0.4);
}

button:active {
  transform: translateY(0);
}

.message {
  margin-top: 20px;
  padding: 12px;
  border-radius: 8px;
  display: none;
}

.message.success {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.message.error {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.message.show {
  display: block;
}

.link {
  text-align: center;
  margin-top: 20px;
  color: #666;
  font-size: 14px;
}

.link a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.link a:hover {
  text-decoration: underline;
}
</style>

