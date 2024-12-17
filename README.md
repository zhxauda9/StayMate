# **Stay Mate: A Hotel Management System** 🏨

## **Project Overview** 🌍
Stay Mate is a sophisticated software solution built to enhance and automate hotel management processes. Designed for hotel staff and management, it simplifies tasks like reservation management, customer tracking, and room availability. By integrating key features, Stay Mate aims to streamline daily operations, reduce manual workload, and increase efficiency.

### **Features** 🚀
- **Booking Management**: Create, update, view, and delete bookings easily.
- **User Management**: Manage users, including their personal information and details.
- **Room Management**: Add, update, view, and delete room details.
- **Efficient Database Handling**: Uses PostgreSQL with GORM ORM for seamless database interactions.

## **Team Members** 👥
- **Temutjin Koszhanov** (SE-2308) 👨🏻‍💻
- **Aida Zhalgassova** (SE-2307) 👩🏻‍💻

## **API Endpoints** 📡

### **Bookings** 📅

| Method | Endpoint          | Description                          | Response                |
|--------|-------------------|--------------------------------------|-------------------------|
| **POST**   | `/bookings`       | Creates a new booking                | 201 Created             |
| **GET**    | `/bookings`       | Retrieves all bookings               | 200 OK                  |
| **GET**    | `/bookings/{id}`  | Retrieves a specific booking by ID   | 200 OK                  |
| **PUT**    | `/bookings/{id}`  | Updates an existing booking          | 200 OK                  |
| **DELETE** | `/bookings/{id}`  | Deletes a booking                    | 204 No Content          |

### **Users** 👤

| Method | Endpoint          | Description                          | Response                |
|--------|-------------------|--------------------------------------|-------------------------|
| **POST**   | `/user`           | Adds a new user                      | 201 Created             |
| **GET**    | `/user`           | Retrieves all users                  | 200 OK                  |
| **GET**    | `/user/{id}`      | Retrieves a specific user by ID      | 200 OK                  |
| **PUT**    | `/user/{id}`      | Updates an existing user             | 200 OK                  |
| **DELETE** | `/user/{id}`      | Deletes a user                       | 204 No Content          |

### **Rooms** 🛏️

| Method | Endpoint          | Description                          | Response                |
|--------|-------------------|--------------------------------------|-------------------------|
| **POST**   | `/rooms`          | Adds a new room                      | 201 Created             |
| **GET**    | `/rooms`          | Retrieves all rooms                  | 200 OK                  |
| **GET**    | `/rooms/{id}`     | Retrieves a specific room by ID      | 200 OK                  |
| **PUT**    | `/rooms/{id}`     | Updates an existing room             | 200 OK                  |
| **DELETE** | `/rooms/{id}`     | Deletes a room                       | 204 No Content          |

## **CRUD Operations** 🛠️

### **Bookings CRUD** 📅
- **Create**: `POST /bookings`  
  Allows you to create a new booking with details like user ID, room ID, check-in and check-out dates.
- **Read**:
    - `GET /bookings`  
      Retrieves a list of all bookings.
    - `GET /bookings/{id}`  
      Retrieves a specific booking by ID.
- **Update**: `PUT /bookings/{id}`  
  Allows you to update an existing booking, including modifying user, room, and date details.
- **Delete**: `DELETE /bookings/{id}`  
  Deletes an existing booking by ID.

### **Users CRUD** 👤
- **Create**: `POST /user`  
  Allows you to add a new user with details like name, email, and other personal information.
- **Read**:
    - `GET /user`  
      Retrieves a list of all users.
    - `GET /user/{id}`  
      Retrieves a specific user by ID.
- **Update**: `PUT /user/{id}`  
  Allows you to update an existing user’s details, such as email, name, or other information.
- **Delete**: `DELETE /user/{id}`  
  Deletes a user by ID.

### **Rooms CRUD** 🛏️
- **Create**: `POST /rooms`  
  Allows you to add a new room with details like room type, price, and availability.
- **Read**:
    - `GET /rooms`  
      Retrieves a list of all rooms.
    - `GET /rooms/{id}`  
      Retrieves a specific room by ID.
- **Update**: `PUT /rooms/{id}`  
  Allows you to update an existing room’s details, such as price, availability, and room type.
- **Delete**: `DELETE /rooms/{id}`  
  Deletes a room by ID.

## **Technologies Used** 🛠️
- **Programming Language**: Go (Golang) 1.22.3 🖥️
- **Database**: PostgreSQL 🗃️
- **Version Control**: GitHub 🧑‍💻
- **Libraries**:
    - **GORM**: For ORM (Object-Relational Mapping) 📦
    - **Zerolog**: For structured logging 📜

## **How it Works** 🔄
Stay Mate integrates a robust system for handling the key functions of a hotel. By using GORM for interacting with the PostgreSQL database, it allows efficient data management for users, rooms, and bookings. The clean and simple API design makes it easy to interact with and manage hotel resources.

### **Booking Process** 📲
1. Users can create, update, or delete bookings.
2. The system automatically ensures room availability before confirming bookings.
3. Staff can review all bookings to manage the hotel's schedule.

### **Room Management** 🏨
Rooms can be added, updated, or removed from the system. Each room is linked to specific details like room type, price, and availability.

### **User Management** 👥
User data such as name, email, and booking history can be managed. The system helps in tracking customer interactions and improves the overall customer experience.

## **Contribution** 📝
If you want to contribute to the project, feel free to fork the repository, make changes, and create pull requests. We appreciate all contributions!

---

For questions or suggestions, feel free to reach out to us! We’d love to hear from you. 💬

---
Stay Mate: Simplifying hotel management with smart technology. 🌟
