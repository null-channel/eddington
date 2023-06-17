<template src="./ory-form.html"></template>
<script lang="ts">
import { oryMapper } from '@helpers';
import { defineComponent } from 'vue';
import { Flow } from '@/core/types'
import { getNodeId } from "@ory/integrations/ui";
export default defineComponent({
    name: "OryForm",
    props: {
        flow: null,
        formId: String
    },
    data() {
        const formSchema = (this.flow as Flow).ui.nodes.map((node) => {
            return oryMapper(node, getNodeId(node))
        });
        return {
            formData: {},
            formSchema,
        }
    },
    methods: {
        submitForm() {
            if (this.formId)
                this.$formkit.submit(this.formId);
        },
        submitFlow() {
            // method contains the form method like: post & get 
            const method = (this.flow as Flow).ui.method.toLowerCase();

            // Set request headers
            const headers = {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            };

            // Make the HTTP request
            window.$axios({
                method: method,
                headers: headers,
                // Add other request parameters such as URL and data
                url: (this.flow as Flow).ui.action,
                data: { ...(this.formData) }
            }).then(( data : any)=>{
                // we need to handel each flow 
            })
        }
    }
})
</script>
<style src="./ory-form.css"></style>