<template src="./home.html"></template>
<script  lang="ts">
import { defineComponent } from 'vue';
import Header from "@components/header/header.vue"
import HeroSection from '@components/hero-section/heroSection.vue';
import { NAVBAR_AFTER_LOGIN, NAVBAR_BEFORE_LOGIN } from '@/core/constants';
import { markRaw } from 'vue';
import { $ory, injectStrict } from '@/core/helpers';
import { useUserStore } from '@/core/stores';
import { storeToRefs } from 'pinia';
import { useRouter } from 'vue-router';
export default defineComponent({
    name: "HomePage",
    components: {
        HeroSection: markRaw(HeroSection),
        Header: markRaw(Header),
    },
    setup() {
        const ory = injectStrict($ory);
        const userStore = useUserStore()
        const { user } = storeToRefs(userStore)
        const router = useRouter();
        return {
            user,
            userStore,
            router,
            ory
        }
    },
    mounted() {
        try {
            this.ory.toSession().then(({ data }) => {
                this.user.session = data
                this.user.authenticated = true
                // If the user is logged in, we want to show a logout link!
                this.ory.createBrowserLogoutFlow().then(({ data }) => {
                    this.user.logoutUrl = data.logout_url
                })
            })
        } catch (err) {
            console.log('no user is logged in')
        }
    },
    computed: {
        isLogged() {
            return this.user.authenticated
        },
        routes: function () {
            return this.isLogged ? NAVBAR_AFTER_LOGIN : NAVBAR_BEFORE_LOGIN
        }
    }
})
</script>
<style src="./home.css"></style>