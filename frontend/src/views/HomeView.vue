<template>
  <div>
    <div class="mb-2">
      <form>
        <label>
          Artist
        </label>
        <input v-model="artistSearch" />

        <button @click.prevent="searchArtist">Search</button>
        <button @click.prevent="getDefaultAlbums">Reset</button>
      </form>
    </div>

    <Album v-for="album in albums" :album=album></Album>

    <RandomAlbum @randomAdded="addNewAlbum"/>
  </div>
</template>

<script>
import Request from "../helpers/Request";
import Album from "./album.vue"
import RandomAlbum from "./RandomAlbum.vue"
import { store } from '../store'

export default {
  components: {RandomAlbum, Album},
  data() {
    return {
      albums: [],
      artistSearch: "",
      store
    }
  },
  methods: {
    addNewAlbum(album) {
      this.albums.push(album)
    },
    searchArtist() {
      if (!this.artistSearch) {
        return
      }

      this.albums = []

      let request = new Request(`albums/artist/${ encodeURIComponent(this.artistSearch) }`)
      request.get().then(data => {
        if (typeof data.errors === "undefined") {
          this.albums = data
        }
      })
    },
    getDefaultAlbums() {
      this.artistSearch = ""
      let request = new Request("albums")
      request.get().then(data => this.albums = data)
    },
  },
  created() {
    this.getDefaultAlbums()
  },
}
</script>

<style scoped>

</style>