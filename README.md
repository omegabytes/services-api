# services-api
Take home challenge for Kong. Versioned API with backing postgres and authorization using JWT.

## Running and testing
This project can be ran and tested using docker-compose and curl (or your API tool of choice).
First, clone this repo and navigate to the project root. To start the project, run:
```
docker-compse up
```

You'll see `app_1 | connected to services` after both the app and database containers successfully start, which means you're ready to go. The database is seeded with ~20 entries to test against. Check out the [init script](/Users/alex/Desktop/Programming/services-api/storage/intidbpsql.sql) to see what you can query on.

You can test the service using another terminal window or an API tool like insomnia. Try the following endpoints:
```
localhost:8080/api/authenticate               // fetch a token
localhost:8080/api/v1/services                // list first page of services
localhost:8080/api/v1/services?offset=3       // list services additional pages
localhost:8080/api/v1/services/{id}           // get service by id
localhost:8080/api/v1/services?search={term}  // fuzzy search
```

Tests can be run with `go test ./...` from the project root.

## Requirements
1. User can see the name, a brief description, and versions available for a given service
2. User can navigate to a given service from its card
3. User can search for a specific service
4. Supports filtering, sorting, pagination.

## Stretch goals
1. Example authorization
2. Example tests

![mockup](https://www.figma.com/file/zeaWiePnc3OCe34I4oZbzN/Service-Card-List?node-id=0%3A1)

### Constraints
- API is read-only
- 12 services shall be displayed per page
- Searches are performed against both name and description fields using fuzzy match

### Assumptions
These assumptions were defined during the technical interview.
- Number of services in the database is less than 100,000
- Service name can not be null
- Service description can be null
- Restricting/filtering search results by User/Organization is not important for this exercise
- Database schema is stable
- Name and description values are english
- Do not need to support live-updating search results as a user types
- "Add New Service" button, notification, help, and user (icons in top right) are out of scope
- Cacheing is out of scope

## Considerations
### Pagination
Pagination is accomplished using offset/limit. This method is straightforward to implement for both the API and database as well as reasonably preformant for our assumed data size. I have set a global config for the limit  to represent a flexble way for engineering to update or even A/B test these page sizes across otherwise identically deployed services. The config limit can also quickly be updated to instead be a user-defined value to support a user defining the number of results on a page, for example.

### Search
I leveraged postgres' extensions to implement a fuzzy search. To support queries against both the name and description fields, I concatenated both fields and built a trigram index. I also concat the name and description when we query. Queries remain snappy, and my understanding is this method of indexing/searching should remain reasonably performant as the amount of data grows (to a certain size, anyway). One tradeoff here is my limited understanding of performance when it comes to indexing/searching - I optimized for time-to-delivery of the project and put optimizing performance as a "follow-up task".

### Datastore
I chose postgres for my data store based on the following considerations:
- It supports extensions for fuzzy searching, an advantage over mysql
- There's several drivers for go's sql package
- Relational (I may want to add a user table for auth + a `createdBy` relationship to the service table)

### Tests
I only wrote super-basic tests for the handlers. One of the trade-offs of using gorilla for my API mux is the need to stand up an entire dummy server in order to run tests. This is something I would do in a production setting. Same with the datastore.

## Examples
See the [examples](./examples.md) file for a quick overview of the project in action.
