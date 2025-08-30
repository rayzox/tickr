<template>
  <div class="p-6 bg-gray-100 min-h-screen flex flex-col items-center">
    <!-- Tasks Card -->
    <div class="w-full max-w-lg bg-white rounded-2xl shadow-2xl p-6">
      <!-- Header -->
      <div class="flex justify-between items-center mb-6">
        <div class="flex items-center gap-3">
          <h2 class="text-xl font-bold text-gray-800">Tasks</h2>
        </div>
      </div>

      <!-- Add Task Form -->
      <form @submit.prevent="addTask" class="flex gap-2 mb-5">
        <input v-model="newTask" type="text" placeholder="New task..." 
               class="flex-1 px-4 py-2 rounded-xl border border-gray-300 focus:ring-2 focus:ring-indigo-400 focus:outline-none"/>
        <select v-model="newPriority" class="px-4 py-2 rounded-xl border border-gray-300 focus:ring-2 focus:ring-indigo-400 focus:outline-none">
          <option value="low">Low</option>
          <option value="medium">Medium</option>
          <option value="high">High</option>
        </select>
        <button type="submit" class="px-4 py-2 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 transition">
          Add
        </button>
      </form>

      <!-- Tasks List -->
      <div class="flex flex-col gap-3">
        <div v-for="task in tasks" :key="task.id" class="flex justify-between items-center p-3 bg-gray-50 rounded-xl shadow-sm hover:shadow-md transition">
          <div class="flex items-center gap-3">
            <input type="checkbox" v-model="task.completed" @change="updateTask(task)" 
                   class="h-5 w-5 rounded border-gray-300 text-indigo-600"/>
            <span :class="{'line-through text-gray-400': task.completed}" class="font-semibold text-gray-800">
              {{ task.title }}
            </span>
          </div>
          <div class="flex items-center gap-2">
            <span v-if="task.priority" 
                  :class="{
                    'bg-red-500': task.priority === 'high',
                    'bg-yellow-400': task.priority === 'medium',
                    'bg-green-500': task.priority === 'low'
                  }" 
                  class="text-white text-xs px-2 py-1 rounded-full font-semibold">
              {{ task.priority.toUpperCase() }}
            </span>
            <button @click="deleteTask(task.id)" class="text-red-500 hover:text-red-700 font-bold">
              ✕
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const tasks = ref([])
const newTask = ref("")
const newPriority = ref("low")

async function fetchTasks() {
  const res = await axios.get('http://localhost:8080/tasks')
  tasks.value = res.data
}

async function addTask() {
  if (!newTask.value.trim()) return
  const res = await axios.post('http://localhost:8080/tasks', {
    title: newTask.value,
    completed: false,
    priority: newPriority.value
  })
  tasks.value.push(res.data)
  newTask.value = ""
  newPriority.value = "low"
}

async function updateTask(task) {
  await axios.put(`http://localhost:8080/tasks/${task.id}`, task)
}

async function deleteTask(id) {
  await axios.delete(`http://localhost:8080/tasks/${id}`)
  tasks.value = tasks.value.filter(t => t.id !== id)
}

onMounted(fetchTasks)
</script>
