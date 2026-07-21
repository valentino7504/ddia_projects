from typing import final, override

from neo4j import Driver, GraphDatabase, Transaction

from model.comment import Comment
from model.like import Like
from model.post import Post
from model.user import User
from repository.abstract_repository import AbstractRepository


@final
class Neo4jRepository(AbstractRepository):
    def __init__(self, password: str) -> None:
        URI = "bolt://localhost:7687"
        AUTH = ("neo4j", password)
        self.driver: Driver = GraphDatabase.driver(URI, auth=AUTH)

    @override
    def seed_users(self, users: list[User]) -> None:
        with self.driver.session() as session:
            tx = session.begin_transaction()
            try:
                for user in users:
                    tx.run(
                        "CREATE (u:User {id: $id, username: $username, email: $email})",
                        id=user.id,
                        username=user.username,
                        email=user.email,
                    )
            except Exception:
                tx.rollback()
                raise
            else:
                tx.commit()

    @override
    def seed_posts(self, posts: list[Post]) -> None:
        with self.driver.session() as session:
            tx = session.begin_transaction()
            try:
                for post in posts:
                    tx.run(
                        """MATCH (u: User {id: $author_id})
                        CREATE (u)-[:WROTE]->(p: Post {id: $id, title: $title, body: $body})""",
                        author_id=post.author_id,
                        id=post.id,
                        title=post.title,
                        body=post.body,
                    )
            except Exception:
                tx.rollback()
                raise
            else:
                tx.commit()

    @override
    def seed_comments(self, comments: list[Comment]) -> None:
        with self.driver.session() as session:
            tx = session.begin_transaction()
            try:
                for comment in comments:
                    if comment.parent_comment_id is None:
                        self._insert_comment_on_post(tx, comment)
                    else:
                        self._insert_comment_on_parent_comment(tx, comment)
            except Exception:
                tx.rollback()
                raise
            else:
                tx.commit()

    def _insert_comment_on_post(self, tx: Transaction, comment: Comment) -> None:
        tx.run(
            """MATCH (u:User {id: $author_id})
            MATCH (p:Post {id: $post_id})
            CREATE (u)-[:WROTE]->(c:Comment {id: $id, body: $body})
            CREATE (c)-[:REPLIES_TO]->(p)""",
            id=comment.id,
            body=comment.body,
            author_id=comment.author_id,
            post_id=comment.post_id,
        )

    def _insert_comment_on_parent_comment(self, tx: Transaction, comment: Comment) -> None:
        tx.run(
            """MATCH (u:User {id: $author_id})
            MATCH (parent:Comment {id: $parent_comment_id})
            CREATE (u)-[:WROTE]->(c:Comment {id: $id, body: $body})
            CREATE (c)-[:REPLIES_TO]->(parent)""",
            id=comment.id,
            body=comment.body,
            author_id=comment.author_id,
            parent_comment_id=comment.parent_comment_id,
        )

    @override
    def seed_likes(self, likes: list[Like]) -> None:
        with self.driver.session() as session:
            tx = session.begin_transaction()
            try:
                for like in likes:
                    tx.run(
                        """
                        MATCH (u: User {id: $user_id})
                        MATCH (p: Post {id: $post_id})
                        CREATE (u)-[:LIKES]->(p)
                        """,
                        user_id=like.user_id,
                        post_id=like.post_id,
                    )
            except Exception:
                tx.rollback()
                raise
            else:
                tx.commit()

    @override
    def clear(self) -> None:
        with self.driver.session() as session:
            session.run("MATCH ()-[r]-() DELETE r")
            session.run("MATCH (n) DELETE n")

    @override
    def create_schema(self) -> None:
        with self.driver.session() as session:
            session.run("CREATE CONSTRAINT IF NOT EXISTS FOR (u:User) REQUIRE u.id IS UNIQUE")
            session.run("CREATE CONSTRAINT IF NOT EXISTS FOR (p:Post) REQUIRE p.id IS UNIQUE")
            session.run("CREATE CONSTRAINT IF NOT EXISTS FOR (c:Comment) REQUIRE c.id IS UNIQUE")

    @override
    def close(self) -> None:
        self.driver.close()

    def __str__(self) -> str:
        return "Neo4j"
