# services-api
Take home challenge for Kong.

## Summary


## Requirements
1. User can see the name, a brief description, and versions available for a given service
2. User can navigate to a given service from its card
3. User can search for a specific service

[!](https://www.figma.com/file/zeaWiePnc3OCe34I4oZbzN/Service-Card-List?node-id=0%3A1)

### Functionality
support filtering, sorting, pagination

### Constraints
- API is read-only
- 12 services sahll be displayed per page
- Searches are performed against both name and description fields using fuzzy match

### Assumptions
- Number of services in the database is less than 100,000
- Service name can not be null
- Service description can be null
- Restricting/filtering search results by User/Organization is not important for this exercise
- Database schema is stable
- Name and description values are english
- Do not need to support live-updating search results as a user types

### Out of Scope
The mockup shows a few components. For clarity, the following shall not be implemented as a part of this exercise:
- "Add New Service" button 
- notification, help, and user (icons in top right)

## Considerations
### Pagination
Pagination is accomplished using offset/limit. This method is straightforward to implement for both the API and database as well as reasonably preformant for our assumed data size. I have set a global config for the limit  to represent a flexble way for engineering to update or even a/b test these page sizes across otherwise identically deployed services. The config limit can also quickly be updated to instead be a user-defined value to support a user defining the number of results on a page, for example.

### Search
I leveraged postgres' extensions to implement a fuzzy search. To support queries against both the name and description fields, I concatenated both fields and built a trigram index. I also concat the name and description when we query. Queries remain pretty quick, and from what I read online this method of indexing and searching should remain reasonably performant as the amount of data grows.

### Datastore
I chose postgres for my data store based on the following considerations:
- Schema is stable, won't need to be changed in the near future
- It supports extensions for fuzzy searching, and advantage over mysql
- There's several drivers for go's sql package