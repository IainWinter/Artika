# API v2

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

