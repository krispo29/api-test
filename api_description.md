## API Endpoint and Request Method

This API is designed to compare data from an uploaded Excel file with data stored in a database.

**HTTP Method:** The typical HTTP method for this type of request, where a file is being sent to the server for processing, is `POST`.

**Endpoint Structure:** A likely and intuitive endpoint structure would be something like:

`/compare/excel`

This clearly indicates the action (compare) and the type of data being sent (excel).

**Request Format:** Since an Excel file is being uploaded, the request format will be `multipart/form-data`. This format allows for sending files along with other data if needed. The request will contain the Excel file itself, typically under a field name like `excel_file`.

**Example Request:**

While the exact implementation can vary, a conceptual example of how a client might make this request (using a tool like `curl`) would be:

```bash
curl -X POST -F "excel_file=@/path/to/your/file.xlsx" http://your-api-domain.com/compare/excel
```

This example demonstrates:
- `-X POST`: Specifies the POST HTTP method.
- `-F "excel_file=@/path/to/your/file.xlsx"`:  Indicates a form field named `excel_file` and provides the path to the local Excel file to be uploaded. The `@` symbol tells `curl` to treat the following string as a file path.
- `http://your-api-domain.com/compare/excel`: The URL of the API endpoint.

**Summary:**

*   **Purpose:** Compare Excel data with a database.
*   **Method:** `POST`
*   **Endpoint:** `/compare/excel` (example)
*   **Request Body Type:** `multipart/form-data` (containing the Excel file)

## Request Payload

The request payload is sent as `multipart/form-data` and primarily consists of the Excel file.

**1. Excel File:**

*   **Key:** `excel_file` (common convention)
*   **Type:** File
*   **Description:** The Excel file (.xlsx, .xls, .csv) containing the data to be compared. The API should be designed to handle common Excel formats and potentially CSV files as well, given their prevalence in data exchange. The internal structure of the Excel (which sheet to read, header row, data rows) might be based on predefined conventions or configurable via other parameters.

**2. Optional Parameters:**

Depending on the API's flexibility, additional form fields might be included to control the comparison process:

*   **Key:** `db_table_name`
    *   **Type:** String
    *   **Description:** Specifies the database table against which the Excel data should be compared. This is crucial if the API can interact with multiple tables.
    *   **Example:** `products`, `sales_records`

*   **Key:** `comparison_columns_excel`
    *   **Type:** JSON string or comma-separated string
    *   **Description:** Specifies which columns from the Excel file to use for comparison. This could be a list of column names or indices.
    *   **Example:** `["ProductID", "ProductName", "Price"]` or `"ID,Name,Cost"`

*   **Key:** `comparison_columns_db`
    *   **Type:** JSON string or comma-separated string
    *   **Description:** Specifies the corresponding columns in the database table to compare against. The order or names should map to `comparison_columns_excel`.
    *   **Example:** `["item_id", "item_name", "unit_price"]` or `"product_identifier,description,value"`

*   **Key:** `sheet_name` (or `sheet_index`)
    *   **Type:** String (or Integer)
    *   **Description:** If the Excel file has multiple sheets, this parameter specifies which sheet to process.
    *   **Example:** `"Sheet1"` or `0`

*   **Key:** `header_row_index`
    *   **Type:** Integer
    *   **Description:** The index of the row in the Excel sheet that contains the column headers. Defaults to `0` (the first row).
    *   **Example:** `1`

**Example `curl` with Optional Parameters:**

```bash
curl -X POST \
  -F "excel_file=@/path/to/your/file.xlsx" \
  -F "db_table_name=products" \
  -F "comparison_columns_excel=[\"ID\", \"Stock\"]" \
  -F "comparison_columns_db=[\"product_id\", \"quantity\"]" \
  -F "sheet_name=Inventory" \
  http://your-api-domain.com/compare/excel
```

The inclusion of these optional parameters makes the API more versatile and configurable by the client. If not provided, the API might rely on default configurations or attempt to infer settings (e.g., use the first sheet, assume the first row is headers).

## Success Response

If the API processes the request successfully and performs the comparison, it should return a JSON response detailing the findings.

*   **HTTP Status Code:** `200 OK`

*   **Response Body (JSON):**

    The JSON response should be structured to provide a clear summary of the comparison.

    ```json
    {
      "summary": {
        "total_excel_rows": 1500,
        "total_db_rows_queried": 1450, // Based on matching criteria if any
        "matched_rows": 1200,
        "mismatched_rows": 50,
        "excel_only_rows": 250, // Rows in Excel not found in DB based on key columns
        "db_only_rows": 200     // Rows in DB (matching query criteria) not found in Excel based on key columns
      },
      "results": {
        "mismatched_details": [
          {
            "excel_row_index": 5, // Original row number in the Excel file
            "excel_data": { "ID": "P100", "Name": "Old Product Name", "Stock": 10 },
            "db_data": { "product_id": "P100", "Name": "New Product Name", "quantity": 15 },
            "differences": [
              { "column_excel": "Name", "excel_value": "Old Product Name", "column_db": "Name", "db_value": "New Product Name" },
              { "column_excel": "Stock", "excel_value": 10, "column_db": "quantity", "db_value": 15 }
            ]
          }
          // ... more mismatched rows
        ],
        "excel_only_details": [
          {
            "excel_row_index": 10,
            "data": { "ID": "P205", "Name": "Extra Product", "Stock": 30 }
          }
          // ... more rows found only in Excel
        ],
        "db_only_details": [
          {
            "data": { "product_id": "P310", "Name": "DB Only Product", "quantity": 5 }
          }
          // ... more rows found only in the database
        ]
      },
      "errors": [] // Optional: For non-critical processing errors, e.g., unparseable rows
    }
    ```

