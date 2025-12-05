<template>
  <div class="flex-1 overflow-y-auto">
    <ul class="space-y-2">
      <li>
        <router-link :to="{ name: 'Dashboard' }" active-class="active" class="font-medium">
          <Icon icon="mdi:data" /> 数据统计
        </router-link>
      </li>
      <li>
        <router-link :to="{ name: 'Client' }" active-class="active" class="font-medium">
          <Icon icon="material-symbols:computer" /> 主机列表
        </router-link>
      </li>
      <li>
        <router-link :to="{ name: 'UserList' }" active-class="active" class="font-medium">
          <Icon icon="mdi:playlist-play" /> 用户列表
        </router-link>
      </li>
      <li>
        <a class="font-medium">
          <Icon icon="mdi:cog-outline" /> 系统设置
        </a>
      </li>
    </ul>
  </div>
</template>
