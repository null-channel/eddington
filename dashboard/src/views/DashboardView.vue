<template>
    <div v-if="loading" class="loading">Loading...</div>

    <div v-if="error" class="error">{{ error }}</div>

    <div v-if="post" class="content">
      <h2>{{ post.title }}</h2>
      <p>{{ post.body }}</p>
    </div>
</template>

<script setup lang=ts>
import { useAuth } from 'vue-clerk'
const { getToken, isLoaded, isSignedIn } = useAuth();
import axios from 'axios';
import { ref, onMounted } from 'vue';

const loading = ref(true);
const error = ref(null);
const post = ref(null);

onMounted(async () => {
  try {
    const response = await axios.get('https://jsonplaceholder.typicode.com/posts/1');
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
