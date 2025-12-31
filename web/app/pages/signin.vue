<template>
  <PageContainer>
    <div class="flex justify-center">
      <div class="bg-white rounded-xl shadow-2xl py-10 px-4 w-full max-w-md">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Sign In</h1>
        <p class="text-gray-600 text-sm mb-8">
          Welcome back! Please sign in to your account
        </p>
        <form id="signinForm" @submit.prevent="handleSubmit">
          <FormInput
            id="email"
            label="Email"
            type="email"
            v-model="email"
            required
          />
          <FormInput
            id="password"
            label="Password"
            type="password"
            v-model="password"
            required
          />
          <button
            type="submit"
            class="w-full py-3.5 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg text-base font-semibold cursor-pointer transition-transform hover:-translate-y-0.5 hover:shadow-lg active:translate-y-0"
          >
            Sign In
          </button>
        </form>

        <div
          v-if="message"
          class="mt-5 p-3 rounded-lg bg-red-100 text-red-800 border border-red-200"
        >
          <div>{{ message }}</div>
        </div>
        
        <div class="text-center mt-5 text-gray-600 text-sm">
          Don't have an account?
          <NuxtLink
            to="/signup"
            class="text-purple-600 no-underline font-medium hover:underline"
            >Sign Up</NuxtLink
          >
        </div>
      </div>
    </div>
  </PageContainer>
</template>

<script setup lang="ts">
const config = useRuntimeConfig();
const apiUrl = config.public.apiUrl;
const { refreshAuth } = useAuth();

// Set page title and meta tags
useHead({
  title: 'GrocerMe | Sign In',
  meta: [
    {
      name: 'description',
      content: 'Sign in to your GrocerMe account to access your grocery lists and manage your shopping.'
    },
    {
      property: 'og:title',
      content: 'GrocerMe | Sign In'
    },
    {
      property: 'og:description',
      content: 'Sign in to your GrocerMe account to access your grocery lists and manage your shopping.'
    },
    {
      property: 'og:type',
      content: 'website'
    },
    {
      name: 'robots',
      content: 'noindex, nofollow'
    }
  ]
});

const email = ref("");
const password = ref("");
const message = ref("");

const handleSubmit = async () => {
  message.value = "";

  try {
    await $fetch<{ token?: string }>(`${apiUrl}/signin`, {
      method: "POST",
      body: {
        email: email.value,
        password: password.value,
      },
      credentials: "include",
    });

    // Refresh auth state to update the flag
    await refreshAuth();

    // Check for redirect parameter
    const route = useRoute();
    const redirectPath = route.query.redirect as string | undefined;

    // Redirect to specified path or dashboard
    await navigateTo(redirectPath || "/dashboard");
  } catch (error: any) {
    message.value =
      "Error: " +
      (error.data?.error || error.message || "Invalid email or password");
  }
};
</script>
