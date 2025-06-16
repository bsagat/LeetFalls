from fastapi import FastAPI, HTTPException
from internal.models import Character
import internal.database as repo
import os
from dotenv import load_dotenv

load_dotenv()

app = FastAPI()
DB_PATH = os.getenv("DATABASE_URL")

if not DB_PATH:
    raise RuntimeError("DB_PATH is not set in .env file")

@app.get("/")
def root():
    return {"Characters": "/characters"}

@app.get("/characters")
def get_characters_count():
    try:
        count = repo.CharacterCount()
        return {"Count": count}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/characters/{character_id}", response_model=Character)
def get_character(character_id: int):

    count = repo.CharacterCount()
    if character_id <= 0 or character_id > count:
        raise HTTPException(status_code=404, detail=f"Character ID must be between 1 and {count}")
        
    character = repo.CharacterByID(character_id)

    if character is None:
        raise HTTPException(status_code=404, detail="Character not found")

    return {
        "Id": character[0],
        "Name": character[1],
        "Species": character[2],
        "Likes": character[3],
        "Quote": character[4],
        "Image": character[5]
    }
