<style src="./header.css"></style>
<template src="./header.html"></template>
<script lang="ts">
import { injectStrict } from '@helpers';
import { useUserStore } from '@stores';
import { defineComponent } from 'vue';
import * as _ from 'lodash';
import { useRouter } from 'vue-router';

export default defineComponent({
    name: "Header",
    setup() {
        const userStore = useUserStore()
        const router = useRouter();
        return {
            userStore,
            router,
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
    methods: {
        toggleMenu: function (_: any) {
            if (this.isLogged)
                this.showMenu = !this.showMenu
        },
        toggleBurgerMenu: function () {
            this.toggled = !this.toggled
        },
        logout() {
        }
    }
})
</script>
