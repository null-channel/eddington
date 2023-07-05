<template src="./container.html"></template>
<style src="./container.css"></style>
<script lang="ts">
import { defineComponent } from "vue";
import SideBar from '@components/side-bar/sideBar.vue'
import { markRaw } from "vue";
import Header from "@/shared/components/header/header.vue";
import { NAVBAR_AFTER_LOGIN } from "@/core/constants";
export default defineComponent({
    name: "ContainerPage",
    components: {
        SideBar: markRaw(SideBar),
        Header: markRaw(Header),
    },
    data() {
        return {
            content: {},
            selected: "Dashboard",
            routes:NAVBAR_AFTER_LOGIN
        }
    },
    methods: {
        isSelected: function (button: string) {
            return button === this.selected
        },
        selectedColor: function (button: string) {
            return this.isSelected(button) ? '#4C1D95' : '#1A202C'
        },
        setSelected: function (option: string) {
            console.log(option)
            this.selected = option
        }
    },
})
</script>