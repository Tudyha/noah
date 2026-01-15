<template>
  <div class="overflow-y-auto">
    <ul class="menu bg-base-100 rounded-box w-full gap-1 p-0">
      <li>
        <router-link :to="{ name: 'Dashboard' }" active-class="bg-primary text-primary-content shadow-md" class="group">
          <Icon icon="mdi:home-outline" class="w-5 h-5 group-hover:scale-110 transition-transform" />
          <span class="font-medium">仪表盘</span>
        </router-link>
      </li>
      <li>
        <router-link :to="{ name: 'Agent' }" active-class="bg-primary text-primary-content shadow-md" class="group">
          <Icon icon="material-symbols:computer-outline" class="w-5 h-5 group-hover:scale-110 transition-transform" />
          <span class="font-medium">主机列表</span>
        </router-link>
      </li>
      <li>
        <router-link :to="{ name: 'UserList' }" active-class="bg-primary text-primary-content shadow-md" class="group">
          <Icon icon="mdi:playlist-play" class="w-5 h-5 group-hover:scale-110 transition-transform" />
          <span class="font-medium">用户列表</span>
        </router-link>
      </li>
      <li>
        <a class="group">
          <Icon icon="mdi:cog-outline" class="w-5 h-5 group-hover:scale-110 transition-transform" />
          <span class="font-medium">系统设置</span>
        </a>
      </li>
    </ul>
  </div>
</template>
