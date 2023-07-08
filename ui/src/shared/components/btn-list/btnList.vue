<template src="./btnList.html"></template>
<script lang="ts">
import { defineComponent } from 'vue'
import { FormButtons } from '@interfaces'

export default defineComponent({
    name: "BtnListComponent",
    props: {
        buttons: {
            type: Array<FormButtons>,
            required: true
        },
        validForm: {
            type: Boolean,
            required: false,
            default: false,
        }
    },
    emits: ['clicked'],
    methods: {
        triggerEvent: function (buttonClicked: FormButtons) {
            this.$emit('clicked', buttonClicked)
        },
        mapButtonForValidity: function () {
            this.buttons.map(btn => {
                if (btn.disableInvalidForm) {
                    btn.disabled = !this.validForm
                }
                return btn
            })
        }
    },
    mounted: function () {
        this.mapButtonForValidity()
    },
    watch: {
        validForm: function () {
            this.mapButtonForValidity()
        }
    }
})
</script>