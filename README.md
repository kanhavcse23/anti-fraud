Project employs a microservices architecture, primarily developed in Go, to ensure scalability, maintainability, and efficient handling of various components.

- Table of Contents:
    - Project Overview
    - Features
    - Architecture
    - Installation
    - Usage
    - Configuration
    - Testing
    - Contributing
    - License


- Project Overview:
    This project consists of multiple services, each responsible for specific functionalities such as account management, transaction processing, and operation handling.

- Features:
    - Account Management: Handles creation and retrieval of account details.
    - Transaction Management: Manages the creation and validation of financial transactions.
    - Operation Handling: Manages different operation types.
    - Mediator Service: Facilitates communication between different services to ensure a decoupled architecture.
    - Database Migrations: Includes scripts for setting up and migrating the database schema.

- Architecture:
    The system follows a microservices architecture with the following components:

    - Account Service: Manages account-related data.
    - Transaction Service: Handles transaction-related data.
    - Operation Service: Manages operation types.
    - Mediator Service: Acts as an intermediary to facilitate communication between services via Mediator Pattern.

    Each service is designed to be independent, promoting scalability and ease of maintenance.


- Automate Installation via run_linux.sh/run_macos.sh:
    - Prerequisites:
        - Ensure that your postgresql database is set up and accessible.
    - Build and Run Services:
        - "git clone https://github.com/kanhavcse23/anti-fraud.git" (Clone the Repository)
        - "cd anti-fraud"
        - For Linux OS:
            - "chmod +x run_linux.sh"
            - "source run_linux.sh"
        - For MacOS OS:
            - "chmod +x run_macos.sh"
            - "source run_macos.sh"

- Manual Installation:
    - To set up the Project System locally, follow these steps:

    - Clone the Repository:
        - "git clone https://github.com/kanhavcse23/anti-fraud.git"
        - "cd anti-fraud"

    - Set Up Environment Variables:
        - Create a .env file in the root directory and configure the necessary environment variables for database connections and service configurations.

    - Run Database Migrations:
        - Ensure that your database is set up and accessible. Then, run the migration scripts located in the database/migration/ directory to set up the necessary tables and schemas.

    - Build and Run Services:
        "docker-compose up --build"

- Usage:
    - Once all services are up and running, you can interact with them using API clients like Postman.

    - Account Service:
        - Create Account: POST /accounts, JSON BODY: {"document_number": <DOCUMENT_NUMBER>}
        - Get Account Details: GET /accounts/{accountId}

    - Transaction Service:
        - Create Transaction: POST /transactions, JSON BODY: {"account_id": <ACC_ID>, "operation_type_id": <OP_ID>, "amount": <AMOUNT>}

- Testing:
    To run tests for the services, use the following command:
        - "go test ./... -v"

- Database Configuration:
    Edit database configuration in following files:
    - config.yml
    - Dockerfile
    - docker-compose.yml