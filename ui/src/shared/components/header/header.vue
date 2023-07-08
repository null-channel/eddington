<style src="./header.css"></style>
<template src="./header.html"></template>
<script lang="ts">
import { $ory, injectStrict } from '@helpers';
import { useUserStore } from '@stores';
import { defineComponent } from 'vue';
import * as _ from 'lodash';
import { useRouter } from 'vue-router';

export default defineComponent({
    name: "Header",
    setup() {
        const ory = injectStrict($ory);
        const userStore = useUserStore()
        const router = useRouter();
        return {
            userStore,
            router,
            ory,
        }
    },
    props: {
        routes: null,
        hideLogo: Boolean
    },
    data() {
        return {
            showMenu: false,
            toggled: false
        }
    },
    computed: {
        isLogged() {
            return !!localStorage.getItem("session")
        }
    },
    methods: {
        toggleMenu: function (_: any) {
            if (this.isLogged)
                this.showMenu = !this.showMenu
        },
        toggleBurgerMenu: function () {
            this.toggled = !this.toggled
        },
        logout() {
            this.ory.toSession().then((_: any) => {
                this.ory.createBrowserLogoutFlow().then(({ data }: any) => {
                    this.userStore.logout(data.logout_url)
                });
            });
        }
    }
})
</script>