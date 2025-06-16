# ☁️ GonIO — Simplified S3 Storage

**GonIO** is a lightweight, simplified version of Amazon S3 (Simple Storage Service), designed for learning and experimentation. It provides a RESTful API that lets you manage **buckets** and **objects** — including creating, uploading, retrieving, and deleting files, as well as storing metadata.

---

## 🔧   Features

Imagine a minimal web service where you can:

- Create virtual containers (buckets)
- Store and retrieve files (objects)
- Access everything via simple HTTP calls

While commercial S3 services are highly scalable and complex, **GonIO** focuses on the **core concepts** of object storage.

---

## 📦 Buckets

Buckets are like folders or containers for your files. Here's how to manage them:

### ✅ Create a Bucket
- **Method:** `PUT`
- **Endpoint:** `/{bucket-name}`
- **Request Body:** _Empty_
- **Constraints:** Bucket names must be 3–63 characters, lowercase, and can contain numbers, hyphens, and periods.

### 📄 List All Buckets
- **Method:** `GET`
- **Endpoint:** `/`

### ❌ Delete a Bucket
- **Method:** `DELETE`
- **Endpoint:** `/{bucket-name}`

---

## 🗂️ Objects

Objects are the actual files stored inside buckets, along with metadata like content type and size.

### 📤 Upload an Object
- **Method:** `PUT`
- **Endpoint:** `/{bucket-name}/{object-key}`
- **Body:** Binary data of the object
- **Headers:**
  - `Content-Type`: MIME type (e.g., `image/png`)
  - `Content-Length`: Size in bytes

### 📄 List All Objects 
- **Method:** `GET`
- **Endpoint:** `/{bucket-name}`

### 📥 Retrieve an Object
- **Method:** `GET`
- **Endpoint:** `/{bucket-name}/{object-key}`

### 🗑️ Delete an Object
- **Method:** `DELETE`
- **Endpoint:** `/{bucket-name}/{object-key}`

---

## 🛠️ Usage

Start the server using:

```bash
$ ./gonIO 