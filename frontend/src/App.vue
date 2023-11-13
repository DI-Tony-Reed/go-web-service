<template>
	<div class="wrapper">
		<div v-show="store.waitingOnAjax">
			Spinner.gif :)
			<br>
		</div>
		<div>
			<div v-if="store.messages" class="mb-2">
				<ServerMessage :key="message" v-for="message in store.messages" :message="message"></ServerMessage>
			</div>
			
			<div v-if="store.errors" class="mb-2">
				<Error :key="error" v-for="error in store.errors" :error="error"></Error>
			</div>
			
			<RouterView/>
		</div>
	</div>
</template>

<script>
import {store} from './store'
import Error from './views/Error.vue'
import ServerMessage from "@/views/ServerMessage.vue";

export default {
  components: {ServerMessage, Error},
  data() {
    return {
      store
    }
  }
}
</script>

<style scoped>
.wrapper {
	padding: 1rem
}
</style>
