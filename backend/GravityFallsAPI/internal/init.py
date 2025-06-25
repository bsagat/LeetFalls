from dotenv import load_dotenv
import pandas as pd
import psycopg2
import requests
import os

load_dotenv()
DATABASE_URL = os.getenv("DATABASE_URL")
S3_URL = os.getenv("S3_URL", 'http://0.0.0.0:9090')  # Лучше с протоколом
EXCEL_PATH = 'data/Characters.xlsx'
IMAGES_PATH = 'data/Character_Images.zip'


def Connect():
    return psycopg2.connect(DATABASE_URL)


def SaveImages():
    # Создать бакет
    requests.put(url=f'{S3_URL}/buckets/characters')

    # Отправить zip-файл
    with open(IMAGES_PATH, 'rb') as f:
        files = {'file': ('Character_Images.zip', f, 'application/zip')}
        resp = requests.put(url=f'{S3_URL}/objects/characters/jar', files=files)
        print(f"Status: {resp.status_code}, Response: {resp.text}")


def CreateTable():
    conn = Connect()
    cursor = conn.cursor()
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS characters (
            Id INTEGER PRIMARY KEY,
            Name TEXT NOT NULL,
            Species TEXT NOT NULL,
            Likes TEXT,
            Quote TEXT,
            Image TEXT NOT NULL    
        )
    ''')
    conn.commit()
    conn.close()


def LoadExcelData():
    conn = Connect()
    cursor = conn.cursor()

    cursor.execute("DELETE FROM characters")

    df = pd.read_excel(EXCEL_PATH)
    for _, row in df.iterrows():
        cursor.execute('''
            INSERT INTO characters (Id, Name, Species, Likes, Quote, Image)
            VALUES (%s, %s, %s, %s, %s, %s)
        ''', (row['ID'], row['Name'], row['Species'], row['Likes'], row['Quote'], row['Image']))

    conn.commit()
    conn.close()


def NukeCharacters():
    conn = Connect()
    cursor = conn.cursor()
    cursor.execute("DROP TABLE IF EXISTS characters")
    conn.commit()
    conn.close()


NukeCharacters()
CreateTable()
LoadExcelData()
SaveImages()
