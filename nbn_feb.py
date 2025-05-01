import requests
from pymongo import MongoClient
import json
import re
from datetime import datetime, timedelta

# TODO: -> May be you can add a key like response["resp_{func_name}"] = function
# passed ? function failed so that in the front end we can have entities to
# know the response of each of these function.
"""
** When defining the workflow, make sure to add the auth token to frontend and
send extraction response with the auth token to our custom api.
File will come for extraction, but for each request check on the basis of the
message thread id that the fields are present in the db for that message id or not.

1- Get auth token
2- Check for valid fields.
3- Check for amounts, convert according to currencies, save the fields in a Mysql
    database with the outlook message ID.
4- By this part you should know if there is an escalation needed or manual
    intervention is needed.
5- Go to the data preparing function(this will have all fields sorted just the
    due date might be needed some work.
6- Send the data to the due date handler function, there the due dates will be
    taken care of and the data will be pushed to adler their only via send
    data to adler function where the api to send the data to adler will be
    integrated and final structure of the response will be made.

Pending points ->
* What to do with the job code?
* Amounts are to be used to cross validate a few things watch the recording.
"""

# def generate_auth_token_adler():
#     """
#     Keep this function in the frontend.
#     This function is to get the job code from adler
#     on this basis of the supplier name.
#     """
#     url = "https://adleridemo.azurewebsites.net/Token"
#     payload = "grant_type=password&username=E42_user&password=vXFAj599Xl%401"
#     headers = {
#         "Content-Type": "application/x-www-form-urlencoded",
#     }
#     response = requests.request("GET", url, headers=headers, data=payload)
#     response = response.json()
#     auth_token = (response.text).get("access_token")
#     return f"Bearer {auth_token}"


def mandatory_field_validation_adler(request):
    """
    This function will receive a dictionary which will contain the response
    from the ICR API.We will then check if the 5 mandatory_fields, need to check
    if these fields are extracted from icr or not.Invoice number, Invoice date,
    net amount, vendor name,amount including taxes.
    """
    invoice_number = request.get("invoice_number", "")
    invoice_date = request.get("invoice_date", "")
    subtotal = request.get("subtotal", "")
    supplier_name = request.get("supplier_name", "")
    total_amount = request.get("total_amount", "")
    auth_token = request.get("auth_token", "")

    missing_fields_dict = {}

    if not invoice_date:
        missing_fields_dict["invoice_date"] = "missing"
    else:
        missing_fields_dict["invoice_date"] = "present"

    if not invoice_number:
        # Send into the exception handling
        missing_fields_dict["invoice_number"] = "missing"
    else:
        missing_fields_dict["invoice_number"] = "present"

    if not subtotal:
        missing_fields_dict["subtotal"] = "missing"
    else:
        missing_fields_dict["subtotal"] = "present"

    if not supplier_name:
        missing_fields_dict["supplier_name"] = "missing"
    else:
        missing_fields_dict["supplier_name"] = "present"

    if not total_amount:
        missing_fields_dict["total_amount"] = "missing"
    else:
        missing_fields_dict["total_amount"] = "present"

    request["missing_fields_dict"] = missing_fields_dict
    # Make sure all the missing fields are grouped together because those
    # missing fields will be going to the email template that will be sent
    # to the email initiator.
    if "missing" in missing_fields_dict.values():
        request["func_resp_mandatory_field_validation"] = "1 or more fields is missing"
    else:
        request["func_resp_mandatory_field_validation"] = "All fields are present"

    return request


def get_currency_and_change_amounts_adler(request):
    subtotal = request.get("subtotal", "")
    total_amount = request.get("total_amount", "")
    auth_token = request.get("auth_token", "")
    currency_code_extracted = ""

    match = re.search(r"([A-Za-z]+)", subtotal)
    if match:
        currency_code_extracted = match.group(1)
        print(f"Alpha part (currency code): {currency_code_extracted}")
    else:
        match = re.search(r"([A-Za-z]+)", total_amount)
        if match:
            currency_code_extracted = match.group(1)
            print(f"Alpha part (currency code): {currency_code_extracted}")

    url = "https://adleridemo.azurewebsites.net/api/public/accounts/currency"
    payload = {}
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {auth_token}",
        "api_key": "HNG37484=",
    }
    exchange_rate = ""
    response = requests.request("GET", url, headers=headers, data=payload)
    currency_list = response.json().get("currency", "")
    for item in currency_list:
        for _, _ in item.items():
            if item["currency_code"] == str(currency_code_extracted):
                exchange_rate = float(item["exchange_rate"])

    if exchange_rate != "":
        subtotal = subtotal * exchange_rate
        total_amount = subtotal * exchange_rate
    request["subtotal"] = subtotal
    request["total_amount"] = total_amount
    if request.get("subtotal") and request.get("total_amount"):
        request["func_resp_get_currency_and_change_amounts_adler"] = "Success"
    else:
        request["func_resp_get_currency_and_change_amounts_adler"] = "Failed"
    return request


