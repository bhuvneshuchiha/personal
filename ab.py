import requests
import time
import os

# API details
url = "https://geniex02.mahindra.com/api/route/api/generic/epof/740/train"
payload = {
    'flag': 'login',
    'query': 'hi',
    'callback': 'https://geniex02.mahindra.com/',
    'user_id': '1'
}
headers = {}

# Folder with PDF files
folder_path = "/Users/bhuvnesh/Downloads/Mahindra_Epod"

# Iterate over files in the folder
for filename in os.listdir(folder_path):
    if filename.lower().endswith(".pdf"):
        file_path = os.path.join(folder_path, filename)
        with open(file_path, 'rb') as f:
            files = [
                ('file_attachment', (filename, f, 'application/pdf'))
            ]

            response = requests.post(url, headers=headers, data=payload, files=files)

            print(f"Sent: {filename}, Response Code: {response.status_code}")
            print(response.text)

        # Wait for 15 seconds before the next request
        time.sleep(15)

