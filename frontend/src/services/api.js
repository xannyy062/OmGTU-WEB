import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Cars API
export const carApi = {
  getAll: () => api.get('/cars'),
  getById: (id) => api.get(`/cars/${id}`),
  create: (carData) => api.post('/cars', carData),
  update: (id, carData) => api.put(`/cars/${id}`, carData),
  delete: (id) => api.delete(`/cars/${id}`),
};

// Dealers API
export const dealerApi = {
  getAll: () => api.get('/dealers'),
  getById: (id) => api.get(`/dealers/${id}`),
  create: (dealerData) => api.post('/dealers', dealerData),
  update: (id, dealerData) => api.put(`/dealers/${id}`, dealerData),
  delete: (id) => api.delete(`/dealers/${id}`),
};

export default api;