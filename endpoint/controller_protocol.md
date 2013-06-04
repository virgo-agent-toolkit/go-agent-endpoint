Communication between controller and an endpoint is through a RESTful web API that a controller shall provide.

POST, PUT, PATCH method should have a body that is a valid json.

## Authentication
Request:
```
method: POST
url:    /ctrl/auth
content:
{
  "agent_name": "name of the agent",
  "agent_id": "agent ID",
  "token": "token of the agent"
}
```
Response:
* Authentication Succeeded:
HTTP Status: 200 (OK)

* Authentication Failed:
HTTP Status 203 (Non-Authoritative Information)
