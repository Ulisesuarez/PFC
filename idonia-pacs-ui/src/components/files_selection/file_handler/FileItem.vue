<template>
  <v-card color="grey lighten-3 file-item" >
    <v-card-title class="pb-0 title-wrapper">
      <div class="name" v-text="file.name" />
      <span class="size grey--text text--darken-3" v-text="size" />
    </v-card-title>
    <v-card-text class="pa-0 text-xs-center file-icon">
      <v-icon color="grey darken-3" v-text="icon" />
    </v-card-text>
    <v-card-actions class="pt-0" >
      <v-tooltip v-if="!isFolder" top>
        <v-btn
          v-if="!file.cantBeReport"
          slot="activator"
          :class="{ 'active grey darken-2': isReport }"
          icon
          @click="$emit('update:is-report', !isReport)"
        >
          <v-icon
            :color="isReport ? 'white' : 'grey darken-2'"
            v-text="'attach_file'"
          />
        </v-btn>
        <span>{{ $tname('report_button_text') }}</span>
      </v-tooltip>
      <v-spacer />
      <v-btn icon @click="$emit('delete')">
        <v-icon color="grey darken-2">delete</v-icon>
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import { humanReadableFileSize } from '@/assets/js/files'

export default {
  name: 'FileItem',
  props: {
    file: {
      type: File,
      required: true
    },
    isReport: {
      type: Boolean,
      required: true
    }
  },
  computed: {
    isFolder () {
      return !this.file.type && this.file.size !== 0 && this.file.size % 4096 === 0
    },
    size () {
      return humanReadableFileSize(this.file.size)
    },
    icon () {
      if (this.isFolder) {
        return 'folder'
      }

      const type = this.file.type.split('/')[0]
      switch (type) {
        case 'text':
          return 'insert_drive_file'
        case 'image':
          return 'image'
        case 'audio':
          return 'audiotrack'
        case 'video':
          return 'play_circle_filled'
        case 'application':
          return 'settings_applications'
        default:
          return 'insert_drive_file'
      }
    }
  }
}
</script>

<style lang="scss" scoped>
  .title-wrapper {
    position: relative;

    .name {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      max-width: 100%;
    }

    .size {
      font-size: .8rem;
      opacity: .9;
      position: absolute;
      top: 0;
      right: 0;
      margin: 3px 5px;
    }
  }

  .file-icon {
    .v-icon {
      font-size: 4rem;
    }
  }
</style>
