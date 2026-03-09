<template>
  <Header @open-appointment="openModal"/>
  <HeroSection @open-appointment="openModal"/>
  <ServicesSection />
  <AppointmentSection @open-appointment="openModal"/>
  <AppointmentModal :open="modalOpen" @close="closeModal"/>
</template>

<script setup>
import { ref } from "vue"

import Header from '../components/layout/Header.vue'
import HeroSection from '../components/hero/HeroSection.vue'
import ServicesSection from '../components/services/ServicesSection.vue'
import AppointmentSection from "../components/appointment/AppointmentSection.vue"
import AppointmentModal from "../components/appointment/AppointmentModal.vue"

const modalOpen = ref(false)

function openModal() {
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
}

</script>