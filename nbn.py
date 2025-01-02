import json
import requests


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
    auth_token = response.text
    return f"Bearer {auth_token}"


def mandatory_field_validation_adler(request):
    """
    This function will receive a dictionary which will contain the response
    from the ICR API.We will then check if the 5 mandatory_fields, need to check
    if these fields are extracted from icr or not.Invoice number, Invoice date,
    net amount, vendor name,amount including taxes.
    """
    return


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

    response_list = response.text.get("supplier", "")
    for item in response_list:
        for key, val in item.items():
            if key == "supplier_name" and val == supplier_name:
                """write the logic to fetch the job code if the supplier name was matched."""

                url = "https://adleridemo.azurewebsites.net/api/public/accounts/supplier_jobs"

                payload = json.dumps({"supplier_name": supplier_name})
                headers = {
                    "Content-Type": "application/json",
                    "Authorization": auth_token,
                    "api_key": "HNG37484=",
                }

                response = requests.request(
                    "GET", url, headers=headers, data=payload)

                print(response.text)
                supplier_jobs = response.text.get("supplier_jobs", [])
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
    final_output = response.get("supplier", [])
    return final_output  # List of dict


def handle_due_date_adler(request):
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

    # If 1 due date was extracted from icr
    due_date = request.get("due_date", "")
    if not len(due_date.split(",")):
        send_data_to_adler(request)

    # If no due date was extracted
    # Need to have the exact wordings of the payment terms for string matching
    # What all types of payment terms will be there?
    elif request.get("due_date", "") == "":
        supplier_info = get_supplier_info_adler(
            request.get("supplier_name", ""))
        payment_terms = ""
        for item in supplier_info:
            payment_terms = item.get("payement_terms")
            if payment_terms == "100% ADVANCE":
                due_date = invoice_date
                # Need to complete this once get the clarification

    # More than 1 due date extracted
    date_list = []
    if len(due_date.split(",")) > 1:
        for item in due_date.split(","):
            date_list.append(item)
        """Now the date_list has multiple dates that we need to send to the payload
        one by one and send the call to push the data for each date"""

        for i in range(len(date_list)):
            request["due_date"] = date_list[i]
            payload = request
            send_data_to_adler(payload)
    else:
        send_data_to_adler(payload)
    return
