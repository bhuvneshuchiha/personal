import json
import requests
import re
from datetime import datetime, timedelta

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


def generate_auth_token_adler():
    """
    This function is to get the job code from adler
    on this basis of the supplier name.
    """
    url = "https://adleridemo.azurewebsites.net/Token"
    payload = "grant_type=password&username=E42_user&password=vXFAj599Xl%401"
    headers = {
        "Content-Type": "application/x-www-form-urlencoded",
    }
    response = requests.request("GET", url, headers=headers, data=payload)
    response = response.json()
    auth_token = (response.text).get("access_token")
    return f"Bearer {auth_token}"


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
    if (
        not invoice_date
        and invoice_number
        and subtotal
        and supplier_name
        and total_amount
    ):
        handle_due_date_send_data_adler(request)
    if not invoice_number:
        # Send into the exception handling
        pass
    if not subtotal:
        # what to do?
        pass
    if not supplier_name:
        # what to do ?
        pass
    if not total_amount:
        # what to do?
        pass

    missing_fields_dict = {}

    # Make sure all the missing fields are grouped together because those
    # missing fields will be going to the email template that will be sent
    # to the email initiator.
    print(missing_fields_dict)
    return


def get_currency_and_change_amounts_adler(request):
    subtotal = request.get("subtotal", "")
    total_amount = request.get("total_amount", "")
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
    auth_token = generate_auth_token_adler()
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

    subtotal = subtotal * exchange_rate
    total_amount = subtotal * exchange_rate
    request["subtotal"] = subtotal
    request["total_amount"] = total_amount
    return request


def send_data_to_adler(request):
    """
    **Check if two due dates are extracted from the invoice, if that's the case
    then the invoice data needs to be send twice with both the due dates one by one.
    """
    return


def get_job_code_info_adler(supplier_name):
    """
    the supplier name will come from the icr extraction and that will be pass
    onto this function.
    """
    url = "https://adleridemo.azurewebsites.net/api/public/procurement/supplier"
    auth_token = generate_auth_token_adler()

    payload = json.dumps({"supplier_name": supplier_name})
    headers = {
        "Content-Type": "application/json",
        "Authorization": auth_token,
        "api_key": "HNG37484=",
    }

    response = requests.request("GET", url, headers=headers, data=payload)
    print(response.text)
    response = response.json()

    response_list = response.text.get("supplier", "")
    for item in response_list:
        for key, val in item.items():
            if key == "supplier_name" and val == supplier_name:
                """write the logic to fetch the job code if the supplier name
                was matched."""

                url = "https://adleridemo.azurewebsites.net/api/public/\
                    accounts/supplier_jobs"

                payload = json.dumps({"supplier_name": supplier_name})
                headers = {
                    "Content-Type": "application/json",
                    "Authorization": auth_token,
                    "api_key": "HNG37484=",
                }

                response = requests.request("GET", url, headers=headers, data=payload)

                print(response.text)
                response = response.json()
                supplier_jobs = response.get("supplier_jobs", [])
                return supplier_jobs
            else:
                return "supplier details was not found in adler"


def get_supplier_info_adler(request):
    """
    Send the supplier name to this api to get the payment terms
    """
    url = "https://adleridemo.azurewebsites.net/api/public/procurement/supplier"
    auth_token = generate_auth_token_adler()

    payload = json.dumps({"supplier_name": request.get("supplier_name", "")})
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {auth_token}",
        "api_key": "HNG37484=",
    }

    response = requests.request("GET", url, headers=headers, data=payload)
    response = response.json()
    final_output = response.get("supplier", [])
    return final_output  # List of dict


def prepare_data_adler(request):
    """
    This function will be called inside the handle due date function just to
    make sure all the parameters are there before we use the send data func
    to finally send the data to adler.
    """
    return


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
        send_data_to_adler(request)

    # If no due date was extracted
    # Need to have the exact wordings of the payment terms for string matching
    # What all types of payment terms will be there?
    elif request.get("due_date", "") == "":
        supplier_info = get_supplier_info_adler(request.get("supplier_name", ""))
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
            send_data_to_adler(payload)
    else:
        send_data_to_adler(payload)
    return
