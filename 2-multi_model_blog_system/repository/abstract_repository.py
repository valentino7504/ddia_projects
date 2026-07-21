from abc import ABC, abstractmethod

from model.comment import Comment
from model.like import Like
from model.post import Post
from model.user import User


class AbstractRepository(ABC):
    @abstractmethod
    def create_schema(self) -> None:
        """should be no op in non schema based databases"""
        ...

    @abstractmethod
    def seed_users(self, users: list[User]) -> None: ...

    @abstractmethod
    def seed_posts(self, posts: list[Post]) -> None: ...

    @abstractmethod
    def seed_comments(self, comments: list[Comment]) -> None: ...

    @abstractmethod
    def seed_likes(self, likes: list[Like]) -> None: ...

    @abstractmethod
    def clear(self) -> None: ...

    @abstractmethod
    def close(self) -> None: ...
