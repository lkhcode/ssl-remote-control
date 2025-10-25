<script setup lang="ts">
import { ref, computed, inject, onMounted, onBeforeUnmount } from 'vue';
import { useQuasar } from 'quasar';
import ErrorMessageBar from './components/ErrorMessageBar.vue';
import StatusBar from './components/StatusBar.vue';
import MatchStateToolbar from './components/MatchStateToolbar.vue';
import MatchInfo from './views/MatchInfo.vue';
import { ApiController } from './services/ApiController';
import { useUiStateStore } from "@/store/uiState/index";

const $q = useQuasar();
const uiStore = useUiStateStore();
const canSubstituteRobot = ref(false);
const api = inject<ApiController>('api');
const backgroundColor = computed(() => {
  if (canSubstituteRobot.value) return 'rgb(249,159,50)';
  return $q.dark.isActive ? 'var(--q-dark-page)' : 'rgba(56,58,58,0.86)';
});

api?.RegisterStateConsumer((s) => {
  canSubstituteRobot.value = s.canSubstituteRobot || false;
});

const rightDrawerOpen = ref(false);
const rightDrawerWidth = ref(360);

let resizing = false;
let lastX = 0;
const onGripMouseDown = (ev: MouseEvent) => {
  ev.preventDefault();
  resizing = true;
  lastX = ev.clientX;
};
const onDocumentMouseMove = (ev: MouseEvent) => {
  if (!resizing) return;
  const dx = ev.clientX - lastX;
  const newW = Math.max(160, Math.min(1000, rightDrawerWidth.value - dx));
  rightDrawerWidth.value = newW;
  lastX = ev.clientX;
};
const onDocumentMouseUp = () => {
  resizing = false;
};

onMounted(() => {
  document.addEventListener('mousemove', onDocumentMouseMove);
  document.addEventListener('mouseup', onDocumentMouseUp);
});
onBeforeUnmount(() => {
  document.removeEventListener('mousemove', onDocumentMouseMove);
  document.removeEventListener('mouseup', onDocumentMouseUp);
});

const darkMode = computed(() => $q.dark.isActive);
const toggleDarkMode = () => {
  $q.dark.toggle();
  uiStore.darkMode = $q.dark.isActive;
};

if (uiStore.darkMode !== undefined) {
  $q.dark.set(uiStore.darkMode);
}

const showShortcuts = computed(() => uiStore.showShortcuts);
const toggleShortcuts = () => {
  uiStore.showShortcuts = !uiStore.showShortcuts;
};

let initialDrawerWidth = 0;
const resizeRightDrawer = (ev: any) => {
  if (ev.isFirst === true) {
    initialDrawerWidth = uiStore.rightDrawerWidth;
  }
  uiStore.rightDrawerWidth = initialDrawerWidth - ev.offset.x;
};
const resizeLeftDrawer = (ev: any) => {
  if (ev.isFirst === true) {
    initialDrawerWidth = uiStore.leftDrawerWidth;
  }
  uiStore.leftDrawerWidth = initialDrawerWidth + ev.offset.x;
};

const dev = computed(() => {
  return import.meta.env.DEV;
});
</script>

<template>
  <div id="app">
    <header class="app-header">
      <div class="left-controls">
        <div class="title">SSL Remote Control</div>
      </div>
      <div class="header-center"></div>
      <div class="right-controls">
        <StatusBar />
        <button class="icon-btn" @click="rightDrawerOpen = !rightDrawerOpen">☰</button>
      </div>
    </header>

    <q-drawer
      v-model="rightDrawerOpen"
      side="right"
      bordered
      :width="rightDrawerWidth"
      class="bg-dark"
    >
      <div class="grip" @mousedown.prevent="onGripMouseDown"></div>
      <div class="drawer-content">
      </div>
    </q-drawer>

    <main class="page-area" :style="{ marginRight: rightDrawerOpen ? rightDrawerWidth + 'px' : '0' }">
      <div class="content-wrap">
        <MatchInfo />
      </div>
    </main>

    <footer class="app-footer">
      <div class="match-toolbar-wrap">
        <MatchStateToolbar />
      </div>
      <div class="right-footer-items">
        <StatusBar />
        <ErrorMessageBar />
      </div>
    </footer>
  </div>
</template>

<style scoped>
:root {
  --header-h: 64px;
}

.app-header {
  height: var(--header-h);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  background: var(--q-primary);
  color: #fff;
  box-sizing: border-box;
}
.left-controls,
.right-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}
.title {
  font-weight: 600;
  font-size: 1.05rem;
}
.icon-btn {
  background: transparent;
  border: 0;
  color: #fff;
  font-size: 1.2rem;
  padding: 8px;
  border-radius: 50%;
  cursor: pointer;
}
.icon-btn:hover {
  background: rgba(255, 255, 255, 0.06);
}

.q-drawer {
  z-index: 40;
}
.drawer .grip {
  position: absolute;
  left: -8px;
  top: 0;
  bottom: 0;
  width: 8px;
  cursor: ew-resize;
}
.drawer-content {
  padding: 12px;
}
.drawer-placeholder {
  color: rgba(255, 255, 255, 0.7);
}

.page-area {
  flex: 1 1 auto;
  padding: 0;
  box-sizing: border-box;
  transition: margin 0.24s cubic-bezier(0.2, 0.8, 0.2, 1);
}
.content-wrap {
  width: 100%;
  height: 100%;
  margin: 0;
}

.app-footer {
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 16px;
  box-sizing: border-box;
  background: var(--q-primary);
  color: #fff;
  border-top: 1px solid rgba(0, 0, 0, 0.12);
}
.match-toolbar-wrap {
  flex: 1;
}
.right-footer-items {
  display: flex;
  gap: 12px;
  align-items: center;
}

@media (max-width: 700px) {
  .content-wrap {
    padding: 8px;
  }
}
</style>

