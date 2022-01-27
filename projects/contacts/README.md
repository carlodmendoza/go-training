# ASSIGNMENT #4: REST API FOR CONTACTS DB
---
# NAME

contacts

# SYNOPSIS

http://localhost:8080/contacts

# DESCRIPTION

The webapp implements REST API for a database of contacts. All data are represented in json format. 

| HTTP Verb | Entire collection /contacts  | Specific item /contacts/{id} |
|-----------|------------|----------------|
| POST      | 201 (Created), creates new crecord; 409 (Conflict), retrieves current record of contact | 405 (Not allowed) |
| GET       | 200 (OK), retrieves all records | 200 (OK), retrieves record; 404 (Not found) |
| PUT      | 405 (Not allowed) | 200 (OK), updates record; 404 (Not found) |
| DELETE   | 405 (Not allowed) | 200 (OK), removes record; 404 (Not found) | 

*CONTACTS DATABASE*

| Field | Data Type | Description |
|-------|-----------|-------------|
| Last  | string    | Last name   |
| First | string    | First name  |
| Company | string  | Company name |
| Address | string | Company address |
| Country | string | Country name |
| Position | string | Position in the Company |

# ERROR HANDLING

Always report errors encountered (e.g., incorrect url)