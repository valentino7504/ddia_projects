from typing import TypedDict


class UserDict(TypedDict):
    id: str
    username: str
    email: str


class PostDict(TypedDict):
    id: str
    title: str
    body: str
    author_id: str


class CommentDict(TypedDict):
    id: str
    post_id: str
    author_id: str
    body: str
    parent_comment_id: str | None


class LikeDict(TypedDict):
    user_id: str
    post_id: str


class BlogData(TypedDict):
    users: list[UserDict]
    posts: list[PostDict]
    comments: list[CommentDict]
    likes: list[LikeDict]
