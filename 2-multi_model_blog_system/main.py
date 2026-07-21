import json
import time
from pathlib import Path
from typing import cast

from config import Config
from datagen.data_generator import DataGenerator
from model.comment import Comment
from model.like import Like
from model.post import Post
from model.user import User
from repository.abstract_repository import AbstractRepository
from repository.mongo import MongoRepository
from repository.neo4j import Neo4jRepository
from repository.postgres import PostgresRepository
from typedefs.json_types import BlogData

DATA_PATH = Path(__file__).parent / "data.json"


def load_data_from_json() -> tuple[list[User], list[Post], list[Comment], list[Like]]:
    with open(DATA_PATH, "r") as f:
        data = cast(BlogData, json.load(f))
    return (
        [User(**user) for user in data["users"]],
        [Post(**post) for post in data["posts"]],
        [Comment(**comment) for comment in data["comments"]],
        [Like(**like) for like in data["likes"]],
    )


postgres_repo = PostgresRepository(Config.POSTGRES_DB, Config.POSTGRES_USER, Config.POSTGRES_PASSWORD)
mongo_repo = MongoRepository(Config.MONGO_USER, Config.MONGO_PASSWORD)
neo4j_repo = Neo4jRepository(Config.NEO4J_PASSWORD)

repos: list[AbstractRepository] = [postgres_repo, mongo_repo, neo4j_repo]


def load_data(use_synthetic: bool = True) -> tuple[list[User], list[Post], list[Comment], list[Like]]:
    if not use_synthetic:
        return load_data_from_json()
    generator = DataGenerator(
        no_users=300,
        no_posts=800,
        viral_post_ratio=0.03,
        viral_comment_range=(300, 1500),
        normal_comment_range=(0, 15),
        viral_likes_range=(100, 280),
        normal_likes_range=(0, 50),
    )
    users = generator.generate_users()
    posts = generator.generate_posts(users)
    comments = generator.generate_comments(users, posts)
    likes = generator.generate_likes(users, posts)
    return users, posts, comments, likes


users, posts, comments, likes = load_data(use_synthetic=True)

for repo in repos:
    start = time.perf_counter()
    repo.clear()
    repo.create_schema()
    repo.seed_users(users)
    repo.seed_posts(posts)
    repo.seed_comments(comments)
    repo.seed_likes(likes)
    repo.close()
    end = time.perf_counter()
    print(f"time to insert data {end - start:.6f} seconds")
