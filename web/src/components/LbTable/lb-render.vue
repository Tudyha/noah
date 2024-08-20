<script>
export default {
  name: 'LbRender',
  functional: true,
  props: {
    scope: {
      type: Object,
      default: () => {
      }
    },
    render: {
      type: Function,
      default: () => {
      }
    }
  },
  render: (h, ctx) => {
    return ctx.props.render ? ctx.props.render(h, ctx.props.scope) : ''
  }
}
</script>
