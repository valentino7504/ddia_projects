from dataclasses import dataclass


@dataclass
class Like:
    user_id: str
    post_id: str
