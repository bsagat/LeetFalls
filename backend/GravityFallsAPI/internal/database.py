from internal.init import Connect

def CharacterByID(id: int):
    with Connect() as conn:
        cursor = conn.cursor()
        cursor.execute('''
            SELECT 
                Id, Name, Species, COALESCE(Likes, '') as Likes, COALESCE(Quote, '') as Quote, Image
            FROM 
                characters 
            WHERE 
                Id = %s
        ''', (id,))
        data = cursor.fetchone()
    return data


def CharacterCount() -> int:
    with Connect() as conn:
        cursor = conn.cursor()
        cursor.execute("SELECT COUNT(*) FROM characters")
        count = cursor.fetchone()[0]
    return count
