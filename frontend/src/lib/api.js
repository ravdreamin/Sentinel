import axios from 'axios';

const api = axios.create({
    // Uses Vite proxy in development, configure for production as needed
    baseURL: '',
});

api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export default api;
