<template>
  <div class="container">
    <filter-group :group="group"></filter-group>
  </div>
</template>

<script>
  import FilterGroup from '@/components/FilterGroup'

  export default {
    components: { FilterGroup },
    data () {
      return {
        group: {
          name: 'Alerts',
          combine: 'and',
          rules: [
            {
              f: 'ends',
              field: 'to',
              arg: '@stuifzand.com'
            },
            {
              f: 'begins',
              field: 'subject',
              arg: 'Alert:'
            }
          ],
          actions: [
            {
              f: 'move_to_folder',
              arg: 'Webwinkels'
            }
          ]
        }
      }
    }
  }
</script>

<style>
.container {
  width: 1024px;
  margin: 0 auto;
}
</style>
