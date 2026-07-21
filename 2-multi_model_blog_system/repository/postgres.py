from typing import final, override

import psycopg

from model.comment import Comment
from model.like import Like
from model.post import Post
from model.user import User
from repository.abstract_repository import AbstractRepository


@final
class PostgresRepository(AbstractRepository):
    def __init__(self, db: str, username: str, password: str) -> None:
        self.connection = psycopg.connect(f"dbname={db} user={username} password={password} host=localhost port=5432")

    @override
    def seed_users(self, users: list[User]) -> None:
        with self.connection.cursor() as cur:
            cur.executemany(
                "INSERT INTO users (id, username, email) VALUES (%s, %s, %s)",
                [(user.id, user.username, user.email) for user in users],
            )
            self.connection.commit()

    @override
    def seed_comments(self, comments: list[Comment]) -> None:
        with self.connection.cursor() as cur:
            cur.executemany(
                """INSERT INTO comments
                (id, body, post_id, author_id, parent_comment_id)
                VALUES (%s, %s, %s, %s, %s)""",
                [
                    (
                        comm.id,
                        comm.body,
                        comm.post_id,
                        comm.author_id,
                        comm.parent_comment_id,
                    )
                    for comm in comments
                ],
            )
            self.connection.commit()

    @override
    def seed_posts(self, posts: list[Post]) -> None:
        with self.connection.cursor() as cur:
            cur.executemany(
                """INSERT INTO posts (id, title, body, author_id)
                VALUES (%s, %s, %s, %s)""",
                [
                    (
                        post.id,
                        post.title,
                        post.body,
                        post.author_id,
                    )
                    for post in posts
                ],
            )
            self.connection.commit()

    @override
    def seed_likes(self, likes: list[Like]) -> None:
        with self.connection.cursor() as cur:
            cur.executemany(
                "INSERT INTO likes (user_id, post_id) VALUES (%s, %s)",
                [(like.user_id, like.post_id) for like in likes],
            )
            self.connection.commit()

    @override
    def clear(self) -> None:
        with self.connection.cursor() as cur:
            cur.execute("TRUNCATE TABLE users, posts, comments, likes CASCADE")
            self.connection.commit()

    @override
    def create_schema(self) -> None:
        with self.connection.cursor() as cur:
            cur.execute("""
                CREATE TABLE IF NOT EXISTS users (
                    id VARCHAR(255) PRIMARY KEY,
                    username VARCHAR(255) UNIQUE NOT NULL,
                    email VARCHAR(255) UNIQUE NOT NULL
                )
            """)
            cur.execute("""
                CREATE TABLE IF NOT EXISTS posts (
                    id VARCHAR(255) PRIMARY KEY,
                    title text NOT NULL,
                    body text NOT NULL,
                    author_id VARCHAR(255) NOT NULL REFERENCES users(id)
                )
            """)
            cur.execute("""
                CREATE TABLE IF NOT EXISTS comments (
                    id VARCHAR(255) PRIMARY KEY,
                    post_id VARCHAR(255) NOT NULL REFERENCES posts(id),
                    author_id VARCHAR(255) NOT NULL REFERENCES users(id),
                    body text NOT NULL,
                    parent_comment_id VARCHAR(255)
                )
            """)
            cur.execute("""
                CREATE TABLE IF NOT EXISTS likes (
                    user_id VARCHAR(255) NOT NULL REFERENCES users(id),
                    post_id VARCHAR(255) NOT NULL REFERENCES posts(id),
                    PRIMARY KEY(user_id, post_id)
                )
            """)
            cur.execute("CREATE INDEX IF NOT EXISTS idx_comments_author_id" + " ON comments(author_id)")
            cur.execute("CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id)")
            cur.execute("CREATE INDEX IF NOT EXISTS idx_likes_user_id ON likes(user_id)")
            cur.execute("CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes(post_id)")
            cur.execute("CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id)")
            self.connection.commit()

    @override
    def close(self) -> None:
        return self.connection.close()

    def __str__(self) -> str:
        return "PostgreSQL"
