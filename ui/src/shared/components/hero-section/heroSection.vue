<template src="./heroSection.html"></template>
<style src="./heroSection.css"></style>
<script lang="ts">
import { defineComponent } from 'vue';
import Header from "@components/header/header.vue"
import { ChevronDoubleRightIcon } from '@heroicons/vue/24/outline';
export default defineComponent({
    name: "HeroSection",
    components: {
        Header,
        ChevronDoubleRightIcon
    }
})
</script>