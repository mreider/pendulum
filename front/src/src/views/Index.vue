<template>
	<div class="index-page">
		<div class="container-fluid">
			<div class="row heading">
				<div class="col-12">
					<logo></logo> <b>{{ path }}</b>

					<div class="actions">
						<button class="btn btn-sm btn-info" @click="create">New file</button>
					</div>
				</div>
			</div>

			<div class="alert alert-danger" v-for="error in errors">{{error.message}}</div>

			<table class="table">
				<thead>
					<th>
						<i class="mdi mdi-folder-open"></i>
						Name
						<a class="back" v-if="path.length > 1" @click="back"><i class="mdi mdi-keyboard-return"></i></a>
					</th>
					<th class="tar">Last modified</th>
				</thead>
				<tbody>
					<tr v-for="item in response.files">
						<td v-if="item.type === 'dir'"><router-link :to="item.path"><i class="mdi mdi-folder-outline"></i> {{ item.name }}/</router-link></td>
						<td v-else><router-link :to="'/edit'+item.path"><i class="mdi mdi-file-document"></i> {{ item.name }}</router-link></td>
						<td class="tar">{{ item.last_modified }}</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</template>

<script>
import axios from 'axios'

export default {
  data () {
    return {
      path: this.$route.path,
      errors: [],
      response: {
        files: []
      }
    }
  },
  beforeRouteLeave (to, from, next) {
    if (to.meta.componentName === 'Index') {
      this.loadContents(to.path, next)
    } else {
      next()
    }
  },
  beforeRouteUpdate (to, from, next) {
    this.loadContents(to.path, next)
  },
  created () {
    this.loadContents(this.$route.path)
  },
  methods: {
    back () {
      this.$router.go(-1)
    },
    create () {
      var filename = prompt('Enter new filename:')
      if (filename.length > 0) {
        this.$router.push('/edit' + this.path + filename)
      }
    },
    loadContents (path, callback) {
      this.path = path
      var params = {}
      var link = '/api/list' + path
      return axios
        .get(link, params)
        .then(response => {
          console.log(response.data)
          if ('error' in response.data) {
            this.errors = [ response.data.error ]
          } else {
            this.response = response.data.response
          }
          if (typeof callback === 'function') {
            callback()
          }
        })
        .catch(err => {
          this.response = {
            files: []
          }
          this.errors = [
            { message: err }
          ]
          if (typeof callback === 'function') {
            callback()
          }
        })
    }
  }
}
</script>

<style lang="scss">
.index-page {
	.container-fluid {
		padding-top: 1em;
		padding-bottom: 1em;
	}
	.heading {
		padding-bottom: 1em;
	}
	.back {
		cursor: pointer;
		margin-left: 1em;
		background: #eee;
		border: 1px solid #ddd;
		border-radius: 3px;
		padding: 1px 3px;
	}
	.actions {
		float: right;
	}
}
</style>