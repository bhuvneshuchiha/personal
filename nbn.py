import requests
import json


def generate_auth_token():
    """This function is to get the job code from adler
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


def get_job_code(supplier_name):
    """
    the supplier name will come from the icr extraction and that will be pass
        onto this function.
    """
    url = "https://adleridemo.azurewebsites.net/api/public/procurement/supplier"
    auth_token = generate_auth_token()

    payload = json.dumps({"supplier_name": supplier_name})
    headers = {
        "Content-Type": "application/json",
        "Authorization": auth_token,
        "api_key": "HNG37484=",
    }

    response = requests.request("GET", url, headers=headers, data=payload)
    print(response.text)

    if response.text == supplier_name:
        """write the logic to fetch the job code if the supplier name was matched."""

        url = "https://adleridemo.azurewebsites.net/api/public/accounts/supplier_jobs"

        payload = json.dumps({"supplier_name": supplier_name})
        headers = {
            "Content-Type": "application/json",
            "Authorization": auth_token,
            "api_key": "HNG37484=",
        }

        response = requests.request("GET", url, headers=headers, data=payload)

        print(response.text)
    else:
        return "supplier name was not found in adler"
