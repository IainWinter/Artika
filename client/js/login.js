/**
 * This is the entry point after a user logs into the application.
 * All signin-with credentials are encoded in JSON web tokens (JWTs).
 * See https://jwt.io/ for more information.
 * 
 * @param {string} jwt 
 */
function user_session_create(jwt) {
    console.log("Authenticating user.");

    arFetch('api/session/create', {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({"jwt": jwt})
    })
    .then(json => {
        Cookies.set('session-id', json['session-id'], { expires: 1 });

        // This actually shouldn't be set on the client
        //Cookies.set('session-expiration', json['session-expiration'], { expires: 1 });
    });
}

function user_session_validate() {
    let sessionId = Cookies.get('session-id');

    arFetch('api/session/validate', {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({"session-id": sessionId})
    })
    .then((json) => {
        let isSessionValid = json['is-session-valid'];
        console.log(isSessionValid);
    });
}

function getSessionId() {
    return Cookies.get('session-id');
}