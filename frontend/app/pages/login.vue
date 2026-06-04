<template>
  <div class="min-h-screen bg-slate-100 dark:bg-slate-950 flex flex-col font-sans">
    
    <!-- Gitea/Forgejo Top Navbar Header -->
    <header class="bg-slate-50 dark:bg-slate-900 border-b border-slate-200 dark:border-slate-800 py-3 px-6 flex justify-between items-center flex-shrink-0">
      <div class="flex items-center gap-3">
        <UIcon name="i-heroicons-sparkles" class="h-6 w-6 text-primary-500" />
        <span class="font-bold tracking-tight text-slate-900 dark:text-white font-mono text-base">Northstar CMDB</span>
      </div>
      <div class="flex items-center gap-4 text-xs font-semibold text-slate-600 dark:text-slate-300">
        <span class="hover:text-primary-500 cursor-pointer">Explore</span>
        <span class="hover:text-primary-500 cursor-pointer">Help</span>
        <span class="text-primary-500 border-b-2 border-primary-500 pb-1 cursor-default font-bold">Sign In</span>
      </div>
    </header>

    <!-- Center Login Box Panel -->
    <div class="flex-grow flex items-center justify-center p-4">
      <div class="max-w-md w-full bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-md shadow-sm overflow-hidden">
        
        <!-- Login Header (Slate Gitea/Forgejo style) -->
        <div class="bg-slate-50 dark:bg-slate-800/80 px-6 py-4 border-b border-slate-200 dark:border-slate-800">
          <h2 class="text-sm font-bold text-slate-800 dark:text-white uppercase tracking-wider flex items-center gap-2">
            <UIcon name="i-heroicons-shield-check" class="h-5 w-5 text-primary-500" />
            Northstar Security Authentication
          </h2>
        </div>

        <div class="p-6 space-y-6">
          <p class="text-xs text-slate-500 dark:text-slate-400 leading-normal">
            Access secure configuration items, network patch cords, and DCIM rack layout configurations. Please authenticate using role-authorized credentials.
          </p>

          <form @submit.prevent="handleLogin" class="space-y-4">
            <UFormGroup label="Username" :error="loginError ? ' ' : ''" class="font-semibold text-slate-700 dark:text-slate-300 text-xs">
              <UInput v-model="username" icon="i-heroicons-user" placeholder="admin, operator, or viewer" required />
            </UFormGroup>

            <UFormGroup label="Password" :error="loginError" class="font-semibold text-slate-700 dark:text-slate-300 text-xs">
              <UInput v-model="password" type="password" icon="i-heroicons-key" placeholder="Enter password..." required />
            </UFormGroup>

            <div class="pt-2">
              <UButton type="submit" block color="primary" size="sm" :loading="isLoggingIn" icon="i-heroicons-arrow-right-end-on-rectangle" class="font-bold">
                Sign In
              </UButton>
            </div>
          </form>

          <!-- Mock Accounts Info (GitLab Site Info Panel style) -->
          <div class="bg-slate-50 dark:bg-slate-800/40 border border-slate-200 dark:border-slate-800 rounded-md p-4">
            <h4 class="text-[10px] font-bold text-slate-500 dark:text-slate-400 uppercase tracking-wider mb-2 font-mono">Pre-seeded Cluster Roles:</h4>
            <ul class="text-[11px] text-slate-600 dark:text-slate-300 space-y-2 font-mono">
              <li class="flex items-center justify-between gap-2">
                <span>🛡️ Admin:</span>
                <span class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 px-1.5 py-0.5 rounded text-[10px]">admin / admin</span>
              </li>
              <li class="flex items-center justify-between gap-2">
                <span>⚙️ Operator:</span>
                <span class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 px-1.5 py-0.5 rounded text-[10px]">operator / operator</span>
              </li>
              <li class="flex items-center justify-between gap-2">
                <span>👁️ Viewer:</span>
                <span class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 px-1.5 py-0.5 rounded text-[10px]">viewer / viewer</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

definePageMeta({
  layout: false // Blank layout for login
})

const { login, loginError, isLoggingIn } = useAuth()

const username = ref('')
const password = ref('')

const handleLogin = async () => {
  const success = await login(username.value, password.value)
  if (success) {
    navigateTo('/')
  }
}
</script>
