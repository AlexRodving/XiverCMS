<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Sidebar -->
    <aside class="fixed inset-y-0 left-0 w-64 bg-white border-r border-gray-200">
      <div class="flex flex-col h-full">
        <!-- Logo -->
        <div class="flex items-center justify-center h-16 border-b border-gray-200">
          <h1 class="text-xl font-bold text-gray-900">XiverCMS</h1>
        </div>

        <!-- Navigation -->
        <nav class="flex-1 px-4 py-4 space-y-1">
          <router-link
            to="/"
            class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors"
            :class="$route.path === '/' ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-100'"
          >
            <HomeIcon class="w-5 h-5 mr-3" />
            {{ $t('navigation.dashboard') }}
          </router-link>

          <router-link
            to="/content-types"
            class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors"
            :class="$route.path.startsWith('/content-types') ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-100'"
          >
            <DocumentIcon class="w-5 h-5 mr-3" />
            {{ $t('navigation.contentTypes') }}
          </router-link>

          <router-link
            v-if="authStore.isSuperAdmin()"
            to="/users"
            class="flex items-center px-4 py-2 text-sm font-medium rounded-lg transition-colors"
            :class="$route.path.startsWith('/users') ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-100'"
          >
            <UsersIcon class="w-5 h-5 mr-3" />
            {{ $t('navigation.users') }}
          </router-link>
        </nav>

        <!-- User section -->
        <div class="p-4 border-t border-gray-200 space-y-3">
          <!-- Language switcher -->
          <div class="flex justify-center">
            <LanguageSwitcher />
          </div>
          
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <div class="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-medium">
                {{ authStore.user?.firstName?.[0] || authStore.user?.email?.[0] || 'U' }}
              </div>
              <div class="ml-3">
                <p class="text-sm font-medium text-gray-900">
                  {{ authStore.user?.firstName || authStore.user?.email }}
                </p>
                <p class="text-xs text-gray-500">{{ authStore.user?.email }}</p>
              </div>
            </div>
            <button
              @click="authStore.logout()"
              class="p-2 text-gray-400 hover:text-gray-600"
              :title="$t('auth.logout')"
            >
              <ArrowRightOnRectangleIcon class="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>
    </aside>

    <!-- Main content -->
    <div class="pl-64">
      <main class="p-8">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth'
import LanguageSwitcher from '../components/LanguageSwitcher.vue'
import {
  HomeIcon,
  DocumentIcon,
  UsersIcon,
  ArrowRightOnRectangleIcon,
} from '@heroicons/vue/24/outline'

const authStore = useAuthStore()
</script>

