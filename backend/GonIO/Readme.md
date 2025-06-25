# ☁️ GonIO — Simplified S3 Storage

**GonIO** is a lightweight alternative to Amazon S3, designed for learning and experimentation. It provides a RESTful API to manage **buckets** and **objects**, including file upload, retrieval, deletion, and metadata storage in a simple way.

---

## 🔧 Features

- ✅ Create virtual containers (buckets)
- ✅ Upload and retrieve objects
- ✅ Store metadata in CSV format
- ✅ Flexible configuration via `.env`, command-line arguments, or default values
- ✅ Support for uploading ZIP archives
- ✅ Name validation rules
- ✅ Proper HTTP status codes
- ✅ Content-Type and Content-Length support

---

## 📦 Buckets

### ✅ Create a Bucket
- **Method:** `PUT`
- **Endpoint:** `/buckets/{bucket-name}`
- **Body:** empty
- **Constraints:** bucket name must be 3–63 characters, lowercase, may include numbers, hyphens, and dots

### 📄 List All Buckets
- **Method:** `GET`
- **Endpoint:** `/buckets`

### ❌ Delete a Bucket
- **Method:** `DELETE`
- **Endpoint:** `/buckets/{bucket-name}`

---

## 🗂️ Objects

### 📤 Upload an Object
- **Method:** `PUT`
- **Endpoint:** `/objects/{bucket-name}/{object-key}`
- **Body:** binary data of the file
- **Headers:**
  - `Content-Type`: MIME type (e.g. `image/png`)
  - `Content-Length`: file size in bytes

### 📄 List All Objects in a Bucket
- **Method:** `GET`
- **Endpoint:** `/objects/{bucket-name}`

### 📥 Retrieve an Object
- **Method:** `GET`
- **Endpoint:** `/objects/{bucket-name}/{object-key}`

### 🗑️ Delete an Object
- **Method:** `DELETE`
- **Endpoint:** `/objects/{bucket-name}/{object-key}`

---

## 📦 Upload ZIP Archive

Allows uploading a ZIP archive with images into a bucket.

### 📤 Upload ZIP
- **Method:** `POST`
- **Endpoint:** `/{bucket-name}/upload-zip`
- **Headers:**
  - `Content-Type`: `application/zip`
- **Body:** A ZIP file containing image files

### Requirements:
- Only image files allowed inside the archive (`.jpg`, `.jpeg`, `.png`)
- Files are automatically extracted and stored as regular objects

---

## ⚙️ Configuration

Configuration is loaded in the following order of priority:

1. `.env` file
2. Command-line arguments (`--port`, `--host`, `--dir`)
3. Default values

### Example `.env` file:
```env
PORT=9090
HOST=localhost
BUCKETPATH=data
