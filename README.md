# OAuth Implementation Example using Golang Fiber

A simple example demonstrating the implementation of OAuth2 for user authentication in a website using Golang Fiber, with the help of the Goth module. This example showcases the use of Google and GitHub as authentication providers.

## Getting Started

### 1. Configure Environment Variables
Create a `.env` file in the root directory of the project and fill it with your configuration settings with basic values from `.example.env`.

### 2. Install Dependencies
Run the following command to ensure all necessary modules are installed:

```bash
go mod tidy
```

### 3. Start the Development Server
To start the development server, run:

```bash
go run main.go
```

This will start the server and automatically load changes when you rerun the command after making changes.

Alternatively, if you prefer using air for live reloading during development, simply run:

```bash
air
```

Make sure to configure air according to your project's needs by adjusting the settings in the .air.toml file.

### 4. Start the Production Server
To start the server in production mode, you can build the binary and run it:

#### On Windows:
```bash
go build -o lorem.exe
lorem.exe
```

#### On Linux/macOS:
```bash
go build -o lorem
./lorem
```
