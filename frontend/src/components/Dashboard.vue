<template>
  <div class="p-6 bg-gray-100 min-h-screen">
    <!-- Dashboard Header -->
    <div class="max-w-7xl mx-auto mb-8">
      <h1 class="text-3xl font-bold text-gray-800 mb-2">Productivity Dashboard</h1>
      <p class="text-gray-600">{{ formatDate(new Date()) }}</p>
    </div>

    <!-- Quick Stats -->
    <div class="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
      <div class="bg-white rounded-2xl shadow-lg p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600">Tasks</p>
            <p class="text-2xl font-bold text-indigo-600">{{ completedTasks }}/{{ totalTasks }}</p>
          </div>
          <div class="text-3xl">✅</div>
        </div>
      </div>
      
      <div class="bg-white rounded-2xl shadow-lg p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600">Habits</p>
            <p class="text-2xl font-bold text-green-600">{{ completedHabits }}/{{ totalHabits }}</p>
          </div>
          <div class="text-3xl">🎯</div>
        </div>
      </div>
      
      <div class="bg-white rounded-2xl shadow-lg p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600">Focus Sessions</p>
            <p class="text-2xl font-bold text-red-600">{{ todayPomodoros }}</p>
          </div>
          <div class="text-3xl">🍅</div>
        </div>
      </div>
      
      <div class="bg-white rounded-2xl shadow-lg p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-600">Events</p>
            <p class="text-2xl font-bold text-purple-600">{{ totalEvents }}</p>
          </div>
          <div class="text-3xl">📅</div>
        </div>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-3 gap-6">
      
      <!-- Recent Tasks -->
      <div class="bg-white rounded-2xl shadow-lg p-6">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-bold text-gray-800">Recent Tasks</h3>
          <button @click="$emit('switchTab', 'tasks')" class="text-indigo-600 hover:text-indigo-800 text-sm font-semibold">
            View All →
          </button>
        </div>
        <div class="space-y-3">
          <div v-for="task in recentTasks" :key="task.ID" 
               class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center gap-3">
              <input type="checkbox" v-model="task.completed" @change="updateTask(task)"
                     class="h-5 w-5 rounded border-gray-300 text-indigo-600"/>
              <div>
                <span :class="{'line-through text-gray-400': task.completed}" class="font-medium">
                  {{ task.title }}
                </span>
                <div class="text-xs text-gray-500" v-if="task.priority">
                  Priority: {{ task.priority }}
                </div>
              </div>
            </div>
            <span :class="getPriorityColor(task.priority)" 
                  class="w-3 h-3 rounded-full"></span>
          </div>
          <div v-if="!recentTasks.length" class="text-center text-gray-500 py-4">
            No tasks yet
          </div>
        </div>
      </div>

      <!-- Active Habits -->
      <div class="bg-white rounded-2xl shadow-lg p-6">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-bold text-gray-800">Today's Habits</h3>
          <button @click="$emit('switchTab', 'habits')" class="text-indigo-600 hover:text-indigo-800 text-sm font-semibold">
            View All →
          </button>
        </div>
        <div class="space-y-3">
          <div v-for="habit in activeHabits" :key="habit.ID" 
               class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center gap-3">
              <input type="checkbox" v-model="habit.completed_today" @change="updateHabit(habit)"
                     class="h-5 w-5 rounded border-gray-300 text-green-600"/>
              <div>
                <span :class="{'line-through text-gray-400': habit.completed_today}" class="font-medium">
                  {{ habit.name }}
                </span>
                <div class="text-xs text-gray-500">
                  🔥{{ habit.streak || 0 }} day streak
                </div>
              </div>
            </div>
            <span class="text-xs bg-gray-200 px-2 py-1 rounded-full">
              {{ habit.frequency }}
            </span>
          </div>
          <div v-if="!activeHabits.length" class="text-center text-gray-500 py-4">
            No habits tracked
          </div>
        </div>
      </div>

      <!-- Recent Events -->
      <div class="bg-white rounded-2xl shadow-lg p-6">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-bold text-gray-800">Upcoming Events</h3>
          <button @click="$emit('switchTab', 'calendar')" class="text-indigo-600 hover:text-indigo-800 text-sm font-semibold">
            View Calendar →
          </button>
        </div>
        <div class="space-y-3">
          <div v-for="event in upcomingEvents" :key="event.ID" 
               class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
            <div class="text-lg">
              {{ getEventEmoji(event.event_type) }}
            </div>
            <div class="flex-1">
              <div class="font-medium">{{ event.title }}</div>
              <div class="text-xs text-gray-500" v-if="event.event_date">
                {{ formatEventTime(event.event_date) }}
              </div>
            </div>
          </div>
          <div v-if="!upcomingEvents.length" class="text-center text-gray-500 py-4">
            No upcoming events
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="max-w-7xl mx-auto mt-8">
      <div class="bg-white rounded-2xl shadow-lg p-6">
        <h3 class="text-lg font-bold text-gray-800 mb-4">Quick Actions</h3>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <button @click="$emit('switchTab', 'tasks')" 
                  class="flex flex-col items-center p-4 bg-indigo-50 rounded-lg hover:bg-indigo-100 transition">
            <div class="text-2xl mb-2">📋</div>
            <span class="text-sm font-semibold text-indigo-700">Add Task</span>
          </button>
          
          <button @click="$emit('switchTab', 'habits')" 
                  class="flex flex-col items-center p-4 bg-green-50 rounded-lg hover:bg-green-100 transition">
            <div class="text-2xl mb-2">🎯</div>
            <span class="text-sm font-semibold text-green-700">Track Habit</span>
          </button>
          
          <button @click="$emit('switchTab', 'pomodoro')" 
                  class="flex flex-col items-center p-4 bg-red-50 rounded-lg hover:bg-red-100 transition">
            <div class="text-2xl mb-2">🍅</div>
            <span class="text-sm font-semibold text-red-700">Start Pomodoro</span>
          </button>
          
          <button @click="$emit('switchTab', 'calendar')" 
                  class="flex flex-col items-center p-4 bg-purple-50 rounded-lg hover:bg-purple-100 transition">
            <div class="text-2xl mb-2">📅</div>
            <span class="text-sm font-semibold text-purple-700">Schedule Event</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'

