<template>
  <div>
    <div v-if="store.waitingOnAjax">
      Spinner.gif :)
    </div>
    <div v-else>
      <RouterView />
    </div>
  </div>
</template>

<script>
  import store from './store'

  export default {
    computed: {
      store() {
        return store
      }
    }
  }
</script>
