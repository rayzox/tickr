<template>
  <div class="p-6 bg-gray-100 min-h-screen flex flex-col items-center">
    <!-- Active Focus Session Banner -->
    <div v-if="activeFocus" class="w-full max-w-lg mb-4 bg-gradient-to-r from-indigo-500 to-purple-600 rounded-xl p-4 text-white">
      <div class="text-sm opacity-90">🎯 Focusing on:</div>
      <div class="font-bold">{{ activeFocus.task?.title }}</div>
      <div class="text-sm opacity-75">Started {{ formatTimeAgo(activeFocus.start_time) }}</div>
    </div>

    <!-- Pomodoro Card -->
    <div class="w-full max-w-lg bg-white rounded-2xl shadow-2xl p-8">
      <!-- Header -->
      <div class="text-center mb-8">
        <h2 class="text-2xl font-bold text-gray-800 mb-2">Pomodoro Timer</h2>
        <p class="text-sm text-gray-500">{{ getCurrentPhase() }}</p>
      </div>

      <!-- Task Selection for Work Sessions -->
      <div v-if="currentPhase === 'work' && !isRunning" class="mb-6">
        <label class="block text-sm font-medium text-gray-700 mb-2">Working on (optional):</label>
        <select v-model="selectedTaskId" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none">
          <option value="">No specific task</option>
          <option v-for="task in incompleteTasks" :key="task.ID" :value="task.ID">
            {{ task.title }} ({{ task.completed_pomodoros || 0 }}/{{ task.estimated_pomodoros || 1 }} 🍅)
          </option>
        </select>
      </div>

      <!-- Timer Display -->
      <div class="text-center mb-8">
        <div class="text-6xl font-mono font-bold text-gray-800 mb-4">
          {{ formatTime(timeLeft) }}
        </div>
        
        <!-- Progress Ring -->
        <div class="relative w-32 h-32 mx-auto mb-6">
          <svg class="w-32 h-32 transform -rotate-90" viewBox="0 0 100 100">
            <!-- Background circle -->
            <circle cx="50" cy="50" r="45" stroke="#e5e7eb" stroke-width="8" fill="none"/>
            <!-- Progress circle -->
            <circle cx="50" cy="50" r="45" 
                    :stroke="getPhaseColor()" 
                    stroke-width="8" 
                    fill="none"
                    stroke-linecap="round"
                    :stroke-dasharray="circumference"
                    :stroke-dashoffset="strokeDashoffset"
                    class="transition-all duration-1000"/>
          </svg>
          <div class="absolute inset-0 flex items-center justify-center">
            <span class="text-2xl">{{ getPhaseEmoji() }}</span>
          </div>
        </div>
      </div>

      <!-- Controls -->
      <div class="flex justify-center gap-4 mb-6">
        <button @click="startTimer" 
                :disabled="isRunning"
                class="px-6 py-3 bg-green-500 text-white rounded-xl hover:bg-green-600 disabled:bg-gray-300 disabled:cursor-not-allowed transition">
          Start
        </button>
        <button @click="pauseTimer" 
                :disabled="!isRunning"
                class="px-6 py-3 bg-yellow-500 text-white rounded-xl hover:bg-yellow-600 disabled:bg-gray-300 disabled:cursor-not-allowed transition">
          Pause
        </button>
        <button @click="resetTimer" 
                class="px-6 py-3 bg-red-500 text-white rounded-xl hover:bg-red-600 transition">
          Reset
        </button>
      </div>

      <!-- Phase Selector -->
      <div class="flex justify-center gap-2 mb-6">
        <button @click="setPhase('work')" 
                :class="{'bg-indigo-600 text-white': currentPhase === 'work', 'bg-gray-200 text-gray-700': currentPhase !== 'work'}"
                class="px-4 py-2 rounded-lg transition">
          Work (25m)
        </button>
        <button @click="setPhase('short')" 
                :class="{'bg-indigo-600 text-white': currentPhase === 'short', 'bg-gray-200 text-gray-700': currentPhase !== 'short'}"
                class="px-4 py-2 rounded-lg transition">
          Short (5m)
        </button>
        <button @click="setPhase('long')" 
                :class="{'bg-indigo-600 text-white': currentPhase === 'long', 'bg-gray-200 text-gray-700': currentPhase !== 'long'}"
                class="px-4 py-2 rounded-lg transition">
          Long (15m)
        </button>
      </div>

      <!-- Stats -->
      <div class="grid grid-cols-3 gap-4 text-center">
        <div class="bg-gray-50 rounded-lg p-3">
          <div class="text-2xl font-bold text-indigo-600">{{ stats.completedSessions }}</div>
          <div class="text-xs text-gray-500">Completed</div>
        </div>
        <div class="bg-gray-50 rounded-lg p-3">
          <div class="text-2xl font-bold text-green-600">{{ stats.totalMinutes }}</div>
          <div class="text-xs text-gray-500">Minutes</div>
        </div>
        <div class="bg-gray-50 rounded-lg p-3">
          <div class="text-2xl font-bold text-purple-600">{{ stats.todaySessions }}</div>
          <div class="text-xs text-gray-500">Today</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import axios from 'axios'

