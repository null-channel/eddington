<script setup>
import { SignIn } from 'vue-clerk'
</script>

<template>
  <SignIn path="/sign-in" routing="path" sign-up-url="/sign-up" />
</template>
