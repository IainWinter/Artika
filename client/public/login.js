/**
 * This is the entry point after a user logs into the application.
 * All signin-with credentials are encoded in JSON web tokens (JWTs).
 * See https://jwt.io/ for more information.
 * 
 * @param {string} jwt 
 */
function user_create_session(jwt) {
    console.log("Authenticating user.");

    arFetch('api/session/create', {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({"jwt": jwt})
    })
    .then(json => {
        console.log(json);
    });
}