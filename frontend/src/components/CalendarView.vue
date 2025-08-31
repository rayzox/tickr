<template>
  <div class="p-6 bg-gray-100 min-h-screen">
    <!-- Calendar Card -->
    <div class="w-full max-w-6xl mx-auto bg-white rounded-2xl shadow-2xl p-6">
      <!-- Header -->
      <div class="flex justify-between items-center mb-6">
        <div class="flex items-center gap-4">
          <h2 class="text-2xl font-bold text-gray-800">Calendar</h2>
          <div class="flex gap-2">
            <button @click="viewMode = 'month'" 
                    :class="{'bg-indigo-600 text-white': viewMode === 'month', 'bg-gray-200 text-gray-700': viewMode !== 'month'}"
                    class="px-4 py-2 rounded-lg transition">
              Month
            </button>
            <button @click="viewMode = 'week'" 
                    :class="{'bg-indigo-600 text-white': viewMode === 'week', 'bg-gray-200 text-gray-700': viewMode !== 'week'}"
                    class="px-4 py-2 rounded-lg transition">
              Week
            </button>
          </div>
        </div>
        
        <div class="flex items-center gap-4">
          <!-- Navigation -->
          <div class="flex items-center gap-2">
            <button @click="previousPeriod" class="p-2 rounded-lg hover:bg-gray-100 transition">
              ←
            </button>
            <h3 class="text-lg font-semibold text-gray-700 min-w-[200px] text-center">
              {{ getCurrentPeriodTitle() }}
            </h3>
            <button @click="nextPeriod" class="p-2 rounded-lg hover:bg-gray-100 transition">
              →
            </button>
          </div>
          
          <!-- Add Event Button -->
          <button @click="showEventModal = true" 
                  class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition">
            + Add Event
          </button>
        </div>
      </div>

      <!-- Month View -->
      <div v-if="viewMode === 'month'" class="grid grid-cols-7 gap-1">
        <!-- Day Headers -->
        <div v-for="day in ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']" :key="day" 
             class="p-3 text-center font-semibold text-gray-600 bg-gray-50 rounded-lg">
          {{ day }}
        </div>
        
        <!-- Calendar Days -->
        <div v-for="day in monthDays" :key="day.date.getTime()" 
             :class="{
               'bg-gray-100 text-gray-400': !day.isCurrentMonth,
               'bg-blue-50 border-2 border-blue-300': isToday(day.date),
               'hover:bg-gray-50': day.isCurrentMonth
             }"
             class="min-h-[120px] p-2 border rounded-lg cursor-pointer transition"
             @click="selectDate(day.date)">
          <div class="font-semibold text-sm mb-1">{{ day.date.getDate() }}</div>
          <div class="space-y-1">
            <div v-for="event in getEventsForDate(day.date)" :key="event.ID" 
                 :class="getEventTypeColor(event.event_type || 'custom')"
                 class="text-xs p-1 rounded text-white truncate cursor-pointer hover:opacity-80"
                 @click.stop="editEvent(event)">
              {{ getEventIcon(event.event_type) }} {{ event.title }}
            </div>
          </div>
        </div>
      </div>

      <!-- Week View -->
      <div v-if="viewMode === 'week'" class="grid grid-cols-8 gap-1">
        <!-- Time Column -->
        <div class="p-2">
          <div class="h-12"></div> <!-- Header spacer -->
          <div v-for="hour in hours" :key="hour" class="h-16 text-xs text-gray-500 text-right pr-2 border-b">
            {{ hour.toString().padStart(2, '0') }}:00
          </div>
        </div>
        
        <!-- Day Columns -->
        <div v-for="day in weekDays" :key="day.date.getTime()" class="border-l">
          <!-- Day Header -->
          <div :class="{'bg-blue-50 text-blue-700': isToday(day.date)}" 
               class="p-3 text-center font-semibold border-b">
            <div class="text-xs text-gray-500">{{ formatDayName(day.date) }}</div>
            <div class="text-lg">{{ day.date.getDate() }}</div>
          </div>
          
          <!-- Hour Slots -->
          <div v-for="hour in hours" :key="hour" 
               class="h-16 border-b border-gray-100 cursor-pointer hover:bg-gray-50 transition relative"
               @click="selectDateTime(day.date, hour)">
            <div v-for="event in getEventsForDateHour(day.date, hour)" :key="event.ID"
                 :class="getEventTypeColor(event.event_type || 'custom')"
                 class="absolute inset-x-0 mx-1 text-xs p-1 rounded text-white truncate cursor-pointer hover:opacity-80"
                 @click.stop="editEvent(event)">
              {{ getEventIcon(event.event_type) }} {{ event.title }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Event Modal -->
    <div v-if="showEventModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="closeModal">
      <div class="bg-white rounded-2xl p-6 w-full max-w-md" @click.stop>
        <h3 class="text-lg font-bold text-gray-800 mb-4">
          {{ editingEvent ? 'Edit Event' : 'Add Event' }}
        </h3>
        
        <form @submit.prevent="saveEvent" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Title</label>
            <input v-model="eventForm.title" type="text" required
                   class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none"/>
          </div>
          
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Date</label>
              <input v-model="eventForm.date" type="date" required
                     class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none"/>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Time</label>
              <input v-model="eventForm.time" type="time" required
                     class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none"/>
            </div>
          </div>
          
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
              <select v-model="eventForm.event_type" 
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none">
                <option value="custom">Custom</option>
                <option value="task">Task</option>
                <option value="habit">Habit</option>
                <option value="pomodoro">Pomodoro</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Priority</label>
              <select v-model="eventForm.priority" 
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none">
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
              </select>
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Duration (minutes)</label>
            <input v-model.number="eventForm.duration" type="number" min="15" max="480"
                   class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none"/>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea v-model="eventForm.description" rows="3"
                      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none"></textarea>
          </div>
          
          <div class="flex gap-3 pt-4">
            <button type="submit" class="flex-1 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition">
              {{ editingEvent ? 'Update' : 'Create' }}
            </button>
            <button type="button" @click="closeModal" 
                    class="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 transition">
              Cancel
            </button>
            <button v-if="editingEvent" type="button" @click="deleteEvent" 
                    class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition">
              Delete
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'

// State
const viewMode = ref('month')
const currentDate = ref(new Date())
const events = ref([])
const showEventModal = ref(false)
const editingEvent = ref(null)

const eventForm = ref({
  title: '',
  date: '',
  time: '',
  event_type: 'custom',
  priority: 'medium',
  duration: 60,
  description: ''
})

const hours = Array.from({ length: 24 }, (_, i) => i)

// Computed properties
const monthDays = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const startDate = new Date(firstDay)
  startDate.setDate(startDate.getDate() - firstDay.getDay())
  
  const days = []
  const current = new Date(startDate)
  
  for (let i = 0; i < 42; i++) {
    days.push({
      date: new Date(current),
      isCurrentMonth: current.getMonth() === month
    })
    current.setDate(current.getDate() + 1)
  }
  
  return days
})