def prepare_and_send_data_to_adler(request):
    resp = {}
    print("Inside the function to prepare and send data")
    # Make sure to correct the structure of the object body and safenames before
    # sending the data.
    return resp


def get_job_code_info_adler(request):
    """
    Handles job code retrieval and processing for Adler ERP.

    Logic:
    1. Single job code with a single service type: Send as-is to Adler.
    2. Single job code with multiple service types: Match total amount
       extracted with the total for that service type and send if matched.
    3. Multiple service types match the invoice amount: Post with all matching
       service types.
    4. Invoice amount is less than total service type amounts: Find combinations
       of service types that sum up to the invoice amount.
    5. No matching combination or invoice amount exceeds all sums: Send for
       manual intervention.
    6. Multiple job codes for a supplier: Send for manual intervention.
    """

    supplier_name = request.get("supplier_name", "")
    auth_token = request.get("auth_token", "")
    invoice_amount = float(request.get("invoice_amount", 0))

    url = "https://adleridemo.azurewebsites.net/api/public/accounts/supplier_jobs"
    payload = json.dumps({"supplier_name": supplier_name})
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {auth_token}",
        "api_key": "HNG37484=",
    }

    response = requests.request("GET", url, headers=headers, data=payload)
    if response.status_code != 200:
        request["func_resp_get_job_code_info_adler"] = "No data found"
        return request

    supplier_jobs = response.json().get("supplier_jobs", [])

    if not supplier_jobs:
        request["func_resp_get_job_code_info_adler"] = "No jobs found for the supplier."
        return request

    # Grouping jobs by job_code
    job_code_groups = {}
    for job in supplier_jobs:
        job_code = job["job_code"]
        amount = float(job["amount"])
        service_type = job["service_type"]

        if job_code not in job_code_groups:
            job_code_groups[job_code] = {
                "service_types": [],
                "total_amount": 0,
            }
        job_code_groups[job_code]["service_types"].append(service_type)
        job_code_groups[job_code]["total_amount"] += amount

    # Check for multiple job codes for a supplier
    if len(job_code_groups) > 1:
        request["func_resp_get_job_code_info_adler"] = (
            "Manual intervention required: Multiple job codes."
        )
        return request

    # Single job code logic
    single_job_code, job_data = next(iter(job_code_groups.items()))
    total_amount = job_data["total_amount"]
    service_types = job_data["service_types"]

    if len(service_types) == 1:
        # Single service type
        if total_amount == invoice_amount:
            request["func_resp_get_job_code_info_adler"] = {
                "job_code": single_job_code,
                "service_types": service_types,
            }
            return request
        else:
            request["func_resp_get_job_code_info_adler"] = (
                "Manual intervention required: Amount mismatch."
            )
            return request

    # Multiple service types logic
    if total_amount == invoice_amount:
        # Exact match with multiple service types
        request["func_resp_get_job_code_info_adler"] = {
            "job_code": single_job_code,
            "service_types": service_types,
        }
        return request
    elif total_amount > invoice_amount:
        # Find combinations of service types matching the invoice amount
        from itertools import combinations

        possible_matches = []
        for r in range(1, len(service_types) + 1):
            for combo in combinations(service_types, r):
                combo_amount = sum(
                    float(job["amount"])
                    for job in supplier_jobs
                    if job["service_type"] in combo
                )
                if combo_amount == invoice_amount:
                    possible_matches.append(combo)

        if possible_matches:
            request["func_resp_get_job_code_info_adler"] = {
                "job_code": single_job_code,
                "matching_combinations": possible_matches,
            }
            return request

    # If no match found or invoice amount exceeds total
    request["func_resp_get_job_code_info_adler"] = "Manual intervention required."
    return request


