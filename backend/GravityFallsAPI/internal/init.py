from dotenv import load_dotenv
import pandas as pd
import psycopg2
import os

load_dotenv()
EXCEL_PATH = os.getenv("EXCEL_PATH")
DATABASE_URL = os.getenv("DATABASE_URL")


def Connect():
    conn = psycopg2.connect(DATABASE_URL)
    return conn


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

    file = pd.read_excel(EXCEL_PATH)

    for _, row in file.iterrows():
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