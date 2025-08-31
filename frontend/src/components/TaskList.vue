<template>
  <div class="p-6 bg-gray-100 min-h-screen flex flex-col items-center">
    <!-- Tasks Card -->
    <div class="w-full max-w-lg bg-white rounded-2xl shadow-2xl p-6">
      <!-- Header -->
      <div class="flex justify-between items-center mb-6">
        <div class="flex items-center gap-3">
          <h2 class="text-xl font-bold text-gray-800">Tasks</h2>
        </div>
        <button @click="showTaskModal = true"
          class="px-3 py-1 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition text-sm">
          + Add
        </button>
      </div>

      <!-- Tasks List -->
      <div class="flex flex-col gap-3">
        <div v-for="task in tasks" :key="task.ID"
          class="flex justify-between items-center p-3 bg-gray-50 rounded-xl shadow-sm hover:shadow-md transition">
          <div class="flex items-center gap-3">
            <input type="checkbox" v-model="task.completed" @change="updateTask(task)"
              class="h-5 w-5 rounded border-gray-300 text-indigo-600" />
            <div>
              <span :class="{ 'line-through text-gray-400': task.completed }" class="font-semibold text-gray-800">
                {{ task.title }}
              </span>
              <div v-if="task.estimated_pomodoros" class="text-xs text-gray-500">
                🍅 {{ task.completed_pomodoros || 0 }}/{{ task.estimated_pomodoros }} pomodoros
              </div>
              <div v-if="task.due_date" class="text-xs text-gray-500">
                📅 Due: {{ formatDate(task.due_date) }}
              </div>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <span v-if="task.priority" :class="{
              'bg-red-500': task.priority === 'high',
              'bg-yellow-400': task.priority === 'medium',
              'bg-green-500': task.priority === 'low'
            }" class="text-white text-xs px-2 py-1 rounded-full font-semibold">
              {{ task.priority.toUpperCase() }}
            </span>
            <button @click="scheduleTask(task)" class="text-blue-500 hover:text-blue-700 text-sm" title="Schedule">
              📅
            </button>
            <button @click="startFocus(task)" class="text-green-500 hover:text-green-700 text-sm" title="Focus">
              🎯
            </button>
            <button @click="deleteTask(task.ID)" class="text-red-500 hover:text-red-700 font-bold">
              ✕
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Task Modal -->
    <div v-if="showTaskModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click="closeTaskModal">
      <div class="bg-white rounded-2xl p-6 w-full max-w-md" @click.stop>
        <h3 class="text-lg font-bold text-gray-800 mb-4">Add New Task</h3>

        <!-- TaskList.vue - Update the modal form -->
        <form @submit.prevent="addTask" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Title</label>
            <input v-model="newTask" type="text" required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none" />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea v-model="newDescription" rows="2"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none"></textarea>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Priority</label>
              <select v-model="newPriority"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none">
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Estimated 🍅</label>
              <input v-model.number="newEstimatedPomodoros" type="number" min="1" max="20"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none" />
            </div>
          </div>

          <!-- Updated Date and Time fields -->
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Due Date</label>
              <input v-model="newDueDate" type="date"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Time</label>
              <input v-model="newDueTime" type="time"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none" />
            </div>
          </div>

          <div class="flex gap-3 pt-4">
            <button type="submit"
              class="flex-1 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition">
              Create Task
            </button>
            <button type="button" @click="closeTaskModal"
              class="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 transition">
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const tasks = ref([]);
const showTaskModal = ref(false)
const newTask = ref("")
const newDescription = ref("")
const newPriority = ref("medium")
const newEstimatedPomodoros = ref(1)
const newDueDate = ref("")
const newDueTime = ref("09:00") // Default time

function formatDate(dateString) {
  if (!dateString) return ""
  return new Date(dateString).toLocaleDateString()
}