**Explanation of Response Fields:**

*   **`summary`**: Provides a high-level overview of the comparison.
    *   `total_excel_rows`: Total number of data rows processed from the Excel file.
    *   `total_db_rows_queried`: Total number of database rows that were considered in the comparison (e.g., if `db_table_name` and some filters were applied).
    *   `matched_rows`: Number of rows where the key identifying columns matched and all compared columns were identical.
    *   `mismatched_rows`: Number of rows where key columns matched, but one or more other compared columns had different values.
    *   `excel_only_rows`: Number of rows present in the Excel file but not found in the database based on the specified key columns.
    *   `db_only_rows`: Number of rows present in the database (within the queried scope) but not found in the Excel file based on the specified key columns.

*   **`results`**: Contains detailed lists of the differences.
    *   `mismatched_details`: An array of objects, each representing a row where data exists in both Excel and the DB but has discrepancies.
        *   `excel_row_index`: Helps locate the row in the original Excel file.
        *   `excel_data`: The data from the Excel row.
        *   `db_data`: The corresponding data from the database row.
        *   `differences`: An array detailing each column that differs, showing the Excel column name/value and the DB column name/value.
    *   `excel_only_details`: An array of objects, each representing a row found only in the Excel file.
    *   `db_only_details`: An array of objects, each representing a row found only in the database (that matched any initial query criteria).

*   **`errors`**: An optional array that can list any non-fatal errors encountered during processing, such as issues with specific rows that didn't halt the entire comparison (e.g., data type conversion issues for a particular cell). Fatal errors would typically result in a different HTTP error status code.

This detailed structure allows the client to not only get a summary but also to programmatically access and act upon the specific differences found. The level of detail in `mismatched_details`, `excel_only_details`, and `db_only_details` can be adjusted based on performance considerations and client requirements (e.g., the API might offer a "summary_only" mode).

## Error Responses

Error responses are returned in a consistent JSON format.

**Common HTTP Status Codes:**

*   **`400 Bad Request`**: The request was malformed or invalid. This can occur due to:
    *   Missing `excel_file` in the `multipart/form-data`.
    *   Unparseable JSON in optional parameters like `comparison_columns_excel`.
    *   Invalid file type (e.g., uploading a text file instead of an Excel sheet).
    *   Missing required parameters like `db_table_name` if it's not optional.
    *   Invalid values for parameters (e.g., `header_row_index` being negative).

*   **`401 Unauthorized`**: If the API requires authentication and the client has not provided valid credentials. (Assuming authentication is outside the scope of this basic description for now but important in a real-world scenario).

*   **`403 Forbidden`**: If the authenticated client does not have permission to access the specified resource or perform the comparison (e.g., trying to access a table they are not allowed to).

*   **`404 Not Found`**:
    *   If the endpoint itself (`/compare/excel`) is incorrect.
    *   Potentially, if a specified `db_table_name` does not exist in the database.

*   **`413 Payload Too Large`**: If the uploaded Excel file exceeds a server-defined size limit.

*   **`415 Unsupported Media Type`**: If the provided file is not of a supported Excel format (e.g., `.xlsx`, `.xls`, `.csv`) and the server cannot process it.

*   **`422 Unprocessable Entity`**: If the request is well-formed, but the server cannot process the instructions due to semantic errors. For example:
    *   Specified column names in `comparison_columns_excel` not found in the Excel header.
    *   Specified column names in `comparison_columns_db` not found in the `db_table_name`.
    *   Inconsistent number of columns specified in `comparison_columns_excel` and `comparison_columns_db`.

*   **`500 Internal Server Error`**: A generic error message for unexpected server-side issues (e.g., database connection failure, unhandled exceptions during processing). These errors should ideally be minimized by robust error handling.

*   **`503 Service Unavailable`**: If the server is temporarily overloaded or down for maintenance.

**Error Response Body (JSON):**

A consistent JSON structure should be used for all error responses to help clients parse them programmatically.

```json
{
  "error": {
    "code": "VALIDATION_ERROR", // Or "DB_CONNECTION_ERROR", "FILE_PROCESSING_ERROR", etc.
    "message": "A human-readable error message detailing what went wrong.",
    "details": [ // Optional: for more granular error reporting, especially for 400/422 errors
      {
        "field": "excel_file",
        "issue": "File is required."
      },
      {
        "field": "comparison_columns_excel",
        "issue": "Column 'Product ID' not found in the uploaded Excel file."
      }
    ]
  }
}
```

**Explanation of Error Fields:**

*   **`error.code`**: A machine-readable error code that clients can use to handle specific error types programmatically.
*   **`error.message`**: A human-readable summary of the error. This should be clear and concise.
*   **`error.details`**: (Optional) An array of objects providing more specific information about the error, particularly useful for validation errors where multiple issues might exist in the request.
    *   `field`: The name of the field that caused the error (e.g., `excel_file`, `db_table_name`).
    *   `issue`: A description of the specific problem with that field.

By providing structured error responses, the API becomes more robust and easier for client applications to integrate with and debug.
```
