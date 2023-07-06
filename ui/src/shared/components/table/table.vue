<template src="./table.html"></template>
<style src="./table.css">
</style>
<script lang="ts">
import { defineComponent } from 'vue';

import * as _ from 'lodash'
import {  TableColumns } from '@interfaces';

export default defineComponent({
    name: "TableComponent",
    props: {
        columns: {
            type: Array<TableColumns>,
            required: true,
        },
        data: {
            type: Array,
            required: true,
        },
        tableTitle: {
            type: String,
            required: true,
        },
        perPage: {
            type: Number,
            required: true
        },
    },
    emits: [],
    components: {},
    data() {
        return {
            page: 1,
            dataList: []
        };
    },
    computed: {
        getRowValue: function () {
            return (row: any, property: string, getFunction?: Function) => !getFunction ? _.get(row, property) : getFunction(row)
        },
        isVisible: function () {
            return (row: any, user: any, visible?: Function) => !visible ? true : visible(row, user)
        }
    },
    methods: {
                
    },
    watch: {
        data: function (values) {
            this.filtredData = _.cloneDeep(values)
        }
    }
})

</script>