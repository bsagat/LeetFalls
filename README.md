# LeetFalls - Gravity Falls Anonymous Imageboard

![LeetFalls Logo](/backend/app/web/images/wallpaper.gif)

> "Stay curious, stay weird, stay kind, and don't let anyone tell you you're not smart enough, brave enough, or worthy enough."
>
> — Alex Hirsch, "Gravity Falls: Journal 3"

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/HTML)
[![CSS3](https://img.shields.io/badge/CSS3-1572B6?style=for-the-badge&logo=css3&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/CSS)
[![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)](https://swagger.io/)
[![REST API](https://img.shields.io/badge/API-REST-007ACC?style=for-the-badge)](https://en.wikipedia.org/wiki/Representational_state_transfer)
[![Hexagonal Architecture](https://img.shields.io/badge/Architecture-Hexagonal-blueviolet?style=for-the-badge)](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)

---

## Table of Contents

* [About the Project](#-about-the-project)
    * [Posts](#-posts)
    * [Comments](#-comments)
    * [Application Screenshots](#-application-screenshots)
    * [Project Struct](#-project-structure)
* [Technical Highlights](#️-technical-highlights)
* [How to Build](#️-how-to-build)
* [API Documentation](#-api-documentation)
* [License](#-license)
* [Future Plans](#-future-plans)
* [Contacts](#-our-contacts)

---

## 🌲 About The Project

Welcome to **LeetFalls**, an anonymous imageboard nestled in the strange and wonderful world of Gravity Falls. Built with Go and clean architecture, it’s your place to share odd discoveries, cryptic theories, and hidden truths without revealing your identity — just like Dipper and Mabel would want.

-   📝 **Post Creation:** Share your thoughts, stories, and images in new posts with title, content, and optional image attachments.
-   💬 **Threaded Discussions:** Reply to posts or comments to create layered, tree-like conversations, perfect for unraveling mysteries together.
-   🍪 **Anonymous Sessions:** Stay truly anonymous with cookie-based user sessions.
-   🖼️ **Image Storage:** Seamless integration with an **S3-compatible storage** for uploading and storing images attached to posts and comments.
-   🗄️ **Automatic Archiving:** Inactive posts quietly fade into the archive like a memory in the woods.
-   🧙‍♂️ **GravityFallsAPI Integration:** Provides **automatic assignment of unique Gravity Falls character avatars and names** to new users.
-   ✏️ **Dynamic User Names:** Users can specify a **custom name**, which dynamically replaces their API-provided persona name in all their messages and comments.

---

## 📝 Posts

Posts are the core content of the imageboard. Each post consists of:

* A **title**, **text content**, and a **unique post ID**.
* An **optional image** attached to the post.
* **User avatar** and **name**. If the user has specified a custom name, it dynamically replaces the API-provided Gravity Falls persona name in all their messages and comments.

---

### 🗑️ Post Deletion Rules

To maintain a dynamic and clean catalog, posts are subject to automated deletion rules from the Main Page (Catalog Page):

* Posts with **no new comments for 15 minutes** are automatically deleted.
* New posts **without any comments** are deleted after **10 minutes**.

---

## 💬 Comments

Users can engage in discussions by adding comments to posts. The comment system supports:

* **Replies:** Comments can reply to the original post or to other comments, creating threaded discussions.
* **Contextual Replies:** Each comment explicitly indicates whether it is a reply to a post or another comment.
* **Interactive Replies:** Users can click on the ID of a post or comment to easily initiate a reply to it.
* **User Identity:** Comments display the comment's unique ID and the user's assigned avatar.

---

## 📸 Application Screenshots

| Page | Screenshot |
|------|------------|
| Catalog | ![Catalog Page Screenshot](/docs/images/CatalogPage.png) |
| Post | ![Post Page Screenshot](/docs/images/CatalogPost.png) |
| Post Creation | ![Post Creation Page Screenshot](/docs/images/CreatePostPage.jpg) |
| Archive | ![Archive Page Screenshot](/docs/images/ArchivePage.png)  |
| Error Page | ![Error Page Screenshot](/docs/images/ErrorPage.png) |
| Swagger UI | ![Swagger UI](/docs/images/swaggerPage.png) |

---

## 🛠️ Technical Highlights

-   🚀 **Go-Powered Backend:** Developed with Go for high performance, concurrency, and reliability.
-   🐳 **Docker Deployment:** Utilizes Docker for easy containerization and consistent deployment across environments.
-   📐 **Hexagonal Architecture:** Ensures clean code, testability, and flexible dependency management.
-   🔗 **RESTful API:** Provides a robust interface, fully documented with **Swagger** for easy integration.
-   🗄️ **PostgreSQL Database:** Reliable and scalable data persistence for all application content.
-   ☁️ **GonIO/MinIO Image Storage:** Efficient S3-compatible object storage for all uploaded images.
-   🌲 **Gravity Falls Character API:** Seamless integration for dynamic and unique user persona assignment.

---

## ⚙️ How to Build

The project can be easily built and run using `make` and Docker Compose.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/bsagat/LeetFalls
    cd LeetFalls
    ```

2.  **Configure environment variables:**
    Before building, you need to create an `.env` file in the project's root directory. Copy the content from `user_friendly.env` and modify values as needed:

    ```bash
    cp user_friendly.env .env 
    ```

3.  **Run the application (Build and Start):**
    Ensure you have Docker and Docker Compose installed.
    ```bash
    make up
    ```
    This command will build Docker images and start all necessary services (LeetFalls application, PostgreSQL database, GonIO storage, and Gravity Falls API server).

4.  **Service stop:**
    ```bash
    make down
    ```



## 📂 Project Structure

The project follows a modular and clean architecture, organized as follows:

```bash
.
├── backend/                   # Main backend services
│   ├── app/                   # LeetFalls Go application
│   │   ├── cmd/               # Application entry points (e.g., main.go)
│   │   ├── Dockerfile         # Dockerfile for LeetFalls app
│   │   ├── init.sql           # Database initialization script
│   │   ├── internal/          # Internal application logic
│   │   │   ├── adapters/      # Implementations of ports (DB, external APIs, storage)
│   │   │   │   ├── dbRepo/    # Database repository implementations (PostgreSQL)
│   │   │   │   ├── external/  # External API adapters
│   │   │   │   ├── handlers/  # HTTP handlers for API endpoints and web pages
│   │   │   │   └── storage/   # Image storage adapter (MinIO/GonIO)
│   │   │   ├── app/           # Application setup and configuration
│   │   │   ├── domain/        # Core business logic, models, and interfaces (ports)
│   │   │   └── service/       # Business logic implementations
│   │   ├── logs/              # Application logs
│   │   └── web/               # Frontend static files (HTML templates, CSS, images)
│   ├── GonIO/                 # S3-compatible image storage service
│   └── GravityFallsAPI/       # Python API for character data
├── docker-compose.yaml        # Orchestrates all services
├── docs/                      # Project documentation (e.g., images for README)
│   └── images/
├── README.md                  # README file
└── user_friendly.env          # Example environment variables file

```

---

## 📖 API Documentation

Interactive API documentation is provided through **Swagger UI**, making it easy to explore and test the API endpoints directly from your browser.

After starting the application, you can access the documentation at the following address:

* **`http://localhost:8080/docs/swagger`** (The port might differ if you've changed your application's configuration.)

---

## 🔮 Future Plans

**Note: All features listed below are currently planned and not yet implemented.**

| Feature | Description | Status |
|---|---|---|
| 👤 **Profile Page** | Introduce a dedicated page for users to **customize their avatar and name**, and view their **author-specific post list**. | (Planned) |
| ⭐ **Favorite Posts** | Implement a feature allowing users to **mark and easily access their favorite posts**. | (Planned) |
| 💬 **Real-time Chat** | Introduce an **anonymous real-time chat** feature for direct communication between users. | (Planned) |
| 🔍 **Search Functionality** | Add a robust search bar to quickly **find posts by keywords, titles, or specific content**. | (Planned) |
| 🏷️ **Tagging System** | Implement a system for posts (e.g., `#mystery`, `#billcipher`) to **improve content discoverability and organization**. | (Planned) |
| 👍 **Reactions/Voting** | Explore adding simple reactions or a voting system (e.g., upvote/downvote, or themed reactions) to posts and comments. | (Planned) |
| 📜 **Advanced Archiving Options** | Provide more granular control or viewing options for archived posts. | (Planned) |
| 🎨 **Themed UI Customization** | Allow users to **choose different Gravity Falls themed UI elements or color schemes**. | (Planned) |
| ⚡ **Real-time Updates** | Enhance the catalog and comment sections with **real-time updates** for new posts and comments without page refresh. | (Planned) |
| 🛡️ **Basic Moderation Tools** | Consider implementing light moderation features (e.g., reporting mechanism) to maintain a healthy community environment. | (Planned) |

---

## 📄 License

This project is distributed under the **MIT License**. For full details, please refer to the [LICENSE](LICENSE) file in this repository.

---

## 📞 Our Contacts

If you have any questions, suggestions, or just want to connect, feel free to reach out to the project contributors:

### 👤 Bsagat

* **GitHub:** [GitHub Profile](https://github.com/bsagat)
* **Email:** [sagatbekbolat854@gmail.com](mailto:sagatbekbolat854@gmail.com)
* **LinkedIn:** [LinkedIn](https://www.linkedin.com/in/bekbolat-sagat/)

### 👤 Mrakhimo

* **GitHub:** [GitHub](https://github.com/zefirkaZirael)
* **Email:** [mansur.cor.tion@gmail.com](mailto:mansur.cor.tion@gmail.com)
* **LinkedIn:** [LinkedIn](https://www.linkedin.com/in/mansur-rakhimov/)
