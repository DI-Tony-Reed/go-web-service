<template>
	<div class="error">
		{{ error }}
	</div>
</template>

<script>
export default {
  name: "Error",
  props: {
    error: String
  }
}
</script>

<style scoped>
.error {
	background-color: indianred;
	border: 2px solid red;
}
</style>