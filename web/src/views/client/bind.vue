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
        <label class="tab">
            <input type="radio" name="my_tabs_4" />
            <Icon :icon="clientOsTypeIconMap[1]" />
            Windows
        </label>
        <div class="tab-content bg-base-100 border-base-300">
            <div class="mockup-code">
                <pre><code>暂不支持</code></pre>
            </div>
        </div>

        <label class="tab">
            <input type="radio" name="my_tabs_4" />
            <Icon :icon="clientOsTypeIconMap[3]" />
            Linux
        </label>
        <div class="tab-content bg-base-100 border-base-300 p-6">
            <div class="mockup-code">
                <pre><code>暂不支持</code></pre>
            </div>
        </div>

        <label class="tab">
            <input type="radio" name="my_tabs_4" checked/>
            <Icon :icon="clientOsTypeIconMap[2]" />
            Mac
        </label>
        <div class="tab-content bg-base-100 border-base-300 p-6">
            在主机终端执行如下命令：
            <div class="mockup-code">
                <pre><code>{{ data?.mac_bind }}</code></pre>
            </div>
        </div>
    </div>
</template>