from dataclasses import dataclass


@dataclass
class Comment:
    id: str
    post_id: str
    author_id: str
    body: str
    parent_comment_id: str | None
