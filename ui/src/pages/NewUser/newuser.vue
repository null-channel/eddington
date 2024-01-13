<template src="./newuser.html"></template>
<script  lang="ts">
import { defineComponent } from 'vue';
import Header from "@components/header/header.vue"
import { NAVBAR_AFTER_LOGIN, NAVBAR_BEFORE_LOGIN } from '@/core/constants';
import { markRaw } from 'vue';
export default defineComponent({
    name: "NewUser",
    components: {
        Header: markRaw(Header),
    },
    computed: {
        routes() {
            return !!localStorage.getItem("session") ? NAVBAR_AFTER_LOGIN : NAVBAR_BEFORE_LOGIN
        }
    },
    methods: {
        sayHey() {
            console.log("HEY!!!")
        }
    },

})
</script>
<style src="./newuser.css"></style>
