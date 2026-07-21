from dataclasses import dataclass


@dataclass
class Post:
    id: str
    title: str
    body: str
    author_id: str
