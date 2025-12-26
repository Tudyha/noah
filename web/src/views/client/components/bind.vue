<script setup lang="ts">
import { getClientBind } from '@/api/client'
import type { ClientBindResponse } from '@/types'
import { clientOsTypeIconMap } from '@/map'
import { useClipboard } from '@vueuse/core'

const data = ref<ClientBindResponse>()
const { copy, copied } = useClipboard()

onMounted(() => {
    getClientBind().then(res => {
        data.value = res
    })
})

const handleCopy = (text?: string) => {
    if (text) {
        copy(text)
    }
}

</script>

<template>
    <div class="tabs tabs-lift">
        <input type="radio" name="bind_tabs" class="tab" :aria-label="'Windows'" />
        <div class="tab-content bg-base-100 border-base-300 rounded-box p-6">
            <div class="flex flex-col gap-4">
                <div class="flex items-center justify-between">
                    <span class="text-sm font-medium flex items-center gap-2">
                        <Icon :icon="clientOsTypeIconMap[1]" class="w-5 h-5 text-blue-500" />
                        Windows 绑定命令
                    </span>
                    <button class="btn btn-xs btn-ghost gap-1" disabled>
                        <Icon icon="mdi:content-copy" class="w-3.5 h-3.5" /> 复制
                    </button>
                </div>
                <div class="mockup-code before:hidden px-4">
                    <pre class="overflow-x-auto"><code>暂不支持</code></pre>
                </div>
            </div>
        </div>

        <input type="radio" name="bind_tabs" class="tab" :aria-label="'Linux'" />
        <div class="tab-content bg-base-100 border-base-300 rounded-box p-6">
            <div class="flex flex-col gap-4">
                <div class="flex items-center justify-between">
                    <span class="text-sm font-medium flex items-center gap-2">
                        <Icon :icon="clientOsTypeIconMap[3]" class="w-5 h-5 text-orange-500" />
                        Linux 绑定命令
                    </span>
                    <button class="btn btn-xs btn-ghost gap-1" disabled>
                        <Icon icon="mdi:content-copy" class="w-3.5 h-3.5" /> 复制
                    </button>
                </div>
                <div class="mockup-code before:hidden px-4">
                    <pre class="overflow-x-auto"><code>暂不支持</code></pre>
                </div>
            </div>
        </div>

        <input type="radio" name="bind_tabs" class="tab" :aria-label="'macOS'" checked />
        <div class="tab-content bg-base-100 border-base-300 rounded-box p-6">
            <div class="flex flex-col gap-4">
                <div class="flex items-center justify-between">
                    <span class="text-sm font-medium flex items-center gap-2">
                        <Icon :icon="clientOsTypeIconMap[2]" class="w-5 h-5 text-base-content" />
                        macOS 绑定命令
                    </span>
                    <button
                        class="btn btn-xs gap-1 transition-all duration-200"
                        :class="copied ? 'btn-success text-white' : 'btn-ghost hover:bg-base-200'"
                        @click="handleCopy(data?.mac_bind)"
                    >
                        <Icon :icon="copied ? 'mdi:check' : 'mdi:content-copy'" class="w-3.5 h-3.5" />
                        {{ copied ? '已复制' : '复制' }}
                    </button>
                </div>
                <div class="mockup-code">
                    <pre class="overflow-x-auto"><code class="text-xs">{{ data?.mac_bind }}</code></pre>
                </div>
            </div>
        </div>
    </div>
</template>
