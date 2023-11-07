<template>
  <div>
  </div>
</template>

<script>
import Request from "../helpers/Request";

export default {
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
