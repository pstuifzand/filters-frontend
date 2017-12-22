<template>
  <div class="container">
    <filter-group :group="group" v-for="group in filters" @save="save" :key="group.id"></filter-group>
    <button @click="add_filter" class="btn btn-primary">Add filter</button>
  </div>
</template>

<script>
  import FilterGroup from '@/components/FilterGroup'
  import axios from 'axios'

  export default {
    components: { FilterGroup },
    data () {
      return {
        filters: {}
      }
    },
    mounted () {
      axios.get('http://localhost:8088/store').then((response) => {
        this.filters = response.data.filters
      })
    },
    methods: {
      save () {
        axios.post('http://localhost:8088/store', { filters: this.filters }).then(() => {
          console.log('ok')
        }).catch(() => {
          console.log('err')
        })
      },
      add_filter () {
        this.filters.push({ name: '', combine: 'and', rules: [{}], actions: [] })
      }
    }
  }
</script>

<style>
</style>
