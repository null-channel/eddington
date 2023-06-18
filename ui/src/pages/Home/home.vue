<template src="./home.html"></template>
<script  lang="ts">
import { defineComponent } from 'vue';
import Header from "@components/header/header.vue"
import HeroSection from '@components/hero-section/heroSection.vue';
import { NAVBAR_BEFORE_LOGIN } from '@/core/constants';
import { markRaw } from 'vue';
export default defineComponent({
    name: "HomePage",
    components: {
        HeroSection: markRaw(HeroSection),
        Header: markRaw(Header),
    },
    data() {
        return {
            routes: NAVBAR_BEFORE_LOGIN
        }
    }
})
</script>
<style src="./home.css"></style>