// Define emits
defineEmits(['switchTab'])

// State
const tasks = ref([])
const habits = ref([])
const events = ref([])
const pomodoroStats = ref({
  completedSessions: 0,
  totalMinutes: 0,
  todaySessions: 0
})

// Computed properties
const totalTasks = computed(() => tasks.value.length)
const completedTasks = computed(() => tasks.value.filter(t => t.completed).length)
const totalHabits = computed(() => habits.value.length)
const completedHabits = computed(() => habits.value.filter(h => h.completed_today).length)
const todayPomodoros = computed(() => pomodoroStats.value.todaySessions || 0)
const totalEvents = computed(() => events.value.length)

const recentTasks = computed(() => tasks.value.slice(-5).reverse())
const activeHabits = computed(() => habits.value.slice(0, 5))
const upcomingEvents = computed(() => {
  const now = new Date()
  return events.value
    .filter(event => {
      const eventDate = new Date(event.event_date || event.date)
      return eventDate >= now
    })
    .sort((a, b) => new Date(a.event_date || a.date) - new Date(b.event_date || b.date))
    .slice(0, 5)
})

// Methods
function formatDate(date) {
  return date.toLocaleDateString('en-US', { 
    weekday: 'long', 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  })
}

function formatEventTime(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit'
  })
}

function getPriorityColor(priority) {
  const colors = {
    high: 'bg-red-500',
    medium: 'bg-yellow-400',
    low: 'bg-green-500'
  }
  return colors[priority] || 'bg-gray-400'
}

function getEventEmoji(type) {
  const emojis = {
    task: '📋',
    habit: '🎯',
    pomodoro: '🍅',
    custom: '📅'
  }
  return emojis[type] || '📅'
}

// API calls
async function fetchTasks() {
  try {
    const res = await axios.get('http://localhost:8080/tasks')
    tasks.value = res.data || []
  } catch (error) {
    console.error('Error fetching tasks:', error)
    tasks.value = []
  }
}

async function fetchHabits() {
  try {
    const res = await axios.get('http://localhost:8080/habits')
    habits.value = res.data || []
  } catch (error) {
    console.error('Error fetching habits:', error)
    habits.value = []
  }
}

async function fetchEvents() {
  try {
    const res = await axios.get('http://localhost:8080/calendar/events')
    events.value = res.data || []
  } catch (error) {
    console.error('Error fetching events:', error)
    events.value = []
  }
}

async function fetchPomodoroStats() {
  try {
    const res = await axios.get('http://localhost:8080/pomodoro/stats')
    pomodoroStats.value = res.data || { completedSessions: 0, totalMinutes: 0, todaySessions: 0 }
  } catch (error) {
    console.error('Error fetching pomodoro stats:', error)
    pomodoroStats.value = { completedSessions: 0, totalMinutes: 0, todaySessions: 0 }
  }
}

async function updateTask(task) {
  try {
    await axios.put(`http://localhost:8080/tasks/${task.ID}`, task)
  } catch (error) {
    console.error('Error updating task:', error)
    // Revert on error
    task.completed = !task.completed
  }
}

async function updateHabit(habit) {
  // Update streak logic
  if (habit.completed_today) {
    habit.streak = (habit.streak || 0) + 1
  } else {
    habit.streak = 0
  }
  
  try {
    await axios.put(`http://localhost:8080/habits/${habit.ID}`, habit)
  } catch (error) {
    console.error('Error updating habit:', error)
    // Revert on error
    habit.completed_today = !habit.completed_today
    if (habit.completed_today) {
      habit.streak = Math.max(0, habit.streak - 1)
    }
  }
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    fetchTasks(),
    fetchHabits(),
    fetchEvents(),
    fetchPomodoroStats()
  ])
})
</script>