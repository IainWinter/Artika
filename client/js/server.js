// This file contains the configuration for the client application

let G_ServerBackend = "http://localhost:3000";

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

/**
 * Wrap a fetch call to the server backend and decode it to json upon completion. 
 * Sends a DELETE http request with a json body
 * 
 * @param {string} url
 * @param {Object} body
 * @returns {Promise<Object>}
 */
function api_delete(url, body) {
    let options = {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body)
    };

    let promise = fetch(`${G_ServerBackend}/${url}`, options)
                 .then(response => response.json());
    
    return promise;
}