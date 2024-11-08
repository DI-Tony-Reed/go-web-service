<template>
	<div>
		<form>
			Artist
			<input v-model="artist"/>
			<br>
			
			Title
			<input v-model="title"/>
			<br>
			
			Price
			<input v-model="price"/>
			<br>
			
			<button @click.prevent="update">Update</button>
		</form>
		
		<br>
		<router-link to="/albums">Back to all</router-link>
	</div>
</template>

<script>
import {store} from "@/store";
import Request from "@/helpers/Request";

export default {
  name: "EditAlbumView",
  data() {
    return {
      id: 0,
      artist: "",
      title: "",
      price: 0,
      store,
    }
  },
  methods: {
    getDefaultAlbums() {
      let request = new Request(`albums/${this.$route.params.id}`)
      request.get().then(data => {
        this.assignProperties(data)
      })
    },
    update() {
      let parameters = {
        price: this.price,
        artist: this.artist,
        title: this.title
      }
			
      let urlParameters = new URLSearchParams(parameters).toString()
			
      let request = new Request(`albums/${this.id}?${urlParameters}`)
			
      request.patch(parameters).then(data => {
        this.assignProperties(data)
      })
    },
    assignProperties(data) {
      this.store.albums = data

      this.id = data[0].id
      this.artist = data[0].artist
      this.title = data[0].title
      this.price = data[0].price
    }
  },
  created() {
    this.getDefaultAlbums()
  }
}
</script>

<style scoped>

</style>