// src/api/http.js

import axios from 'axios';

// Create an instance
const myhttp = axios.create({
  baseURL: 'http://localhost:9000', // Your API base URL
  timeout: 1000,
});


export default myhttp;
