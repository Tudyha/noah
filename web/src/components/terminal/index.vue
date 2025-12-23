<template>
  <div ref="terminalEl" />
</template>

<script setup lang="ts">
import { Terminal } from "@xterm/xterm";
import { FitAddon } from "@xterm/addon-fit";
import "@xterm/xterm/css/xterm.css";
import { useUserStore } from "@/stores/auth";
const { VITE_WS_API_BASE_URL } = import.meta.env;

const props = defineProps<{
  id: string
}>()

// xterm DOM 引用
const terminalEl = ref();

// WebSocket 实例
let webSocket: WebSocket | null = null;
// xterm 实例
let xterm: Terminal | null = null;

// xterm 配置
const options = {
  cursorBlink: true,
  theme: {
    foreground: '#ECECEC',
    background: '#000000',
  },
  rows: 40
};

const fitAddon = new FitAddon();

// 窗口尺寸变化处理
const sendSize = () => {
  const windowSize = {
    type: "size",
    high: xterm?.rows,
    width: xterm?.cols
  };
  const blob = new Blob([JSON.stringify(windowSize)], {
    type: "application/json"
  });
  webSocket?.send(blob);
};

const resizeScreen = () => {
  fitAddon.fit();
  sendSize();
};

const initXterm = () => {
  xterm = new Terminal(options);
  xterm.loadAddon(fitAddon);

  // 挂载到 DOM
  xterm.open(terminalEl.value);
  fitAddon.fit();

  const u = useUserStore();

  // 建立 WebSocket 连接
  const wsUrl = `${VITE_WS_API_BASE_URL}/client/${props.id}/pty?Authorization=${u.token}`;

  webSocket = new WebSocket(wsUrl);

  // 接收消息并显示在终端上
  webSocket.onmessage = event => {
    xterm?.write(event.data);
  };

  // 发送输入内容到服务端
  xterm.onData(data => {
    webSocket?.send(JSON.stringify({ type: "data", data }));
  });

  webSocket.onopen = sendSize;
  window.addEventListener("resize", resizeScreen, false);
}

const closeTerminal = () => {
  console.log("close terminal");
  if (webSocket) {
    webSocket.close();
  }
  if (xterm) {
    xterm.dispose();
  }
}

defineExpose({
  close: closeTerminal
});

onMounted(() => {
  initXterm();
});

onUnmounted(() => {
  closeTerminal();
  window.removeEventListener("resize", resizeScreen, false);
});

</script>
