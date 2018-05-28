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
      var parts = this.$route.path.split('/')
      var project = parts[2].replace('.md', '')

      var storyTitle = prompt('Enter a new story title:')
      if (storyTitle == null | storyTitle.length === 0) {
        this.$router.go(-1)
        return
      }
      var link = '/api/addStory?title=' + encodeURIComponent(storyTitle) + '&project=' + encodeURIComponent(project)
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