function closeTaskModal() {
  showTaskModal.value = false
  newTask.value = ""
  newDescription.value = ""
  newPriority.value = "medium"
  newEstimatedPomodoros.value = 1
  newDueDate.value = ""
  newDueTime.value = "09:00" // Reset to default
}

async function fetchTasks() {
  try {
    const res = await axios.get('http://localhost:8080/tasks')
    tasks.value = res.data || []
  } catch (error) {
    console.error('Error fetching tasks:', error)
  }
}

async function addTask() {
  if (!newTask.value.trim()) return
  
  try {
    // Combine date and time with timezone offset
    let dueDateTime = ""
    if (newDueDate.value && newDueTime.value) {
      // Create a date object with the selected date and time
      const date = new Date(`${newDueDate.value}T${newDueTime.value}`);
      // Convert to ISO string which includes timezone information
      dueDateTime = date.toISOString();
      console.log('Sending dueDateTime with timezone:', dueDateTime);
    }
    
    const taskData = {
      title: newTask.value,
      description: newDescription.value,
      completed: false,
      priority: newPriority.value,
      due_date: dueDateTime,
      estimated_pomodoros: newEstimatedPomodoros.value,
      completed_pomodoros: 0
    }
    
    console.log('Sending task data:', taskData);
    
    const res = await axios.post('http://localhost:8080/tasks', taskData);
    
    console.log('Task created successfully:', res.data);
    
    // Automatically schedule the task on the calendar
    if (dueDateTime) {
      try {
        // Use the correct field name (ID instead of id)
        const scheduleRes = await axios.post('http://localhost:8080/productivity/schedule/task', {
          task_id: res.data.ID, // Use ID (uppercase) as returned from the task creation
          event_date: dueDateTime,
          duration: (newEstimatedPomodoros.value || 1) * 25
        });
        
        console.log('Task scheduled successfully:', scheduleRes.data);
        
        // Update the task with the calendar event ID if needed
        if (scheduleRes.data.ID) {
          const updateRes = await axios.put(`http://localhost:8080/tasks/${res.data.ID}`, {
            ...res.data,
            calendar_event_id: scheduleRes.data.ID
          });
          console.log('Task updated with calendar event ID:', updateRes.data);
        }
      } catch (scheduleError) {
        console.error('Error scheduling task:', scheduleError);
        // Don't fail the entire operation if scheduling fails
      }
    }
    
    tasks.value.push(res.data);
    closeTaskModal();
  } catch (error) {
    console.error('Error adding task:', error);
    if (error.response) {
      console.error('Server response:', error.response.data);
    }
    alert('Failed to create task. Please try again.');
  }
}

async function updateTask(task) {
  try {
    await axios.put(`http://localhost:8080/tasks/${task.ID}`, task)
  } catch (error) {
    console.error('Error updating task:', error)
  }
}

async function deleteTask(id) {
  try {
    await axios.delete(`http://localhost:8080/tasks/${id}`)
    tasks.value = tasks.value.filter(t => t.ID !== id)
  } catch (error) {
    console.error('Error deleting task:', error)
  }
}

async function scheduleTask(task) {
  try {
    const eventDate = new Date()
    eventDate.setHours(eventDate.getHours() + 1) // Schedule 1 hour from now

    await axios.post('http://localhost:8080/productivity/schedule/task', {
      task_id: task.ID,
      event_date: eventDate.toISOString(),
      duration: (task.estimated_pomodoros || 1) * 25 // 25 min per pomodoro
    })

    alert(`Task "${task.title}" scheduled on calendar!`)
  } catch (error) {
    console.error('Error scheduling task:', error)
  }
}

async function startFocus(task) {
  try {
    await axios.post('http://localhost:8080/productivity/focus/start', {
      task_id: task.ID,
      planned_duration: (task.estimated_pomodoros || 1) * 25
    })

    alert(`Started focus session for "${task.title}"! Switch to Pomodoro tab.`)
  } catch (error) {
    console.error('Error starting focus session:', error)
  }
}

onMounted(fetchTasks)
</script>