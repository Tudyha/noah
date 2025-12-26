<template>
  <div class="overflow-y-auto">
    <ul class="menu bg-base-100 rounded-box w-full gap-1 p-0">
      <li>
        <router-link :to="{ name: 'Dashboard' }" active-class="active" class="group">
          <Icon icon="mdi:data" class="w-5 h-5 group-hover:scale-110 transition-transform" />
          <span class="font-medium">数据统计</span>
        </router-link>
      </li>
      <li>
        <router-link :to="{ name: 'Client' }" active-class="active" class="group">
          <Icon icon="material-symbols:computer" class="w-5 h-5 group-hover:scale-110 transition-transform" />
          <span class="font-medium">主机列表</span>
        </router-link>
      </li>
      <li>
        <router-link :to="{ name: 'UserList' }" active-class="active" class="group">
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
