# **InfluxUI-UG10 Group Project**

## **Setup Instructions**

### **Prerequisites**

Ensure you have the following installed and configured:

- **Windows Subsystem for Linux (WSL)**
- **Visual Studio Code (VSC)**
- **Docker** (Ensure the Windows Docker is running and configured for WSL)
- **Python**
- **Go (v1.23.1)**
- **Node Version Manager (NVM)**
- **Docker inside WSL**

## **Starting the Server**

1. **Build the frontend:**
   ```sh
   cd frontend
   npm run build
   ```
2. **Run the backend server:**
   ```sh
   cd backend
   go run .
   ```

---

## **Making Changes**

### **Frontend Changes (While the Server is Running)**

1. **Rebuild the frontend:**
   ```sh
   npm run build
   ```
2. **Refresh the page to see the changes.**

### **Backend Changes (While the Server is Running)**

1. **Stop the current server:**
   - Use `Ctrl+C` (DO NOT use `Ctrl+Z`).
2. **Restart the backend server:**
   ```sh
   go run .
   ```

---

## **Grafana Setup**

If you are running a container that has **Grafana**, follow these steps:

### **1. Open Grafana**

- Navigate to: **[http://localhost:3000](http://localhost:3000)**
- Sign in with the credentials:
  ```
  Username: admin
  Password: adminpassword
  ```

### **2. Create a Data Source**

1. Click **"Add your first data source"**.
2. Select **InfluxDB**.
3. Change the query language to **Flux**.
4. Set the URL to:
   ```
   http://localhost:8086
   ```
5. Fill in the **auth details** using your InfluxDB admin user:
   ```
   Username: admin
   Password: <your-password>
   ```
6. Configure the InfluxDB connection:
   - **Organization**: `admin`
   - **Token**: Use the Influx API token from the `.env` file.
   - **Default Bucket**: `default`
7. Click **"Save & Test"**.
8. Retrieve the **Datasource UID** from the URL:
   - Example: If the URL is:
     ```
     http://localhost:3000/datasources/edit/f6ddea66-df79-48e1-a050-dfa6ed19ef0e
     ```
     Then add this to your `.env` file:
     ```
     GRAFANA_INFLUXDB_DATASOURCE_UID="f6ddea66-df79-48e1-a050-dfa6ed19ef0e"
     ```

### **3. Set Up an API Token**

1. Go to **Settings → Administration → Service Accounts**.
2. Create two service accounts:
   - **Viewer** (`role: viewer`)
   - **Main** (`role: editor`)
3. Save the **editor API token** in your `.env` file:
   ```
   GRAFANA_API_TOKEN="glsa_bjViL3WpMiBM40rNEfsj1Qu7bmIAKCsh_2dad91d3"
   ```
4. Save the **Grafana URL** in the `.env` file:
   ```
   GRAFANA_URL="http://grafana:3000"
   ```

### **4. Create a Dashboard**

1. Navigate to the Grafana homepage.
2. Click **"Create your first dashboard"**.
3. Click **"Add Visualization"**.
4. Configure the query and click **"Apply"**.
5. Click **"Save"**, give it a name, and retrieve the **Dashboard UID**:
   - Example: If the URL is:
     ```
     http://localhost:3000/d/d97f7b04-0c2b-42bf-8c51-6e83373d4db3/a-brand-new-dashboard?orgId=1
     ```
     Then store this in your `.env` file:
     ```
     GRAFANA_DASHBOARD_UID="d97f7b04-0c2b-42bf-8c51-6e83373d4db3"
     ```
6. Update the **Grafana test script** to use this UID and verify the script runs successfully.

---

## **Troubleshooting**

### **Docker Issues**

- Ensure Docker is running and configured for WSL.
- Restart Docker if you experience build issues.

### **Backend or Frontend Failing to Start**

- Run `npm install` again inside `/frontend`.
- Run `go run .` inside `/backend`.

---

## **License**

This project is for **educational purposes** and is part of the **InfluxUI-UG10** group project.**fluxUI-UG10 Group Project**

### 🎉 **Enjoy working on InfluxUI-UG10!** 🚀

