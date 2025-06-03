#### A Starter Kit for writing http services using standard library for the most parts (except GORM for talking to the database)

- uses GO's new routing feature
- uuid as session token stored in Redis
- custom types for input sanitization, validator package to validate shape of data
- custom Rule interface for various validation (exists, unique, etc)
- reduces a bunch of boilterplate by extensively using generics
- becoming like a pseudo-framework but tries not to (aim to be as simple as I can)
