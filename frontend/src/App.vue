<template>
	<div class="wrapper">
		<div v-show="store.waitingOnAjax">
			Spinner.gif :)
			<br>
		</div>
		
		<div>
			<div v-if="store.errors" class="mb-2">
				<Error v-for="error in store.errors" :error="error"></Error>
			</div>
			
			<RouterView/>
		</div>
	</div>
</template>

<script>
import {store} from './store'
import Error from './views/Error.vue'

export default {
  components: {Error},
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
