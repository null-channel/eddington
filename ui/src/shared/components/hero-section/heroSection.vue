<template src="./heroSection.html"></template>
<style src="./heroSection.css"></style>
<script lang="ts">
import { defineComponent } from 'vue';
import { ChevronDoubleRightIcon } from '@heroicons/vue/24/outline';
import NullCloudTitle from '@components/null-cloud-title/nullCloudTitle.vue';
export default defineComponent({
    name: "HeroSection",
    components: {
        ChevronDoubleRightIcon,
        NullCloudTitle
    },

    
})
</script>