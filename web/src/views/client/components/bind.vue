<script setup lang="ts">
import { getClientBind } from '@/api/client'
import type { ClientBindResponse } from '@/types'
import { clientOsTypeIconMap } from '@/map'

const data = ref<ClientBindResponse>()

onMounted(() => {
    getClientBind().then(res => {
        data.value = res
    })
})

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
                    <CopyButton class="btn-xs" :text="data?.mac_bind" />
                </div>
                <div class="mockup-code">
                    <pre class="overflow-x-auto"><code class="text-xs">{{ data?.mac_bind }}</code></pre>
                </div>
            </div>
        </div>
    </div>
</template>
