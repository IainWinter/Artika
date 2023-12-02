arFetch("api/session/is-session-valid", {
    method: "GET",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ "session-id": getSessionId() })
})
.then((json) => {
    let isSessionValid = json['is-session-valid'];
    console.log(isSessionValid);

    // If the session is valid, load the account header
    // if the session is invalid, load the logic page

    if (isSessionValid) {
                
    }
})