<style src="./header.css"></style>
<template src="./header.html"></template>
<script lang="ts">
import { useUserStore } from '@/core/stores/user.store';
import { storeToRefs } from 'pinia';
import { defineComponent } from 'vue';
import { useCookies } from 'vue3-cookies';
export default defineComponent({
    name: "Header",
    setup() {
        const userStore = useUserStore()
        const { user } = storeToRefs(userStore)
        const { cookies } = useCookies();
        return {
            cookies,
            user,
            userStore
        }
    },
    props: {
        routes: null
    },
    data() {
        return {
            showMenu: false,
            toggled: false
        }
    },
    computed: {
        isLogged() {
            return this.cookies.get('user-token') == null
        }
    },
    methods: {
        togelMenu: function (_: any) {
            if (this.cookies.get('user-token') == null)
                this.showMenu = !this.showMenu
            else
                this.toggled = !this.toggled

        }
    }
})
</script>