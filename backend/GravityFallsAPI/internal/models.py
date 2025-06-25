from pydantic import BaseModel
from typing import Optional

class Character(BaseModel):
    Id: int
    Name: str
    Species: str
    Likes: str
    Quote: str
    Image: str
