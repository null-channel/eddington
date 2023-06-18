<style src="./log-in.css"></style>
<template src="./log-in.html"></template>
<script lang="ts">
import { NAVBAR_BEFORE_LOGIN } from '@constants';
import { defineComponent } from 'vue';
import Header from "@components/header/header.vue"
import NullCloudTitle from '@components/null-cloud-title/nullCloudTitle.vue';
import { $ory, injectStrict, oryErrorHandler } from '@helpers';
import { useRoute, useRouter } from 'vue-router';
import { ref } from 'vue';
import { LoginFlow } from '@ory/client';
import  OryFlow  from '@components/flow/login/login-flow.vue'
export default defineComponent({
    name: "LogInPage",
    components: {
        Header,
        NullCloudTitle,
        OryFlow
    },
    setup() {
        const ory = injectStrict($ory);
        const route = useRoute();
        const router = useRouter();
        const loginFlow = ref<LoginFlow | undefined>();
        const handleGetFlowError = oryErrorHandler(router);
        return {
            ory,
            route,
            router,
            loginFlow,
            handleGetFlowError
        }
    },
    mounted() {
        const { flow, refresh, aal, returnTo } = this.route.query;

        if (typeof flow !== 'string') {
            this.loginFlowForBrowsers(Boolean(refresh), aal?.toString(), returnTo?.toString());
        } else {
            this.ory.getLoginFlow({ id: flow })
                .then((response) => {
                    this.loginFlow = response.data;
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
        loginFlowForBrowsers: function (refresh: boolean, aal?: string, returnTo?: string) {
            this.ory.createBrowserLoginFlow(
                {
                    refresh,
                    aal,
                    returnTo
                }
            )
                .then((response) => {
                    this.loginFlow = response.data;
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