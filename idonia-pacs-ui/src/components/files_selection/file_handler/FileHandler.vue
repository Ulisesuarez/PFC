<template>
  <div
    class="file-handler"
    @dragover.prevent.stop="isOver = true"
    @dragenter.prevent.stop="isOver = true"
    @dragleave.prevent.stop="isOver = false"
    @dragend.prevent.stop="isOver = false"
    @mouseout.prevent.stop="isOver = false"
    @drop.prevent.stop="event => {
      isOver = false
      supportDirectoryUpload? parseDropped(event): parseFiles(event)
    }"
  >
    <v-layout v-if="files.length || dirfiles.length">
      <v-flex>
        <h2>{{ $tc('common.files_added', files.length + nonstudyfiles.length, { count: files.length + nonstudyfiles.length }) | title }}</h2>
      </v-flex>
      <v-flex class="text-xs-right">
        <v-btn color="primary" @click="$refs.fileInput.click()" >
          {{ $tcommon('select_more_files') }}
        </v-btn>
        <v-btn v-if="supportDirectoryUpload" color="primary" @click="$refs.dirInput.click()" >
          {{ $tcommon('select_another_directory') }}
        </v-btn>
      </v-flex>
    </v-layout>
    <input
      ref="dirInput"
      webkitdirectory
      mozdirectory
      style="display: none"
      type="file"
      accept="*"
      multiple
      @change="parseDirectory">
    <input
      ref="fileInput"
      style="display: none"
      type="file"
      accept="*"
      multiple
      @change="parseFiles">
    <div
      :class="{
        active: files.length === 0 && dirfiles.length === 0 || isOver,
        border: isOver
      }"
      class="drop-zone">
      <v-icon>move_to_inbox</v-icon>
      <strong class="grey--text text--darken-3 drop-box_text" >
        {{ $tname('drop_box_text') }}
      </strong>
      <transition name="slide-y-transition">
        <div v-if="files.length + dirfiles.length === 0 && !isOver" class="upload-button">
          <v-btn color="primary" @click="$refs.fileInput.click()" >
            {{ $tname('upload_button_file') }}
          </v-btn>
          <v-btn v-if="supportDirectoryUpload" color="primary" @click="$refs.dirInput.click()" >
            {{ $tname('upload_button_dir') }}
          </v-btn>
        </div>
      </transition>
    </div>

    <div class="scroll-box">
      <v-container
        grid-list-md
        fluuid>
        <v-layout
          row
          wrap>
          <v-flex
            v-for="file in files"
            :key="file.uuid"
            xs6
            sm4
            md2>
            <FileItem
              :file="file"
              :is-report="typeof reports[file.uuid] !== 'undefined'"
              @delete="() => {
                $store.dispatch('payload/deleteFile',file.uuid)
                isFileReport(file, false)
              }"
              @update:is-report="isFileReport(file, $event)"
            />
          </v-flex>
          <v-flex
            v-for="file in nonstudyfiles"
            :key="file.uuid"
            xs6
            sm4
            md2>
            <FileItem
              :file="file"
              :is-report="typeof reports[file.uuid] !== 'undefined'"
              @delete="() => {
                $store.dispatch('payload/deleteDirFile',file.uuid)
                isFileReport(file, false)
              }"
              @update:is-report="isFileReport(file, $event)"
            />
          </v-flex>
        </v-layout>
      </v-container>
    </div>
    <v-layout v-if="studyfiles.length">
      <h2>{{ $tc('common.studies_added', studyfiles.length, { count: studyfiles.length }) + ' ' +
      $tc('common.studies_added_from', numberOfStudies, { count: numberOfStudies })| title }}</h2>
    </v-layout>
    <div class="scroll-box">
      <v-container
        grid-list-md
        fluuid>
        <v-layout
          row
          wrap
        >

          <v-flex
            v-for="file in studyfiles"
            :key="file.uuid"
            xs6
            sm4
            md2
            style="margin-top: 10px;">
            <FileItem
              :file="file"
              :is-report="typeof reports[file.uuid] !== 'undefined'"
              @delete="() => {
                $store.dispatch('payload/deleteDirFile',file.uuid)
                isFileReport(file, false)
              }"
              @update:is-report="isFileReport(file, $event)"
            />
          </v-flex>
          <v-flex
            v-for="index in 6"
            :key="'dummy' + index"
            xs6
            sm4
            md2
          />
        </v-layout>
      </v-container>
    </div>
    <v-dialog v-model="showDialog" persistent max-width="600px" >
      <Card>
        <template slot="title">
          {{ $tc('FileHandler.DICOM_without_study_dialog_title', dialog, { count: dialog, name:DICOMname }) }}
        </template>
        <template slot="content" >
          <div v-html="$tc('FileHandler.DICOM_without_study_dialog_content', dialog, { count: dialog, name:DICOMname })"/>
        </template>
        <template slot="actions">
          <v-spacer />
          <v-btn
            color="grey"
            class="white--text"
            @click="dialog= 0; putInStudy= false">
            {{ $tcommon('no') }}
          </v-btn>
          <v-btn
            color="primary"
            @click="dialog = 0; putInStudy= true">
            {{ $tcommon('yes') }}
          </v-btn>
        </template>
      </Card>
    </v-dialog>
  </div>
