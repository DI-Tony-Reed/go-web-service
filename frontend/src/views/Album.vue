<template>
  <div class="card">
    ID: {{ album.ID }}
    <br>
    Title: {{ album.Title }}
    <br>
    Artist: {{ album.Artist }}
    <br>
    Price: {{ album.Price }}
  </div>
</template>

<script>
export default {
  name: "Album",
  props: {
    album: Object,
  },
}
</script>

<style scoped>
.card {
  background-color: white;
  border: 2px solid black;
  padding: 15px;
  margin: 0 0 15px 0;
}
</style>