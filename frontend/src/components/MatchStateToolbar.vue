<script setup lang="ts">
import {ref, onMounted} from 'vue'

const now = ref(new Date().toLocaleTimeString())
onMounted(() => {
  setInterval(() => { now.value = new Date().toLocaleTimeString() }, 1000)
})
</script>

<template>
  <div class="match-toolbar">
    <div>Match toolbar — time: {{ now }}</div>
  </div>
</template>

<style scoped>
.match-toolbar { color: white; font-weight: 600 }
</style>
