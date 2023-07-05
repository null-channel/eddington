<template src="./sign-up.html"></template>
<script lang="ts">
import { oryErrorHandler, oryMapper } from '@helpers';
import { defineComponent } from 'vue';
import { Flow } from '@types'
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
            this.userStore.signUp((this.flow as Flow).ui.action, headers, this.formData).then((session: any) => {
                localStorage.setItem('session', session)
                this.router.push('/')
            }).catch((err: any) => {
                this.handleGetFlowError(err);
                if (this.formId) {
                    this.$formkit.setErrors(
                        this.formId,
                        err.response.data.ui.messages.map((el: any) => el.text)
                    )
                }
            })
        }
    }
})
</script>
<style src="./sign-up.css"></style>