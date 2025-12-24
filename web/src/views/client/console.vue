<script setup lang="ts">
import SystemInfo from './components/system-info.vue'
import Terminal from './components/terminal.vue';

const id = useRoute().params.id + '';

const tabs = [
    {
        name: "系统资源",
        icon: "mdi:home",
        active: true,
        component: SystemInfo,
    },
    {
        name: "在线终端",
        icon: "mdi:terminal",
        component: Terminal,
    },
];
</script>

<template>
    <div class="tabs tabs-lift lg:tabs-lg">
        <template v-for="item in tabs">
            <label class="tab">
                <input type="radio" name="tab" :checked="item.active" />
                <Icon :icon="item.icon" />
                {{ item.name }}
            </label>
            <div class="tab-content bg-base-100 border-base-300 p-6">
                <KeepAlive>
                    <component :is="item.component" :id="id" />
                </KeepAlive>
            </div>
        </template>
    </div>
</template>