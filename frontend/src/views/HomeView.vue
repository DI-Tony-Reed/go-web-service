<template>
  <div class="wrapper">
    <Album v-for="album in albums" :album=album></Album>

    <RandomAlbum @randomAdded="addNewAlbum" />
  </div>
</template>

<script>
import Request from "../helpers/Request";
import Album from './album.vue'
import RandomAlbum from './RandomAlbum.vue'

export default {
  components: {RandomAlbum, Album},
  data() {
    return {
      albums: [],
    }
  },
  methods: {
    addNewAlbum(album) {
      this.albums.push(album)
    }
  },
  created() {
    let request = new Request('albums')
    request.get().then(data => this.albums = data)
  },
}
</script>

<style scoped>
  .wrapper {
    padding: 1rem
  }
</style>