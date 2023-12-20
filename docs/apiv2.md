# API v2

There are two levels to the API. One is the data API which returns JSON, and the other is the hypermedia API which returns HTML. For our initial build, we are creating a website which will use the hypermedia API. The data API allows for future development of apps that do not depend on the web.

POST /session
    JWT

DELETE /session
    SessionID

DELETE /user
    SessionID

UPDATE /user
    SessionID
    Diff list [
        field: value
    ]

GET /user/designers

POST /work
    SessionID
    Work Name
    Work Description
    Work Thumbnail

DELETE /work
    SessionID
    WorkID

UPDATE /work
    SessionID
    WorkID
    Diff list [
        field: value
    ]

POST /cloth
    SessionID

