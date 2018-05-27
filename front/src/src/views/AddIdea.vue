<template>
</template>
<script>
  import axios from 'axios'

  export default {
    data () {
      return {
      }
    },
    created () {
      var ideaTitle = prompt('Enter a new idea title:')
      if (ideaTitle.length === 0) {
        this.$router.go(-1)
        return
      }
      var link = '/api/addIdea?idea=' + encodeURIComponent(ideaTitle)
      var params = {}
      return axios
        .get(link, params)
        .then(response => {
          console.log(response)
          this.$router.replace({ path: '/edit' + response.data.response.path, query: { response: JSON.stringify(response.data.response) } })
        })
        .catch(err => {
          console.log(err)
          this.$router.go(-1)
        })
    }
  }
</script>

<style lang="scss">
</style>
