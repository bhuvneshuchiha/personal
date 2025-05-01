import re

def read_email_body_adler(email_body):
    field_names = [
        "invoice_number", "invoice_date", "subtotal", "supplier_name", "total_amount"
    ]
    extracted_data = {}

    # Create variations of field names to match against
    for field in field_names:
        value = None

        # Generate different variations of the field name
        field_variations = [
            field,                          # invoice_number
            field.replace('_', ' '),        # invoice number
            field.replace('_', ''),         # invoicenumber
            field.upper(),                  # INVOICE_NUMBER
            field.replace('_', ' ').upper() # INVOICE NUMBER
        ]

        # Create a regex pattern that matches any of the variations
        pattern_parts = []
        for variation in field_variations:
            # Escape special regex characters if any
            escaped_variation = re.escape(variation)
            pattern_parts.append(escaped_variation)

        # Join all variations with OR operator
        combined_pattern = '|'.join(pattern_parts)

        # Final pattern matches any variation followed by optional spaces, colon, and captures the value
        pattern = rf"(?:{combined_pattern})\s*:?\s*(.*?)(?:\n|$)"

        match = re.search(pattern, email_body, re.IGNORECASE | re.MULTILINE)
        if match:
            value = match.group(1).strip()

            # Additional cleaning: remove any trailing colons or excess whitespace
            value = re.sub(r':+\s*$', '', value)

        extracted_data[field] = value

    return extracted_data

# Example usage
email_body = """
Hello,
Here are the invoice details:
INVOICE NUMBER: 123ABC
Invoice Date: 10th Feb 2025
SubTotal: 1000.00
supplierName: ABC Supplies Ltd.
Total_Amount: 1200.00
Best regards,
Company XYZ
"""

result = read_email_body_adler(email_body)
print(result)
# import re
#
# def extract_fields_from_html(email_body):
#     # Define a dictionary to hold the field names and values
#     field_names = [
#         "invoice_number", "invoice_date", "subtotal", "supplier_name", "total_amount"
#     ]
#     extracted_data = {}
#
#     # Create a regex pattern for each field
#     for field in field_names:
#         # Convert field name to the format it appears in the email (replace underscore with space)
#         field_display = field.replace('_', ' ')
#
#         # The regex pattern now handles:
#         # - Case insensitive matching
#         # - Optional spaces before and after the colon
#         # - Captures everything until the next line break or end of string
#         pattern = rf"{field_display}\s*:\s*(.*?)(?:\n|$)"
#
#         match = re.search(pattern, email_body, re.IGNORECASE | re.MULTILINE)
#         if match:
#             # Store the extracted value, removing any trailing/leading whitespace
#             extracted_data[field] = match.group(1).strip()
#         else:
#             # If not found, store None
#             extracted_data[field] = None
#
#     return extracted_data
#
# # Example usage
# email_body = """
# Hello,
# Please find the details of the invoice below:
# Invoice number: 123ABC
# Invoice date: 10th Feb 2025
# Subtotal: 1000.00
# Supplier name: ABC Supplies Ltd.
# Total amount: 1200.00
# Best regards,
# Company XYZ
# """
#
# extracted_data = extract_fields_from_html(email_body)
# print(extracted_data)
