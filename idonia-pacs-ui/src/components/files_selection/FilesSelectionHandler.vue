<template>
  <Card
    :value="expanded"
    icon="attach_file"
    @input="expanded = !expanded" >
    <template slot="title">
      {{ $tname('title') }}
    </template>
    <template slot="content">
      <FileHandler />
    </template>
  </Card>
</template>

<script>
import FileHandler from '@/components/files_selection/file_handler/FileHandler'
import Card from '@/components/common/Card'

export default {
  name: 'FilesSelectionHandler',
  components: { Card, FileHandler },
  props: {
    expand: {
      type: Boolean,
      default: true
    }
  },
  data () {
    return {
      expanded: this.expand
    }
  }
}
</script>

<style scoped>

</style>
