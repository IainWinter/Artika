// This file contains the configuration for the client application

let G_ServerBackend = "http://localhost:3001";

/**
 * Wrap a fetch call to the server backend and decode it to json upon completion. 
 * 
 * @param {string} url
 * @param {Object} options
 * @returns {Promise<Object>}
 */
function arFetch(url, options) {
    return fetch(`${G_ServerBackend}/${url}`, options)
        .then(response => response.json());
}