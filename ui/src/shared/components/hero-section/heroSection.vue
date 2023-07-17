<template src="./heroSection.html"></template>
<style src="./heroSection.css"></style>
<script lang="ts">
import { defineComponent } from 'vue';
import { ChevronDoubleRightIcon } from '@heroicons/vue/24/outline';
import NullCloudTitle from '@components/null-cloud-title/nullCloudTitle.vue';
import { env } from '@constants';

export default defineComponent({
    name: "HeroSection",
    components: {
        ChevronDoubleRightIcon,
        NullCloudTitle
    },
    data() {
        return {
            subForm:{}
        }
    },
    methods: {
        submit() {
            this.$formkit.submit('sub')
        },
        sendMail() {
            window.$axios.post(`${env.BACKEND_BASE_URL.toString()}api/v1/marketing/email`, this.subForm,
                {
                    headers: { 'Content-Type': 'application/x-www-form-urlencoded', }
                }
            ).then(()=>{
                console.log("email added")
            }).catch((err)=>{console.error(err);
            })
        }
    }

})
</script>