const weekDays = computed(() => {
  const startOfWeek = new Date(currentDate.value)
  startOfWeek.setDate(currentDate.value.getDate() - currentDate.value.getDay())
  
  const days = []
  for (let i = 0; i < 7; i++) {
    const day = new Date(startOfWeek)
    day.setDate(startOfWeek.getDate() + i)
    days.push({ date: day })
  }
  
  return days
})

// Methods
function getCurrentPeriodTitle() {
  if (viewMode.value === 'month') {
    return currentDate.value.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
  } else {
    const startOfWeek = new Date(currentDate.value)
    startOfWeek.setDate(currentDate.value.getDate() - currentDate.value.getDay())
    const endOfWeek = new Date(startOfWeek)
    endOfWeek.setDate(startOfWeek.getDate() + 6)
    
    return `${startOfWeek.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })} - ${endOfWeek.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}`
  }
}

function previousPeriod() {
  if (viewMode.value === 'month') {
    currentDate.value.setMonth(currentDate.value.getMonth() - 1)
  } else {
    currentDate.value.setDate(currentDate.value.getDate() - 7)
  }
  currentDate.value = new Date(currentDate.value)
}

function nextPeriod() {
  if (viewMode.value === 'month') {
    currentDate.value.setMonth(currentDate.value.getMonth() + 1)
  } else {
    currentDate.value.setDate(currentDate.value.getDate() + 7)
  }
  currentDate.value = new Date(currentDate.value)
}

