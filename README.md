# services-api
Take home challenge for Kong.

## Summary

## Requirements
1. User can see the name, a brief description, and versions available for a given service
2. User can navigate to a given service from its card
3. User can search for a specific service

### Functionality
support filtering, sorting, pagination

### Constraints
- API is read-only
- 12 services sahll be displayed per page

### Assumptions
- Number of services in the database is less than 100,000
- Service name can not be null
- Service description can be null
- Restricting/filtering search results by User/Organization is not important for this exercise
- database schema is stable 

### Out of Scope
The mockup shows a few components. For clarity, the following shall not be implemented as a part of this exercise:
- "Add New Service" button 
- notification, help, and user (icons in top right)

## Considerations
### Pagination
Pagination is accomplished using offset/limit. This method is straightforward to implement for both the API and database as well as reasonably preformant for our assumed data size.

### Datastore
I chose mysql for my data store based on the following considerations:
- There is a relationship between tables
- I wont need to change the schema in the foreseeable future
- I've worked with it in the past and am comfortable with it