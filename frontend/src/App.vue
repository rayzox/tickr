<template>
  <div class="flex h-screen bg-gray-50">
    <!-- Sidebar -->
    <aside class="w-20 bg-white border-r border-gray-200 flex flex-col items-center py-4 shadow-lg">
      <button @click="tab = 'dashboard'" 
              :class="tabClass('dashboard')" 
              class="mb-6 p-2 rounded-lg hover:bg-gray-100 transition">
        <ChartBarIcon class="w-7 h-7 font-bold"/>
      </button>

      <button @click="tab = 'tasks'" 
              :class="tabClass('tasks')" 
              class="mb-6 p-2 rounded-lg hover:bg-gray-100 transition">
        <HomeIcon class="w-7 h-7 font-bold"/>
      </button>

      <button @click="tab = 'calendar'" 
              :class="tabClass('calendar')" 
              class="mb-6 p-2 rounded-lg hover:bg-gray-100 transition">
        <CalendarIcon class="w-7 h-7 font-bold"/>
      </button>

      <button @click="tab = 'pomodoro'" 
              :class="tabClass('pomodoro')" 
              class="mb-6 p-2 rounded-lg hover:bg-gray-100 transition">
        <ClockIcon class="w-7 h-7 font-bold"/>
      </button>

      <button @click="tab = 'habits'" 
              :class="tabClass('habits')" 
              class="p-2 rounded-lg hover:bg-gray-100 transition">
        <CheckCircleIcon class="w-7 h-7 font-bold"/>
      </button>
    </aside>

    <!-- Main content -->
    <main class="flex-1 p-6 overflow-auto">
      <div class="w-full h-full bg-white rounded-2xl shadow-xl p-6">
        <Dashboard v-if="tab==='dashboard'" @switch-tab="switchTab"/>
        <TaskList v-if="tab==='tasks'"/>
        <CalendarView v-if="tab==='calendar'"/>
        <Pomodoro v-if="tab==='pomodoro'"/>
        <HabitTracker v-if="tab==='habits'"/>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Dashboard from './components/Dashboard.vue'
import TaskList from './components/TaskList.vue'
import CalendarView from './components/CalendarView.vue'
import Pomodoro from './components/Pomodoro.vue'
import HabitTracker from './components/HabitTracker.vue'

import { HomeIcon, CalendarIcon, ClockIcon, CheckCircleIcon, ChartBarIcon } from '@heroicons/vue/24/outline'

const tab = ref('dashboard')

function tabClass(t) {
  return tab.value === t
    ? 'text-blue-500 shadow-md'
    : 'text-gray-400 hover:text-blue-500'
}

function switchTab(newTab) {
  tab.value = newTab
}
</script>