function isToday(date) {
  const today = new Date()
  return date.toDateString() === today.toDateString()
}

function formatDayName(date) {
  return date.toLocaleDateString('en-US', { weekday: 'short' })
}

function getEventsForDate(date) {
  return events.value.filter(event => {
    if (!event.event_date) return false;
    
    const eventDate = new Date(event.event_date);
    return eventDate.toDateString() === date.toDateString();
  });
}

function getEventsForDateHour(date, hour) {
  return events.value.filter(event => {
    const eventDate = new Date(event.event_date || event.date)
    const eventHour = eventDate.getHours()
    return eventDate.toDateString() === date.toDateString() && eventHour === hour
  })
}

function getEventTypeColor(type) {
  const colors = {
    task: 'bg-blue-500',
    habit: 'bg-green-500',
    pomodoro: 'bg-red-500',
    custom: 'bg-purple-500'
  }
  return colors[type] || 'bg-gray-500'
}

function getEventIcon(type) {
  const icons = {
    task: '📋',
    habit: '🎯',
    pomodoro: '🍅',
    custom: '📅'
  }
  return icons[type] || '📅'
}

function selectDate(date) {
  const dateString = date.toISOString().split('T')[0]
  eventForm.value.date = dateString
  eventForm.value.time = '09:00'
  showEventModal.value = true
}

function selectDateTime(date, hour) {
  const dateString = date.toISOString().split('T')[0]
  const timeString = hour.toString().padStart(2, '0') + ':00'
  eventForm.value.date = dateString
  eventForm.value.time = timeString
  showEventModal.value = true
}

function editEvent(event) {
  editingEvent.value = event
  const eventDate = new Date(event.event_date || event.date)
  eventForm.value = {
    title: event.title,
    description: event.description || '',
    date: eventDate.toISOString().split('T')[0],
    time: eventDate.toTimeString().split(' ')[0].substring(0, 5),
    event_type: event.event_type || 'custom',
    priority: event.priority || 'medium',
    duration: event.duration || 60
  }
  showEventModal.value = true
}

function closeModal() {
  showEventModal.value = false
  editingEvent.value = null
  eventForm.value = {
    title: '',
    date: '',
    time: '',
    event_type: 'custom',
    priority: 'medium',
    duration: 60,
    description: ''
  }
}

async function saveEvent() {
  try {
    const eventDateTime = new Date(`${eventForm.value.date}T${eventForm.value.time}:00`)
    
    const eventData = {
      title: eventForm.value.title,
      description: eventForm.value.description,
      date: eventForm.value.date, // Keep string format for compatibility
      event_date: eventDateTime.toISOString(),
      event_type: eventForm.value.event_type,
      priority: eventForm.value.priority,
      duration: eventForm.value.duration,
      all_day: false
    }

    if (editingEvent.value) {
      // Update existing event
      await axios.put(`http://localhost:8080/calendar/events/${editingEvent.value.ID}`, eventData)
      const index = events.value.findIndex(e => e.ID === editingEvent.value.ID)
      if (index !== -1) {
        events.value[index] = { ...editingEvent.value, ...eventData }
      }
    } else {
      // Create new event
      const res = await axios.post('http://localhost:8080/calendar/events', eventData)
      events.value.push(res.data)
    }
    
    closeModal()
  } catch (error) {
    console.error('Error saving event:', error)
  }
}

async function deleteEvent() {
  if (!editingEvent.value) return
  
  try {
    await axios.delete(`http://localhost:8080/calendar/events/${editingEvent.value.ID}`)
    events.value = events.value.filter(e => e.ID !== editingEvent.value.ID)
    closeModal()
  } catch (error) {
    console.error('Error deleting event:', error)
  }
}

async function fetchEvents() {
  try {
    const res = await axios.get('http://localhost:8080/calendar/events');
    events.value = res.data || [];
    console.log('Fetched events:', events.value); // Debug log
  } catch (error) {
    console.error('Error fetching events:', error);
  }
}

onMounted(() => {
  fetchEvents()
})
</script>