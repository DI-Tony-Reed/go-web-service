<template>
  <div>
    <Album v-for="album in albums" :album=album></Album>
  </div>
  <RandomAlbum @randomAdded="addNewAlbum" />
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
