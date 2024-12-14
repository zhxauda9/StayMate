**Stay Mate. A Hotel Management System**

## Project Overview
Stay Mate is a software solution designed to improve and automate hotel management processes. By integrating key features such as reservation management, customer tracking, and room availability, the system aims to simplify daily operations and enhance efficiency for hotel staff and management.

## Team Members
- **Temutjin Koszhanov**
- **Aida Zhalgasova**


## How to Run the Project

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-repository-url/stalemate.git
   cd stalemate
   ```

2. **Build the programm**:
   ```bash
   go build -o main cmd/main.go
   ```

3. **Start the Server**:
   ```bash
   ./main
   ```

### Now your server started on localhost on port 8080

### API Endpoints

#### Bookings

| Method | Endpoint          | Description                       | Response                   |
|--------|-------------------|-----------------------------------|----------------------------|
| POST   | `/bookings`         | Creates a new booking.             | 201 Created                |
| GET    | `/bookings`         | Retrieves all bookings.             | 200 OK                     |
| GET    | `/bookings/{id}`    | Retrieves a specific order by ID. | 200 OK     |
| PUT    | `/bookings/{id}`    | Updates an existing order.        | 200 OK     |
| DELETE | `/bookings/{id}`    | Deletes an order.                 | 204 No Content  |

#### Rooms

| Method | Endpoint          | Description                        | Response                   |
|--------|-------------------|------------------------------------|----------------------------|
| POST   | `/rooms`           | Adds a new rooms item.             | 201 Created                |
| GET    | `/rooms`           | Retrieves all rooms items.          | 200 OK                     |
| GET    | `/rooms/{id}`      | Retrieves a specific rooms item.    | 200 OK     |
| PUT    | `/rooms/{id}`      | Updates an existing rooms item.     | 200 OK     |
| DELETE | `/rooms/{id}`      | Deletes a rooms item.               | 204 No Content  |

## Tools and Resources Used
- **Programming Language**: Golang 1.22.3
- **Database**: PostgreSQL
- **Version Control**: Git
- **Libraries**:
  - `gorm` for ORM

---
Feel free to reach out to us for any questions or suggestions!
