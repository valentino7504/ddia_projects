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

### Embedding at Scale

This is where we start to see the frailties of the document database style. If one chooses to
embed comments entirely into the post documents then we tend to have the following problems:

1. Constant overwriting of post documents to add new comments. Even though this is a problem whether or not
   you normalize/use references, it is worse with embedding because the document grows wildly in size, and the
   larger the document gets, the more heavy this operation is. Having separate comment documents also does not really solve
   eliminate the cost, it just moves it elsewhere because then you have to simulate a JOIN in your application code, just
   to get all the comments for a single post. Of course the cost of this join is not as efficient as say a relational database,
   because relational databases are optimized for this kind of operation, while most document databases, even if they have something
   similar eg `$lookup`, are not built for j oin-optimized structures.
2. If you have a mega viral post, all of a sudden the document size limit becomes a real issue. Storing entire comment data,
   including the author_id, the body of the comment and the id of the comment in the same document on a post with say 1.5M
   comments becomes an issue. Yes this is unlikely to happen (the record for comments on X is about 707k), but even with 707k,
   the limit for document size is still reached. And this is without considering the performance degradation when documents are
   hitting 1MB, imagine how it becomes when they are 16MB.
3. Concurrency. With a relatively popular post, another issue that immediately rears its head is concurrent comments competing
   for the same document lock at the same time, which might lead to higher throughput and latency.

### Adding Likes

1. MongoDB: The easiest approach will probably be to add likes to the posts, but then in the context of adding likes after posts and
   comments already exists, you have to go back and alter existing posts to add the likes field. Of course this poses an issue.
   Generally introducing new mandatory properties on a "schema" requires going back to rewrite all the existing documents of that
   schema to introduce the new property, which depending on the number of documents, is a very heavy series of operations. Depending
   on the approach, there will still be query time costs because if you want to fetch the users that liked a post then you have to
   do an application join to find the users who liked it, if you need it.
2. Neo4j: The major cost is at query time, especially with introducing the new relationship. Fetching posts with their likes involves
   tracing the [:LIKES] relationships backwards (i.e <-[:LIKES]) and finding the users who liked a post.
3. PostgreSQL: Also pushes the cost to query time, introducing a new likes table with two foreign keys for user_id and post_id. Fetching
   posts with their likes is just using the foreign key reference to get the count. Getting the individual users involves a join on user_id.

## Benchmark Validity Caveat

I will be using two datasets, a small hand authored one for correctness verification across the three databases, then a larger synthetically
generated one for timing benchmarks. Timing results in [the query analysis](#query-analysis) are based on the synthetic larger dataset.
Everything is being run on localhost, and not under realistic concurrent load, and also with no replicas or shards (obviously). My hardware
is an Intel Core i5 11th Gen laptop - not exactly benchmark-grade hardware.

## Conclusion

Based off of what I have studied so far and the results, I think it is safe to say that there is no "best" database. It simply depends on the
use case and what kind of operations your application does. If your application is primarily stuff like many to many relationships and your
it prioritises relationships over data itself, graph databases might be good. If your application needs locality, and ease of access to data,
with less joins, then a document database might serve you well, as documents are self-contained (depending on how you structure them). Of
course, when it comes to pure JOIN operations, simple CRUD and analytics, nothing beats a relational database for establishing this (hence
why they are so ubiquitous). Relational DBs bet on structural integrity and ad hoc querying via JOINs. Of course, this often comes at the cost
of locality, even though DBs like PostgreSQL have introduced extensive support for JSON. Graph databases tend to focus more on traversal via
pointers which can be very useful for establishing interconnected webs of relationships that might not be easily modelled in relational and
document databases. They also offer a lot of freedom to establish relationships because they are somewhat "schemaless" compared to relational
databases though this has its drawbacks as it provides a lot less guarantee of validity compared to regular relational databases.
Document databases also offer this freedom of "schemalessness". Technically in most cases, it is more so a transition to schema on read vs schema
on write (relational databases). They also offer superior locality, and are very useful in cases where you need self contained database objects where
all the data pertaining to an entity is stored with it for easy access. If you tend to do single reads more often, this is quite useful as just
one visit to the database gets you the information you need. Of course, this has its costs on write, as adding new information involves overwriting
documents in most cases which is obviously expensive. There is also the question of embedding, and how that exponentially increases the cost of reading
and appending over time.
