# **Stay Mate: A Hotel Management System** 🏨

## **Project Overview** 🌍
Stay Mate is a sophisticated software solution built to enhance and automate hotel management processes. Designed for hotel staff and management, it simplifies tasks like reservation management, customer tracking, and room availability. By integrating key features, Stay Mate aims to streamline daily operations, reduce manual workload, and increase efficiency.

## ✨ Key Features  

- **📅 Booking Management**:  
  Effortlessly create, update, view, and delete bookings.  

- **👥 User Management**:  
  Manage user information, including personal details and account data.  

- **🏨 Room Management**:  
  Add, update, view, and delete room details with ease.  

- **⚙️ Efficient Database Handling**:  
  Powered by PostgreSQL and GORM ORM for smooth and reliable database operations.  

- **📧 Email Functionality**:  
  - Send emails to users.  
  - Attach files and images effortlessly.  
  - Includes email verification for secure account setup.  

- **📊 Data Handling**:  
  - **Sorting**: Sort data based on various parameters for better insights.  
  - **Pagination**: Navigate large datasets with ease using pagination.  
  - **Filtering**: Apply filters to display only relevant data.  

- **🔐 Registration and Authorization**:  
  - User registration for new accounts.  
  - Secure login for existing users.  

- **👤 User Profile Management**:  
  Users can view and update their profile information effortlessly.
  
- **💬 Support Chat**:  
  Real-time chat between users and admin using WebSockets for seamless communication.

## **Team Members** 👥
- **[Temutjin Koszhanov](https://github.com/Temutjin2k)** (SE-2308) 👨🏻‍💻
- **[Aida Zhalgassova](https://github.com/zhxauda9)** (SE-2307) 👩🏻‍💻

## **Home Page** 🏠
![Image alt](https://github.com/zhxauda9/StayMate/blob/main/assets/home.png)

## **Login** 👤
![Image alt](https://github.com/zhxauda9/StayMate/blob/main/assets/login.png)

## **Sign Up** 👤
![Image alt](https://github.com/zhxauda9/StayMate/blob/main/assets/signup.png)

## **Verify** 🖥️
![Image alt](https://github.com/zhxauda9/StayMate/blob/main/assets/verify.png)

## **Profile** 👤
![Image alt](https://github.com/zhxauda9/StayMate/blob/main/assets/profile.png)

## **Admin Panel** 🛠️
### **Bookings** 📅
![Image alt](https://github.com/zhxauda9/StayMate/blob/main/assets/bookings.png)
### **Users** 👤
![Image alt](https://github.com/zhxauda9/StayMate/blob/main/assets/image.png)
### **Rooms** 🛏️
![Image alt](https://github.com/zhxauda9/StayMate/blob/main/assets/rooms.png)
### **Email sending** 📧
![Image alt](https://github.com/zhxauda9/StayMate/raw/main/assets/email.png)
### **Support chat** ❓
### User page
![image](https://github.com/user-attachments/assets/e1525677-6dca-4659-9f49-2f99012d576a)
### Admin page
![image](https://github.com/user-attachments/assets/866a9d41-3755-42fe-9a24-84724ff041a3)
### Admin can see active chats
![image](https://github.com/user-attachments/assets/52435e05-463d-46f4-b0e8-ef8dcc6963c0)

## **Technologies Used** 🛠️
- **Programming Language**: Go (Golang) 1.22.3 🖥️
- **Database**: PostgreSQL 🗃️
- **Version Control**: Git 🧑‍💻
- **Libraries**:
    - **[GORM](https://github.com/go-gorm/gorm)**: For ORM (Object-Relational Mapping) 📦
    - **[Zerolog](https://github.com/rs/zerolog)**: For structured logging 📜
    - **[go-mail/mail/v2](https://github.com/go-gomail/gomail)**: For sending emails 📧  
    - **[x/time/rate](https://pkg.go.dev/golang.org/x/time/rate)**: For rate-limiting ⏱️
    - **[gorilla/websocket](https://pkg.go.dev/github.com/gorilla/websocket@v1.5.3)**: For support chat ❓

## **How it Works** 🔄
Stay Mate integrates a robust system for handling the key functions of a hotel. By using GORM for interacting with the PostgreSQL database, it allows efficient data management. The clean and simple API design makes it easy to interact with and manage hotel resources.

### **Booking Process** 📲
1. Users can create, update, or delete bookings.
2. The system automatically ensures room availability before confirming bookings.
3. Staff can review all bookings to manage the hotel's schedule.

### **Room Management** 🏨
Rooms can be added, updated, or removed from the system. Each room is linked to specific details like room type, price, and availability.

### **User Management** 👥
User data such as name, email, and booking history can be managed. The system helps in tracking customer interactions and improves the overall customer experience.


## How to Run the Project

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/zhxauda9/StayMate.git
   cd StayMate
   ```
2. **Define enviromental variables. (You can use .env file)**:
   ```bash
    # Variables for the database
    DB_HOST=
    DB_PORT=
    DB_USER=
    DB_PASSWORD=
    DB_NAME=

    # Variables for email sending
    SMTP_HOST=
    SMTP_PORT=
    EMAIL=
    PASSWORD=
   ```
3. **Build the programm from the root directory**:
   ```bash
   go build -o main cmd/main.go
   ```

4. **Start the Server**:
   ```bash
   ./main
   ```

## **Contribution** 📝
If you want to contribute to the project, feel free to fork the repository, make changes, and create pull requests. We appreciate all contributions!

---

For questions or suggestions, feel free to reach out to us! We’d love to hear from you. 💬

---
Stay Mate: Simplifying hotel management for you. 🌟
