export type ClientBindResponse = {
    mac_bind: string;
}

export type ClientResponse = {
  id: number;
  app_id: number;
  device_id: string;
  os_type: 1 | 2 | 3; // 1: windows, 2: mac, 3: linux
  hostname: string;
  username: string;
  gid: string;
  uid: string;
  os_name: string;
  os_arch: string;
  remote_ip: string;
  remote_ip_country: string;
  local_ip: string;
  port: string;
  uptime: number;
  boot_time: number;
  os: string;
  platform: string;
  platform_family: string;
  platform_version: string;
  kernel_version: string;
  kernel_arch: string;
  host_id: string;
  cpu_num: number;
  cpu_info: string;
  mem_total: number;
  disk_total: number;
  conn_id: number;
  status: 1 | 2; // 1: online, 2: offline
  last_online_time: Date;
  created_at: Date;
  updated_at: Date;
}

export type ClientSystemInfoResponse = {
  mem_available: number;
  mem_used: number;
  mem_used_percent: number;
  mem_free: number;
  cpu_percent: number;
  disk_usage: null;
  bytesSent: number;
  bytesRecv: number;
  created_at: Date;
}