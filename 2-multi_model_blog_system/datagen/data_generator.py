import random
import uuid

from faker import Faker

from model.comment import Comment
from model.like import Like
from model.post import Post
from model.user import User


class DataGenerator:
    def __init__(
        self,
        no_users: int,
        no_posts: int,
        viral_post_ratio: float,
        viral_comment_range: tuple[int, int],
        normal_comment_range: tuple[int, int],
        viral_likes_range: tuple[int, int],
        normal_likes_range: tuple[int, int],
    ) -> None:
        self.fake = Faker()
        self._no_users = no_users
        self._no_posts = no_posts
        self._viral_post_ratio = viral_post_ratio
        self._viral_comment_range = viral_comment_range
        self._normal_comment_range = normal_comment_range
        self._viral_likes_range = viral_likes_range
        self._normal_likes_range = normal_likes_range
        self._viral_post_ids: list[str] = []

    def generate_users(self) -> list[User]:
        return [
            User(id=str(uuid.uuid4()), username=self.fake.unique.user_name(), email=self.fake.unique.email())
            for _ in range(self._no_users)
        ]

    def generate_posts(self, users: list[User]) -> list[Post]:
        posts: list[Post] = []
        for _ in range(self._no_posts):
            posts.append(
                Post(
                    id=str(uuid.uuid4()),
                    title=self.fake.catch_phrase(),
                    body=self.fake.paragraph(),
                    author_id=self.fake.random_element(users).id,
                )
            )
        viral_posts = random.sample(posts, int(self._viral_post_ratio * len(posts)))
        self._viral_post_ids = [p.id for p in viral_posts]
        return posts

    def generate_comments(self, users: list[User], posts: list[Post]) -> list[Comment]:
        comments: list[Comment] = []
        REPLY_RATIO = 0.02
        for post in posts:
            current_post_comments: list[Comment] = []
            # default case for if it is not a viral post
            no_comments = random.randint(*self._normal_comment_range)
            if post.id in self._viral_post_ids:
                # viral post with larger number of comments
                no_comments = random.randint(*self._viral_comment_range)
            for _ in range(no_comments):
                new_comment = Comment(
                    id=str(uuid.uuid4()),
                    author_id=self.fake.random_element(users).id,
                    body=self.fake.paragraph(2),
                    post_id=post.id,
                    parent_comment_id=None,
                )
                current_post_comments.append(new_comment)
            no_replies = int(len(current_post_comments) * REPLY_RATIO)
            for comment in random.sample(current_post_comments, no_replies):
                comment.parent_comment_id = random.choice([c for c in current_post_comments if c.id != comment.id]).id
            comments.extend(current_post_comments)
        return comments

    def generate_likes(self, users: list[User], posts: list[Post]) -> list[Like]:
        likes: list[Like] = []
        for post in posts:
            no_likes = random.randint(*self._normal_likes_range)
            if post.id in self._viral_post_ids:
                no_likes = random.randint(*self._viral_likes_range)
            for user in random.sample(users, no_likes):
                likes.append(Like(user_id=user.id, post_id=post.id))
        return likes