def get_supplier_info_adler(request):
    """
    Send the supplier name to this api to get the payment terms
    """
    url = "https://adleridemo.azurewebsites.net/api/public/procurement/supplier"
    auth_token = request.get("auth_token", "")

    payload = json.dumps({"supplier_name": request.get("supplier_name", "")})
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {auth_token}",
        "api_key": "HNG37484=",
    }

    response = requests.request("GET", url, headers=headers, data=payload)
    response = response.json()
    final_output = response.get("supplier", [])
    # The supplier name the needs to be posted in the adler erp needs to come
    # from the supplier master.
    if final_output:
        request["func_resp_get_supplier_info_adler"] = final_output
    else:
        request["func_resp_get_supplier_info_adler"] = "Supplier info not found"
    return request


def handle_due_date_send_data_adler(request):
    # THIS SHOULD BE THE LAST FUNCTION TO GET CALLED BECAUSE THIS WILL SEND THE
    # DATA BACK TO ADLER.....

    # Request dict should contain all the mandatory fields
    # Keep this function at last in your code because on the basis of due date
    # we will send the call to push data to adler
    """If the due date was extracted from the invoice then pass it to adler as
    it is else we need to calculate the due date on the basis of some
    logic present in the DRG

    **Check if two due dates are extracted from the invoice, if that's the case
    then the invoice data needs to be send twice with both the due dates one by one.
    Call the send_data_to_adler function twice for both the dates and push the
    data along with it.

    **If no due date is mentioned on the invoice then fetch the payment terms
    from the supplier master.
    a) If payment term is 100% advance then due date == invoice date.
    b) If payment term is in future eg. 30 days then due date == invoice date + 30
    {FYI if you add any number to the date make sure to change the month and year
    if it gets to the next month or the year.}
    """

    invoice_date = request.get("invoice_date", "")
    due_date = request.get("due_date", "")
    date_format = "%dd/%mm/%YYYY"
    due_date_obj = datetime.strptime(due_date, date_format)
    final_due_date = ""

    # If 1 due date was extracted from icr
    due_date = request.get("due_date", "")
    # Dont use date_obj to split as split works only on strings.
    if not len(due_date.split(",")) or not len(due_date.split()):
        request["due_date"] = due_date_obj
        prepare_and_send_data_to_adler(request)

    # If no due date was extracted
    # Need to have the exact wordings of the payment terms for string matching
    # What all types of payment terms will be there?
    elif request.get("due_date", "") == "":
        supplier_info = get_supplier_info_adler(request)
        payment_terms = ""
        for item in supplier_info:
            payment_terms = item.get("payement_terms")
            if payment_terms == "100% ADVANCE":
                due_date = invoice_date
            elif payment_terms == "100% upon completion":
                pass
                # Manual Intervention
            elif "days" in payment_terms:
                match = re.search(r"(\d+)\s+days", payment_terms)
                if match:
                    final_due_date = match.group(1)

        if final_due_date:
            days_to_add = int(final_due_date)
            date_obj = datetime.strptime(final_due_date, "%dd/%mm/%YYYY")
            final_due_date_obj = date_obj + timedelta(days=days_to_add)
            final_due_date = final_due_date_obj.strftime("%dd/%mm/%YYYY")

    # More than 1 due date extracted
    date_list = []
    # List to store multiple response for API calls incase when more than 1
    # due date was extracted.
    final_response = []
    payload = request
    # What's the separator symbol for multiple due dates comma or space?
    if len(due_date.split(",")) > 1:
        for item in due_date.split(","):
            item_date_obj = datetime.strptime(item, date_format)
            date_list.append(item_date_obj)
        """Now the date_list has multiple dates that we need to send to the
        payload one by one and send the call to push the data for each date"""

        for i in range(len(date_list)):
            payload["due_date"] = date_list[i]
            final_response.append(prepare_and_send_data_to_adler(payload))
        return final_response
    else:
        final_response.append(prepare_and_send_data_to_adler(payload))
    return final_response



