from pydantic import BaseModel

class Character(BaseModel):
    Id: int
    Name: str
    Species: str
    Likes: str
    Quote: str
    Image: str
