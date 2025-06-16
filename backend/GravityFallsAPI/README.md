# 🔮 GravityFallsAPI

**GravityFallsAPI** is a lightweight Python web application that provides access to character data from the *Gravity Falls* universe. Built with FastAPI and powered by PostgreSQL.

The API returns information about characters including their name, quote, interests, and image.

---

## 🚀 Technologies

* **Python 3**
* **FastAPI**
* **PostgreSQL**

---

## 📡 API Endpoints

### `/`

Returns a list of available routes:

```json
{
  "Characters": "/characters"
}
```

---

### `/characters`

Returns the total number of characters in the database:

```json
{
  "Count": 114
}
```

---

### `/characters/{id}`

Returns detailed information about a character by their `id`:

```json
{
  "id": 85,
  "name": "Mabel Pines",
  "image": "https://static.wikia.nocookie.net/gravityfalls/images/b/b2/S1e3_mabel_new_wax_figure.png/",
  "quote": "When life gives you lemons, draw faces on those lemons and wrap them in a blanket. Ta-daaa! Now you have lemon babies."
}
```

---

## 📌 Example Request

```bash
curl http://localhost:9090/characters/85
```

---

## 🔮 Roadmap

* 🧼 Pagination and filtering
* 🕵️ Episode details integration
* 🖼️ Fast image support (upload/URL)
* 🌐 Web hosting with visual API documentation

---

## 🛠 Installation

```bash
git clone https://github.com/bsagat/GravityFallsAPI.git
cd GravityFallsAPI
pip install fastapi
uvicorn main:app
```
