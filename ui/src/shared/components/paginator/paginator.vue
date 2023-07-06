<template src="./paginator.html"></template>
<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
    name: "PaginatorComponent",
    props: {
        page: { type: Number, required: true },
        perPage: { type: Number, required: true },
        dataList: { type: Array, required: true },
        paginatedData: { type: Array, required: true },
    },
    emits: ["update:paginatedData", "update:page"],
    data: function () {
        return {
            pageNumber: 1,
        };
    },
    watch: {
        dataList: function () {
            this.$emit(
                "update:paginatedData",
                this.dataList.length
                    ? this.dataList.slice(
                        (this.page - 1) * this.perPage,
                        this.page * this.perPage
                    )
                    : []
            );
        },
    },
    mounted: function () {
        this.$emit(
            "update:paginatedData",
            this.dataList.length
                ? this.dataList.slice(
                    (this.pageNumber - 1) * this.perPage,
                    this.pageNumber * this.perPage
                )
                : []
        );
    },
    methods: {
        nextPage() {
            if (this.page !== Math.ceil(this.dataList.length / this.perPage)) {
                this.pageNumber += 1;

                this.$emit("update:page", this.pageNumber);
                this.$emit(
                    "update:paginatedData",
                    this.dataList.length
                        ? this.dataList.slice(
                            (this.pageNumber - 1) * this.perPage,
                            this.pageNumber * this.perPage
                        )
                        : []
                );
            }
        },
        backPage() {
            if (this.page !== 1) {
                this.pageNumber -= 1;
                this.$emit("update:page", this.pageNumber);
                this.$emit(
                    "update:paginatedData",
                    this.dataList.length
                        ? this.dataList.slice(
                            (this.pageNumber - 1) * this.perPage,
                            this.pageNumber * this.perPage
                        )
                        : []
                );
            }
        },
        goToPage(numPage: number) {
            this.pageNumber = numPage;

            this.$emit("update:page", this.pageNumber);
            this.$emit(
                "update:paginatedData",
                this.dataList.length
                    ? this.dataList.slice(
                        (this.pageNumber - 1) * this.perPage,
                        this.pageNumber * this.perPage
                    )
                    : []
            );
        },
    },
});
</script>