export let baseUrl = 'http://localhost:3000/api/v1';

if (window.location.protocol === 'https:') {
  baseUrl = 'https://rating-party.onrender.com/api/v1'; // TODO: change to production url
}
