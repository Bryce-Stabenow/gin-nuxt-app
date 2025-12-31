<template>
  <PageContainer>
    <div class="flex items-start justify-center">
      <div class="bg-white rounded-xl shadow-2xl py-10 px-4 w-full max-w-md">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Sign Up</h1>
        <p class="text-gray-600 text-sm mb-8">
          Create a new account to get started
        </p>
        <form id="signupForm" @submit.prevent="handleSubmit">
          <div class="mb-5">
            <label
              for="email"
              class="block text-gray-900 mb-2 font-medium text-sm"
              >Email</label
            >
            <input
              type="email"
              id="email"
              v-model="email"
              required
              class="w-full px-3 py-3 border-2 border-gray-200 rounded-lg text-base transition-colors focus:outline-none focus:border-purple-500"
            />
          </div>
          <div class="mb-5">
            <label
              for="password"
              class="block text-gray-900 mb-2 font-medium text-sm"
              >Password</label
            >
            <input
              type="password"
              id="password"
              v-model="password"
              required
              minlength="6"
              class="w-full px-3 py-3 border-2 border-gray-200 rounded-lg text-base transition-colors focus:outline-none focus:border-purple-500"
            />
          </div>
          <div class="mb-5">
            <label
              for="firstName"
              class="block text-gray-900 mb-2 font-medium text-sm"
              >First Name</label
            >
            <input
              type="text"
              id="firstName"
              v-model="firstName"
              required
              class="w-full px-3 py-3 border-2 border-gray-200 rounded-lg text-base transition-colors focus:outline-none focus:border-purple-500"
            />
          </div>
          <div class="mb-5">
            <label
              for="lastName"
              class="block text-gray-900 mb-2 font-medium text-sm"
              >Last Name</label
            >
            <input
              type="text"
              id="lastName"
              v-model="lastName"
              required
              class="w-full px-3 py-3 border-2 border-gray-200 rounded-lg text-base transition-colors focus:outline-none focus:border-purple-500"
            />
          </div>
          <div class="mb-5">
            <label
              for="avatarUrl"
              class="block text-gray-900 mb-2 font-medium text-sm"
              >Avatar URL
              <span class="text-gray-500 font-normal">(optional)</span></label
            >
            <input
              type="url"
              id="avatarUrl"
              v-model="avatarUrl"
              class="w-full px-3 py-3 border-2 border-gray-200 rounded-lg text-base transition-colors focus:outline-none focus:border-purple-500"
            />
          </div>
          <button
            type="submit"
            class="w-full py-3.5 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg text-base font-semibold cursor-pointer transition-transform hover:-translate-y-0.5 hover:shadow-lg active:translate-y-0"
          >
            Sign Up
          </button>
        </form>
        <div
          v-if="message"
          :class="[
            'mt-5 p-3 rounded-lg',
            messageType === 'success'
              ? 'bg-green-100 text-green-800 border border-green-200'
              : 'bg-red-100 text-red-800 border border-red-200',
          ]"
        >
          <div v-html="message"></div>
        </div>
        <div class="text-center mt-5 text-gray-600 text-sm">
          Already have an account?
          <NuxtLink
            to="/signin"
            class="text-purple-600 no-underline font-medium hover:underline"
            >Sign In</NuxtLink
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

const email = ref("");
const password = ref("");
const firstName = ref("");
const lastName = ref("");
const avatarUrl = ref("");
const message = ref("");
const messageType = ref<"success" | "error">("success");

const handleSubmit = async () => {
  message.value = "";

  try {
    const body: {
      email: string;
      password: string;
      first_name: string;
      last_name: string;
      avatar_url?: string;
    } = {
      email: email.value,
      password: password.value,
      first_name: firstName.value,
      last_name: lastName.value,
    };

    if (avatarUrl.value.trim()) {
      body.avatar_url = avatarUrl.value.trim();
    }

    const response = await $fetch<{ token?: string }>(`${apiUrl}/signup`, {
      method: "POST",
      body,
      credentials: "include",
    });

    messageType.value = "success";
    message.value =
      'Account created successfully! Cookie set. <br><button onclick="testMe()" class="mt-2.5 px-4 py-2 bg-green-600 text-white border-none rounded cursor-pointer">Test /me endpoint</button>';

    if (response.token) {
      if (typeof window !== "undefined") {
        localStorage.setItem("token", response.token);
      }
    }

    // Refresh auth state to update the flag
    await refreshAuth();

    // Check for redirect parameter
    const route = useRoute();
    const redirectPath = route.query.redirect as string | undefined;

    // Redirect to specified path or dashboard
    await navigateTo(redirectPath || "/dashboard");
  } catch (error: any) {
    messageType.value = "error";
    message.value =
      "Error: " +
      (error.data?.error || error.message || "Something went wrong");
  }
};

const testMe = async () => {
  try {
    const response = await $fetch(`${apiUrl}/me`, {
      method: "GET",
      credentials: "include",
    });

    messageType.value = "success";
    message.value =
      'Success! User info: <pre class="mt-2.5 bg-gray-100 p-2.5 rounded overflow-x-auto">' +
      JSON.stringify(response, null, 2) +
      "</pre>";
  } catch (error: any) {
    messageType.value = "error";
    message.value =
      "Error: " +
      (error.data?.error || error.message || "Failed to fetch user info");
  }
};

// Expose testMe to window for onclick handler
if (typeof window !== "undefined") {
  (window as any).testMe = testMe;
}
</script>
