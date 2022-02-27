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
- Searches are performed against both name and description fields using fuzzy match

### Assumptions
- Number of services in the database is less than 100,000
- Service name can not be null
- Service description can be null
- Restricting/filtering search results by User/Organization is not important for this exercise
- Database schema is stable
- Name and description values are english

### Out of Scope
The mockup shows a few components. For clarity, the following shall not be implemented as a part of this exercise:
- "Add New Service" button 
- notification, help, and user (icons in top right)

## Considerations
### Pagination
Pagination is accomplished using offset/limit. This method is straightforward to implement for both the API and database as well as reasonably preformant for our assumed data size.

### Search
- I do not need to support "live" search results, only return results once for a single query

### Datastore
I chose postgres for my data store based on the following considerations:
- There is a relationship between tables
- I wont need to change the schema in the foreseeable future
- It supports extensions for fuzzy searching, and advantage over mysql

```
services=# explain analyse select * from servicetable where similarity(description, 'serv') > 0.05;
                                                QUERY PLAN

----------------------------------------------------------------------------------------
------------------
 Seq Scan on servicetable  (cost=0.00..22.75 rows=283 width=68) (actual time=0.072..0.22
7 rows=4 loops=1)
   Filter: (similarity(description, 'serv'::text) > '0.05'::double precision)
   Rows Removed by Filter: 14
 Planning Time: 0.092 ms
 Execution Time: 0.474 ms
(5 rows)


services=# explain analyse select * from servicetable where similarity((name || ' ' || description), 'not') > 0.05;
                                                QUERY PLAN

----------------------------------------------------------------------------------------
------------------
 Seq Scan on servicetable  (cost=0.00..27.00 rows=283 width=68) (actual time=0.046..0.36
0 rows=3 loops=1)
   Filter: (similarity(((name || ' '::text) || description), 'not'::text) > '0.05'::doub
le precision)
   Rows Removed by Filter: 15
 Planning Time: 0.115 ms
 Execution Time: 0.470 ms
(5 rows)


services=# explain analyse select * from servicetable where name || ' ' || description ilike '%serv%';
                                               QUERY PLAN

----------------------------------------------------------------------------------------
----------------
 Seq Scan on servicetable  (cost=0.00..24.88 rows=7 width=68) (actual time=0.158..0.215
rows=3 loops=1)
   Filter: (((name || ' '::text) || description) ~~* '%serv%'::text)
   Rows Removed by Filter: 15
 Planning Time: 0.174 ms
 Execution Time: 0.322 ms
(5 rows)
```