<style src="./header.css"></style>
<template src="./header.html"></template>
<script lang="ts">
import { $ory, injectStrict } from '@helpers';
import { useUserStore } from '@stores';
import { storeToRefs } from 'pinia';
import { defineComponent } from 'vue';
import * as _ from 'lodash';
import { useRouter } from 'vue-router';

export default defineComponent({
    name: "Header",
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
        try{
            this.ory.toSession().then(({ data }) => {
                this.user.session = data
                this.user.authenticated=true
                // If the user is logged in, we want to show a logout link!
                this.ory.createBrowserLogoutFlow().then(({ data }) => {
                    this.user.logoutUrl = data.logout_url
                })
            })
        }catch(err){
            console.log('no user is logged in')
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
            return this.user.authenticated 
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
            if (this.user.logoutUrl)
                window.$axios.get(this.user.logoutUrl, {
                    headers: {
                        'Content-Type': 'application/json',
                        'Accept': 'application/json'
                    }
                }).then(()=>{
                    this.userStore.$reset()
                    this.router.push('/')
                }
                )
        }
    }
})
</script>