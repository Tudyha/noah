<template>
  <div>
    <!-- Hero Section -->
    <section class="hero min-h-[80vh] bg-base-200">
      <div class="hero-content text-center">
        <div class="max-w-3xl">
          <h1
            class="text-4xl md:text-6xl font-bold mb-6 bg-linear-to-r from-primary to-secondary bg-clip-text text-transparent">
            {{ t('home.hero.title') }}
          </h1>
          <p class="text-lg md:text-xl mb-8 opacity-80">
            {{ t('home.hero.subtitle') }}
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <router-link :to="{ name: 'Dashboard' }" class="btn btn-primary btn-lg">
              {{ t('home.hero.cta') }}
            </router-link>
            <a href="#features" class="btn btn-outline btn-lg">
              {{ t('home.hero.learnMore') }}
            </a>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
</script>
