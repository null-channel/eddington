<style src="./header.css"></style>
<template src="./header.html"></template>
<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
    name: "Header",
    data() {
        return { showMenue: false }
    },
    methods: {
        TogelMenu: function (_: any) {
            this.showMenue = !this.showMenue
        }
    }
})
</script>