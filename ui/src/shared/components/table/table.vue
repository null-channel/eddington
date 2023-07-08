<template src="./table.html"></template>
<style src="./table.css"></style>
<script lang="ts">
import { defineComponent, markRaw } from 'vue';

import * as _ from 'lodash'
import { TableColumns } from '@interfaces';
import PaginatorComponent from '@components/paginator/paginator.vue'
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
        selected: {
            type: Array,
            required: true
        },
        tableTitle: {
            type: String,
            required: true,
        },
        tableSubtitle: {
            type: String,
            required: false
        },
        perPage: {
            type: Number,
            required: true
        },
    },
    emits: ["update:selected"],
    components: {
        PaginatorComponent: markRaw(PaginatorComponent)
    },
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
    },
    // watch: {
    //     data: function (values) {
    //         this.dataList = _.cloneDeep(values)
    //     }
    // },
    methods: {
        allChecked(isChecked: Boolean) {

            this.$emit('update:selected', isChecked ? this.dataList : [])
        },
        elementChecked(isChecked: Boolean, row: any) {
            if (isChecked) {
                this.selected.push(row);
            } else {
                _.remove(this.selected, { id: row.id });
            }
            this.$emit('update:selected', this.selected);
        },
        isChecked(id: Number) {
            return _.some(this.selected, { id });
        }
    }

})

</script>