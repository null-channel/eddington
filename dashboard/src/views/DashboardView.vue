<template>
    <div v-if="loading" class="loading">Loading...</div>

    <div v-if="error" class="error">{{ error }}</div>

    <div v-if="post" class="content">
      <p>{{}}</p>
    </div>
</template>

<script setup>
import { useAuth } from 'vue-clerk'
import http from '../../src/axios/axios'
import { ref, onMounted, inject } from 'vue';

const { getToken, userId, isLoaded, isSignedIn } = useAuth();
const loading = ref(true);
const error = ref(null);
const post = ref(null);

const $http = inject('$http');

onMounted(async () => {
  try {
    const token = await getToken.value(); 
    const config = {
      headers: {
        'Authorization': 'Bearer ' + token
      }
    };
    const response = await $http.get('http://localhost:9000/api/v1/users/id', config);
    post.value = response.data;
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
});

const data = ref(post.value);


</script>

<style>
@media (min-width: 1024px) {
    .about {
        min-height: 100vh;
        display: flex;
        align-items: center;
    }
}
</style>
