<template>
  <div v-if="open"
  class="fixed inset-0 z-50 flex items-center justify-center bg-black/40"
  @click="onOverlayClick"
  >

    <div class="bg-white w-[1000px] h-[520px] rounded-[24px] p-12 overflow-hidden">

      <button
        @click="close"
        class="absolute top-6 right-6 text-gray-400 hover:text-black"
      >
        ✕
      </button>

      <StepDoctor v-if="step === 1" @next="next" />

      <StepPatient
        v-if="step === 2"
        @next="next"
        @back="back"
      />

      <StepDate
        v-if="step === 3"
        @next="next"
        @back="back"
      />

      <StepPhone
        v-if="step === 4"
        @next="next"
        @back="back"
      />

      <StepSuccess
        v-if="step === 5"
        @close="close"
      />

    </div>

  </div>
</template>

<script setup>
import { ref } from "vue"

import StepDoctor from "@/components/appointment/steps/StepDoctor.vue"
import StepPatient from "@/components/appointment/steps/StepPatient.vue"
import StepDate from '@/components/appointment/steps/StepDate.vue'
import StepPhone from "@/components/appointment/steps/StepPhone.vue"
import StepSuccess from "@/components/appointment/steps/StepSuccess.vue"

const props = defineProps({
  open: Boolean
})

const emit = defineEmits(["close"])

const step = ref(1)

function next() {
  step.value++
}

function back() {
  step.value--
}

function close() {
  step.value = 1
  emit("close")
}

function onOverlayClick(event) {
  // Закрываем только если кликнули именно по затемнению, а не по содержимому модалки
  if (event.target === event.currentTarget) {
    close()
  }
}
</script>