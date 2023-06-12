<style src="./sign-up.css"></style>
<template src="./sign-up.html"></template>
<script lang="ts">
import { NAVBAR_BEFORE_LOGIN } from '@constants';
import { useUserStore } from '@/core/stores/user.store';
import { storeToRefs } from 'pinia';
import { defineComponent } from 'vue';
import Header from "@components/header/header.vue"
import NullCloudTitle from '@components/null-cloud-title/nullCloudTitle.vue';
export default defineComponent({
    name: "SignUpPage",
    components: {
        Header,
        NullCloudTitle
    },
    setup() {
        const userStore = useUserStore()
        const { user } = storeToRefs(userStore)
        return {
            user,
            userStore
        }
    },
    data() {
        return {
            routes: NAVBAR_BEFORE_LOGIN,
            formSchema: [
                {
                    $cmp: 'FormKit',
                    props: {
                        type:"email",
                        // 'outer-class': 'flex flex-col items-center justify-center mx-auto',
                        name: 'email',
                        id: 'email',
                        placeholder: 'Email',
                        validation: 'required',
                    },
                },
                {
                    $cmp: 'FormKit',
                    props: {
                        type: 'password',
                        name: 'password',
                        id: 'password',
                        prefixIconClass: '!border-0 !bg-gray-200',
                        'prefix-icon': 'password',
                        placeholder: '**********',
                        validation: 'required',
                        messageClass:'flex items-start'
                    },
                },
                {
                    $cmp: 'FormKit',
                    props: {
                        type: 'submit',
                        label: 'Log-in',
                        ignore:"false",
                    },
                }
            ]
        }
    },
    computed: {},
    methods: {}
})
</script>