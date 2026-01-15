export type DashboardResponse = {
    sys_info?: {
        hostname: string;
        username: string;
        gid: string;
        uid: string;
        os_name: string;
        os_arch: string;
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
    };
    agent_stats?: {
        online: number;
        offline: number;
    };
}