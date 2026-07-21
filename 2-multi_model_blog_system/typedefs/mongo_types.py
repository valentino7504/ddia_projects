from typing import TypedDict


class UserDocument(TypedDict):
    id: str
    username: str
    email: str


class CommentDocument(TypedDict):
    id: str
    author_id: str
    body: str
    parent_comment_id: str | None


class PostDocument(TypedDict):
    id: str
    title: str
    body: str
    author_id: str
    comments: list[CommentDocument]
    likes: list[str]
