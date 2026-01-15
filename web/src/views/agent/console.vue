<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';
import SystemStat from './components/system-stat.vue'
import Terminal from './components/terminal.vue';
import Tunnel from './components/tunnel.vue';
import SystemInfo from './components/system-info.vue';

const route = useRoute();
const id = route.params.id + '';

const activeTab = ref(0);

const tabs = [
    {
        name: "系统资源",
        icon: "mdi:home",
        component: SystemStat,
    },
    {
        name: "在线终端",
        icon: "mdi:terminal",
        component: Terminal,
    },
    {
        name: "隧道管理",
        icon: "mdi:tunnel",
        component: Tunnel,
    },
];
const currentComponent = computed(() => tabs[activeTab.value]?.component);
</script>

<template>
    <div class="flex flex-col h-[calc(100vh-theme(spacing.16))] -m-4 sm:-m-6 overflow-hidden">
        <!-- 极简顶部栏 -->
        <SystemInfo :id="id" />

        <!-- 选项卡导航 (紧凑型) -->
        <div class="bg-base-100 border-b border-base-200 px-4 shrink-0 flex items-center gap-1">
            <button
                v-for="(item, index) in tabs"
                :key="index"
                class="px-4 py-2 text-sm font-medium transition-all relative"
                :class="activeTab === index ? 'text-primary' : 'text-base-content/50 hover:text-base-content'"
                @click="activeTab = index"
            >
                <div class="flex items-center gap-2">
                    <Icon :icon="item.icon" class="w-4 h-4" />
                    {{ item.name }}
                </div>
                <!-- 激活状态下划线 -->
                <div v-if="activeTab === index" class="absolute bottom-0 left-0 right-0 h-0.5 bg-primary"></div>
            </button>
        </div>

        <!-- 内容区域：自动填充剩余高度 -->
        <div class="flex-1 overflow-hidden p-3 bg-base-200/30">
            <div class="h-full bg-base-100 rounded-lg border border-base-200 shadow-sm overflow-hidden">
                <KeepAlive>
                    <component :is="currentComponent" :id="id" class="h-full" />
                </KeepAlive>
            </div>
        </div>
    </div>
</template>
