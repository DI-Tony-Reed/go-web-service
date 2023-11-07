<template>
  <div>
    <Album v-for="album in albums" :album=album />
  </div>
</template>

<script>
import Request from "../helpers/Request";
import Album from './album.vue'

export default {
  components: {Album},
	data() {
		return {
			albums: []
		}
	},
	methods: {

	},
	created() {
		let request = new Request('albums')
		request.call().then(data => this.albums = data)
	}
}
</script>
