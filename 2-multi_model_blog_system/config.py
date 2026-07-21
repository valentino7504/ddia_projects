import os
from typing import final

from dotenv import load_dotenv

load_dotenv()


@final
class Config:
    POSTGRES_USER = os.getenv("POSTGRES_USER") or ""
    POSTGRES_PASSWORD = os.getenv("POSTGRES_PASSWORD") or ""
    POSTGRES_DB = os.getenv("POSTGRES_DB") or ""
    MONGO_USER = os.getenv("MONGO_USER") or ""
    MONGO_PASSWORD = os.getenv("MONGO_PASSWORD") or ""
    NEO4J_PASSWORD = os.getenv("NEO4J_PASSWORD") or ""
