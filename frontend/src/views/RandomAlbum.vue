<template>
  <div>
    <button @click="addRandom">Add random</button>
  </div>
</template>

<script>
import Request from '../helpers/Request'

export default {
  name: "RandomAlbum",
  methods: {
    addRandom() {
      let request = new Request('albums/random')
      request.put().then(data => this.$emit('randomAdded', data))
    },
  },
}
</script>

<style scoped>

</style>