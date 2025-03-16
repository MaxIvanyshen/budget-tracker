(function() {
  document.addEventListener('DOMContentLoaded', function() {
    if (typeof htmx !== 'undefined') {
      htmx.on('htmx:beforeRequest', function(evt) {
        const token = localStorage.getItem('accessToken');
        if (token) {
          evt.detail.headers['Authorization'] = 'Bearer ' + token;
        }
      });
      
      htmx.on('htmx:responseError', function(evt) {
        if (evt.detail.xhr.status === 401 || evt.detail.xhr.status === 403) {
          localStorage.removeItem('authToken');
          window.location.href = '/login';
        }
      });
      
      console.log('HTMX Auth module initialized');
    } else {
      console.error('HTMX not loaded. Auth module not initialized.');
    }
  });
})();