def check_update_data_mongo_adler(request_data):
    MONGO_URL = "mongodb://localhost:27017/"
    client = MongoClient(MONGO_URL)
    db = client.nbn
    collection = db.nbn_details

    required_fields = [
        "invoice_number", "invoice_date", "subtotal", "supplier_name",
        "total_amount", "auth_token", "message_id", "sender_email", "email_body"
    ]

    req_data = {field: request_data.get(field, "") for field in required_fields}
    message_id = req_data["message_id"]

    if not message_id:
        return {"error": "message_id is required"}

    existing_doc = collection.find_one({"message_id": message_id})

    if existing_doc:
        missing_fields = [field for field in required_fields if field not in existing_doc]
        if not missing_fields:
            return {"status": "Duplicate entry"}
        update_data = {field: req_data[field] for field in missing_fields if req_data[field]}

        if update_data:
            collection.update_one({"message_id": message_id}, {"$set": update_data})
            return {"status": "Updated missing fields", "updated_fields": list(update_data.keys())}
        else:
            return {"status": "No new data to update"}

    else:
        missing_fields = [field for field in required_fields if field not in request_data]
        collection.insert_one(req_data)
        if missing_fields:
            return {"status": "Inserted with missing fields", "missing_fields": missing_fields}
        else:
            return {"status": "Inserted successfully"}



def read_email_body_adler(request):
    pass


# MAIN DJANGO ENTRY POINT VIEW
def main_api_view_adler(request):
    try:
        message = "success"
        status = 200
        request_data = json.loads(request.body.decode(encoding="utf-8"))
        resp_data = {}
        req_data = {}
        req_data["invoice_number"] = request_data.get("invoice_number", "")
        req_data["invoice_date"] = request_data.get("invoice_date", "")
        req_data["subtotal"] = request_data.get("subtotal", "")
        req_data["supplier_name"] = request_data.get("supplier_name", "")
        req_data["total_amount"] = request_data.get("total_amount", "")
        req_data["auth_token"] = request_data.get("auth_token", "")
        req_data["message_id"] = request_data.get("message_id", "")
        req_data["sender_email"] = request_data.get("sender_email", "")
        req_data["email_body"] = request_data.get("email_body", "")

        missing_list = []
        missing_dict = {}
        missing_fields_data = mandatory_field_validation_adler(req_data)
        missing_fields = missing_fields_data.get("missing_fields_dict", "")
        resp_data["mandatory_field_validation"] = missing_fields_data.get(
            "func_resp_mandatory_field_validation"
        )

        # Check if missing fields are present
        for k, _ in missing_fields.items():
            if missing_fields[k] == "missing":
                missing_list.append(k)

        if len(missing_list) != 0:
            """
            Store whatever fields are present in the database and return this
            function so that we can know what fields are missing in the frontend
            to be filled in the email template.
            """

            response_from_mongo = check_update_data_mongo_adler(req_data)
            # This means parameters are missing
            # Create and send template for the missing field
            print(missing_list)
            resp_data["missing_fields"] = missing_list
            missing_dict["missing_fields"] = missing_list
            # You need to return the missing dict to the front end so that
            # these fields can be used to create template.
            content = {"data": resp_data, "message": message}
            return make_response(content, status)

        if len(missing_list) == 0:
            """
            Means all the mandatory fields are present.
            Store the fields in the DB with the message ID and continue the process.
            and push the data to adler.
            (Happy path)
            """
            response_from_mongo = check_update_data_mongo_adler(req_data)

            req_data.pop("missing_fields_dict")
            resp_data["missing_fields"] = missing_list
            # All keys are present
            req_data = get_currency_and_change_amounts_adler(req_data)
            resp_data["currency_and_change_amounts"] = req_data.get(
                "func_resp_get_currency_and_change_amounts_adler"
            )
            supplier_info = get_supplier_info_adler(req_data)
            resp_data["func_resp_get_supplier_info_adler"] = req_data.get(
                "func_resp_get_supplier_info_adler"
            )
            if supplier_info:
                job_code_info = get_job_code_info_adler(req_data)
                resp_data["func_resp_get_job_code_info_adler"] = req_data.get(
                    "func_resp_get_job_code_info_adler"
                )
                print(job_code_info)
            """
            Check if the same message id is stored in the database.
            Update the missing fields with the ones recieved and continue on
            with the process.
            """

            # Write the part to handle the amount combinations which are used for calc.

            # final_api_response_to_send_data = handle_due_date_send_data_adler(
            #    req_data
            # )
    except Exception as e:
        status = 400
        message = str(e)
        logger.error("Exception is", e)
    content = {"data": resp_data, "message": message}
    return make_response(content, status)
