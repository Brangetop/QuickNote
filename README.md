# 🔒QuickNote
**QuickNote** is a web service designed for sharing notes securely using one-time links. With **end-to-end encryption**, your messages are protected from unauthorized access, ensuring that only the intended recipient can read them.

Key Features: 
* **One-Time Links**: Share notes that can only be accessed once. After reading, the note is removed from the database— *gone, reduced to atoms*. 
* **End-to-End Encryption**: Your notes and files are encrypted, providing maximum security and privacy. 
* **File Attachment**: Easily attach files to your messages, making it convenient to share documents alongside your notes.

With **QuickNote**, you can share sensitive information with confidence, knowing that your data is **secure and ephemeral**.
# 📝How does it work
<p align="center">
  <img width="1088" height="700" alt="canvas" src="https://github.com/user-attachments/assets/e337b7d7-beba-43fd-aac1-7dd0a32d1734" />
</p>
<p align="center">
  <em>Simple schematic of server-side code</em>
</p>


When a user creates a note, they fill out a form that includes the note content and can optionally upload a file. Upon submission, the **saveHandler function** is triggered. This function first checks if the request method is POST and validates the content length, ensuring it does not exceed 200 characters(example value). If the content is valid, it **generates a unique 10-character link** for the note.

If a file is uploaded, the handler creates a directory named after the generated link to store the file. It saves the uploaded file in this directory and records the creation time in a hidden file. After successfully saving the file, the note content is stored in the database along with the generated link.

When a user wants to view a note, the **viewHandler function** retrieves the note content using the unique link. After displaying the note, it **deletes the message** from the database to ensure that it is no longer accessible. This process maintains the privacy and ephemeral nature of the notes, as they are removed after being viewed.

<p align="center">
  <img width="700" height="700" alt="Encryption process" src="https://github.com/user-attachments/assets/ec9468f9-2322-45c8-ab67-b765e57cd284" />
</p>
<p align="center">
  <em>Simple schematic of e2e encryption</em>
</p>
This service ensures privacy by incorporating a unique key into the link format, e.g. example.com/ID#Key. This design guarantees that access to data is restricted, even from the service provider. By utilizing this method, QuickNote prioritize user confidentiality and data protection.

# ⚙Documentation
| File Name   | Description                                                                 |
|-------------|-----------------------------------------------------------------------------|
| main.go     | The main application file that initializes the server, handles routes, and manages message creation, viewing, and downloading. |
| create.html | HTML template for the message creation form, including a character counter and file upload option. |
| link.html   | HTML template displayed after saving a message, showing the unique link for retrieval. |
| view.html   | HTML template for displaying the saved message and any associated files, with a download option. |
| create.js   | JavaScript file for handling character count updates and toggling additional information in the create form. |
| db.go       | Database management file that initializes the connection to the MySQL database and provides functions to save, retrieve, and delete messages. |
| .env        | Environment file for storing database configuration variables (user, password, database name, and host). |

## Database Configuration
To **configure the database**, edit the .env file with the following parameters: 
1. **DB_USER**: Your MySQL username (e.g., root).
2. **DB_PASSWORD**: Your MySQL password (e.g., 1290).
3. **DB_NAME**: The name of the database to use (e.g., messages).
4. **DB_HOST**: The host where the MySQL server is running (e.g., localhost).
## Dependencies
The project uses the following indirect **dependencies**: 
1. filippo.io/edwards25519 v1.1.0: A cryptographic library for secure operations.
2. github.com/go-sql-driver/mysql v1.9.3: MySQL driver for Go, enabling database interactions.
3. github.com/joho/godotenv v1.5.1: A library for loading environment variables from a .env file.
## Running the Application
Ensure you have Go installed and set up on your machine. Install the required dependencies using go mod tidy. Create the MySQL database and the messages table with appropriate fields. Run the application using go run main.go. **Access the application at** http://localhost:8080.
# To-Do list
1. Fix SQL injection.
2. Finish end-to-end encryption.
3. Improve frontend(make it more usable and customizable).
4. Finish another project for server auto-cleanup.
