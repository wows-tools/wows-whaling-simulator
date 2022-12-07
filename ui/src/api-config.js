let backendHost;

const hostname = window && window.location && window.location.hostname;

if(hostname === 'localhost') {
  backendHost = process.env.REACT_APP_BACKEND_HOST || 'http://localhost:8080';
}

export const API_ROOT = `${backendHost}`;


