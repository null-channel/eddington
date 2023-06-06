<template src="./home.html"></template>
<script  lang="ts">
import { defineComponent } from 'vue';
import HeroSection from '@components/hero-section/heroSection.vue';
export default defineComponent({
    name: "HomePage",
    components: {
        HeroSection
    }
})
</script>
<style src="./home.css"></style>