<style src="./sign-up.css"></style>
<template src="./sign-up.html"></template>
<script lang="ts">
import { NAVBAR_BEFORE_LOGIN } from '@constants';
import { defineComponent } from 'vue';
import Header from "@components/header/header.vue"
import NullCloudTitle from '@components/null-cloud-title/nullCloudTitle.vue';
import { $ory, injectStrict, oryErrorHandler } from '@helpers';
import { useRoute, useRouter } from 'vue-router';
import { ref } from 'vue';
import { RegistrationFlow } from '@ory/client';
import OryFlow from '@components/flow/sign-up/sign-up.vue'
export default defineComponent({
    name: "SignUpPage",
    components: {
        Header,
        NullCloudTitle,
        OryFlow
    },
    setup() {
        const ory = injectStrict($ory);
        const route = useRoute();
        const router = useRouter();
        const registrationFlow = ref<RegistrationFlow | undefined>();
        const handleGetFlowError = oryErrorHandler(router);
        return {
            ory,
            route,
            router,
            registrationFlow,
            handleGetFlowError
        }
    },
    mounted() {
        const { flow, returnTo } = this.route.query;

        if (typeof flow !== 'string') {
            this.registrationFlowForBrowsers( returnTo?.toString());
        } else {
            this.ory.getRegistrationFlow({ id: flow })
                .then((response) => {
                    this.registrationFlow = response.data;
                })
                .catch(this.handleGetFlowError);
        }
    },
    data() {
        return {
            routes: NAVBAR_BEFORE_LOGIN,
        }
    },
    computed: {},
    methods: {
        registrationFlowForBrowsers: function (returnTo?: string) {
            this.ory.createBrowserRegistrationFlow(
                {
                    returnTo
                }
            )
                .then((response) => {
                    this.registrationFlow = response.data;
                    this.router.replace({
                        query: {
                            flow: response.data.id,
                        },
                    });
                })
                .catch(this.handleGetFlowError);
        }
    }
})
</script>