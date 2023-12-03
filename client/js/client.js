class Client {
    /**
     * Create a session from a jwt token from a third party authentication service.
     * Sets a cookie with the SessionID upon success.
     * 
     * @param {string} jwt
     * @returns {Promise<Object>}
     */
    createSession(jwt) {
        return this.#post("api/session", {"JWT": jwt})
            .then((json) => {
                Cookies.set("SessionID", json["SessionID"], { expires: 1 });
                location.reload();
            });
    }

    /**
     * Deletes the session from the server and removes the SessionID cookie.
     * 
     * @returns {Promise<Object>}
     */
    deleteSession() {
        let sessionID = Cookies.get("SessionID");

        return this.#delete("api/session", {"SessionID": sessionID})
            .then((json) => {
                Cookies.remove("SessionID");
                location.reload();
            });
    }

//private:
    #fetch(url, method, body) {
        let options = {
            method: method,
            headers: { "Content-Type": "application/json" }
        };

        if (body !== null) {
            options["body"] = JSON.stringify(body);
        }

        let promise = fetch(`${this.#backendURL}/${url}`, options)
                     .then(response => response.json());
        
        return promise;
    }

    #get(url) {
        return this.#fetch(url, "GET", null);
    }

    #post(url, body) {
        return this.#fetch(url, "POST", body);
    }

    #delete(url, body) {
        return this.#fetch(url, "DELETE", body);
    }

    #backendURL = "http://localhost:3000";
};

G_Client = new Client();