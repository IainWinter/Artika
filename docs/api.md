# Artika API Specification

Outline:

* users
    * /user/create
    * /user/delete
    * /user/edit
* designers
    * /designer/list
* sessions
    * /session/create
    * /session/delete
* requests
    * /request/create
    * /request/delete
    * /request/edit
    * /request/designer/respond
    * /request/user/respond

## Users

There are two types of users, a 'designer' and a normal 'user'. I will refer to normal users as 'user' and 
designers as 'designer'. The designer account has a superset of the permissions of a user account.

### POST /user/create

| Parameters ||| 
|-|-|-|
| `fullName`   | `string`  | Required |
| `address`    | `string`  | Required if a designer |
| `isDesigner` | `boolean` | Required |
| `avatar`     | `binary`  | If this is not set, a default will be used. |
| `oAuth2Id`   | `string`  | Required. The only supported login schemes will be through other services. |

| Returns ||| 
| ------- | ------ | -- |
| `error` | `string` | If there is an error, this will contain the reason for UI display. |

### DELETE /user/delete

| Parameters ||| 
|-|-|-|
| `session`   | `string`  | Required |

| Returns ||| 
| ------- | ------ | -- |
| `error` | `string` | If there is an error, this will contain the reason for UI display. |

### PATCH /user/edit

Edit the fields of a user profile. All edits have to be valid for the update to apply.

```
UserProfileField {
    field : string
    value : string
}
```

| Parameters ||| 
|-|-|-|
| `session` | `string`   | Required |
| `fields`  | `UserProfileField[]` | Required |

| Returns ||| 
|-|-|-|
| `error` | `string` | If there is an error, this will contain the reason for UI display. |

## Sessions

### POST /session/create

| Parameters ||| 
|-|-|-|
| `oAuth2Id` | `string` | Required |

| Returns ||| 
|-|-|-|
| `session` | `string` | If there is no error, this will contain the session token to store in a cookie. |
| `error`   | `string` | If there is an error, this will contain the reason for UI display. |

### DELETE /session/delete

| Parameters ||| 
|-|-|-|
| `session` | `string` | Required |

| Returns ||| 
|-|-|-|
| `error` | `string` | If there is an error, this will contain the reason for UI display. |

## Designers

### GET /designer/list

List all designers, this allows users to select one to request work from.

```
Designer {
    DesignerId : string
    FullName : string
    AvatarURI : string
    ThumbnailURI: string
}
```

| Parameters ||| 
|-|-|-|
| `offset` | `integer` | Defaults to 0 |

| Returns ||| 
|-|-|-|
| `designers` | `Designer[]` | Could be empty. Will always be empty if there is an error. |
| `error`     | `string`     | If there is an error, this will contain the reason for UI display. |

## Work Requests

A work request is a request from a user to a designer. The user can upload a description and some photos for the designer to get an idea of what the user wants.
The designer can respond to the request by sending a message back.

### POST /request/create

| Parameters ||| 
|-|-|-|
| `session`     | `string`   | Required |
| `designerId`  | `string`   | Required |
| `title`       | `string`   | Required |
| `description` | `string`   | Required |
| `images`      | `binary[]` | Required, but can be an empty list. The first image will become the thumbnail. |

| Returns ||| 
|-|-|-|
| `error` | `string` | If there is an error, this will contain the reason for UI display. |

### DELETE /request/delete

Delete a request. If its response state is set to Open or Negotiated, it will just be deleted. If its response is anything else
there will be a request to delete and both parties have to agree. 

| Parameters ||| 
|-|-|-|
| `session`     | `string`   | Required |
| `requestId`   | `string`   | Required |

| Returns ||| 
|-|-|-|
| `error` | `string` | If there is an error, this will contain the reason for UI display. |

### GET /request/list

A designer can get a list of all their requests

```
WorkRequest {
    UserId : string
    UserFullName : string
    UserAvatarURI : string
    Title : string
    Description: string
    Images: string[]
}
```

| Parameters ||| 
|-|-|-|
| `session` | `string` | Required, must map to a designer account. |

| Returns ||| 
|-|-|-|
| `requests` | `WorkRequest[]` | Could be empty. Will always be empty if there is an error. |
| `error` | `string` | If there is an error, this will contain the reason for UI display. |

### POST /request/designer/response

A designer can respond to a users request for work by either denying, accepting, or a custom message.

A request can be in 5 differnt states.

| State ||
|-|-|
| Open | No designer has taked this request |
| Negotiated | A designer has responded with a custom message, not denying or accepting. (implement in the future) |
| Accepted | A designer has accepted to work on the request |
| Denied | The asked designer has denied to work on the request |
| Completed | The designer completed the work |

```
WorkReponse {
    Open,
    Negotiated,
    Accepted,
    Denied,
    Completed
}
```

| Parameters ||| 
|-|-|-|
| `response` | `WorkReponse` | Required |
| `message` | `string` | Required, but can be empty. |

| Returns ||| 
|-|-|-|
| `error` | `string` | If there is an error, this will contain the reason for UI display. |

### POST /request/user/response

After a designer has responded to a request, the user can respond.

