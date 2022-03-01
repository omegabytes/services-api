# Examples

## Authorizing
I only support one user for this demo project.
```
curl localhost:8080/api/authenticate -F name=user -F password=pass
{
	Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidXNlciIsInJvbGUiOiJtZW1iZXIifQ.XBqsh1PM4Ne9eETglL7aSve-YWlCzUUvp0evHQxAxN0
}

curl localhost:8080/api/authenticate -F name=user -F password=badpass
{
    "message": "Unable to sign in with that user or password",
    "status": 401
}
```

## Unauthorized Request
```
 curl localhost:8080/api/v1/services -v
< HTTP/1.1 401 Unauthorized
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Mon, 28 Feb 2022 22:54:35 GMT
< Content-Length: 29
<
Missing authorization header
* Connection #0 to host localhost left intact
```

## Get Service by ID
```
curl localhost:8080/api/v1/services/1 -H 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidXNlciIsInJvbGUiOiJtZW1iZXIifQ.XBqsh1PM4Ne9eETglL7aSve-YWlCzUUvp0evHQxAxN0' | json_pp
[
   {
      "Description" : "Cloud storage monitoring",
      "Id" : 1,
      "Name" : "Notifications",
      "Versions" : [
         {
            "semver" : "0.0.0"
         },
         {
            "semver" : "0.1.0"
         },
         {
            "semver" : "2.1.0"
         }
      ]
   }
]
```

## Searching
```
curl localhost:8080/api/v1/services?search=notification -H 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidXNlciIsInJvbGUiOiJtZW1iZXIifQ.XBqsh1PM4Ne9eETglL7aSve-YWlCzUUvp0evHQxAxN0' -v | json_pp

< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Mon, 28 Feb 2022 02:12:34 GMT
< Content-Length: 489
<
[
   {
      "Description" : "Customizable notifications",
      "Id" : 0,
      "Name" : "Notifications",
      "Versions" : [
         {
            "semver" : "0.0.0"
         },
         {
            "semver" : "0.1.0"
         },
         {
            "semver" : "2.1.0"
         }
      ]
   },
   {
      "Description" : "Cloud storage monitoring",
      "Id" : 1,
      "Name" : "Notifications",
      "Versions" : [
         {
            "semver" : "0.0.0"
         },
         {
            "semver" : "0.1.0"
         },
         {
            "semver" : "2.1.0"
         }
      ]
   },
   {
      "Description" : "Multi-channel communication tools",
      "Id" : 4,
      "Name" : "Contact Us",
      "Versions" : [
         {
            "semver" : "0.0.0"
         }
      ]
   },
   {
      "Description" : "Suite of notification tools",
      "Id" : 6,
      "Name" : "Contact Us",
      "Versions" : null
   }
]
```