// Timer state
const currentPhase = ref('work')
const timeLeft = ref(25 * 60)
const isRunning = ref(false)
const timerId = ref(null)
const selectedTaskId = ref('')

// Data
const stats = ref({
  completedSessions: 0,
  totalMinutes: 0,
  todaySessions: 0
})
const incompleteTasks = ref([])
const activeFocus = ref(null)

// Phase durations in seconds
const phaseDurations = {
  work: 25 * 60,
  short: 5 * 60,
  long: 15 * 60
}

// Computed properties
const circumference = 2 * Math.PI * 45
const strokeDashoffset = computed(() => {
  const progress = timeLeft.value / phaseDurations[currentPhase.value]
  return circumference * (1 - progress)
})

// Methods
function formatTime(seconds) {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

function formatTimeAgo(dateString) {
  const diff = Date.now() - new Date(dateString).getTime()
  const minutes = Math.floor(diff / 60000)
  return `${minutes}m ago`
}

function getCurrentPhase() {
  const phases = {
    work: 'Focus Time - Stay concentrated!',
    short: 'Short Break - Relax a bit',
    long: 'Long Break - Take your time'
  }
  return phases[currentPhase.value]
}

function getPhaseColor() {
  const colors = {
    work: '#dc2626',
    short: '#16a34a',
    long: '#2563eb'
  }
  return colors[currentPhase.value]
}

function getPhaseEmoji() {
  const emojis = {
    work: '🍅',
    short: '☕',
    long: '🛋️'
  }
  return emojis[currentPhase.value]
}

function startTimer() {
  if (isRunning.value) return
  
  isRunning.value = true
  timerId.value = setInterval(() => {
    if (timeLeft.value > 0) {
      timeLeft.value--
    } else {
      completeSession()
    }
  }, 1000)
}

function pauseTimer() {
  isRunning.value = false
  if (timerId.value) {
    clearInterval(timerId.value)
    timerId.value = null
  }
}

function resetTimer() {
  pauseTimer()
  timeLeft.value = phaseDurations[currentPhase.value]
}

function setPhase(phase) {
  pauseTimer()
  currentPhase.value = phase
  timeLeft.value = phaseDurations[phase]
}

async function completeSession() {
  pauseTimer()
  
  try {
    // Save session to backend with task linking
    const sessionData = {
      phase: currentPhase.value,
      duration: phaseDurations[currentPhase.value] / 60,
      completed_at: new Date().toISOString(),
      productive: true
    }

    // Link to task if selected
    if (selectedTaskId.value && currentPhase.value === 'work') {
      sessionData.task_id = parseInt(selectedTaskId.value)
      sessionData.notes = `Worked on task during pomodoro session`
    }

    const res = await axios.post('http://localhost:8080/pomodoro/sessions', sessionData)
    
    // Update task pomodoro count if linked
    if (selectedTaskId.value && currentPhase.value === 'work') {
      await updateTaskPomodoros(selectedTaskId.value)
    }
    
    // Update local stats
    stats.value.completedSessions++
    stats.value.totalMinutes += phaseDurations[currentPhase.value] / 60
    stats.value.todaySessions++
    
    // Auto-transition to next phase
    if (currentPhase.value === 'work') {
      if (stats.value.completedSessions % 4 === 0) {
        setPhase('long')
      } else {
        setPhase('short')
      }
    } else {
      setPhase('work')
    }
    
    // Browser notification
    if (Notification.permission === 'granted') {
      new Notification('Pomodoro Complete!', {
        body: `${getCurrentPhase()} session finished!`,
        icon: '/favicon.ico'
      })
    }
    
  } catch (error) {
    console.error('Error saving session:', error)
  }
}

async function updateTaskPomodoros(taskId) {
  try {
    // Fetch current task
    const res = await axios.get(`http://localhost:8080/tasks`)
    const task = res.data.find(t => t.ID == taskId)
    if (task) {
      task.completed_pomodoros = (task.completed_pomodoros || 0) + 1
      await axios.put(`http://localhost:8080/tasks/${taskId}`, task)
    }
  } catch (error) {
    console.error('Error updating task pomodoros:', error)
  }
}

async function fetchStats() {
  try {
    const res = await axios.get('http://localhost:8080/pomodoro/stats')
    stats.value = res.data
  } catch (error) {
    console.error('Error fetching stats:', error)
  }
}

async function fetchTasks() {
  try {
    const res = await axios.get('http://localhost:8080/tasks')
    incompleteTasks.value = res.data.filter(task => !task.completed)
  } catch (error) {
    console.error('Error fetching tasks:', error)
  }
}

async function fetchActiveFocus() {
  try {
    const res = await axios.get('http://localhost:8080/productivity/focus/active')
    activeFocus.value = res.data
  } catch (error) {
    console.error('Error fetching active focus:', error)
  }
}

// Request notification permission
onMounted(async () => {
  if (Notification.permission === 'default') {
    await Notification.requestPermission()
  }
  await fetchStats()
  await fetchTasks()
  await fetchActiveFocus()
})

onUnmounted(() => {
  if (timerId.value) {
    clearInterval(timerId.value)
  }
})
</script>