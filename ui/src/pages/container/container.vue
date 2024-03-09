<template src="./container.html"></template>
<style src="./container.css"></style>
<script lang="ts">
import { defineComponent } from "vue";
import SideBar from '@components/side-bar/sideBar.vue'
import { markRaw } from "vue";
import Header from "@components/header/header.vue";
import { NAVBAR_AFTER_LOGIN } from "@/core/constants";
import Table from "@components/table/table.vue";
import ButtonListComponent from '@components/btn-list/btnList.vue'
export default defineComponent({
    name: "AppPage",
    components: {
        SideBar: markRaw(SideBar),
        Header: markRaw(Header),
        Table: markRaw(Table),
        BtnList: markRaw(ButtonListComponent)
    },
    data() {
        return {
            elements:{
                Dashboard:'Dashboard',
                Setting:'Setting'
            },
            content: {},
            currentElement: "Dashboard",
            routes: NAVBAR_AFTER_LOGIN,
            buttonsList: [
                {
                    label: "Actions",
                    type: "primary",
                    class: 'h-10'
                },
                {
                    label: "Add",
                    type: "secondary",
                    class:'h-10'
                }
            ],
            appsColumn: [
                { label: 'App name', propertyKey: "name" },
                {
                    label: 'Date & time', propertyKey: "date", type: 'date', getValueFunction: (row:any) => {
                        var date = new Date(row.date)
                        return `${date.toLocaleString('default', { month: 'short', day: 'numeric', year: 'numeric' })}`
                    }
                },
                { label: 'Access URL', propertyKey: "access_url" },
                { label: 'Status', propertyKey: "status", type: 'progress' },
            ],
            selectedList: [],
            data: [
                {
                    id: 1,
                    name: 'app1',
                    date: Date(),
                    access_url: 'https://auth.nullcloud.io',
                    status: 'completed'
                },
                {
                    id: 2,
                    name: 'app2',
                    date: Date(),
                    access_url: 'N/A',
                    status: 'cancelled'
                },
                {
                    id: 0,
                    name: 'app3',
                    date: Date(),
                    access_url: 'N/A',
                    status: 'In Progress'
                }
            ]
        }
    },
    methods: {
        isSelected: function (button: string) {
            return button === this.currentElement
        },
        selectedColor: function (button: string) {
            return this.isSelected(button) ? '#4C1D95' : '#1A202C'
        },
        setSelected: function (option: string) {
            this.currentElement = option
        }
    },
})
</script>
