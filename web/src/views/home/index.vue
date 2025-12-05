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
            <router-link to="/register" class="btn btn-primary btn-lg">
              {{ t('home.hero.cta') }}
            </router-link>
            <a href="#features" class="btn btn-outline btn-lg">
              {{ t('home.hero.learnMore') }}
            </a>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="py-20">
      <div class="container-custom">
        <h2 class="text-3xl md:text-4xl font-bold text-center mb-16">
          {{ t('home.features.title') }}
        </h2>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          <div v-for="i in 4" :key="i" class="card bg-base-200 shadow-xl hover:shadow-2xl transition-all duration-300">
            <div class="card-body items-center text-center">
              <div class="w-16 h-16 rounded-full bg-primary/20 flex items-center justify-center mb-4">
                <svg class="w-8 h-8 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <h3 class="card-title">{{ t(`home.features.feature${i}.title`) }}</h3>
              <p class="opacity-80">{{ t(`home.features.feature${i}.desc`) }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-20 bg-base-200">
      <div class="container-custom text-center">
        <h2 class="text-3xl md:text-4xl font-bold mb-6">
          准备开始了吗？
        </h2>
        <p class="text-lg mb-8 opacity-80">
          立即注册，开启您的企业级应用之旅
        </p>
        <router-link to="/register" class="btn btn-primary btn-lg">
          {{ t('nav.register') }}
        </router-link>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
</script>
