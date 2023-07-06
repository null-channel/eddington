<template src="./container.html"></template>
<style src="./container.css"></style>
<script lang="ts">
import { defineComponent } from "vue";
import SideBar from '@components/side-bar/sideBar.vue'
import { markRaw } from "vue";
import Header from "@/shared/components/header/header.vue";
import { NAVBAR_AFTER_LOGIN } from "@/core/constants";
import Table from "../../shared/components/table/table.vue";
export default defineComponent({
    name: "ContainerPage",
    components: {
        SideBar: markRaw(SideBar),
        Header: markRaw(Header),
        Table:markRaw(Table)
    },
    data() {
        return {
            content: {},
            selected: "Dashboard",
            routes:NAVBAR_AFTER_LOGIN,
            containerColum: [
                { label: 'App name', propertyKey: "name", isFiltredBy: false },
                { label: 'Date & time', propertyKey: "urgencyLevel", isFiltredBy: false },
                { label: 'Access URL', propertyKey: "qty", isFiltredBy: false },
                { label: 'Status', propertyKey: "status", isFiltredBy: true },
            ]
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