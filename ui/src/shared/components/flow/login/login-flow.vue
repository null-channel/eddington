<template src="./login-flow.html"></template>
<script lang="ts">
import { oryErrorHandler, oryMapper } from '@helpers';
import { defineComponent } from 'vue';
import { Flow } from '@/core/types'
import { getNodeId } from "@ory/integrations/ui";
import { useRoute, useRouter } from 'vue-router';
import { useUserStore } from '@stores';
export default defineComponent({
    name: "LoginForm",
    props: {
        flow: null,
        formId: String
    },
    data() {
        const formSchema = (this.flow as Flow).ui.nodes.map((node) => {
            return oryMapper(node, getNodeId(node))
        });
        const userStore = useUserStore()
        const router = useRouter();
        const route = useRoute();
        const handleGetFlowError = oryErrorHandler(router);

        return {
            formData: {},
            formSchema,
            handleGetFlowError,
            route,
            router,
            userStore,
        }
    },
    methods: {
        submitForm() {
            if (this.formId)
                this.$formkit.submit(this.formId);
        },
        submitFlow() {
            const headers = {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            };
            this.userStore.login((this.flow as Flow).ui.action, headers, this.formData).then((session: any) => {
                localStorage.setItem('session', JSON.stringify(session))
                this.router.push('/')
            }).catch((err) => {
                this.handleGetFlowError(err);
                if (this.formId) {
                    this.$formkit.setErrors(
                        this.formId, // You must use the ID of the form
                        err.response.data.ui.messages.map((el:any) => el.text) // (optional) An array of form level errors
                    )
                }
            })
        }
    }

})
</script>
<style src="./login-flow.css"></style>