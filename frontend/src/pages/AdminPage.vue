<script setup>
import { ref, onMounted } from 'vue'

const stats = ref({
  todayAppointments: 0,
  activeDoctors: 0,
  totalClients: 0
})

const appointments = ref([])
const loading = ref(false)

async function loadStats() {
  try {
    const res = await fetch('/api/admin/stats')
    stats.value = await res.json()
  } catch (err) {
    console.error(err)
  }
}

async function loadAppointments() {
  loading.value = true
  try {
    const res = await fetch('/api/admin/appointments?limit=10')
    appointments.value = await res.json()
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

function getStatusClass(status) {
  const classes = {
    'pending': 'bg-yellow-100 text-yellow-800',
    'confirmed': 'bg-green-100 text-green-800',
    'cancelled': 'bg-red-100 text-red-800',
    'completed': 'bg-blue-100 text-blue-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

function getStatusText(status) {
  const texts = {
    'pending': 'Ожидание',
    'confirmed': 'Подтверждено',
    'cancelled': 'Отменено',
    'completed': 'Завершено'
  }
  return texts[status] || status
}

onMounted(() => {
  loadStats()
  loadAppointments()
})
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <!-- Шапка -->
    <header class="bg-white shadow">
      <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <h1 class="text-xl font-bold text-primary">МедЦентр Админ</h1>
        <div class="flex items-center gap-4">
          <span class="text-sm text-gray-600">{{ new Date().toLocaleDateString('ru-RU') }}</span>
          <router-link to="/" class="text-sm text-gray-500 hover:text-primary">На сайт</router-link>
        </div>
      </div>
    </header>

    <div class="container mx-auto px-4 py-6">
      <!-- Статистика -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-gray-500 text-sm">Сегодня записей</div>
          <div class="text-3xl font-bold">{{ stats.todayAppointments }}</div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-gray-500 text-sm">Активных врачей</div>
          <div class="text-3xl font-bold">{{ stats.activeDoctors }}</div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-gray-500 text-sm">Всего клиентов</div>
          <div class="text-3xl font-bold">{{ stats.totalClients }}</div>
        </div>
      </div>

      <!-- Последние записи -->
      <div class="bg-white rounded-lg shadow">
        <div class="p-4 border-b flex justify-between items-center">
          <h2 class="font-semibold">Последние записи</h2>
          <router-link to="/admin/appointments" class="text-sm text-primary hover:underline">
            Все записи →
          </router-link>
        </div>
        
        <div class="p-4">
          <div v-if="loading" class="text-center py-4">
            Загрузка...
          </div>
          <table v-else class="w-full">
            <thead>
              <tr class="text-left text-sm text-gray-600">
                <th class="pb-2">Пациент</th>
                <th class="pb-2">Врач</th>
                <th class="pb-2">Дата</th>
                <th class="pb-2">Время</th>
                <th class="pb-2">Статус</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="apt in appointments" :key="apt.id" class="border-t">
                <td class="py-2">{{ apt.last_name }} {{ apt.first_name }}</td>
                <td class="py-2">{{ apt.doctor_name }}</td>
                <td class="py-2">{{ apt.appointment_date }}</td>
                <td class="py-2">{{ apt.appointment_time.slice(0,5) }}</td>
                <td class="py-2">
                  <span class="px-2 py-1 text-xs rounded-full" :class="getStatusClass(apt.status)">
                    {{ getStatusText(apt.status) }}
                  </span>
                </td>
              </tr>
              <tr v-if="appointments.length === 0">
                <td colspan="5" class="text-center py-4 text-gray-400">
                  Нет записей
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Быстрые ссылки -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-6">
        <router-link to="/admin/doctors" class="bg-white rounded-lg shadow p-4 hover:shadow-md transition">
          <div class="text-lg font-semibold mb-2">👨‍⚕️ Управление врачами</div>
          <p class="text-sm text-gray-600">Добавление, редактирование, расписание</p>
        </router-link>
        
        <router-link to="/admin/clients" class="bg-white rounded-lg shadow p-4 hover:shadow-md transition">
          <div class="text-lg font-semibold mb-2">👤 Управление клиентами</div>
          <p class="text-sm text-gray-600">Просмотр и редактирование данных клиентов</p>
        </router-link>
      </div>
    </div>
  </div>
</template>