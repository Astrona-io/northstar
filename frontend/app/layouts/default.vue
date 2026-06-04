<template>
  <div class="flex h-screen bg-slate-100 dark:bg-slate-950 overflow-hidden font-sans">
    
    <!-- Left Sidebar (GitLab / Gitea Vertical Navigation Sidebar) -->
    <aside 
      :class="isSidebarCollapsed ? 'w-20' : 'w-72'" 
      class="bg-slate-50 dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 flex flex-col justify-between h-screen p-4 flex-shrink-0 shadow-sm transition-all duration-300 ease-in-out"
    >
      
      <!-- Top Portion: Logo & Standard Navigation -->
      <div class="space-y-6">
        <!-- Logo Branding (Forgejo Style) -->
        <div class="flex items-center justify-between px-2 py-1.5 border-b border-slate-200 dark:border-slate-800 pb-4">
          <div class="flex items-center gap-3">
            <UIcon name="i-heroicons-sparkles" class="h-6 w-6 text-primary-500 shrink-0" />
            <h1 v-if="!isSidebarCollapsed" class="text-xl font-bold tracking-tight text-slate-900 dark:text-white font-mono shrink-0">
              Northstar
            </h1>
          </div>
          <!-- Toggle Collapse Button -->
          <UButton
            :icon="isSidebarCollapsed ? 'i-heroicons-chevron-double-right' : 'i-heroicons-chevron-double-left'"
            color="gray"
            variant="ghost"
            size="xs"
            @click="isSidebarCollapsed = !isSidebarCollapsed"
            :class="isSidebarCollapsed ? 'mx-auto' : ''"
            title="Toggle Sidebar"
          />
        </div>

        <!-- Standard Operational Navigation List -->
        <nav class="space-y-1">
          <UTooltip text="Assets" :prevent="!isSidebarCollapsed" :popper="{ placement: 'right' }" class="block">
            <UButton 
              to="/" 
              variant="ghost" 
              color="gray" 
              block 
              :active-class="isSidebarCollapsed 
                ? 'bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-white font-bold rounded-md' 
                : 'bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-white font-bold border-l-4 border-primary-500 rounded-none'" 
              :class="[
                isSidebarCollapsed ? 'justify-center px-0' : 'justify-start text-left px-3 border-l-4 border-transparent',
                'py-2.5 rounded-md transition-all text-xs'
              ]"
              icon="i-heroicons-rectangle-stack"
            >
              <span v-if="!isSidebarCollapsed">Assets</span>
            </UButton>
          </UTooltip>
          
          <UTooltip text="Datacenters" :prevent="!isSidebarCollapsed" :popper="{ placement: 'right' }" class="block">
            <UButton 
              to="/dcim" 
              variant="ghost" 
              color="gray" 
              block 
              :active-class="isSidebarCollapsed 
                ? 'bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-white font-bold rounded-md' 
                : 'bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-white font-bold border-l-4 border-primary-500 rounded-none'" 
              :class="[
                isSidebarCollapsed ? 'justify-center px-0' : 'justify-start text-left px-3 border-l-4 border-transparent',
                'py-2.5 rounded-md transition-all text-xs'
              ]"
              icon="i-heroicons-server"
            >
              <span v-if="!isSidebarCollapsed">Datacenters</span>
            </UButton>
          </UTooltip>
          
          <UTooltip text="Topology" :prevent="!isSidebarCollapsed" :popper="{ placement: 'right' }" class="block">
            <UButton 
              to="/topology" 
              variant="ghost" 
              color="gray" 
              block 
              :active-class="isSidebarCollapsed 
                ? 'bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-white font-bold rounded-md' 
                : 'bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-white font-bold border-l-4 border-primary-500 rounded-none'" 
              :class="[
                isSidebarCollapsed ? 'justify-center px-0' : 'justify-start text-left px-3 border-l-4 border-transparent',
                'py-2.5 rounded-md transition-all text-xs'
              ]"
              icon="i-heroicons-globe-alt"
            >
              <span v-if="!isSidebarCollapsed">Topology</span>
            </UButton>
          </UTooltip>
          
          <UTooltip text="Device Catalog" :prevent="!isSidebarCollapsed" :popper="{ placement: 'right' }" class="block">
            <UButton 
              to="/devices" 
              variant="ghost" 
              color="gray" 
              block 
              :active-class="isSidebarCollapsed 
                ? 'bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-white font-bold rounded-md' 
                : 'bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-white font-bold border-l-4 border-primary-500 rounded-none'" 
              :class="[
                isSidebarCollapsed ? 'justify-center px-0' : 'justify-start text-left px-3 border-l-4 border-transparent',
                'py-2.5 rounded-md transition-all text-xs'
              ]"
              icon="i-heroicons-cpu-chip"
            >
              <span v-if="!isSidebarCollapsed">Device Catalog</span>
            </UButton>
          </UTooltip>
        </nav>
      </div>

      <!-- Bottom Portion: Separated Site Admin & User Controls (Phase 1 Admin Re-org) -->
      <div class="space-y-4 border-t border-slate-200 dark:border-slate-800 pt-4">
        
        <!-- GitOps Exporter Button (Phase 1 GitOps YAML Generator) -->
        <UTooltip text="GitOps YAML Exporter" :prevent="!isSidebarCollapsed" :popper="{ placement: 'right' }" class="block">
          <UButton 
            variant="ghost" 
            color="gray" 
            block 
            @click="isGitOpsDrawerOpen = true"
            :class="[
              isSidebarCollapsed ? 'justify-center px-0' : 'justify-start text-left px-3 border-l-4 border-transparent',
              'py-2.5 rounded-md transition-all text-xs text-blue-600 dark:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-950/20'
            ]"
            icon="i-heroicons-code-bracket-square"
          >
            <span v-if="!isSidebarCollapsed">GitOps YAML Exporter</span>
          </UButton>
        </UTooltip>

        <!-- Site Admin button (Anchored at the bottom separately, GitLab style) -->
        <UTooltip text="Site Admin" :prevent="!isSidebarCollapsed" :popper="{ placement: 'right' }" class="block">
          <UButton 
            v-if="isAdmin" 
            to="/settings" 
            variant="ghost" 
            color="gray" 
            block 
            :active-class="isSidebarCollapsed 
              ? 'bg-red-50 dark:bg-red-950/30 text-red-600 dark:text-red-400 font-bold rounded-md' 
              : 'bg-red-50 dark:bg-red-950/30 text-red-600 dark:text-red-400 font-bold border-l-4 border-red-500 rounded-none'" 
            :class="[
              isSidebarCollapsed ? 'justify-center px-0' : 'justify-start text-left px-3 border-l-4 border-transparent',
              'py-2.5 rounded-md transition-all text-xs text-red-600 dark:text-red-500 hover:bg-red-50 dark:hover:bg-red-950/20'
            ]"
            icon="i-heroicons-shield-check"
          >
            <span v-if="!isSidebarCollapsed">Site Admin</span>
          </UButton>
        </UTooltip>

        <!-- Session Profile Card -->
        <div v-if="isAuthenticated" :class="isSidebarCollapsed ? 'p-1' : 'p-3'" class="space-y-3 bg-white dark:bg-slate-900 rounded-md border border-slate-200 dark:border-slate-800 shadow-sm transition-all duration-300">
          <div v-if="!isSidebarCollapsed" class="flex items-center justify-between gap-2">
            <div class="truncate flex-1">
              <span class="text-[9px] text-slate-400 uppercase font-bold block leading-none mb-1 font-mono">Signed Operator</span>
              <span class="text-xs font-semibold text-slate-800 dark:text-white truncate block">{{ user.username }}</span>
            </div>
            <UBadge :color="isAdmin ? 'red' : (isOperator ? 'blue' : 'gray')" size="xs" variant="subtle" class="uppercase text-[8px] font-bold font-mono">{{ user.role }}</UBadge>
          </div>
          
          <div :class="isSidebarCollapsed ? 'flex-col gap-2 align-center' : 'flex-row justify-between pt-2 border-t border-slate-200 dark:border-slate-800'" class="flex items-center">
            <!-- Theme Toggle -->
            <ClientOnly>
              <UTooltip :text="isDark ? 'Switch to Light Theme' : 'Switch to Dark Theme'" :prevent="!isSidebarCollapsed" :popper="{ placement: 'right' }">
                <UButton
                  :icon="isDark ? 'i-heroicons-moon-20-solid' : 'i-heroicons-sun-20-solid'"
                  color="gray"
                  variant="ghost"
                  size="xs"
                  @click="isDark = !isDark"
                />
              </UTooltip>
              <template #fallback>
                <div class="w-6 h-6" />
              </template>
            </ClientOnly>
            
            <UTooltip text="Sign Out" :prevent="!isSidebarCollapsed" :popper="{ placement: 'right' }">
              <UButton 
                size="xs" 
                variant="ghost" 
                color="red" 
                icon="i-heroicons-arrow-right-start-on-rectangle" 
                :label="isSidebarCollapsed ? '' : 'Sign Out'" 
                :class="[
                  isSidebarCollapsed ? 'p-1' : 'text-[10px] px-2 py-1',
                  'hover:bg-red-50 dark:hover:bg-red-950/20'
                ]"
                @click="logout" 
              />
            </UTooltip>
          </div>
        </div>

        <UTooltip v-else text="Sign In" :prevent="!isSidebarCollapsed" :popper="{ placement: 'right' }" class="block">
          <UButton 
            to="/login" 
            block 
            size="sm" 
            variant="outline" 
            icon="i-heroicons-arrow-right-end-on-rectangle"
            :class="isSidebarCollapsed ? 'justify-center px-0' : ''"
          >
            <span v-if="!isSidebarCollapsed">Sign In</span>
          </UButton>
        </UTooltip>

        <!-- Muted Footer (Astrona Branding Requirements) -->
        <div class="pt-3 border-t border-slate-200 dark:border-slate-800 text-[10px] text-slate-400 dark:text-slate-500 font-sans space-y-1 mt-2 shrink-0 text-center">
          <p v-if="!isSidebarCollapsed" class="leading-relaxed">
            Northstar &ndash; Owned by <a href="https://astrona.io" target="_blank" class="hover:underline hover:text-primary-500 font-semibold font-mono">Astrona.io</a>
          </p>
          <p v-if="!isSidebarCollapsed" class="text-[9px] font-semibold text-slate-500 font-mono leading-none">
            {{ appVersion }}
          </p>
          <div :class="isSidebarCollapsed ? 'flex flex-col gap-1.5 items-center text-center' : 'flex items-center gap-2 justify-between'" class="flex">
            <!-- License Link (Triggering Modal) -->
            <UTooltip text="GNU AGPLv3 License" :prevent="!isSidebarCollapsed" :popper="{ placement: 'right' }">
              <button @click="openLicenseModal" class="hover:underline hover:text-primary-500 flex items-center gap-1 text-left">
                <UIcon name="i-heroicons-document-text" class="h-3.5 w-3.5 shrink-0 text-slate-400 dark:text-slate-500" />
                <span v-if="!isSidebarCollapsed">GNU AGPLv3</span>
              </button>
            </UTooltip>
            <!-- GitHub Link (Corrected URL) -->
            <UTooltip text="GitHub Repository" :prevent="!isSidebarCollapsed" :popper="{ placement: 'right' }">
              <a href="https://github.com/Astrona-io/northstar" target="_blank" class="hover:underline hover:text-primary-500 flex items-center gap-1">
                <UIcon name="i-heroicons-code-bracket" class="h-3.5 w-3.5 shrink-0 text-slate-400 dark:text-slate-500" />
                <span v-if="!isSidebarCollapsed">GitHub Repo</span>
              </a>
            </UTooltip>
          </div>
        </div>
      </div>
    </aside>

    <!-- Right Column Main Content Port -->
    <main 
      :class="isSidebarCollapsed ? 'px-3 pt-3 pb-2' : 'px-6 pt-6 pb-3'" 
      class="flex-1 overflow-y-auto h-screen bg-slate-100 dark:bg-slate-950 transition-all duration-300"
    >
      <div 
        :class="isSidebarCollapsed ? 'max-w-none' : 'max-w-[1500px]'" 
        class="mx-auto transition-all duration-300"
      >
        
        <!-- Global Northstar Onboarding Learning Checklist (Phase 5 Global Checklist) -->
        <UCard v-if="showOnboarding" class="mb-4 border border-dashed border-primary-500 bg-primary-50/10 dark:bg-primary-950/10 rounded-md relative shadow-sm">
          <div class="absolute top-4 right-4">
            <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark" @click="dismissOnboarding" />
          </div>
          <div class="flex gap-4 items-start pr-8">
            <UIcon name="i-heroicons-academic-cap" class="h-10 w-10 text-primary-500 animate-bounce flex-shrink-0" />
            <div class="space-y-2 w-full">
              <h3 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2 font-mono">
                Northstar Quickstart: Your Guided Learning Path
              </h3>
              <p class="text-xs text-slate-500 dark:text-slate-400 leading-normal">
                Welcome to Northstar! Follow this interactive onboarding checklist to successfully model your environment. Steps are checked off automatically in real-time as you config the platform across all pages.
              </p>
              
              <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 pt-3">
                <!-- Step 1: DC & Rack -->
                <button 
                  type="button" 
                  @click="navigateTo('/settings/dcim')" 
                  class="flex items-center gap-3 px-3 py-2 bg-white dark:bg-slate-900 rounded-md shadow-sm border border-slate-200 dark:border-slate-800 hover:border-primary-500 hover:bg-slate-50 dark:hover:bg-slate-800/55 text-left w-full transition-all cursor-pointer group"
                >
                  <UIcon :name="step3Done ? 'i-heroicons-check-circle' : 'i-heroicons-minus-circle'" :class="step3Done ? 'text-green-500' : 'text-slate-400 group-hover:text-primary-500'" class="h-5 w-5 flex-shrink-0 transition-colors" />
                  <span class="text-xs leading-none" :class="step3Done ? 'line-through text-slate-400' : 'font-semibold text-slate-700 dark:text-slate-200 group-hover:text-primary-500 transition-colors'">1. Define DC & Rack</span>
                </button>

                <!-- Step 2: Register Brand -->
                <button 
                  type="button" 
                  @click="navigateTo('/devices')" 
                  class="flex items-center gap-3 px-3 py-2 bg-white dark:bg-slate-900 rounded-md shadow-sm border border-slate-200 dark:border-slate-800 hover:border-primary-500 hover:bg-slate-50 dark:hover:bg-slate-800/55 text-left w-full transition-all cursor-pointer group"
                >
                  <UIcon :name="step1Done ? 'i-heroicons-check-circle' : 'i-heroicons-minus-circle'" :class="step1Done ? 'text-green-500' : 'text-slate-400 group-hover:text-primary-500'" class="h-5 w-5 flex-shrink-0 transition-colors" />
                  <span class="text-xs leading-none" :class="step1Done ? 'line-through text-slate-400' : 'font-semibold text-slate-700 dark:text-slate-200 group-hover:text-primary-500 transition-colors'">2. Register Brand</span>
                </button>

                <!-- Step 3: Define Spec -->
                <button 
                  type="button" 
                  @click="navigateTo('/devices')" 
                  class="flex items-center gap-3 px-3 py-2 bg-white dark:bg-slate-900 rounded-md shadow-sm border border-slate-200 dark:border-slate-800 hover:border-primary-500 hover:bg-slate-50 dark:hover:bg-slate-800/55 text-left w-full transition-all cursor-pointer group"
                >
                  <UIcon :name="step2Done ? 'i-heroicons-check-circle' : 'i-heroicons-minus-circle'" :class="step2Done ? 'text-green-500' : 'text-slate-400 group-hover:text-primary-500'" class="h-5 w-5 flex-shrink-0 transition-colors" />
                  <span class="text-xs leading-none" :class="step2Done ? 'line-through text-slate-400' : 'font-semibold text-slate-700 dark:text-slate-200 group-hover:text-primary-500 transition-colors'">3. Define Spec Model</span>
                </button>

                <!-- Step 4: Asset Mount -->
                <button 
                  type="button" 
                  @click="navigateTo('/')" 
                  class="flex items-center gap-3 px-3 py-2 bg-white dark:bg-slate-900 rounded-md shadow-sm border border-slate-200 dark:border-slate-800 hover:border-primary-500 hover:bg-slate-50 dark:hover:bg-slate-800/55 text-left w-full transition-all cursor-pointer group"
                >
                  <UIcon :name="step4Done ? 'i-heroicons-check-circle' : 'i-heroicons-minus-circle'" :class="step4Done ? 'text-green-500' : 'text-slate-400 group-hover:text-primary-500'" class="h-5 w-5 flex-shrink-0 transition-colors" />
                  <span class="text-xs leading-none" :class="step4Done ? 'line-through text-slate-400' : 'font-semibold text-slate-700 dark:text-slate-200 group-hover:text-primary-500 transition-colors'">4. Mount First Asset</span>
                </button>
              </div>
            </div>
          </div>
        </UCard>

        <slot />
      </div>
    </main>

    <!-- GitOps YAML Exporter Drawer (Phase 1 GitOps Reconciliation) -->
    <USlideover v-model="isGitOpsDrawerOpen" :overlay="true" class="z-50">
      <div class="h-full flex flex-col bg-slate-50 dark:bg-slate-900 shadow-xl overflow-hidden font-sans">
        <div class="p-4 border-b border-slate-200 dark:border-slate-800 flex justify-between items-center bg-white dark:bg-slate-950 flex-shrink-0">
          <div>
            <h3 class="font-bold text-slate-900 dark:text-white font-mono flex items-center gap-2">
              <UIcon name="i-heroicons-code-bracket-square" class="h-5 w-5 text-blue-500" />
              GitOps Spec Generator
            </h3>
            <span class="text-[10px] text-slate-500">Generate valid declarative YAML manifests</span>
          </div>
          <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark" @click="isGitOpsDrawerOpen = false" />
        </div>
        
        <div class="p-4 flex-1 overflow-y-auto space-y-6">
          <p class="text-xs text-slate-600 dark:text-slate-400 leading-relaxed">
            Northstar's Declarative Reconciliation Engine automatically syncs YAML specifications from Git repositories or local PVC mounts to the GORM database. Use these templates to declare state safely.
          </p>

          <UFormGroup label="Resource Target Kind" class="text-xs font-semibold">
            <USelect v-model="gitOpsKind" :options="['Asset', 'Datacenter']" />
          </UFormGroup>

          <UFormGroup v-if="gitOpsKind === 'Asset'" label="Specific Hardware Category" class="text-xs font-semibold">
            <USelect v-model="gitOpsAssetCategory" :options="['Server', 'Router', 'Switch', 'Container']" />
          </UFormGroup>
          
          <div class="relative group mt-2">
            <div class="absolute top-2 right-2 z-10 flex gap-2">
              <UButton size="xs" color="gray" variant="solid" icon="i-heroicons-clipboard-document" @click="copyYaml">Copy YAML</UButton>
            </div>
            <pre class="bg-slate-950 text-slate-300 p-4 rounded-md text-[10px] sm:text-xs overflow-x-auto border border-slate-800 font-mono shadow-inner leading-relaxed"><code>{{ generatedYamlSpec }}</code></pre>
          </div>
        </div>

        <div class="p-4 border-t border-slate-200 dark:border-slate-800 bg-slate-100 dark:bg-slate-950 text-[10px] text-slate-500 font-mono flex-shrink-0">
          Mount to <code>/etc/astrona/gitops/*.yaml</code> for background sync
        </div>
      </div>
    </USlideover>

    <!-- GNU AGPLv3 License Modal Dialog -->
    <UModal v-model="isLicenseModalOpen" :ui="{ width: 'sm:max-w-4xl' }" :prevent-close="showAgreementButton">
      <UCard :ui="{ body: { padding: 'p-0' } }">
        <template #header>
          <div class="flex items-center justify-between">
            <div class="flex flex-col gap-1">
              <h3 class="text-sm font-bold text-slate-900 dark:text-white font-mono flex items-center gap-2">
                <UIcon name="i-heroicons-document-text" class="text-primary-500 h-5 w-5" />
                GNU Affero General Public License v3 (AGPLv3)
              </h3>
              <span class="text-[10px] text-slate-400 dark:text-slate-500 font-mono pl-7">
                Last updated: June 4, 2026 (Adopted AGPLv3)
              </span>
            </div>
            <UButton v-if="!showAgreementButton" color="gray" variant="ghost" icon="i-heroicons-x-mark" @click="isLicenseModalOpen = false" />
          </div>
        </template>
        
        <div 
          @scroll="onLicenseScroll"
          class="p-6 overflow-y-auto max-h-[60vh] bg-slate-950 text-slate-300 font-mono text-[11px] leading-relaxed rounded-b-md whitespace-pre-wrap"
        >
          <div v-if="isLoadingLicense" class="flex justify-center p-8">
            <UIcon name="i-heroicons-arrow-path" class="animate-spin h-6 w-6 text-primary-500" />
          </div>
          <div v-else>
            {{ licenseText }}
          </div>
        </div>

        <template v-if="showAgreementButton" #footer>
          <div class="p-4 bg-slate-50 dark:bg-slate-900 border-t border-slate-200 dark:border-slate-800 flex flex-col sm:flex-row items-center justify-between gap-4 rounded-b-md">
            <div class="flex items-start gap-2 max-w-xl">
              <UIcon name="i-heroicons-information-circle" class="h-5 w-5 text-amber-500 shrink-0 mt-0.5" />
              <div class="space-y-1.5 text-[11px] text-slate-500 dark:text-slate-400 leading-normal font-sans text-left">
                <p>
                  As an Administrator, you are signing this binding copyleft agreement on behalf of your deployment. Access to the CMDB administration systems requires active AGPLv3 compliance registration.
                </p>
                <p class="text-[10px] text-slate-400 dark:text-slate-500 font-semibold">
                  Commercial licensing exceptions and enterprise SLAs are available. Contact <a href="https://astrona.io" target="_blank" class="text-primary-500 hover:underline">Astrona.io</a> or email <span class="font-mono">contact@astrona.io</span> for custom corporate licensing.
                </p>
              </div>
            </div>
            <UButton 
              color="green" 
              icon="i-heroicons-check-circle"
              :loading="isSavingAgreement"
              :disabled="!hasReadLicense"
              @click="submitAcceptLicense"
              class="font-bold shrink-0 shadow-sm"
              :class="hasReadLicense ? 'animate-pulse' : 'opacity-50 cursor-not-allowed'"
            >
              I Agree & Accept AGPLv3 Terms
            </UButton>
          </div>
        </template>
      </UCard>
    </UModal>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'

const colorMode = useColorMode()
const isSidebarCollapsed = ref(false)

const isLicenseModalOpen = ref(false)
const isLoadingLicense = ref(false)
const licenseText = ref('')
const hasReadLicense = ref(false)

const showAgreementButton = ref(false)
const isSavingAgreement = ref(false)

const appVersion = computed(() => {
  return import.meta.dev ? 'v0.0.0-dev' : 'v0.0.1'
})

const onLicenseScroll = (event) => {
  const el = event.target
  const scrollHeight = el.scrollHeight - el.clientHeight
  if (scrollHeight <= 0) {
    hasReadLicense.value = true
    return
  }
  const scrollPercentage = (el.scrollTop / scrollHeight) * 100
  if (scrollPercentage >= 80) {
    hasReadLicense.value = true
  }
}

const openLicenseModal = async () => {
  isLicenseModalOpen.value = true
  if (showAgreementButton.value) {
    hasReadLicense.value = false
  } else {
    hasReadLicense.value = true
  }
  
  if (!licenseText.value) {
    isLoadingLicense.value = true
    try {
      licenseText.value = await $fetch('/LICENSE', { parseResponse: txt => txt })
    } catch (err) {
      console.error('Failed to load license text:', err)
      licenseText.value = 'GNU Affero General Public License v3 (AGPLv3) - Copyright (C) 2026 Astrona (astrona.io)'
    } finally {
      isLoadingLicense.value = false
    }
  }
}

const { user, isAuthenticated, isAdmin, isOperator, logout, getAuthHeader } = useAuth()
const route = useRoute()

const checkLicenseAgreementStatus = async () => {
  if (!isAuthenticated.value || !isAdmin.value) return

  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  try {
    const res = await $fetch(`${apiBase}/license/status`, {
      headers: getAuthHeader()
    })
    if (res && !res.accepted) {
      showAgreementButton.value = true
      await openLicenseModal()
    } else {
      showAgreementButton.value = false
    }
  } catch (err) {
    console.error('Failed to check license agreement status:', err)
  }
}

const submitAcceptLicense = async () => {
  isSavingAgreement.value = true
  const runtimeConfig = useRuntimeConfig()
  const apiBase = runtimeConfig.public.apiBase

  try {
    await $fetch(`${apiBase}/license/accept`, {
      method: 'POST',
      headers: getAuthHeader()
    })
    showAgreementButton.value = false
    isLicenseModalOpen.value = false
    
    const toast = useToast()
    toast.add({
      title: 'Agreement Registered',
      description: 'AGPLv3 License Agreement accepted and logged successfully.',
      color: 'green',
      icon: 'i-heroicons-check-circle'
    })
  } catch (err) {
    console.error('Failed to accept license:', err)
    const toast = useToast()
    toast.add({
      title: 'Registration Failed',
      description: 'Failed to register license agreement. Please try again.',
      color: 'red',
      icon: 'i-heroicons-x-circle'
    })
  } finally {
    isSavingAgreement.value = false
  }
}

onMounted(() => {
  checkLicenseAgreementStatus()
})

watch([isAuthenticated, isAdmin], () => {
  checkLicenseAgreementStatus()
})

// Global Checklist Logic (Phase 5)
const { fetchAssets } = useAssets()
const { fetchManufacturers, fetchDevices } = useDevices()
const { fetchRacks } = useDCIM()

const showOnboarding = ref(true)
const step1Done = ref(false)
const step2Done = ref(false)
const step3Done = ref(false)
const step4Done = ref(false)

const checkOnboardingProgress = async () => {
  if (!isAuthenticated.value) return

  try {
    const { data: brands } = await fetchManufacturers()
    step1Done.value = brands.value && brands.value.length > 0

    const { data: models } = await fetchDevices()
    step2Done.value = models.value && models.value.length > 0

    const { data: racks } = await fetchRacks()
    step3Done.value = racks.value && racks.value.length > 0

    const { data: assets } = await fetchAssets({ limit: 1 })
    step4Done.value = assets.value && assets.value.length > 0

    if (step1Done.value && step2Done.value && step3Done.value && step4Done.value) {
      showOnboarding.value = false
    }
  } catch (err) {
    console.error('Failed to parse onboarding checklist state:', err)
  }
}

onMounted(() => {
  checkOnboardingProgress()
})

// Re-check progress whenever the route changes (simulating a global event listener)
watch(() => route.path, () => {
  if (showOnboarding.value) {
    checkOnboardingProgress()
  }
})

const dismissOnboarding = () => {
  showOnboarding.value = false
}

const isDark = computed({
  get () {
    return colorMode.value === 'dark'
  },
  set () {
    colorMode.preference = colorMode.value === 'dark' ? 'light' : 'dark'
  }
})

// GitOps YAML Spec Generator Logic
const isGitOpsDrawerOpen = ref(false)
const gitOpsKind = ref('Asset')
const gitOpsAssetCategory = ref('Server')

const generatedYamlSpec = computed(() => {
  if (gitOpsKind.value === 'Datacenter') {
    return `apiVersion: network.astrona.io/v1alpha1
kind: Datacenter
metadata:
  name: dublin-hq
spec:
  type: homelab
  country: Ireland
  city: Dublin
  properties:
    uplink_speed: "2.5 Gbps Fiber"
    public_ip: "85.12.85.112"
  floors:
    - name: "Basement Level"
      level: -1
      width: 800
      depth: 1200
  racks:
    - name: "Rack-A01"
      height_u: 42
      placement_zone: "Aisle 1"
      floor_level: -1
      x: 100
      y: 200`
  }

  // Asset configurations based on category
  if (gitOpsAssetCategory.value === 'Router') {
    return `apiVersion: network.astrona.io/v1alpha1
kind: Asset
metadata:
  name: edge-router-01
spec:
  type: Network
  status: active
  ip_address: 10.0.0.1
  description: "Primary Edge BGP Router"
  properties:
    network_subtype: "Router"
    bgp_asn: 65001
    throughput: "10 Gbps"
  interfaces:
    - name: "GigabitEthernet0/0"
      type: "fiber"
      mac_address: "AA:BB:CC:DD:EE:01"
      ip_address: "85.12.85.112"
      link_speed: "10000"
      mtu: 9000
      vlan_id: 100
  relationships:
    downlinks:
      - "core-switch-01"`
  }

  if (gitOpsAssetCategory.value === 'Container') {
    return `apiVersion: network.astrona.io/v1alpha1
kind: Asset
metadata:
  name: redis-cache-cluster
spec:
  type: Container
  status: active
  description: "Ephemeral Redis cache nodes"
  properties:
    image: "redis:7.2-alpine"
    runtime: "docker"
    ports: "6379"
  host_asset_name: "docker-worker-01" # Mounts container to host server`
  }

  // Default Server
  return `apiVersion: network.astrona.io/v1alpha1
kind: Asset
metadata:
  name: esxi-host-01
spec:
  type: Server
  status: active
  ip_address: 192.168.1.10
  description: "Hypervisor virtualization host"
  properties:
    cpu_model: "AMD EPYC 7003"
    ram_gb: "256"
    os: "VMware ESXi 8.0"
    warranty_expiry: "2029-01-01"
  rack_name: "Rack-A01"
  rack_position_u: 15`
})

const copyYaml = () => {
  navigator.clipboard.writeText(generatedYamlSpec.value)
  alert('GitOps YAML spec copied to clipboard!')
}
</script>
>
