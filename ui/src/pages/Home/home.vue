<template src="./home.html"></template>
<script  lang="ts">
import { defineComponent } from 'vue';
import Header from "@components/header/header.vue"
import HeroSection from '@components/hero-section/heroSection.vue';
import { NAVBAR_BEFORE_LOGIN } from '@/core/constants';
export default defineComponent({
    name: "HomePage",
    components: {
        HeroSection,
        Header,
    },
    data() {
        return {
            routes: NAVBAR_BEFORE_LOGIN
        }
    }
})
</script>
<style src="./home.css"></style>