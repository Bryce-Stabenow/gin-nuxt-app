<template>
  <PageContainer>
    <div class="flex items-start justify-center">
      <div class="bg-white rounded-xl shadow-2xl py-10 px-4 w-full max-w-md">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Sign Up</h1>
        <p class="text-gray-600 text-sm mb-8">
          Create a new account to get started
        </p>
        <form id="signupForm" @submit.prevent="handleSubmit">
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
            :minlength="6"
          />
          <FormInput
            id="firstName"
            label="First Name"
            type="text"
            v-model="firstName"
            required
          />
          <FormInput
            id="lastName"
            label="Last Name"
            type="text"
            v-model="lastName"
            required
          />
          <FormInput
            id="avatarUrl"
            label="Avatar URL"
            type="url"
            v-model="avatarUrl"
          >
            <template #label-suffix>
              <span class="text-gray-500 font-normal">(optional)</span>
            </template>
          </FormInput>
          <button
            type="submit"
            class="w-full py-3.5 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg text-base font-semibold cursor-pointer transition-transform hover:-translate-y-0.5 hover:shadow-lg active:translate-y-0"
          >
            Sign Up
          </button>
        </form>

        <div
          v-if="message"
          class="mt-5 p-3 rounded-lg bg-red-100 text-red-800 border border-red-200"
        >
          <div>{{ message }}</div>
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

// Set page title and meta tags
useHead({
  title: "GrocerMe | Sign Up",
  meta: [
    {
      name: "description",
      content:
        "Create a new GrocerMe account to start organizing your grocery lists and simplify your shopping experience.",
    },
    {
      property: "og:title",
      content: "GrocerMe | Sign Up",
    },
    {
      property: "og:description",
      content:
        "Create a new GrocerMe account to start organizing your grocery lists and simplify your shopping experience.",
    },
    {
      property: "og:type",
      content: "website",
    },
    {
      name: "robots",
      content: "noindex, nofollow",
    },
  ],
});

const email = ref("");
const password = ref("");
const firstName = ref("");
const lastName = ref("");
const avatarUrl = ref("");
const message = ref("");

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

    await $fetch<{ token?: string }>(`${apiUrl}/signup`, {
      method: "POST",
      body,
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
      (error.data?.error || error.message || "Something went wrong");
  }
};
</script>
