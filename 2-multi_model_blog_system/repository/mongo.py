import urllib.parse
from typing import final, override

from pymongo import MongoClient
from pymongo.collection import Collection

from model.comment import Comment
from model.like import Like
from model.post import Post
from model.user import User
from repository.abstract_repository import AbstractRepository
from typedefs.mongo_types import CommentDocument, PostDocument, UserDocument

# pyright: reportUnknownMemberType=false


@final
class MongoRepository(AbstractRepository):
    def __init__(self, username: str, password: str) -> None:
        self.client = MongoClient(
            f"mongodb://{urllib.parse.quote_plus(username)}:{urllib.parse.quote_plus(password)}@localhost:27017/?authSource=admin"
        )
        self.db = self.client.get_database("blog")
        self._users: Collection[UserDocument] = self.db["users"]
        self._posts: Collection[PostDocument] = self.db["posts"]

    @override
    def seed_users(self, users: list[User]) -> None:
        self._users.insert_many([UserDocument(username=u.username, email=u.email, id=u.id) for u in users])

    @override
    def seed_posts(self, posts: list[Post]) -> None:
        self._posts.insert_many(
            [
                PostDocument(
                    id=p.id,
                    title=p.title,
                    body=p.body,
                    author_id=p.author_id,
                    comments=[],
                    likes=[],
                )
                for p in posts
            ]
        )

    @override
    def seed_comments(self, comments: list[Comment]) -> None:
        for comment in comments:
            query_filter: dict[str, str] = {"id": comment.post_id}
            update_operation = {
                "$push": {
                    "comments": CommentDocument(
                        id=comment.id,
                        author_id=comment.author_id,
                        body=comment.body,
                        parent_comment_id=comment.parent_comment_id,
                    )
                }
            }
            self._posts.update_one(query_filter, update_operation)

    @override
    def seed_likes(self, likes: list[Like]) -> None:
        for like in likes:
            query_filter: dict[str, str] = {"id": like.post_id}
            update_operation = {"$push": {"likes": like.user_id}}
            self._posts.update_one(query_filter, update_operation)

    @override
    def create_schema(self) -> None:
        # no schema to create - Mongo collections are created on first insert
        pass

    @override
    def clear(self) -> None:
        self._users.delete_many({})
        self._posts.delete_many({})

    @override
    def close(self) -> None:
        self.client.close()

    def __str__(self) -> str:
        return "MongoDB"
