# **Azure Container Instances Job Manager**

## **Environment Variables**

The RESTful service requires a few environment variable configurations to be able to run.

### **Logging**

<table>
    <hr>
        <th>Variable Name</th>
        <th>Required</th>
        <th>Details</th>
    </hr>
    <tr>
        <td>LOGGER_LOG_LEVEL</td>
        <td>true</td>
        <td>
            <ul>
                <li>INFO</li>
                <li>DEBUG</li>
                <li>ERROR</li>
                <li>WARNING</li>
            </ul>
        </td>
    </tr>
    <tr>
        <td>LOGGER_FILE_PATH</td>
        <td>true</td>
        <td>
            Path to where the JSON log file is stored.
        </td>
    </tr>
</table>

### **GIN WEB FRAMEWORK**

<table>
    <hr>
        <th>Variable Name</th>
        <th>Required</th>
        <th>Details</th>
    </hr>
    <tr>
        <td>GIN_MODE</td>
        <td>true</td>
        <td>
            <ul>
                <li>release</li>
                <li>debug</li>
            </ul>
        </td>
    </tr>
    <tr>
        <td>GIN_PORT</td>
        <td>true</td>
        <td>Port number for HTTP traffic</td>
    </tr>
</table>

### **Azure Resource Manager Authentication**

<table>
    <hr>
        <th>Variable Name</th>
        <th>Required</th>
        <th>Details</th>
    </hr>
    <tr>
        <td>AZURE_AUTH_TYPE</td>
        <td>true</td>
        <td>
            <ul>
                <li>USER_ASSIGNED_MANAGED_IDENTITY</li>
                <li>SERVICE_PRINCIPAL</li>
            </ul>
        </td>
    </tr>
    <tr>
        <td>AZURE_AUTH_CLIENT_ID</td>
        <td>false</td>
        <td>
            Required when using either USER_ASSIGNED_MANAGED_IDENTITY
            or SERVICE_PRINCIPAL.
        </td>
    </tr>
    <tr>
        <td>AZURE_AUTH_CLIENT_SECRET</td>
        <td>false</td>
        <td>
            Required when using either SERVICE_PRINCIPAL.
        </td>
    </tr>
    <tr>
        <td>AZURE_AUTH_TENANT_ID</td>
        <td>false</td>
        <td>
            Required when using either SERVICE_PRINCIPAL.
        </td>
    </tr>
</table>