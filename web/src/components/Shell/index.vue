<template>
  <div id="terminal" ref="terminal" />
</template>

<script>
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'
import { getToken } from '@/utils/auth'

export default {
  name: 'Shell',
  props: {
    id: {
      type: Number,
      default: null
    },
    shellType: {
      type: Number,
      default: 2
    }
  },
  data() {
    return {
      terminal: null,
      webSocket: null
    }
  },
  mounted() {
    this.initXterm()
  },
  methods: {
    initXterm() {
      var options = {
        cursorBlink: true,
        theme: {
          background: '#000',
          foreground: '#fff'
        },
        rows: 40
      }

      this.terminal = new Terminal(options)
      const fitAddon = new FitAddon()
      this.terminal.loadAddon(fitAddon)
      this.terminal.open(document.getElementById('terminal'))
      fitAddon.fit()

      if (this.shellType === 1) {
        this.webSocket = new WebSocket(process.env.VUE_APP_WS_ADDR + '/shell/ws/' + this.id + "?token=" + getToken())
      }
      if (this.shellType === 2) {
        this.webSocket = new WebSocket(process.env.VUE_APP_WS_ADDR + '/pty/ws/' + this.id + "?token=" + getToken())
      }

      this.webSocket.onmessage = (event) => {
        this.terminal.write(event.data)
      }

      this.terminal.onData((data) => {
        const blob = new Blob([JSON.stringify({type: 'data', data: data})], { type: 'application/json' })
        this.webSocket.send(blob)
      })

      const sendSize = () => {
        const windowSize = { type: 'size', high: this.terminal.rows, width: this.terminal.cols }
        const blob = new Blob([JSON.stringify(windowSize)], { type: 'application/json' })
        this.webSocket.send(blob)
      }

      this.webSocket.onopen = sendSize

      const resizeScreen = () => {
        fitAddon.fit()
        sendSize()
      }
      window.addEventListener('resize', resizeScreen, false)
    },
    close() {
      if (this.webSocket !== null) {
        this.webSocket.close()
      }
    }
  }
}
</script>

<style lang="scss" scoped>
html,
body,
#app {
    height: 100%;
}
</style>
