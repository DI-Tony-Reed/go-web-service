<template>
	<div>
		<div class="mb-2">
			<form>
				<label>
					Artist
				</label>
				<input v-model="artistSearch"/>
				
				<button @click.prevent="searchArtist">Search</button>
				<button @click.prevent="getDefaultAlbums">Reset</button>
			</form>
		</div>

		<Album
			v-for="album in store.albums"
			v-bind:key="album.id"
			:album=album
			@deleteAlbum="deleteAlbum"
		></Album>
		
		<RandomAlbum @randomAdded="addNewAlbum"/>
	</div>
</template>

<script>
import Request from "../helpers/Request";
import Album from "./Album.vue"
import RandomAlbum from "./RandomAlbum.vue"
import {store} from '@/store'

export default {
  components: {RandomAlbum, Album},
  data() {
    return {
      artistSearch: "",
      store,
    }
  },
  methods: {
    addNewAlbum(album) {
      this.store.albums.push(album)
    },
    searchArtist() {
      if (!this.artistSearch) {
        return
      }
			
      this.store.albums = []
			
      let request = new Request(`albums/artist/${encodeURIComponent(this.artistSearch)}`)
      request.get().then(data => {
        if (typeof data.errors === "undefined") {
          this.store.albums = data
        }
      })
    },
    deleteAlbum(id) {
      let request = new Request(`albums/${id}`)
      request.delete().then(() => {
        this.store.albums = this.store.albums.filter(album => album.id !== id)
      })
    },
    getDefaultAlbums() {
      this.artistSearch = ""
      let request = new Request("albums")
      request.get().then(data => this.store.albums = data)
    },
  },
  created() {
    this.getDefaultAlbums()
  },
}
</script>

<style scoped>

</style>