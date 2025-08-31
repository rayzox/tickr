<template>
  <div class="p-6 bg-gray-100 min-h-screen flex flex-col items-center">
    <!-- Habits Card -->
    <div class="w-full max-w-lg bg-white rounded-2xl shadow-2xl p-6">
      <!-- Header -->
      <div class="flex justify-between items-center mb-6">
        <div class="flex items-center gap-3">
          <h2 class="text-xl font-bold text-gray-800">Habits</h2>
        </div>
        <button @click="scheduleAllHabits" class="px-3 py-1 bg-green-600 text-white rounded-lg hover:bg-green-700 transition text-sm">
          📅 Schedule All
        </button>
      </div>

      <!-- Add Habit Form -->
      <form @submit.prevent="addHabit" class="flex gap-2 mb-5">
        <input v-model="newHabit" type="text" placeholder="New habit..." 
               class="flex-1 px-4 py-2 rounded-xl border border-gray-300 focus:ring-2 focus:ring-indigo-400 focus:outline-none"/>
        <select v-model="newFrequency" class="px-4 py-2 rounded-xl border border-gray-300 focus:ring-2 focus:ring-indigo-400 focus:outline-none">
          <option value="daily">Daily</option>
          <option value="weekly">Weekly</option>
          <option value="monthly">Monthly</option>
        </select>
        <button type="submit" class="px-4 py-2 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 transition">
          Add
        </button>
      </form>

      <!-- Habits List -->
      <div class="flex flex-col gap-3">
        <div v-for="habit in habits" :key="habit.ID" class="flex justify-between items-center p-3 bg-gray-50 rounded-xl shadow-sm hover:shadow-md transition">
          <div class="flex items-center gap-3">
            <input type="checkbox" v-model="habit.completed_today" @change="updateHabit(habit)" 
                   class="h-5 w-5 rounded border-gray-300 text-indigo-600"/>
            <div>
              <span :class="{'line-through text-gray-400': habit.completed_today}" class="font-semibold text-gray-800">
                {{ habit.name }}
              </span>
              <div v-if="habit.target_time" class="text-xs text-gray-500">
                ⏰ Best time: {{ habit.target_time }}
              </div>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <span v-if="habit.frequency" 
                  :class="{
                    'bg-purple-500': habit.frequency === 'daily',
                    'bg-blue-500': habit.frequency === 'weekly',
                    'bg-green-500': habit.frequency === 'monthly'
                  }" 
                  class="text-white text-xs px-2 py-1 rounded-full font-semibold">
              {{ habit.frequency.toUpperCase() }}
            </span>
            <span class="text-xs text-gray-500 bg-gray-200 px-2 py-1 rounded-full">
              🔥 {{ habit.streak || 0 }}
            </span>
            <button @click="scheduleHabit(habit)" class="text-blue-500 hover:text-blue-700 text-sm" title="Schedule">
              📅
            </button>
            <button @click="deleteHabit(habit.ID)" class="text-red-500 hover:text-red-700 font-bold">
              ✕
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue"
import axios from "axios"

const habits = ref([])
const newHabit = ref("")
const newFrequency = ref("daily")

async function fetchHabits() {
    try {
        const res = await axios.get("http://localhost:8080/habits")
        habits.value = res.data?.map(h => ({ ...h, streak: h.streak || 0 })) || []
    } catch (error) {
        console.error("Error fetching habits:", error)
    }
}

async function addHabit() {
    if (!newHabit.value.trim()) return
    try {
        const res = await axios.post("http://localhost:8080/habits", {
            name: newHabit.value,
            frequency: newFrequency.value,
            completed_today: false,
            streak: 0,
            color: getRandomColor(),
            target_time: "09:00" // Default time
        })
        habits.value.push({ ...res.data, streak: res.data.streak || 0 })
        newHabit.value = ""
        newFrequency.value = "daily"
    } catch (error) {
        console.error("Error adding habit:", error)
    }
}

async function updateHabit(habit) {
    // Update streak and completion time
    if (habit.completed_today) {
        habit.streak = (habit.streak || 0) + 1
        habit.last_completed_at = new Date().toISOString()
    } else {
        habit.streak = 0
    }
    
    try {
        await axios.put(`http://localhost:8080/habits/${habit.ID}`, habit)
    } catch (error) {
        console.error("Error updating habit:", error)
        // Revert the change if update fails
        habit.completed_today = !habit.completed_today
        if (habit.completed_today) {
            habit.streak = Math.max(0, habit.streak - 1)
        }
    }
}

async function deleteHabit(id) {
    try {
        await axios.delete(`http://localhost:8080/habits/${id}`)
        habits.value = habits.value.filter(h => h.ID !== id)
    } catch (error) {
        console.error("Error deleting habit:", error)
    }
}

async function scheduleHabit(habit) {
    try {
        const startDate = new Date()
        await axios.post('http://localhost:8080/productivity/schedule/habit', {
            habit_id: habit.ID,
            start_date: startDate.toISOString(),
            days: 7 // Schedule for next 7 days
        })
        
        alert(`Habit "${habit.name}" scheduled for the next 7 days!`)
    } catch (error) {
        console.error('Error scheduling habit:', error)
    }
}

async function scheduleAllHabits() {
    try {
        const startDate = new Date()
        for (const habit of habits.value) {
            await axios.post('http://localhost:8080/productivity/schedule/habit', {
                habit_id: habit.ID,
                start_date: startDate.toISOString(),
                days: 7
            })
        }
        alert('All habits scheduled for the next 7 days!')
    } catch (error) {
        console.error('Error scheduling habits:', error)
    }
}

function getRandomColor() {
    const colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7', '#DDA0DD', '#98D8C8']
    return colors[Math.floor(Math.random() * colors.length)]
}

onMounted(fetchHabits)
</script>