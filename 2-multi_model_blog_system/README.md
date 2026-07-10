# Multi Model Blog Mini Project To Analyse Different DBs

This is just a mini technical experiment to compare three database paradigms -
relational, document and graph using PostgreSQL, MongoDB and Neo4j respectively.

## Schema Mapping

| Concept  | PostgreSQL                                            | MongoDB                                                            | Neo4j                                                                           |
| -------- | ----------------------------------------------------- | ------------------------------------------------------------------ | ------------------------------------------------------------------------------- |
| Users    | `users` table, columns: id, username, email           | `users` collection, one document per user                          | `(:User)` node, properties: id, username, email                                 |
| Posts    | `posts` table, columns: id, title, body, user_id      | `posts` collection, fields: id, title, body, user_id, comments etc | `(:Post)` node, properties: id, title, body with WROTE relationship             |
| Comments | `comments` table, columns: post_id, user_id, id, body | embedded in `posts`                                                | `(:Comment)` node, properties: id, body. Relationships are WROTE and REPLIES_TO |
| Likes    | `likes` table, columns: post_id, user_id              | Array of user ids on posts                                         | `[:LIKES]` relationship                                                         |

## Query Analysis

## Structural Cost of Change

## Benchmark Validity Caveat

## Conclusion
