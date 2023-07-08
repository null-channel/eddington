<style src="./sideBar.css"></style>
<template src="./sideBar.html"></template>
<script lang="ts">

import { defineComponent } from 'vue';
import * as _ from 'lodash';
import NullCloudTitle from '@components/null-cloud-title/nullCloudTitle.vue';
import { markRaw } from 'vue';

export default defineComponent({
    name: "SideBar",
    components: {
        NullCloudTitle: markRaw(NullCloudTitle)
    },
    props: {
        title: String
    }
})
</script>