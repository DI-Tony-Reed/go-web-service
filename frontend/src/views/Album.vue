<template>
  <div class="card d-flex justify-between align-center">
    <div>
      ID: {{ album.ID }}
      <br>
      Title: {{ album.Title }}
      <br>
      Artist: {{ album.Artist }}
      <br>
      Price: {{ album.Price }}
    </div>
    <div>
      <button @click.prevent="deleteAlbum(album.ID)">Delete</button>
    </div>
  </div>
</template>

<script>
export default {
  name: "Album",
  props: {
    album: Object,
  },
  methods: {
    deleteAlbum(id) {
      if (confirm("Are you sure?")) {
        this.$emit('deleteAlbum', id)
      }
    },
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