</template>

<script>
import { fileParser, droppedFilesParser } from '@/assets/js/files'
import FileItem from './FileItem'
import uuidv4 from 'uuid/v4'
import Card from '@/components/common/Card'

export default {
  name: 'FileHandler',
  components: { FileItem, Card },
  data () {
    return {
      isOver: false,
      files: this.$store.state.payload.data.files,
      dirfiles: this.$store.state.payload.data.dirfiles,
      dialog: 0,
      putInStudy: undefined,
      DICOMname: '',
      droppedFiles: [],
      supportDirectoryUpload:
        typeof window.DataTransferItem !== 'undefined' &&
        typeof window.DataTransferItem.prototype.webkitGetAsEntry !== 'undefined',
      reports: {},
      BANNED_EXTENSIONS: ['exe', 'dll', 'log', 'inf', 'html', 'htm', 'cab', 'appimage', 'dat', 'db', 'ocx', 'vbs']
    }
  },
  computed: {
    studyfiles () {
      return this.dirfiles.filter(file => file.isDicom)
    },
    nonstudyfiles () {
      return this.dirfiles.filter(file => !file.isDicom)
    },
    numberOfStudies () {
      let paths = []
      for (let study in this.studyfiles) {
        if (!paths.includes(this.studyfiles[study].directory)) {
          paths.push(this.studyfiles[study].directory)
        }
      }
      return paths.length
    },
    showDialog () {
      return this.dialog > 0
    }
  },
  watch: {
    putInStudy (val) {
      if (typeof val !== 'undefined') {
        this.$emit('include', val)
      }
    },
    droppedFiles: {
      deep: true,
      handler (val) {
        for (let array in val) {
          this.parseDirectory(null, val[array])
        }
      }
    }
  },
  beforeDestroy () {
    this.$off('delete')
  },
  updated () {
    window.dispatchEvent(new CustomEvent('resize'))
  },
  beforeMount () {
    (function () {
      if (typeof window.CustomEvent === 'function') return false

      function CustomEvent (event, params) {
        params = params || { bubbles: false, cancelable: false, detail: undefined }
        var evt = document.createEvent('CustomEvent')
        evt.initCustomEvent(event, params.bubbles, params.cancelable, params.detail)
        return evt
      }

      CustomEvent.prototype = window.Event.prototype

      window.CustomEvent = CustomEvent
    })()
  },
  methods: {
    filterHiddenFiles (files) {
      if (typeof files[0].webkitRelativePath === 'undefined') {
        return files.filter(file => file.name.substring(0, 1) !== '.')
      } else if (typeof files[0].manuallySetPath !== 'undefined') {
        return files.filter(file => file.name.substring(0, 1) !== '.' && !file.manuallySetPath.includes('/.'))
      } else {
        return files.filter(file => file.name.substring(0, 1) !== '.' && !file.webkitRelativePath.includes('/.'))
      }
    },
    markDICOM (file, isDicom) {
      file.uuid = uuidv4()
      const splittedName = file.name.split('.')
      let fileExtension = ''
      fileExtension = splittedName.length > 1 ? splittedName.slice(-1)[0] : ''
      if (!fileExtension || fileExtension.toLowerCase() === 'dcm') {
        file.isDicom = isDicom
        file.cantBeReport = true
      }
      return file
    },
    isFileReport (file, isReport) {
      if (isReport) {
        this.$set(this.reports, file.uuid, file)
      } else {
        this.$delete(this.reports, file.uuid)
      }

      this.$store.dispatch('payload/setReports', Object.keys(this.reports))
    },
    parseDirectory (event, array) {
      let newFiles
      if (array) {
        newFiles = array
      } else {
        newFiles = fileParser(event)
      }
      let filesNotInAStudy = []
      let filesInStudy = []
      newFiles = this.filterHiddenFiles(newFiles)
      let studies = newFiles.filter(file => file.name.toLowerCase() === 'dicomdir')
      if (studies.length > 0 && typeof studies[0].webkitRelativePath !== 'undefined') {
        for (let study in studies) {
          let directory = studies[study].manuallySetPath
            ? studies[study].manuallySetPath.split('/').slice(0, studies[study].manuallySetPath.split('/').length - 1)
              .join('/')
            : studies[study].webkitRelativePath.split('/').slice(0, studies[study].webkitRelativePath.split('/').length - 1)
              .join('/')
          filesInStudy.push(...newFiles
            .filter(file => file.manuallySetPath
              ? file.manuallySetPath.includes(directory)
              : file.webkitRelativePath.includes(directory))
            .map(file => {
              file.directory = directory
              return file
            }))
        }
        filesNotInAStudy = newFiles.filter(file => !filesInStudy.includes(file))
        newFiles = filesInStudy.filter(file => {
          const splittedName = file.name.split('.')
          let fileExtension = ''
          fileExtension = splittedName.length > 1 ? splittedName.slice(-1)[0] : ''
          return (!this.BANNED_EXTENSIONS.includes(fileExtension) || splittedName.length === 1) &&
            splittedName[0].toLowerCase() !== 'dicomdir'
        }).map(file => {
          return this.markDICOM(file, true)
        })
        filesNotInAStudy = filesNotInAStudy.map(file => {
          return this.markDICOM(file, false)
        })
        newFiles.push(...filesNotInAStudy)
      } else {
        newFiles.map(file => {
          return this.markDICOM(file, false)
        })
      }
      let DICOMOutOfStudy = newFiles.reduce((accumulator, current) =>
        current.cantBeReport && !current.isDicom
          ? ++accumulator
          : accumulator, 0)
      if (DICOMOutOfStudy > 0) {
        if (DICOMOutOfStudy === 1) {
          this.DICOMname = newFiles[newFiles.findIndex(file => file.cantBeReport && !file.isDicom)].name
        }
        this.dialog = DICOMOutOfStudy

        this.$on('include', function (include) {
          if (include) {
            newFiles.map(file => {
              if (file.cantBeReport && !file.isDicom) {
                file.isDicom = true
              }
              return file
            })
          }
          this.dirfiles = this.dirfiles.concat(newFiles)
          this.$store.dispatch('payload/setDirFiles', this.dirfiles)
          if (event && event.target === this.$refs.dirInput) {
            this.$refs.dirInput.value = ''
          }
          this.putInStudy = undefined
          this.$off('include')
        })
      } else {
        this.dirfiles = this.dirfiles.concat(newFiles)
        this.$store.dispatch('payload/setDirFiles', this.dirfiles)
        if (event && event.target === this.$refs.dirInput) {
          this.$refs.dirInput.value = ''
        }
      }
    },
    parseFiles (event) {
      const newFiles = fileParser(event).map(file => {
        return this.markDICOM(file, false)
      })
      let DICOMOutOfStudy = newFiles.reduce((accumulator, current) =>
        current.cantBeReport && !current.isDicom
          ? ++accumulator
          : accumulator, 0)
      if (DICOMOutOfStudy > 0) {
        if (DICOMOutOfStudy === 1) {
          this.DICOMname = newFiles[newFiles.findIndex(file => file.cantBeReport && !file.isDicom)].name
        }
        this.dialog = DICOMOutOfStudy
        this.$on('include', function (include) {
          if (include) {
            newFiles.map(file => {
              if (file.cantBeReport && !file.isDicom) {
                file.isDicom = true
              }
              return file
            })
          }
          this.dirfiles = this.dirfiles.concat(newFiles)
          this.$store.dispatch('payload/setDirFiles', this.dirfiles)
          if (event.target === this.$refs.fileInput) {
            this.$refs.fileInput.value = ''
          }
          this.putInStudy = undefined
          this.$off('include')
        })
      } else {
        this.files = this.files.concat(newFiles)
        this.$store.dispatch('payload/setFiles', this.files)

        if (event.target === this.$refs.fileInput) {
          this.$refs.fileInput.value = ''
        }
      }
    },
    parseDropped (event) {
      droppedFilesParser(event).then((files) => {
        this.droppedFiles = files
      })
    }
  }
}
</script>

<style lang="scss" scoped>
  .scroll-box {
    max-height: 400px;
    overflow-x: hidden;
    overflow-y: auto;
  }

  .file-handler {
    position: relative;
    min-height: 250px;

    .drop-zone {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      width: 100%;
      height: 100%;
      font-size: 2rem;
      display: flex;
      flex-flow: row wrap;
      align-items: center;
      justify-content: center;
      opacity: 0;
      transition: all .2s ease-in-out;
      z-index: 100;
      pointer-events: none;

      &.active {
        opacity: 1;
        .upload-button {
          pointer-events: auto;
        }
      }

      &.border {
        background-color: rgba(69, 169, 212, 0.9);
        border: #45a9d4 3px solid;
        border-radius: 12px;
        padding: 2rem;

        .v-icon,
        .drop-box_text
         {
          text-shadow: 0px 2px 12px #3480a0;
          color: white !important;
        }
      }

      .v-icon,
      .drop-box_text,
      .upload-button {
        text-align: center;
        flex-basis: 100%;
        line-height: 5rem;
      }

      .v-icon {
        line-height: 8rem;
        font-size: 10rem;
      }
    }
  }
</style>
