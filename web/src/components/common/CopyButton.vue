<script setup lang="ts">
const { copy, copied } = useClipboard()

const props = defineProps<{
    text?: string
}>()

const handleCopy = () => {
    if (props.text) {
        copy(props.text)
    }
}
</script>

<template>
    <button class="btn btn-primary gap-1 transition-all duration-200"
        :class="copied ? 'btn-success text-white' : 'btn-ghost'" @click="handleCopy">
        <Icon :icon="copied ? 'mdi:check' : 'mdi:content-copy'" class="w-3.5 h-3.5" />
        {{ copied ? '已复制' : '复制' }}
    </button>
</template>